package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirBlinov/AuthService/pkg/authservice"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/handler"
	store2 "github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/store"

	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/model"
	authservicefake "github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/pkg/authservice/fake"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/service"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/store/teststore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

var sessManager authservice.AuthServiceClient
var testSessionID string = "sessionid"

func TestServerHandleSignOut(t *testing.T) {
	store := teststore.New()
	services := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(services, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		context            *model.User
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:               "valid",
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:               "not valid cookie",
			context:            u,
			serviceCookieValue: "",
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/private/signout", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotEqual(t, tc.coockieValue, rec.Result().Header["Set-Cookie"])
		})
	}
}

func TestServer_HandleProductCreate(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		payload            interface{}
		context            *model.User
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"product_name":    "product",
				"category_id":     105,
				"pieces_in_pack":  1,
				"material_id":     1,
				"weight":          500,
				"lenght":          500,
				"width":           300,
				"height":          20,
				"description":     "descript",
				"wildberries_sku": 1234,
				"ozon_sku":        1234567,
			},
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "valid_empty_sku",
			payload: map[string]interface{}{
				"product_name":    "product",
				"category_id":     105,
				"material_id":     1,
				"wildberries_sku": 0,
				"ozon_sku":        0,
			},
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},

		{
			name: "valid_minimum_params",
			payload: map[string]interface{}{
				"product_name": "product",
				"category_id":  105,
				"material_id":  1,
			},
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid_less_params",
			payload: map[string]interface{}{
				"product_name": "product",
				"material_id":  1,
			},
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/private/product/product", b)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleProductGetProduct(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	p := model.TestProduct(t)
	p.UserID = u.ID
	mpiList := &model.MarketPlaceItemsList{}
	mpiList.GetMPIList(p)
	store.Product().Create(p, mpiList)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		context            *model.User
		productId          interface{}
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:               "valid",
			context:            u,
			productId:          p.ProductID,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:               "invalid req",
			context:            u,
			productId:          "",
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:               "not found",
			productId:          p.ProductID + 1,
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/private/product/product/%v", tc.productId), nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotEqual(t, 0, rec.Result().ContentLength)
		})
	}
}

func TestServer_HandleDeleteProduct(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	p := model.TestProduct(t)
	p.UserID = u.ID
	mpiList := &model.MarketPlaceItemsList{}
	mpiList.GetMPIList(p)
	store.Product().Create(p, mpiList)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		context            *model.User
		productId          interface{}
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:               "valid",
			context:            u,
			productId:          p.ProductID,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:               "invalid",
			context:            u,
			productId:          "",
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:               "invalid_id",
			context:            u,
			productId:          "a",
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(`/api/v1/private/product/product/%v`, tc.productId), nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			pId, ok := tc.productId.(int)
			if ok {
				p1, err := store.Product().GetProductById(pId)
				assert.Error(t, store2.ErrRecordNotFound, err)
				assert.Nil(t, p1)
			}
		})
	}
}

func TestServer_HandleProductFindByUserId(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	p1 := model.TestProduct(t)
	p1.UserID = u.ID
	mpi1 := &model.MarketPlaceItemsList{}
	mpi1.GetMPIList(p1)
	store.Product().Create(p1, mpi1)

	p2 := model.TestProductWOSKU(t)
	p2.UserID = u.ID
	mpi2 := &model.MarketPlaceItemsList{}
	mpi2.GetMPIList(p2)
	store.Product().Create(p2, mpi2)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		context            *model.User
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:               "valid",
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/private/product/product", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotEqual(t, 0, rec.Result().ContentLength)
		})
	}
}

func TestServer_HandleProductUpdate(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)

	u := model.TestUser(t)
	store.User().Create(u)

	p := model.TestProduct(t)
	p.UserID = u.ID
	mpiList := &model.MarketPlaceItemsList{}
	mpiList.GetMPIList(p)
	store.Product().Create(p, mpiList)

	p1 := model.TestProduct(t)
	p1.UserID = u.ID
	p1.OzonSKU = 0
	p1.WildberriesSKU = 0
	mpiList1 := &model.MarketPlaceItemsList{}
	mpiList1.GetMPIList(p1)
	store.Product().Create(p, mpiList1)

	p.Description = "new description"

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		context            *model.User
		payload            interface{}
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:    "valid",
			context: u,
			payload: map[string]interface{}{
				"product_id":      p.ProductID,
				"product_name":    p.ProductName,
				"category_id":     p.CategoryID,
				"pieces_in_pack":  p.PiecesInPack,
				"material_id":     p.MaterialID,
				"weight":          p.Weight,
				"lenght":          p.Lenght,
				"width":           p.Width,
				"height":          p.Height,
				"description":     p.Description,
				"wildberries_sku": 2222222222,
				"ozon_sku":        1111111111,
			},
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:    "valid_empty_sku",
			context: u,
			payload: map[string]interface{}{
				"product_id":      p1.ProductID,
				"product_name":    p1.ProductName,
				"category_id":     p1.CategoryID,
				"pieces_in_pack":  p1.PiecesInPack,
				"material_id":     p1.MaterialID,
				"weight":          p1.Weight,
				"lenght":          p1.Lenght,
				"width":           p1.Width,
				"height":          p1.Height,
				"description":     p1.Description,
				"wildberries_sku": 2222222222,
				"ozon_sku":        1111111111,
			},
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/private/product/product/%v", p.ProductID), b)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleProductGetCategories(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	c1 := model.TestCategory(t)
	c2 := model.TestCategory(t)
	store.Product().CreateCategory(c1)
	store.Product().CreateCategory(c2)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		context            *model.User
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:               "valid",
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/private/product/category/get_categories", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotEqual(t, 0, rec.Result().ContentLength)
		})
	}
}

func TestServer_HandleProductGetMaterials(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	m := model.TestMaterial(t)
	m.MaterialName = "Дерево"
	store.Product().CreateMaterial(m)

	m1 := model.TestMaterial(t)
	m1.MaterialName = "Пластик"
	store.Product().CreateMaterial(m1)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name               string
		context            *model.User
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:               "valid",
			context:            u,
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/private/product/material/get_materials", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotEqual(t, 0, rec.Result().ContentLength)
		})
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	sessManager := authservicefake.NewAuthServiceClientFake()
	sessionS, _ := sessManager.Create(context.Background(), &authservice.Session{
		UserID: int32(u.ID),
	})

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey), sessManager)
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)
	handl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		name               string
		serviceCookieValue string
		coockieValue       map[interface{}]interface{}
		expectedCode       int
	}{
		{
			name:               "authenticated",
			serviceCookieValue: sessionS.ID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:               "not authenticated",
			serviceCookieValue: "",
			coockieValue:       nil,
			expectedCode:       http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			req.Header.Add("Cookie", fmt.Sprintf("%s=%s", handler.SessionIDKey, tc.serviceCookieValue))
			handlers.AuthenticateUser(handl).ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleRegister(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore([]byte("secret_key")), sessManager)
	handlers.InitHandler()
	testCases := []struct {
		name               string
		serviceCookieValue string
		payload            interface{}
		expectedCode       int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", b)
			handlers.Router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleSignIn(t *testing.T) {
	u := model.TestUser(t)
	store := teststore.New()
	srvc := service.NewService(store)
	store.User().Create(u)

	sessManager := authservicefake.NewAuthServiceClientFake()

	handlers := handler.NewHandler(srvc, sessions.NewCookieStore([]byte("secret_key")), sessManager)
	handlers.InitHandler()
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/signin", b)
			req.Header.Set("Origin", "http://localhost:3000")
			handlers.Router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

			cookieSessionID := rec.Result().Cookies()
			for _, c := range cookieSessionID {
				if c.Name == handler.SessionIDKey {
					sessManager.Delete(context.Background(), &authservice.SessionID{
						ID: c.Value,
					})
				}
			}
		})
	}
}

package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/VladimirBlinov/MarketPlace/internal/handler"
	store2 "github.com/VladimirBlinov/MarketPlace/internal/store"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/service"
	"github.com/VladimirBlinov/MarketPlace/internal/store/teststore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestServerHandleSignOut(t *testing.T) {
	store := teststore.New()
	services := service.NewService(store)
	u := model.TestUser(t)
	store.User().Create(u)

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(services, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		context      *model.User
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:    "valid",
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/private/signout", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
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

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		payload      interface{}
		context      *model.User
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"product_name":    "product",
				"category_id":     "105",
				"pieces_in_pack":  "1",
				"material_id":     "1",
				"weight":          "500",
				"lenght":          "500",
				"width":           "300",
				"height":          "20",
				"description":     "descript",
				"wildberries_sku": "1234",
				"ozon_sku":        "1234567",
			},
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "valid_empty_sku",
			payload: map[string]string{
				"product_name":    "product",
				"category_id":     "105",
				"material_id":     "1",
				"wildberries_sku": "0",
				"ozon_sku":        "0",
			},
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},

		{
			name: "valid_minimum_params",
			payload: map[string]string{
				"product_name": "product",
				"category_id":  "105",
				"material_id":  "1",
			},
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid_less_params",
			payload: map[string]string{
				"product_name": "product",
				"material_id":  "1",
			},
			context: u,
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
			req, _ := http.NewRequest(http.MethodPost, "/private/product/product", b)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
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

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		context      *model.User
		productId    interface{}
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:      "valid",
			context:   u,
			productId: p.ProductID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:      "invalid req",
			context:   u,
			productId: "",
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:      "not found",
			productId: p.ProductID + 1,
			context:   u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/private/product/product/%v", tc.productId), nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
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

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		context      *model.User
		productId    interface{}
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:      "valid",
			context:   u,
			productId: p.ProductID,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:      "invalid",
			context:   u,
			productId: "",
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:      "invalid_id",
			context:   u,
			productId: "a",
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(`/private/product/product/%v`, tc.productId), nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
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

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		context      *model.User
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:    "valid",
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/private/product/product", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotEqual(t, 0, rec.Result().ContentLength)
		})
	}
}

func TestServer_HandleProductUpdate(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

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

	testCases := []struct {
		name         string
		context      *model.User
		payload      interface{}
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:    "valid",
			context: u,
			payload: map[string]interface{}{
				"product_id":      p.ProductID,
				"product_name":    p.ProductName,
				"category_id":     strconv.Itoa(p.CategoryID),
				"pieces_in_pack":  strconv.Itoa(p.PiecesInPack),
				"material_id":     strconv.Itoa(p.MaterialID),
				"weight":          fmt.Sprintf("%f", p.Weight),
				"lenght":          fmt.Sprintf("%f", p.Lenght),
				"width":           fmt.Sprintf("%f", p.Width),
				"height":          fmt.Sprintf("%f", p.Height),
				"description":     p.Description,
				"wildberries_sku": "2222222222",
				"ozon_sku":        "1111111111",
			},
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
				"category_id":     strconv.Itoa(p1.CategoryID),
				"pieces_in_pack":  strconv.Itoa(p1.PiecesInPack),
				"material_id":     strconv.Itoa(p1.MaterialID),
				"weight":          fmt.Sprintf("%f", p1.Weight),
				"lenght":          fmt.Sprintf("%f", p1.Lenght),
				"width":           fmt.Sprintf("%f", p1.Width),
				"height":          fmt.Sprintf("%f", p1.Height),
				"description":     p1.Description,
				"wildberries_sku": "2222222222",
				"ozon_sku":        "1111111111",
			},
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
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/private/product/product/%v", p.ProductID), b)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
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

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		context      *model.User
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:    "valid",
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/private/product/category/get_categories", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
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

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		context      *model.User
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:    "valid",
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/private/product/material/get_materials", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			s.ServeHTTP(rec, req.WithContext(ctx))
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

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	s := newServer(*handlers)
	sc := securecookie.New(secretKey, nil)
	handl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		name         string
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			coockieValue: nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			s.handler.AuthenticateUser(handl).ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleRegister(t *testing.T) {
	store := teststore.New()
	srvc := service.NewService(store)
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore([]byte("secret_key")))
	s := newServer(*handlers)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
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
			req, _ := http.NewRequest(http.MethodPost, "/register", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleSignIn(t *testing.T) {
	u := model.TestUser(t)
	store := teststore.New()
	srvc := service.NewService(store)
	store.User().Create(u)

	handlers := handler.NewHandler(srvc, sessions.NewCookieStore([]byte("secret_key")))
	s := newServer(*handlers)
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
			req, _ := http.NewRequest(http.MethodPost, "/signin", b)
			req.Header.Set("Origin", "http://localhost:3000")
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}

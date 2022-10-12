package handler

import (
	"encoding/json"
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) handleProductOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.respond(w, r, http.StatusOK, nil)
	}
}

func (h *Handler) handleProductCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//req := &service.InputProduct{}
		req := &model.Product{}
		req.UserID = r.Context().Value(CtxKeyUser).(*model.User).ID
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.service.ProductService.CreateProduct(req); err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusCreated, nil)
	}
}

func (h *Handler) handleProductUpdate() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqVars := mux.Vars(r)
		productId, err := strconv.Atoi(reqVars["id"])
		if err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &model.Product{}
		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		req.UserID = r.Context().Value(CtxKeyUser).(*model.User).ID

		if err = h.service.ProductService.UpdateProduct(productId, req); err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusCreated, nil)
	}
}

func (h *Handler) handleProductGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqVars := mux.Vars(r)
		productId, err := strconv.Atoi(reqVars["id"])
		if err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		product, err := h.service.ProductService.GetProductById(productId)
		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusOK, product)
	}
}

func (h *Handler) handleProductDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqVars := mux.Vars(r)
		productId, err := strconv.Atoi(reqVars["id"])
		if err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(CtxKeyUser).(*model.User)

		if err = h.service.ProductService.DeleteProduct(productId, u.ID); err != nil {
			h.error(w, r, http.StatusInternalServerError, nil)
			return
		}

		h.respond(w, r, http.StatusOK, nil)
	}
}

func (h *Handler) handleProductList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(CtxKeyUser).(*model.User)

		products, err := h.service.ProductService.GetProductsByUserId(u.ID)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, products)
	}
}

func (h *Handler) handleProductCategoryGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := h.service.ProductService.GetProductCategories()
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, categories)
	}
}

func (h *Handler) handleProductMaterialGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		materials, err := h.service.ProductService.GetProductMaterials()
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, materials)
	}
}

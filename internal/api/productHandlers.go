package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"bybarcode/internal/products"
)

func (h *handlers) findProductByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode, err := strconv.ParseInt(chi.URLParam(r, "barcode"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	product, err := h.db.FindProductByBarcode(r.Context(), barcode)
	if errors.As(err, &pgx.ErrNoRows) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusNotFound, []byte("product not found"))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	b, err := product.Encode()
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.response.json(w, http.StatusOK, b)
}

func (h *handlers) addProduct(w http.ResponseWriter, r *http.Request) {
	var p products.Product
	if err := p.Decode(r.Body); err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte("bad request"))
	}

	productId, err := h.db.CreateProduct(r.Context(), p)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	p.ID = productId
	b, err := p.Encode()
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.response.json(w, http.StatusOK, b)
}

func (h *handlers) updateProduct(w http.ResponseWriter, r *http.Request) {
	var p products.Product
	if err := p.Decode(r.Body); err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte("bad request"))
	}

	updP, err := h.db.UpdateProduct(r.Context(), p)
	if errors.As(err, &pgx.ErrNoRows) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusNotFound, []byte("product not found"))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	b, err := updP.Encode()
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.response.json(w, http.StatusOK, b)
}

func (h *handlers) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	err = h.db.DeleteProduct(r.Context(), id)
	if errors.As(err, &pgx.ErrNoRows) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusNotFound, []byte("product not found"))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

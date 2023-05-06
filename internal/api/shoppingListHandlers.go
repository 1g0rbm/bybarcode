package api

import (
	"bybarcode/internal/db"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"bybarcode/internal/products"
)

func (h *handlers) addShoppingList(w http.ResponseWriter, r *http.Request) {
	var sl products.ShoppingList
	if err := sl.Decode(r.Body); err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte("bad request"))
		return
	}

	slId, err := h.db.CreateShoppingList(r.Context(), sl)
	fmt.Println(err)
	if errors.Is(err, db.ErrDuplicateKey) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	sl.ID = slId
	b, err := sl.Encode()
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.response.json(w, http.StatusOK, b)
}

func (h *handlers) updateShoppingList(w http.ResponseWriter, r *http.Request) {
	var sl products.ShoppingList
	if err := sl.Decode(r.Body); err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte("bad request"))
		return
	}

	updSl, err := h.db.UpdateShoppingList(r.Context(), sl)
	fmt.Println(err)
	if errors.As(err, &pgx.ErrNoRows) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusNotFound, []byte("shopping list not found"))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	b, err := updSl.Encode()
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.response.json(w, http.StatusOK, b)
}

func (h *handlers) getShoppingListsByAccount(w http.ResponseWriter, r *http.Request) {
	accountId, err := strconv.ParseInt(chi.URLParam(r, "account_id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	lists, err := h.db.GetShoppingListsByAccount(r.Context(), accountId)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	if err := json.NewEncoder(w).Encode(lists); err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) deleteShoppingList(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	err = h.db.DeleteShoppingList(r.Context(), id)
	if errors.As(err, &pgx.ErrNoRows) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusNotFound, []byte("shopping list not found"))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) addProductToShoppingList(w http.ResponseWriter, r *http.Request) {
	slId, err := strconv.ParseInt(chi.URLParam(r, "sl_id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	pId, err := strconv.ParseInt(chi.URLParam(r, "barcode_or_id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	err = h.db.AddProductToShoppingListByIds(r.Context(), pId, slId)
	if errors.As(err, &db.ErrDuplicateKey) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.listener.Notify(r.Context(), slId)

	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) deleteProductFromShoppingList(w http.ResponseWriter, r *http.Request) {
	slId, err := strconv.ParseInt(chi.URLParam(r, "sl_id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	pId, err := strconv.ParseInt(chi.URLParam(r, "barcode_or_id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	err = h.db.DeleteProductFromShoppingList(r.Context(), slId, pId)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.listener.Notify(r.Context(), slId)

	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) getShoppingListProducts(w http.ResponseWriter, r *http.Request) {
	slId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	list, err := h.db.GetShoppingListProducts(r.Context(), slId)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) toggleProductStateInShoppingList(w http.ResponseWriter, r *http.Request) {
	slId, err := strconv.ParseInt(chi.URLParam(r, "sl_id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	pId, err := strconv.ParseInt(chi.URLParam(r, "product_id"), 10, 64)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	err = h.db.ToggleProductStateInShoppingList(r.Context(), slId, pId)
	if errors.As(err, &pgx.ErrNoRows) {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusNotFound, []byte(err.Error()))
		return
	}
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	h.listener.Notify(r.Context(), slId)

	w.WriteHeader(http.StatusNoContent)
}

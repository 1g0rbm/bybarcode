package api

import (
	"bybarcode/internal/products"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"bybarcode/internal/config"
	"bybarcode/internal/db"
)

type AppApi struct {
	db     db.Connect
	cfg    *config.ApiConfig
	logger zerolog.Logger
	router *chi.Mux
	server *http.Server
}

func NewAppApi(db db.Connect, cfg *config.ApiConfig, logger zerolog.Logger) *AppApi {
	r := chi.NewRouter()

	api := &AppApi{
		db:     db,
		cfg:    cfg,
		logger: logger,
		router: r,
		server: &http.Server{
			Addr:    cfg.Address,
			Handler: r,
		},
	}

	api.router.Use(middleware.RequestID)
	api.router.Use(middleware.RealIP)
	api.router.Use(middleware.Logger)
	api.router.Use(middleware.Recoverer)
	api.router.Use(api.authMiddleware)

	api.router.Get("/api/v1/ping", api.pingHandler)

	api.router.Get("/api/v1/product/{barcode}", api.findProductByBarcode)
	api.router.Post("/api/v1/product", api.addProduct)
	api.router.Put("/api/v1/product", api.updateProduct)
	api.router.Delete("/api/v1/product/{id}", api.deleteProduct)

	api.router.Post("/api/v1/shopping-list", api.addShoppingList)
	api.router.Put("/api/v1/shopping-list", api.updateShoppingList)

	return api
}

func (aa *AppApi) Run() error {
	return aa.server.ListenAndServe()
}

func (aa *AppApi) pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`{"status": "ok"}`))
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte(err.Error()))
	}
}

func (aa *AppApi) findProductByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode, err := strconv.ParseInt(chi.URLParam(r, "barcode"), 10, 64)
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	product, err := aa.db.FindProductByBarcode(r.Context(), barcode)
	if errors.As(err, &pgx.ErrNoRows) {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusNotFound, []byte("product not found"))
		return
	}
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	b, err := product.Encode()
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	aa.sendJson(w, http.StatusOK, b)
}

func (aa *AppApi) addProduct(w http.ResponseWriter, r *http.Request) {
	var p products.Product
	if err := p.Decode(r.Body); err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusBadRequest, []byte("bad request"))
	}

	productId, err := aa.db.CreateProduct(r.Context(), p)
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	p.ID = productId
	b, err := p.Encode()
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	aa.sendJson(w, http.StatusOK, b)
}

func (aa *AppApi) updateProduct(w http.ResponseWriter, r *http.Request) {
	var p products.Product
	if err := p.Decode(r.Body); err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusBadRequest, []byte("bad request"))
	}

	updP, err := aa.db.UpdateProduct(r.Context(), p)
	if errors.As(err, &pgx.ErrNoRows) {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusNotFound, []byte("product not found"))
		return
	}
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	b, err := updP.Encode()
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	aa.sendJson(w, http.StatusOK, b)
}

func (aa *AppApi) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	err = aa.db.DeleteProduct(r.Context(), id)
	if errors.As(err, &pgx.ErrNoRows) {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusNotFound, []byte("product not found"))
		return
	}
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (aa *AppApi) addShoppingList(w http.ResponseWriter, r *http.Request) {
	var sl products.ShoppingList
	if err := sl.Decode(r.Body); err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusBadRequest, []byte("bad request"))
		return
	}

	slId, err := aa.db.CreateShoppingList(r.Context(), sl)
	fmt.Println(err)
	if errors.Is(err, db.ErrDuplicateKey) {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	sl.ID = slId
	b, err := sl.Encode()
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	aa.sendJson(w, http.StatusOK, b)
}

func (aa *AppApi) updateShoppingList(w http.ResponseWriter, r *http.Request) {
	var sl products.ShoppingList
	if err := sl.Decode(r.Body); err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusBadRequest, []byte("bad request"))
		return
	}

	updSl, err := aa.db.UpdateShoppingList(r.Context(), sl)
	fmt.Println(err)
	if errors.As(err, &pgx.ErrNoRows) {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusNotFound, []byte("shopping list not found"))
		return
	}
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	b, err := updSl.Encode()
	if err != nil {
		aa.logger.Error().Msg(err.Error())
		aa.sendJson(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	aa.sendJson(w, http.StatusOK, b)
}

func (aa *AppApi) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := aa.db.FindNotExpiredSession(r.Context(), token)
		if err != nil {
			aa.logger.Error().Msg(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (aa *AppApi) sendJson(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		aa.logger.Fatal().Msg(err.Error())
	}

	aa.logger.Debug().Msgf("Send response with headers %s and body %s", w.Header(), string(body))
}

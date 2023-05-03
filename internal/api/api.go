package api

import (
	"context"
	"net/http"
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

type handlers struct {
	db       db.Connect
	logger   zerolog.Logger
	response response
}

func NewAppApi(db db.Connect, cfg *config.ApiConfig, logger zerolog.Logger) *AppApi {
	r := chi.NewRouter()
	h := handlers{
		db:     db,
		logger: logger,
		response: response{
			logger: logger,
		},
	}

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

	api.router.Get("/api/v1/ping", h.pingHandler)

	api.router.Route("/api/v1/product", func(r chi.Router) {
		r.Get("/{barcode}", h.findProductByBarcode)
		r.Post("/", h.addProduct)
		r.Put("/", h.updateProduct)
		r.Delete("/{id}", h.deleteProduct)
	})

	api.router.Route("/api/v1/shopping-list", func(r chi.Router) {
		r.Get("/{account_id}", h.getShoppingListsByAccount)
		r.Post("/", h.addShoppingList)
		r.Put("/", h.updateShoppingList)
		r.Delete("/{id}", h.deleteShoppingList)
		r.Get("/{id}/product", h.getShoppingListProducts)
		r.Post("/{sl_id}/product/{barcode_or_id}", h.addProductToShoppingList)
		r.Delete("/{sl_id}/product/{barcode_or_id}", h.deleteProductFromShoppingList)
		r.Post("/{sl_id}/product/{product_id}/check", h.toggleProductStateInShoppingList)
	})

	api.router.Route("/api/v1/statistic", func(r chi.Router) {
		r.Get("/{date_from}/{date_to}", h.getStatistic)
	})

	return api
}

func (aa *AppApi) Run() error {
	aa.logger.Info().Msgf("Api was started on host %s \n", aa.cfg.Address)
	return aa.server.ListenAndServe()
}

func (aa *AppApi) ShoutDown(ctx context.Context) error {
	if err := aa.db.Close(); err != nil {
		return err
	}

	return aa.server.Shutdown(ctx)
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

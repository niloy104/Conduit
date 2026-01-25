package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var r *chi.Mux

func RegisterRoutes(h *handler) *chi.Mux {
	r = chi.NewRouter()

	r.Route("/products", func(r chi.Router) {
		r.Post("/", h.createProduct)
		r.Get("/", h.listProducts)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getProduct)
			r.Patch("/", h.updateProduct)
			r.Delete("/", h.deleteProduct)
		})
	})

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", h.createOrder)
		r.Get("/", h.listOrders)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getOrder)
			r.Delete("/", h.deleteOrder)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/", h.listUsers)
		r.Patch("/", h.updateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", h.deleteUser)
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", h.loginUser)
		})

		r.Route("/logout", func(r chi.Router) {
			r.Post("/", h.logoutUser)
		})
	})

	r.Route("/tokens", func(r chi.Router) {
		r.Route("/renew", func(r chi.Router) {
			r.Post("/", h.renewAccessToken)
		})

		r.Route("/revoke/{id}", func(r chi.Router) {
			r.Post("/", h.revokeSession)
		})
	})

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, r)
}

package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/xRebellion/MessageDistributor/message"
	"github.com/xRebellion/MessageDistributor/recipient"
)

type ChiRouter struct {
	recipientController *recipient.Controller
	messageController   *message.Controller
}

func (cr ChiRouter) InitRouter() *chi.Mux {
	r := chi.NewRouter()

	// CORS
	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	r.Use(corsHandler)

	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/recipients", func(r chi.Router) {
			r.Get("/", cr.recipientController.ListRecipients)
			r.Get("/{recipientID}", cr.recipientController.GetRecipientByID)
			r.Post("/", cr.recipientController.CreateRecipient)
			// r.Delete("/", recipient.DeleteRecipient)
		})
		r.Route("/message", func(r chi.Router) {
			r.Post("/send", cr.messageController.SendMessage)
		})
	})
	return r
}

func NewChiRouter(recipientController *recipient.Controller, messageController *message.Controller) *ChiRouter {
	return &ChiRouter{
		recipientController: recipientController,
		messageController:   messageController,
	}
}

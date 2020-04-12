/*
All Routes with handler functions goes here
Router specific context also defined here
*/
package main

import (
	"github.com/go-chi/chi"
)

// Routes creates a REST router for the groupsets API server
func (a *APIServer) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", a.Get)

	return r
}

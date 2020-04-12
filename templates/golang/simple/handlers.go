package main

import (
	"net/http"

	"github.com/go-chi/render"
)

func (a *APIServer) Get(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{
		"msg": "Glue Initialized",
	}
	render.JSON(w, r, res)
}

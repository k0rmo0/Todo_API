package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/unrolled/render"
)

// Renderer ->
var Renderer *RendererCtrl

// RendererCtrl ->
type RendererCtrl struct {
	r   *render.Render
	xml *template.Template
}

// Render ->
func (rend *RendererCtrl) Render(res http.ResponseWriter, status int, v interface{}) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if rend == nil {
		log.Println("REND ctrlr is NIL")
		return
	}

	if rend.r == nil {
		log.Println("REND is NIL")
		return
	}

	rend.r.JSON(res, status, v)
}

//WriteJSON encodes interface to json and writes it to writer
func WriteJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

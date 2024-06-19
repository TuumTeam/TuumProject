package handlers

import (
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Gestion des erreurs ici
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

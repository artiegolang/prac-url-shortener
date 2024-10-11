package handler

import "net/http"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/" {
			HandleShortenURL(w, r)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	case http.MethodGet:
		HandleRedirect(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

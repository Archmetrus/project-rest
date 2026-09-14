package transport

import (
	"example.com/project-rest/internal/store"
	"net/http"
)

func UserHandler(s store.UserStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		values, err := s.List(r.Context())
		if err != nil {
			failure(w, err)
			return
		}
		respond(w, 200, values)
	})
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(w, r)
		if !ok {
			return
		}
		value, err := s.Get(r.Context(), id)
		if err != nil {
			failure(w, err)
			return
		}
		respond(w, 200, value)
	})
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var value store.User
		if !decode(w, r, &value) {
			return
		}
		value, err := s.Create(r.Context(), value)
		if err != nil {
			failure(w, err)
			return
		}
		respond(w, 201, value)
	})
	mux.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(w, r)
		if !ok {
			return
		}
		var value store.User
		if !decode(w, r, &value) {
			return
		}
		value.ID = id
		value, err := s.Update(r.Context(), value)
		if err != nil {
			failure(w, err)
			return
		}
		respond(w, 200, value)
	})
	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(w, r)
		if !ok {
			return
		}
		if err := s.Delete(r.Context(), id); err != nil {
			failure(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}

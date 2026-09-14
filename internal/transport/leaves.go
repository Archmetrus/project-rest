package transport

import (
	"example.com/project-rest/internal/store"
	"net/http"
)

func LeaveHandler(s store.LeaveStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /leaves", func(w http.ResponseWriter, r *http.Request) {
		values, err := s.List(r.Context())
		if err != nil {
			failure(w, err)
			return
		}
		respond(w, 200, values)
	})
	mux.HandleFunc("GET /leaves/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("POST /leaves", func(w http.ResponseWriter, r *http.Request) {
		var value store.Leave
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
	mux.HandleFunc("PUT /leaves/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(w, r)
		if !ok {
			return
		}
		var value store.Leave
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
	mux.HandleFunc("DELETE /leaves/{id}", func(w http.ResponseWriter, r *http.Request) {
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

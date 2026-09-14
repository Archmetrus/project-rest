package transport

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

func respond(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}

func failure(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		respond(w, http.StatusNotFound, map[string]string{"error": "record not found"})
		return
	}
	log.Printf("database error: %v", err)
	respond(w, http.StatusInternalServerError, map[string]string{"error": "database operation failed"})
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	// Inspect raw JSON to reject null before decoding the object fields.
	var raw json.RawMessage
	if err := decoder.Decode(&raw); err != nil || string(raw) == "null" {
		respond(w, 400, map[string]string{"error": "invalid JSON object"})
		return false
	}
	// Preserve strict field/type checking after validating the top-level object.
	if err := strictJSON(raw, v); err != nil {
		respond(w, 400, map[string]string{"error": "invalid JSON object"})
		return false
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		respond(w, 400, map[string]string{"error": "expected one JSON object"})
		return false
	}
	return true
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respond(w, 400, map[string]string{"error": "id must be an int64"})
		return 0, false
	}
	return id, true
}

func strictJSON(raw []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

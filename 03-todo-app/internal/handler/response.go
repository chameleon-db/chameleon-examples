package handler

import (
	"encoding/json"
	"net/http"
)

// Response is the standard JSON response structure
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// respondJSON writes JSON response with status code
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Data: data})
}

// respondError writes error JSON response
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Error: message})
}

// decodeJSON decodes JSON body into target
func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// queryIntParam gets integer query parameter with default
func queryIntParam(r *http.Request, name string, defaultVal int) int {
	val := r.URL.Query().Get(name)
	if val == "" {
		return defaultVal
	}

	var result int
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return defaultVal
	}

	return result
}

// queryBoolParam gets boolean query parameter
func queryBoolParam(r *http.Request, name string) *bool {
	val := r.URL.Query().Get(name)
	if val == "" {
		return nil
	}

	var result bool
	switch val {
	case "true", "1":
		result = true
	case "false", "0":
		result = false
	default:
		return nil
	}

	return &result
}

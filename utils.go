package main

import (
	"net/http"
	"encoding/json"
)

//
// Encoding/Serialization: go data -> json
// Encoder() method takes the data container
// and encodes it to JSON
//
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	//
	// Writing header has to be done
	// before writing body
	//
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

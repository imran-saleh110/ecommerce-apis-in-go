package main

import (
	"encoding/json"
	"log"
	"net/http"
	"uuid"
)

type server struct {
	addr string
}

// GET users
func (s *server) getUsersListHandler(w http.ResponseWriter, r *http.Request) {
	usersList := []string{}
	for _, user := range users {
		usersList = append(usersList, user.Name)
	}

	WriteJSON(w, http.StatusOK, usersList)
}

// POST user
func (s *server) createNewUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	//
	// Decoding/Deserialization: json -> go data
	// NOTE: Decode() method takes the address of the struct
	// where you wish to store the JSON decoded data
	//
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	/* A UUID (Universally Unique Identifier) is a 128-bit value used in
	computer systems to uniquely identify information, typically shown
	as a 36-character string like 123e4567-e89b-12d3-a456-426614174000.
	UUIDs are generated locally (without a central authority) and are
	designed so that the chance of duplicates is extremely small.
	*/
	newID := uuid.New().String()

	// add a new user
	newUser := User{
		ID:   newID,
		Name: req.Name,
	}
	users[newID] = newUser

	//
	// w.Write([]byte("User Created"))
	// using fmt.Fprintln is more cleaner
	//
	w.Header().Set("Content-Type", "application/json") // Mandatory!
    w.WriteHeader(http.StatusCreated) // 201

	json.NewEncoder(w).Encode(newUser)
}

// GET /users/{id} -> { ID, Name }
func (s *server) getSingleUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// PathValue parses the relevant part from the handler

	user, exists := users[id]
	if !exists {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	WriteJSON(w, http.StatusCreated, user)

	/*
		  json.Marshal converts a value into JSON bytes ([]byte),
			while json.Encoder.Encode writes JSON directly to an
			io.Writer, such as an HTTP response or file. Use Marshal
			when you need the bytes in memory; use Encoder for streaming
			output, which can avoid holding the entire result in memory.

			data, err := json.Marshal(user)
			if err != nil {
			log.Println("Json Marshaling failed")
			return
			}
			fmt.Fprintln(w, string(data))
	 */
}


func main() {
	srv := &server{addr: ":8080"}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", srv.getUsersListHandler)
	mux.HandleFunc("POST /users", srv.createNewUserHandler)
	mux.HandleFunc("GET /users/{id}", srv.getSingleUserHandler)

	log.Println("Server is starting...")
	http.ListenAndServe(srv.addr, mux)
}

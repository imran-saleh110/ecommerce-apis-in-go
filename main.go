package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"uuid"
)

type server struct {
	addr string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// GET users
	if r.Method == "GET" && r.URL.Path == "/users" {
		//
		// Writing header has to be done
		// before writing body
		//
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		usersList := []string{}
		for _, user := range users {
			usersList = append(usersList, user.Name)
		}

		json.NewEncoder(w).Encode(usersList)
		return
	}

	// POST user
	if r.Method == "POST" && r.URL.Path == "/users" {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201

		json.NewEncoder(w).Encode(newUser)
		return
	}

	// GET /users/{id} -> { ID, Name }
	if r.Method == "GET" {
		//
		// parsing rotues
		// /users/42 -> ["", "users", "42"]
		//
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) == 3 && parts[1] == "users" {
			id := parts[2]

			user, exists := users[id]
			if !exists {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			//
			// Encoding/Serialization: go data -> json
			// Encoder() method takes the data container
			// and encodes it to JSON
			//
			json.NewEncoder(w).Encode(user)
			//
			// data, err := json.Marshal(user)
			// if err != nil {
			// 	log.Println("Json Marshaling failed")
			// 	return
			// }
			// fmt.Fprintln(w, string(data))
			//
		}

		return
	}

	http.Error(w, "the route does not exist", http.StatusNotFound)
}

func main() {
	srv := &server{
		addr: ":8080",
	}

	log.Println("Server is starting...")

	http.ListenAndServe(srv.addr, srv)
}

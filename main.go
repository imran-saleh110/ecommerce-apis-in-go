package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}

		if string(body) == "" {
			http.Error(w, "Username cannot be empty.", http.StatusBadRequest)
			return
		}

		// add a new user
		newID := uuid.New().String()
		newUser := User{
			ID:   newID,
			Name: string(body),
		}
		users[newID] = newUser

		//
		// w.Write([]byte("User Created"))
		// using fmt.Fprintln is more cleaner
		//
		w.WriteHeader(http.StatusCreated) // 201
		fmt.Fprintln(w, "User Created")
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

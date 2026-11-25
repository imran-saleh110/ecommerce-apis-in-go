package main

import (
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
		w.WriteHeader(http.StatusOK)

		for _, user := range users {
			fmt.Fprintln(w, user)
		}
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

			fmt.Fprintln(w, user)
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

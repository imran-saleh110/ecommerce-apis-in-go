package main

import (
	"io"
	"log"
	"net/http"
)

type server struct {
	addr string
}

var users = []string{"Abdul Hamid", "Haider Ali"}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// GET users
	if r.Method == "GET" && r.URL.Path == "/users" {
		//
		// Writing header has to be done
		// before writing body
		//
		w.WriteHeader(http.StatusOK)

		for _, user := range users {
			w.Write([]byte(user))
			w.Write([]byte("\n"))
		}
	}

	// POST user
	if r.Method == "POST" && r.URL.Path == "/users" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}

		// add a new user
		users = append(users, string(body))
		w.WriteHeader(http.StatusCreated)  // 201

		w.Write([]byte("User Created"))
	}
}

func main() {
	srv := &server {
		addr: ":8080",
	}

	log.Println("Server is starting...")

	http.ListenAndServe(srv.addr, srv)
}

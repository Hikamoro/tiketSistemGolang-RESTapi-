package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Issue struct { //структура тикета
	Subject  string `json:"subject"`
	Text     string `json:"text"`
	Priority string `json:"priority"`
}

type Response struct { //структура ответа JSON
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type Store struct { //структура хранилища
	mu     sync.Mutex
	data   map[int]Issue
	nextID int
}

func NewStore() *Store { //конструктор хранилища
	return &Store{
		data:   make(map[int]Issue),
		nextID: 1,
	}
}

func (s *Store) Create(issue Issue) int { //добавление обьекта в хранилище
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	s.data[id] = issue
	return id
}

func (s *Store) GetAll() map[int]Issue { //получить все тикеты
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := make(map[int]Issue)
	for k, v := range s.data {
		copy[k] = v
	}

	return copy
}

func (s *Store) Get(id int) (Issue, bool) { //получить тикет по id
	s.mu.Lock()
	defer s.mu.Unlock()

	issue, ok := s.data[id]
	return issue, ok
}

func (s *Store) Update(id int, issue Issue) bool { //обновление тикета
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return false
	}

	s.data[id] = issue
	return true
}

func (s *Store) Delete(id int) bool { //удаление тикета
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return false
	}

	delete(s.data, id)
	return true
}

func writeJSON(w http.ResponseWriter, status int, data any) { //выход JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func parseID(path string) (int, error) { //парсинг ID
	parts := strings.Split(path, "/")

	if len(parts) != 3 {
		return 0, http.ErrNoLocation
	}

	return strconv.Atoi(parts[2])
}

type Server struct { //структура сервера
	store *Store
}

func NewServer(store *Store) *Server { //конструктор сервера
	return &Server{store: store}
}

func (s *Server) IssuesHandler(w http.ResponseWriter, r *http.Request) { //функция обоработки методов

	// /issues
	if r.URL.Path == "/issues" {

		switch r.Method {

		case http.MethodGet:
			data := s.store.GetAll()

			writeJSON(w, http.StatusOK, Response{
				Data: data,
			})

		case http.MethodPost:
			var issue Issue

			if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
				writeJSON(w, http.StatusBadRequest, Response{
					Error: "invalid json",
				})
				return
			}

			id := s.store.Create(issue)

			writeJSON(w, http.StatusCreated, Response{
				Data: map[string]int{"id": id},
			})

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		return
	}

	// /issues/{id}
	if strings.HasPrefix(r.URL.Path, "/issues/") {

		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Error: "invalid id",
			})
			return
		}

		switch r.Method {

		case http.MethodGet:

			issue, ok := s.store.Get(id)
			if !ok {
				writeJSON(w, http.StatusNotFound, Response{
					Error: "not found",
				})
				return
			}

			writeJSON(w, http.StatusOK, Response{
				Data: issue,
			})

		case http.MethodPut:

			var issue Issue

			if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
				writeJSON(w, http.StatusBadRequest, Response{
					Error: "invalid json",
				})
				return
			}

			if !s.store.Update(id, issue) {
				writeJSON(w, http.StatusNotFound, Response{
					Error: "not found",
				})
				return
			}

			writeJSON(w, http.StatusOK, Response{
				Data: "updated",
			})

		case http.MethodDelete:

			if !s.store.Delete(id) {
				writeJSON(w, http.StatusNotFound, Response{
					Error: "not found",
				})
				return
			}

			writeJSON(w, http.StatusOK, Response{
				Data: "deleted",
			})

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		return
	}

	w.WriteHeader(http.StatusNotFound)
}


func main() {

	store := NewStore()
	server := NewServer(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/issues", server.IssuesHandler)
	mux.HandleFunc("/issues/", server.IssuesHandler)

	log.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", mux) // поднятие сервера
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func main() {
	client := NewClient("http://localhost:9200")
	if err := client.CheckHealth(); err != nil {
		log.Fatal("failed to health check: ", err)
	}

	if err := client.CreateIndex(); err != nil {
		log.Fatal("failed create index: ", err)
	}

	server := Server{client: client}
	http.HandleFunc("/insert", server.InsertDataHandler)
	http.HandleFunc("/update", server.UpdateDataHandler)
	http.HandleFunc("/delete", server.DeleteDataHandler)
	http.HandleFunc("/search", server.SearchDataHandler)
	http.HandleFunc("/health", server.HealthCheckHandler)

	log.Println("listening server on port 8181")
	http.ListenAndServe(":8181", nil)
}

type Server struct {
	client *Client
}

func (s *Server) InsertDataHandler(w http.ResponseWriter, r *http.Request) {
	var employee *Employee
	json.NewDecoder(r.Body).Decode(&employee)
	if err := s.client.InsertData(employee); err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, employee)
}

func (s *Server) UpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	var employee *Employee
	json.NewDecoder(r.Body).Decode(&employee)
	if err := s.client.UpdateData(employee); err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, employee)
}

func (s *Server) DeleteDataHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	if err := s.client.DeleteData(id); err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, Employee{Id: id})
}

func (s *Server) SearchDataHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.FormValue("keyword")
	employees, err := s.client.SearchData(keyword)
	if err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, employees)
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.client.CheckHealth(); err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, map[string]string{
		"status": "OK",
	})
}

func writeResponseOK(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeResponse(w, response)
}

func writeResponseInternalError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	writeResponse(w, map[string]interface{}{
		"error": err,
	})
}

func writeResponse(w http.ResponseWriter, response interface{}) {
	json.NewEncoder(w).Encode(response)
}

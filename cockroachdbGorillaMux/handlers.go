package main

import (
	"encoding/json"

	"net/http"

	"github.com/carlosm27/blog/cockroachdb-gorillamux/model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (s *Server) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/expenses", s.getExpenses)
	router.HandleFunc("/expense/{id}", s.getExpense).Methods("GET")
	router.HandleFunc("/expense", s.createExpense).Methods("POST")
	router.HandleFunc("/expense/{id}", s.updateExpense).Methods("PUT")
	router.HandleFunc("/expense/{id}", s.deleteExpense).Methods("DELETE")

}

func (s *Server) getExpenses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expenses []model.Expenses

	if err := s.db.Find(&expenses).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expenses)
	}
}

func (s *Server) createExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.db.Create(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expense)
	}
}

func (s *Server) getExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.Where("id= ?", id).First(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expense)
	}
}
func (s *Server) updateExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses

	vars := mux.Vars(r)
	id := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.db.Where("id=?", id).Updates(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expense)
	}
}

func (s *Server) deleteExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.Where("id = ?", id).Delete(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}

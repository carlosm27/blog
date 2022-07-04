package main

import (
	"encoding/json"

	"net/http"

	"github.com/carlosm27/blog/cockroachdb-gorillamux/model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UpdateExpense struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

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
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)

}

func (s *Server) createExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newExpense := &model.Expenses{Amount: expense.Amount, Description: expense.Description, Category: expense.Category}
	if err := s.db.Create(newExpense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newExpense)

}

func (s *Server) getExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.Where("id = ?", id).First(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expense)

}

func (s *Server) updateExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var updateExpense UpdateExpense
	var expense model.Expenses

	vars := mux.Vars(r)
	id := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&updateExpense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.db.Where("id = ?", id).First(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := s.db.Model(&expense).Updates(&model.Expenses{Amount: updateExpense.Amount, Description: updateExpense.Description, Category: updateExpense.Category}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expense)

}

func (s *Server) deleteExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.Where("id = ?", id).First(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := s.db.Delete(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Expense Deleted")

}

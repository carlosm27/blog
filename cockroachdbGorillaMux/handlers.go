package main

import (
	"encoding/json"

	"net/http"

	"github.com/carlosm27/blog/cockroachdb-gorillamux/model"
	"github.com/google/uuid"
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
	router.HandleFunc("/expense/{id_expense}", s.getExpense).Methods("GET")
	router.HandleFunc("/expense", s.createExpense).Methods("POST")
	router.HandleFunc("/expense/{id_expense}", s.updateExpense).Methods("PUT")
	router.HandleFunc("/expense/{id_expense}", s.deleteExpense).Methods("DELETE")

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

	expense_id := uuid.New().String()
	if err := s.db.Create(&model.Expenses{IdExpense: expense_id, Amount: expense.Amount, Description: expense.Description, Category: expense.Category}).Error; err != nil {
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
	id := vars["id_expense"]

	if err := s.db.Where("id_expense = ?", id).First(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expense)
	}
}
func (s *Server) updateExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var updateExpense UpdateExpense
	var expense model.Expenses

	vars := mux.Vars(r)
	id := vars["id_expense"]

	if err := json.NewDecoder(r.Body).Decode(&updateExpense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.db.Where("id_expense = ?", id).First(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err := s.db.Model(&expense).Updates(&model.Expenses{Amount: updateExpense.Amount, Description: updateExpense.Description, Category: updateExpense.Category}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expense)
	}

}

func (s *Server) deleteExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var expense model.Expenses
	vars := mux.Vars(r)
	id := vars["id_expense"]

	if err := s.db.Where("id_expense = ?", id).First(&expense).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		s.db.Delete(&expense)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Deleted")
	}

}

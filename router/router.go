package router

import (
	"github.com/gorilla/mux"
	"github.com/to4to/itr-blw/middleware"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/itr/{id}", middleware.GetITR).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/itr", middleware.GetAllITR).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newitr", middleware.CreateITR).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/itr/{id}", middleware.UpdateITR).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteitr/{id}", middleware.DeleteITR).Methods("DELETE", "OPTIONS")

	return router
}
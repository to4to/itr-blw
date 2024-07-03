package router

import (
	"github.com/gorilla/mux"
	"github.com/to4to/itr-blw/middleware"
)


// Router returns a new instance of mux.Router configured with various API endpoints and their corresponding handlers.
func Router() *mux.Router {

    router := mux.NewRouter()

    router.HandleFunc("/api/itr/{id}", middleware.HandlerGetITR).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/itr", middleware.HandlerGetAllITR).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/newitr", middleware.HandlerCreateITR).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/itr/{id}", middleware.HandlerUpdateITR).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/deleteitr/{id}", middleware.HandlerDeleteITR).Methods("DELETE", "OPTIONS")

    return router
}

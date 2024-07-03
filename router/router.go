package router


import (
	"github.com/gorilla/mux"
	"github.com/to4to/itr-blw/middleware"
)


// Router returns a new instance of mux.Router configured with various API endpoints and their corresponding handlers.
// Router returns a new instance of mux.Router with predefined routes for handling ITR (Individual Tax Return) related operations.
func Router() *mux.Router {
    router := mux.NewRouter()

    // Handle GET request for retrieving a specific ITR by ID
    router.HandleFunc("/api/itr/{id}", middleware.HandlerGetITR).Methods("GET", "OPTIONS")
    
    // Handle GET request for retrieving all ITRs
    router.HandleFunc("/api/itr", middleware.HandlerGetAllITR).Methods("GET", "OPTIONS")
    
    // Handle POST request for creating a new ITR
    router.HandleFunc("/api/newitr", middleware.HandlerCreateITR).Methods("POST", "OPTIONS")
    
    // Handle PUT request for updating an existing ITR by ID
    router.HandleFunc("/api/itr/{id}", middleware.HandlerUpdateITR).Methods("PUT", "OPTIONS")
    
    // Handle DELETE request for deleting an ITR by ID
    router.HandleFunc("/api/deleteitr/{id}", middleware.HandlerDeleteITR).Methods("DELETE", "OPTIONS")

    return router
}

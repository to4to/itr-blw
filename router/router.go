package router

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/to4to/itr-blw/handler"
)

func Router() chi.Router {
	godotenv.Load()
	port := os.Getenv("PORT")
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/create", handler.CreateUser)
	v1Router.Get("/find/{id}", handler.FindUser)
	v1Router.Patch("/update", handler.UpdateUser)
	v1Router.Delete("/delete/{id}", handler.DeleteUser)
	v1Router.Get("/findall", handler.FindAllUser)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal((err))

	}

	return router
}

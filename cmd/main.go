package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	postgres "github.com/horzu/golang/picus-security-bootcamp/homework-4-week-5-horzu/pkg/db"
	"github.com/horzu/golang/picus-security-bootcamp/homework-4-week-5-horzu/pkg/models/repos"
	"github.com/joho/godotenv"
)

func main() {
	// Set Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db, err := postgres.NewPsqlDB()
	if err != nil {
		log.Fatalf("Postgres cannot init: %s", err)
	}
	log.Printf("Connected to Postgres Database.")

	// Initialize Repositories
	authorRepo := repos.NewAuthorRepository(db)
	authorRepo.Migration()
	bookRepo := repos.NewBookRepository(db)
	bookRepo.Migration()
	// bookRepo.InsertSampleData()
	// authorRepo.InsertSampleData()

	// fmt.Println(authorRepo.GetAuthorWithBooksById(12))
	// fmt.Println(bookRepo.GetBookWithAuthorsById(12))

	// fmt.Println(authorRepo.GetAuthorWithBooks2())

	r := mux.NewRouter()

	handlers.AllowedOrigins([]string{"https://localhost"})
	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})

	b := r.PathPrefix("/books").Subrouter()

    b.HandleFunc("/", bookRepo.GetAllBooks).Methods(http.MethodGet)
	b.HandleFunc("/withauthors", bookRepo.GetAllBooksWithAuthorById).Methods(http.MethodGet)
    b.HandleFunc("/{id}", bookRepo.GetBookByID).Methods(http.MethodGet)
    b.HandleFunc("/{id}/withauthors", bookRepo.GetBooksWithAuthorById).Methods(http.MethodGet)
	b.HandleFunc("/", bookRepo.AddBook).Methods(http.MethodPost)
    b.HandleFunc("/find/{name}", bookRepo.FindBookByName).Methods(http.MethodGet)
	b.HandleFunc("/{id}", bookRepo.UpdateBook).Methods(http.MethodPut)
	b.HandleFunc("/buy/{id}/{quantity}", bookRepo.BuyBookByID).Methods(http.MethodPatch)
	b.HandleFunc("/{id}", bookRepo.DeleteBook).Methods(http.MethodDelete)
    r.HandleFunc("/bookcount", bookRepo.GetBooksCount).Methods(http.MethodGet)
    b.HandleFunc("/lessthen/{pages}", bookRepo.GetBooksByPagesLessThenWithAuthorInformation).Methods(http.MethodGet)

	a := r.PathPrefix("/authors").Subrouter()
	
	a.HandleFunc("/", authorRepo.GetAllAuthors).Methods(http.MethodGet)
	a.HandleFunc("/withbooks", authorRepo.GetAllAuthorsWithBooksById).Methods(http.MethodGet)
	a.HandleFunc("/{id}", authorRepo.GetAuthorByID).Methods(http.MethodGet)
    a.HandleFunc("/{id}/withbooks", authorRepo.GetAuthorWithBooksById).Methods(http.MethodGet)
	a.HandleFunc("/", authorRepo.AddAuthor).Methods(http.MethodPost)
    a.HandleFunc("/find/{name}", authorRepo.FindAuthorByName).Methods(http.MethodGet)
	a.HandleFunc("/{id}", authorRepo.UpdateAuthor).Methods(http.MethodPut)
	a.HandleFunc("/{id}", authorRepo.DeleteAuthor).Methods(http.MethodDelete)
    r.HandleFunc("/authorcount", authorRepo.GetAuthorsCount).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         "127.0.0.1:4000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Println("API is running!")

		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(srv, time.Second*10)
}

func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
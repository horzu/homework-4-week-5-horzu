package main

import (
	"log"
	"net/http"

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

    log.Println("API is running!")
	http.ListenAndServe(":4000", r)


	// Initialize methods
	// fmt.Println(authorRepo.GetAuthorByID(1))
	// fmt.Println(authorRepo.FindAuthorByName("es"))
	// fmt.Println(authorRepo.GetAuthorsWithBookInformation())
	// fmt.Println(authorRepo.GetDeletedAuthorsWithBookInformation())
	// fmt.Println(authorRepo.GetAuthorWithBookInformationByID(3))
	// authorRepo.DeleteById(2)
	// authorRepo.GetAuthorByID(1)
	// fmt.Println(authorRepo.GetAllBooks())
	// fmt.Println(authorRepo.GetAuthorsCount())

	// bookRepo.InsertSampleData()
	// fmt.Println(bookRepo.GetBookByID(1))
	// fmt.Println(bookRepo.FindBookByName("decoder"))
	// result, _ := bookRepo.GetAllBooksWithAuthorInformation()
	// for _, v := range result {
	// 	fmt.Printf("Book: %s, Author: %s\n",v.Title, v.Authors[0].Name)
	// }
	// bookRepo.DeleteById(2)
	// bookRepo.GetBookByID(1)
	// fmt.Println(bookRepo.GetAllBooks())
	// fmt.Println(bookRepo.GetBooksCount())
	// fmt.Println(bookRepo.GetBooksWithAuthorInformation())
	// fmt.Println(bookRepo.GetDeletedBooksWithAuthorInformation())
	// fmt.Println(bookRepo.GetBookWithAuthorInformationByID(3))
	// fmt.Println(bookRepo.GetBooksByPagesLessThenWithAuthorInformation(500))
	// bookRepo.GetStockCodeByTitle("the")
	// fmt.Println(bookRepo.BuyBookByID(3,50))

}

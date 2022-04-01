package repos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	models "github.com/horzu/golang/picus-security-bootcamp/homework-4-week-5-horzu/pkg/models/entities"
	http_errors "github.com/horzu/golang/picus-security-bootcamp/homework-4-week-5-horzu/pkg/models/errors"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (b *BookRepository) Migration() {
	b.db.AutoMigrate(&models.Book{})
}


// InsertSampleData inserts sample data to database
func (b *BookRepository) InsertSampleData() {
	jsonFile, err := os.Open("./pkg/mocks/books.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	values, _ := ioutil.ReadAll(jsonFile)
	books := []models.Book{}
	json.Unmarshal(values, &books)

	for _, book := range books {
		b.db.FirstOrCreate(&book)
	}
}

// swagger:route GET /books books GetAllBooks
// Returns a list of books
// responses:
//  200: booksResponseSlice

// GetAllBooks lists all available books
func (b *BookRepository) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.Book

	if result := b.db.Find(&books); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(books)
}

// swagger:route GET /books/{id} books GetBookByID
// Returns the book of given id
// responses:
//  200: bookResponse

// GetBookByID returns book information according to given id
func (b *BookRepository) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err!=nil{
		fmt.Println(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	// Iterate over all the books
	var book models.Book

	if result := b.db.First(&book, id); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(book)
	}
}

// swagger:route POST /books/ books AddBook
// Creates the book of given body
// responses:
//  201: bookResponse

// AddBook creates a new book
func (b *BookRepository) AddBook(w http.ResponseWriter, r *http.Request) {
	// Read to request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	var book models.Book
	json.Unmarshal(body, &book)

	// Append to the Book
	if result := b.db.Create(&book); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}
	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// UpdateBook updates the given book
func (b *BookRepository) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}
	// Read to request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	var updatedBook models.Book
	updatedBook.ID = uint(id)
	json.Unmarshal(body, &updatedBook)

	// Iterate over all the Books
	var book models.Book

	// Append to the Book
	if result := b.db.First(&book, id); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	b.db.Save(&updatedBook)
	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

// swagger:route DELETE /books/{id} books DeleteBook
// Deletes and returns the book of given id
// responses:
//  201: bookResponse

// DeleteBook deletes given book according to given id
func (b *BookRepository) DeleteBook(w http.ResponseWriter, r *http.Request) {
	// Read dynamic parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	// Find the book by id
	var book models.Book

	if result := b.db.First(&book, id); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}
	// Delete that book
	b.db.Delete(&book)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Deleted")
}

// FindBookByName returns books found according to given search query
func (b *BookRepository) FindBookByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var books []models.Book

	if result := b.db.Where("title ILIKE ? ", "%"+vars["name"]+"%").Find(&books); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	fmt.Println(vars["name"])

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(books)
}

// BuyBookByID buys book and returns the new state of the given book
func (b *BookRepository) BuyBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}
	quantity, err := strconv.Atoi(vars["quantity"])
	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}
	// Find the book by id
	var book models.Book

	if result := b.db.First(&book, id); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	// Update that book
	book.Stock -= quantity
	b.db.Save(&book)

	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// GetBooksCount returns number of books
func (b *BookRepository) GetBooksCount(w http.ResponseWriter, r *http.Request) {
	var count int

	b.db.Raw("SELECT COUNT(books.title)	FROM books WHERE books.deleted_at is null").Scan(&count)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}

// GetBooksWithAuthorById returns book with its author information
func (b *BookRepository) GetBooksWithAuthorById(w http.ResponseWriter, r *http.Request) {
	// Read dynamic parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	var Book models.Books

	if result := b.db.Preload("Authors").First(&Book, id); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Book)
}

// GetAllBooksWithAuthorById returns all books with their author information
func (b *BookRepository) GetAllBooksWithAuthorById(w http.ResponseWriter, r *http.Request) {
	var Books []models.Books

	if result := b.db.Preload("Authors").Find(&Books); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Books)
}

// GetBooksByPagesLessThenWithAuthorInformation returns all books which have less pages then given page number
func (b *BookRepository) GetBooksByPagesLessThenWithAuthorInformation(w http.ResponseWriter, r *http.Request) {
	// Read dynamic parameter
	vars := mux.Vars(r)
	pages, err := strconv.Atoi(vars["pages"])

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	var Books []models.Books

	if result := b.db.
		Table("books").
		Select("*").
		Where("books.page < ? ", pages).
		Joins("left join authors on authors.id = books.author_id").
		Scan(&Books); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Books)
}

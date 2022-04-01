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

func (b *BookRepository) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.Book

	if result := b.db.Find(&books); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(books)
}

func (b *BookRepository) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Iterate over all the books
	var book models.Book

	if result := b.db.First(&book, id); result.Error != nil {
		fmt.Println(result.Error)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(book)
	}
}

func (b *BookRepository) AddBook(w http.ResponseWriter, r *http.Request) {
	// Read to request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var book models.Book
	json.Unmarshal(body, &book)

	// Append to the Book
	result := b.db.Create(&book)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (b *BookRepository) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// Read to request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var updatedBook models.Book
	updatedBook.ID = uint(id)
	json.Unmarshal(body, &updatedBook)

	// Iterate over all the Books
	var book models.Book

	// Append to the Book
	if result := b.db.First(&book, id); result.Error != nil {
		fmt.Println(result.Error)
	}

	b.db.Save(&updatedBook)
	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

func (b *BookRepository) DeleteBook(w http.ResponseWriter, r *http.Request){
	// Read dynamic parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	// Find the book by id
	var book models.Book

	if result := b.db.First(&book, id); result.Error != nil{
		fmt.Println(result.Error)
	}
	// Delete that book
	b.db.Delete(&book)

	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}

func (b *BookRepository) FindBookByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var books []models.Book

	if result := b.db.Where("title ILIKE ? ", "%" + vars["name"] + "%").Find(&books); result.Error != nil {
		fmt.Println(result.Error)
	}

	fmt.Println(vars["name"])

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(books)
}

func (b *BookRepository) BuyBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	quantity, _ := strconv.Atoi(vars["quantity"])
	// Find the book by id
	var book models.Book

	if result := b.db.First(&book, id); result.Error != nil{
		fmt.Println(result.Error)
	}
	
	// Update that book
	book.Stock -= quantity
	b.db.Save(&book)

	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (b *BookRepository) GetBooksCount(w http.ResponseWriter, r *http.Request) {
	var count int

	b.db.Raw("SELECT COUNT(books.title)	FROM books WHERE books.deleted_at is null").Scan(&count)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}

func (b *BookRepository) GetBooksWithAuthorById(w http.ResponseWriter, r *http.Request) {
	// Read dynamic parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var Book models.Books

	if result := b.db.Preload("Authors").First(&Book, id); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Book)
}

func (b *BookRepository) GetAllBooksWithAuthorById(w http.ResponseWriter, r *http.Request) {
	var Books []models.Books

	if result := b.db.Preload("Authors").Find(&Books); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Books)
}

func (b *BookRepository) GetBooksByPagesLessThenWithAuthorInformation(w http.ResponseWriter, r *http.Request) {
	// Read dynamic parameter
	vars := mux.Vars(r)
	pages, _ := strconv.Atoi(vars["pages"])

	var Books []models.Books
	
	if result := b.db.
			Table("books").
			Select("*").
			Where("books.page < ? ", pages).
			Joins("left join authors on authors.id = books.author_id").
			Scan(&Books); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Books)
}


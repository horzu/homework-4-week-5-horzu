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

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (a *AuthorRepository) Migration() {
	a.db.AutoMigrate(&models.Author{})
}

// InsertSampleData inserts sample data to database
func (a *AuthorRepository) InsertSampleData() {
	jsonFile, err := os.Open("./pkg/mocks/authors.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	values, _ := ioutil.ReadAll(jsonFile)
	authors := []models.Author{}
	json.Unmarshal(values, &authors)

	for _, author := range authors {
		a.db.FirstOrCreate(&author)
	}
}

// GetAllAuthors lists all available authors
func (a *AuthorRepository) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	var author []models.Author

	if result := a.db.Find(&author); result.Error != nil {
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(author)
}

// GetAllAuGetAuthorByID returns author information according to given id
func (a *AuthorRepository) GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Iterate over all the books
	var author models.Author

	if result := a.db.First(&author, id); result.Error != nil {
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(author)
}

// AddAuthor creates a new author
func (a *AuthorRepository) AddAuthor(w http.ResponseWriter, r *http.Request) {
	// Read to request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	var author models.Author
	json.Unmarshal(body, &author)

	// Append to the Book
	result := a.db.Create(&author)
	if result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}
	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

// UpdateAuthor updates the given author
func (a *AuthorRepository) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// Read to request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatalln(err)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(err))
		return
	}

	var updatedAuthor models.Author
	updatedAuthor.ID = uint(id)
	json.Unmarshal(body, &updatedAuthor)

	// Iterate over all the Books
	var author models.Author

	// Append to the Book
	if result := a.db.First(&author, id); result.Error != nil {
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	a.db.Save(&updatedAuthor)
	// Send a 201 created response
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAuthor)
}

// DeleteAuthor deletes given author according to given id
func (a *AuthorRepository) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	// Read dynamic parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Find the author by id
	var author models.Author

	if result := a.db.First(&author, id); result.Error != nil {
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}
	// Delete that book
	a.db.Delete(&author)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}

// FindAuthorByName returns authors found according to given search query
func (a *AuthorRepository) FindAuthorByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var author []models.Author

	if result := a.db.Where("name ILIKE ? ", "%"+vars["name"]+"%").Find(&author); result.Error != nil {
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	fmt.Println(vars["name"])

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(author)
}

// GetAuthorsCount returns number of authors
func (a *AuthorRepository) GetAuthorsCount(w http.ResponseWriter, r *http.Request) {
	var count int

	a.db.Raw("SELECT COUNT(authors.name) FROM authors WHERE authors.deleted_at is null").Scan(&count)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}

// GetAuthorWithBooksById returns author with its book information
func (a *AuthorRepository) GetAuthorWithBooksById(w http.ResponseWriter, r *http.Request) {
	// Read dynamic parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var Author models.Author

	if result := a.db.Preload("Books").First(&Author, id); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Author)
}

// GetAllAuthorsWithBooksById returns all authors with their book information
func (a *AuthorRepository) GetAllAuthorsWithBooksById(w http.ResponseWriter, r *http.Request) {
	var Authors []models.Author

	if result := a.db.Preload("Books").Find(&Authors); result.Error != nil {
		fmt.Println(result.Error)
		json.NewEncoder(w).Encode(http_errors.ParseErrors(result.Error))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Authors)
}

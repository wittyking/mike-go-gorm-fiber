package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	PublisherID uint
	Publisher   Publisher
	Authors     []Author `gorm:"many2many:author_books;"`
}

type Publisher struct {
	gorm.Model
	Details string
	Name    string
}

type Author struct {
	gorm.Model
	Name  string
	Books []Book `gorm:"many2many:author_books;"`
}

type AuthorBook struct {
	AuthorID uint
	Author   Author
	BookID   uint
	Book     Book
}

func main() {
	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Book{}, &Publisher{}, &Author{}, &AuthorBook{})

	// publisher := Publisher{
	// 	Details: "Mikelopster",
	// 	Name:    "Mike",
	// }
	// _ = createPublisher(db, &publisher)

	// author1 := Author{
	// 	Name: "Mike1",
	// }
	// _ = createAuthor(db, &author1)

	// author2 := Author{
	// 	Name: "Mike2",
	// }
	// _ = createAuthor(db, &author2)

	// book := Book{
	// 	Name:        "Mike Book",
	// 	Author:      "0000",
	// 	Description: "Book Description",
	// 	PublisherID: publisher.ID,
	// 	Authors:     []Author{author1, author2},
	// }
	// _ = createBookWithAuthor(db, &book)

	book, err := listBooksOfAuthor(db, 2)
	// book, err := getBookWithAuthors(db, 2)
	// book, err := getBookWithPublisher(db, 2)
	fmt.Println("=========")
	fmt.Println(book)

}

func createPublisher(db *gorm.DB, publisher *Publisher) error {
	result := db.Create(publisher)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func createAuthor(db *gorm.DB, author *Author) error {
	result := db.Create(author)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func createBookWithAuthor(db *gorm.DB, book *Book) error {
	if err := db.Create(book).Error; err != nil {
		return err
	}

	return nil
}

func getBookWithPublisher(db *gorm.DB, bookID uint) (*Book, error) {
	var book Book
	result := db.Preload("Publisher").First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func getBookWithAuthors(db *gorm.DB, bookID uint) (*Book, error) {
	var book Book
	result := db.Preload("Authors").First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func listBooksOfAuthor(db *gorm.DB, authorID uint) ([]Book, error) {
	var books []Book
	result := db.Joins("JOIN author_books on author_books.book_id = books.id").
		Where("author_books.author_id = ?", authorID).
		Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

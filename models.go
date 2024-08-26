package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string
	Author      string
	Description string
	Price       uint
}

func createBook(db *gorm.DB, book *Book) {
	result := db.Create(book)

	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}

	fmt.Println("Create Book Successful")
}

func getBook(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id)

	if result.Error != nil {
		log.Fatalf("Error get book: %v", result.Error)
	}
	return &book
}

func updateBook(db *gorm.DB, book *Book) {
	result := db.Save(&book)

	if result.Error != nil {
		log.Fatalf("Update book failed: %v", result.Error)
	}

	fmt.Println("Update Book Successful")
}

func deleteBook(db *gorm.DB, id uint) {
	var book Book
	// result := db.Delete(&book, id) //soft delete
	result := db.Unscoped().Delete(&book, id) //Permanent delete

	if result.Error != nil {
		log.Fatalf("Delete book failed: %v", result.Error)
	}

	fmt.Println("Delete Book Successful")
}

func searchBook(db *gorm.DB, bookName string) []Book {
	var books []Book

	result := db.Where("name = ?", bookName).Order("price").Find(&books)
	if result.Error != nil {
		log.Fatalf("Search book failed: %v", result.Error)
	}
	return books
}

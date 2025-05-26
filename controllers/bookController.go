package controllers

import (
	"net/http"

	"minha-api/models"

	"github.com/gin-gonic/gin"
)

var books = []models.Book{
	{ID: "1", Title: "Dom Casmurro", Author: "Machado de Assis"},
	{ID: "2", Title: "O Hobbit", Author: "J.R.R. Tolkien"},
}

// GET /books
func GetBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

// GET /books/:id
func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	for _, b := range books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Livro não encontrado"})
}

// POST /books
func CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

// DELETE /books/:id
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Livro deletado"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Livro não encontrado"})
}

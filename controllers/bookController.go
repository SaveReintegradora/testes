package controllers

import (
	"minha-api/models"
	"minha-api/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookController struct {
	repo *repositories.BookRepository
}

func NewBookController(repo *repositories.BookRepository) *BookController {
	return &BookController{repo: repo}
}

// GetBooks - Lista todos os livros
func (c *BookController) GetBooks(ctx *gin.Context) {
	books, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar livros"})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

// GetBookByID - Busca livro por ID
func (c *BookController) GetBookByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	book, err := c.repo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

// CreateBook - Cria um novo livro
func (c *BookController) CreateBook(ctx *gin.Context) {
	var newBook models.Book

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if newBook.Title == "" || newBook.Author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Título e autor são obrigatórios"})
		return
	}

	newBook.ID = uuid.New().String()
	newBook.CreatedAt = time.Now()

	if err := c.repo.Create(&newBook); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar livro"})
		return
	}

	ctx.JSON(http.StatusCreated, newBook)
}

// UpdateBook - Atualiza um livro
func (c *BookController) UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedBook models.Book
	if err := ctx.ShouldBindJSON(&updatedBook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	updatedBook.ID = id

	if err := c.repo.Update(&updatedBook); err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar livro"})
		}
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

// DeleteBook - Remove um livro
func (c *BookController) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.repo.Delete(id); err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar livro"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

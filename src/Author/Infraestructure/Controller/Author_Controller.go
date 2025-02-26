package controller

import (
	application "Client/src/Author/Application"
	entities "Client/src/Author/Domain/Entities"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)


type AuthorController struct {
	service       *application.AuthorService
	longPollMutex sync.Mutex
	longPollChans []chan []entities.Author
}

func NewAuthorController(service *application.AuthorService) *AuthorController {
	return &AuthorController{
		service:       service,
		longPollChans: []chan []entities.Author{},
	}
}

// Crear un autor
func (c *AuthorController) CreateAuthor(ctx *gin.Context) {
	var author entities.Author
	if err := ctx.ShouldBindJSON(&author); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := c.service.CreateAuthor(&author)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar a los clientes sobre el cambio
	c.NotifyAuthorChanges()

	ctx.JSON(http.StatusCreated, gin.H{"message": "Author Created"})
}

// Obtener todos los autores (Short Polling)
func (c *AuthorController) GetAllAuthors(ctx *gin.Context) {
	authors, err := c.service.GetAllAuthor()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, authors)
}

// Obtener un autor por ID (Short Polling)
func (c *AuthorController) GetAuthorByID(ctx *gin.Context) {
	id := ctx.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}
	author, err := c.service.GetAuthorByID(int16(num))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, author)
}

// Actualizar un autor
func (c *AuthorController) UpdateAuthor(ctx *gin.Context) {
	id := ctx.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var author entities.Author
	if err := ctx.ShouldBindJSON(&author); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Entrada inválida"})
		return
	}
	author.ID = authorID

	err = c.service.UpdateAuthor(&author)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar a los clientes sobre el cambio
	c.NotifyAuthorChanges()

	ctx.JSON(http.StatusOK, gin.H{"message": "Autor actualizado"})
}

// Eliminar un autor
func (c *AuthorController) DeleteAuthor(ctx *gin.Context) {
	id := ctx.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	err = c.service.DeleteAuthor(int16(num))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar a los clientes sobre el cambio
	c.NotifyAuthorChanges()

	ctx.JSON(http.StatusOK, gin.H{"message": "Author Deleted"})
}

// --- SHORT POLLING ---

// Short Polling para obtener la lista de autores
func (c *AuthorController) ShortPollingAuthors(ctx *gin.Context) {
	authors, err := c.service.GetAllAuthor()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, authors)
}

// Short Polling para verificar si existe un autor por ID
func (c *AuthorController) ShortPollingAuthorByID(ctx *gin.Context) {
	id := ctx.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	author, err := c.service.GetAuthorByID(int16(authorID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	ctx.JSON(http.StatusOK, author)
}

// --- LONG POLLING ---

// Long Polling para obtener autores cuando haya cambios
func (c *AuthorController) LongPollingAuthors(ctx *gin.Context) {
	ch := make(chan []entities.Author)

	c.longPollMutex.Lock()
	c.longPollChans = append(c.longPollChans, ch)
	c.longPollMutex.Unlock()

	select {
	case authors := <-ch:
		ctx.JSON(http.StatusOK, authors)
	case <-time.After(30 * time.Second): // Expira si no hay cambios en 30s
		ctx.JSON(http.StatusNoContent, nil)
	}
}

// Notificar cambios a los clientes que están esperando con Long Polling
func (c *AuthorController) NotifyAuthorChanges() {
	c.longPollMutex.Lock()
	defer c.longPollMutex.Unlock()

	authors, err := c.service.GetAllAuthor()
	if err != nil {
		return
	}

	for _, ch := range c.longPollChans {
		ch <- authors
	}
	c.longPollChans = []chan []entities.Author{} // Limpiar canales
}

// Long Polling para obtener información de un autor específico cuando haya cambios
func (c *AuthorController) LongPollingAuthorByID(ctx *gin.Context) {
	id := ctx.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	ch := make(chan *entities.Author)

	go func() {
		for {
			author, err := c.service.GetAuthorByID(int16(authorID))
			if err == nil {
				ch <- author
				return
			}
			time.Sleep(2 * time.Second) 
		}
	}()

	select {
	case author := <-ch:
		ctx.JSON(http.StatusOK, author)
	case <-time.After(30 * time.Second): // Expira si no hay cambios en 30s
		ctx.JSON(http.StatusNoContent, nil)
	}
}

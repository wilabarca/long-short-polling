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

	c.NotifyAuthorChanges()

	ctx.JSON(http.StatusCreated, gin.H{"message": "Author Created"})
}


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
	c.longPollChans = nil 
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

// Obtener un autor por ID
func (c *AuthorController) GetAuthorByID(ctx *gin.Context) {
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

// Actualizar un autor
func (c *AuthorController) UpdateAuthor(ctx *gin.Context) {
	id := ctx.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author entities.Author
	if err := ctx.ShouldBindJSON(&author); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	author.ID = authorID

	err = c.service.UpdateAuthor(&author)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar cambios
	c.NotifyAuthorChanges()

	ctx.JSON(http.StatusOK, gin.H{"message": "Author updated"})
}

// Eliminar un autor
func (c *AuthorController) DeleteAuthor(ctx *gin.Context) {
	id := ctx.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	err = c.service.DeleteAuthor(int16(authorID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar cambios
	c.NotifyAuthorChanges()

	ctx.JSON(http.StatusOK, gin.H{"message": "Author deleted"})
}



// Obtener lista de autores (Short Polling)
func (c *AuthorController) ShortPollingAuthors(ctx *gin.Context) {
	authors, err := c.service.GetAllAuthor()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, authors)
}

// Verificar si existe un autor por ID (Short Polling)
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
		
		ctx.JSON(http.StatusOK, gin.H{"message": "No hay cambios detectados después de 30 segundos"})
	}
}

// Long Polling para obtener información de un autor específico cuando haya cambios
func (c *AuthorController) LongPollingAuthorByID(ctx *gin.Context) {
	id := ctx.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de autor inválido"})
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
		
		ctx.JSON(http.StatusOK, gin.H{"message": "No hay cambios detectados después de 30 segundos"})
	}
}

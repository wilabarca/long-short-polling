package router

import (
	controller "Client/src/Author/Infraestructure/Controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorRoutes(router *gin.Engine, AuthorController *controller.AuthorController) {
	AuthorGroup := router.Group("/Author")
	{
		// Rutas para operaciones CRUD
		AuthorGroup.GET("/", AuthorController.GetAllAuthors)
		AuthorGroup.GET("/:id", AuthorController.GetAuthorByID)
		AuthorGroup.POST("/", AuthorController.CreateAuthor)
		AuthorGroup.PUT("/:id", AuthorController.UpdateAuthor)
		AuthorGroup.DELETE("/:id", AuthorController.DeleteAuthor)

		// Short Polling para obtener la lista de autores
		AuthorGroup.GET("/shortPolling", AuthorController.ShortPollingAuthors)

		// Short Polling para verificar un autor por ID
		AuthorGroup.GET("/shortPolling/:id", AuthorController.ShortPollingAuthorByID)

		// Long Polling para obtener la lista de autores cuando haya cambios
		AuthorGroup.GET("/longPolling", AuthorController.LongPollingAuthors)

		// Long Polling para obtener un autor específico cuando haya cambios
		AuthorGroup.GET("/longPolling/:id", AuthorController.LongPollingAuthorByID)
	}
}
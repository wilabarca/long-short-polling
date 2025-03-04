package router

import (
	controller "Client/src/Author/Infraestructure/Controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorRoutes(router *gin.Engine, authorController *controller.AuthorController) {
	authorGroup := router.Group("/author")
	{
		// Rutas para operaciones CRUD
		authorGroup.GET("/", authorController.GetAllAuthors)
		authorGroup.GET("/:id", authorController.GetAuthorByID)
		authorGroup.POST("/", authorController.CreateAuthor)
		authorGroup.PUT("/:id", authorController.UpdateAuthor)
		authorGroup.DELETE("/:id", authorController.DeleteAuthor)

		// Short Polling para obtener la lista de autores
		authorGroup.GET("/shortPolling", authorController.ShortPollingAuthors)

		// Short Polling para verificar un autor por ID
		authorGroup.GET("/shortPolling/:id", authorController.ShortPollingAuthorByID)

		// Long Polling para obtener la lista de autores cuando haya cambios
		authorGroup.GET("/longPolling", authorController.LongPollingAuthors)

		// Long Polling para obtener un autor específico cuando haya cambios
		authorGroup.GET("/longPolling/:id", authorController.LongPollingAuthorByID)
	}
}

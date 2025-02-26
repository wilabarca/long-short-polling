package main

import (
	
	AuthorService "Client/src/Author/Application"
	AuthorControlller "Client/src/Author/Infraestructure/Controller"
	AuthorDb "Client/src/Author/Infraestructure/Database"
	AuthorRouter "Client/src/Author/Infraestructure/Router"
     core "Client/Core"
	BookService "Client/src/Book/Application"
	BookController "Client/src/Book/Infraestructure/Controller"
	BookDb "Client/src/Book/Infraestructure/Database"
	BookRouter "Client/src/Book/Infraestructure/Router"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
    db, err := core.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
        
        return
        
    }
    defer db.Close()

	

	AuthorRepo :=  AuthorDb.NewsqlAuthorRepository(db)
	AuthorService :=  AuthorService.NewAuthorService(AuthorRepo)
	AuthorControlller := AuthorControlller.NewAuthorController(AuthorService)

	BookRepo := BookDb.NewsqlBookRepository(db)
    BookService := BookService.NewBookService(BookRepo)
    BookController := BookController.NewBookController(BookService)


	router := gin.Default()


	AuthorRouter.RegisterAuthorRoutes(router , AuthorControlller)
	BookRouter.RegisterBookRoutes(router, BookController)


	
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
        println(err)
		
	}

	
}
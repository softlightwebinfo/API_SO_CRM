package main

import (
	"so-crm/src/apps"
	"so-crm/src/controller"
	"so-crm/src/middleware"
	"so-crm/src/service"
)

func main() {
	apps.
		NewApp().
		Routes(func(app *apps.App) {
			directoryController := controller.DirectoryController{}
			todoController := controller.TodoController{}

			authController := controller.LoginController{
				JWtService:   service.JWTAuthService(),
				LoginService: service.NewLoginService(),
			}

			//ROUTES Directory
			app.GetApp().Group("/user/directory").
				Use(middleware.AuthorizeJWT()).
				GET("/", directoryController.GetAll()).
				GET("/:id", directoryController.GetAllNotes()).
				POST("/", directoryController.CreateDirectory()).
				POST("/:id", directoryController.CreateNote()).
				DELETE("/:id/note/:note", directoryController.DeleteNote()).
				PUT("/:id/note/:note", directoryController.UpdateNote())

			// Routes TODO
			app.GetApp().Group("/user/todo").
				Use(middleware.AuthorizeJWT()).
				GET("/", todoController.GetAll())
			app.GetApp().POST("/login", authController.Auth())
		}).
		Start()
}

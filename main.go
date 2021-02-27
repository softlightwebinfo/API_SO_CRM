package main

import (
	"so-crm/src/apps"
	"so-crm/src/controller"
	"so-crm/src/middleware"
	"so-crm/src/service"
)

func main() {
	app := apps.NewApp()
	app.Routes(func(app *apps.App) {
		directoryController := controller.DirectoryController{}
		authController := controller.LoginController{
			JWtService:   service.JWTAuthService(),
			LoginService: service.NewLoginService(),
		}
		//ROUTES
		app.GetApp().GET("/user/directory", middleware.AuthorizeJWT(), directoryController.GetAll())
		app.GetApp().GET("/user/directory/:id", middleware.AuthorizeJWT(), directoryController.GetAllNotes())
		app.GetApp().POST("/login", authController.Auth())
	})
	app.Start()
}

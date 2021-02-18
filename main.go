package main

import (
	"github.com/gin-gonic/gin"
	"so-crm/src/controller"
)

func main() {
	r := gin.Default()

	directoryController := controller.DirectoryController{}

	//ROUTES
	r.GET("/user/directory", directoryController.GetAll())
	r.GET("/user/directory/:id", directoryController.GetAllNotes())
	_ = r.Run()
}

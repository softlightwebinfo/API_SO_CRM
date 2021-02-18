package controller

import "github.com/gin-gonic/gin"

type DirectoryController struct {
}

func (s *DirectoryController) GetAll() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (s *DirectoryController) GetAllNotes() gin.HandlerFunc {

}

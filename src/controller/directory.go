package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"so-crm/src/entity"
	"so-crm/src/models"
	"so-crm/src/response"
	"so-crm/src/service"
)

type DirectoryController struct {
	User models.AuthUser
}

func (s *DirectoryController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := service.AuthServiceLogin(c)
		var notesDirectory []response.NoteDirectory

		sqlNotes, _ := models.DB.Query().
			Select("count(id)").
			From("notes").
			WhereX("fk_directory", "=", "note_directories.id").
			Build()

		sql, params := models.DB.Query().
			Select("id, title, updated_at").
			SelectSubQuery(sqlNotes, "notes").
			From("note_directories").
			Where("fk_user_id", "=", user.ID).
			Build()

		query, err := models.DB.Query().DB.Query(sql, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		for query.Next() {
			directory := response.NoteDirectory{}
			_ = query.Scan(
				&directory.ID,
				&directory.Title,
				&directory.UpdatedAt,
				&directory.Notes,
			)
			notesDirectory = append(notesDirectory, directory)
		}
		c.JSON(http.StatusOK, notesDirectory)
	}
}

func (s *DirectoryController) GetAllNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		user := service.AuthServiceLogin(c)
		var notes []entity.Note

		sql, params := models.DB.Query().
			Select("notes.id, notes.title, notes.description, notes.updated_at, notes.fk_directory").
			From("notes").
			JOIN("note_directories", "notes.fk_directory=note_directories.id").
			Where("notes.fk_directory", "=", param).
			AndWhere("note_directories.fk_user_id", "=", user.ID).
			Build()

		query, err := models.DB.Query().DB.Query(sql, params...)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		for query.Next() {
			note := entity.Note{}
			_ = query.Scan(
				&note.ID,
				&note.Title,
				&note.Description,
				&note.UpdatedAt,
				&note.FkDirectory,
			)
			notes = append(notes, note)
		}
		c.JSON(http.StatusOK, notes)
	}
}

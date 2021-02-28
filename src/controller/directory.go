package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"so-crm/src/entity"
	"so-crm/src/models"
	"so-crm/src/models/request"
	"so-crm/src/response"
	"so-crm/src/service"
	"strconv"
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
			Where("fk_user_id", "=", user.Id).
			Build()

		query, err := models.DB.Query().DB.Query(sql, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		for query.Next() {
			directory := response.NoteDirectory{}
			_ = query.Scan(
				&directory.Id,
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
		var notes []entity.Notes

		sql, params := models.DB.Query().
			Select("notes.id, notes.title, notes.description, notes.updated_at, note_directories.id").
			From("notes").
			JOIN("note_directories", "notes.fk_directory=note_directories.id").
			Where("notes.fk_directory", "=", param).
			AndWhere("note_directories.fk_user_id", "=", user.Id).
			Build()
		query, err := models.DB.Query().DB.Query(sql, params...)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		for query.Next() {
			note := entity.Notes{}
			_ = query.Scan(
				&note.Id,
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

func (s *DirectoryController) CreateDirectory() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := service.AuthServiceLogin(c)
		directory := entity.NoteDirectories{}
		err := c.BindJSON(&directory)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		directory.FkUserId = user.Id
		errSave := models.DB.Mutation().Return(true).Save(&directory)
		if errSave != nil {
			c.JSON(http.StatusInternalServerError, errSave)
			return
		}
		c.JSON(http.StatusOK, directory)
	}
}

func (s *DirectoryController) CreateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		fkDirectoryId, _ := strconv.Atoi(c.Param("id"))
		note := entity.Notes{
			FkDirectory: int64(fkDirectoryId),
			Title:       "Escribe un titulo",
			Description: "Escribe una descripcion",
		}
		_ = models.DB.Mutation().Return(true).Save(&note)
		c.JSON(http.StatusOK, note)
	}
}

func (s *DirectoryController) DeleteNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		directoryParam := c.Param("id")
		noteParam := c.Param("note")
		directoryId, _ := strconv.Atoi(directoryParam)
		noteId, _ := strconv.Atoi(noteParam)
		note := entity.Notes{FkDirectory: int64(directoryId)}
		note.Id = int64(noteId)

		_ = models.DB.Mutation().Delete(&note)

		c.JSON(http.StatusOK, note)
	}
}

func (s *DirectoryController) UpdateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		updateNote := request.NoteUpdate{}
		if err := c.BindJSON(&updateNote); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		directoryParam := c.Param("id")
		noteParam := c.Param("note")

		directoryId, _ := strconv.Atoi(directoryParam)
		noteId, _ := strconv.Atoi(noteParam)

		note := entity.Notes{FkDirectory: int64(directoryId)}
		note.Id = int64(noteId)

		noteFind, err := models.DB.Query().
			Select("id, created_at, updated_at, deleted_at, title, description, fk_directory").
			From("notes").
			Where("id", "=", note.Id).
			AndWhere("fk_directory", "=", note.FkDirectory).
			Find(&note)

		newNote := noteFind.(*entity.Notes)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if updateNote.Column == "title" {
			newNote.Title = updateNote.Value
		} else if updateNote.Column == "description" {
			newNote.Description = updateNote.Value
		}

		_ = models.DB.Mutation().Return(true).
			Where("FkDirectory", "=", newNote.FkDirectory).
			Save(newNote)

		c.JSON(http.StatusOK, newNote)
	}
}

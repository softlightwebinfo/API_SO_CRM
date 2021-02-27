package service

import (
	"so-crm/src/entity"
	"so-crm/src/models"
)

type LoginService interface {
	LoginUser(email string, password string) (bool, *entity.User)
}

type loginInformation struct{}

func NewLoginService() LoginService {
	return &loginInformation{}
}

func (info *loginInformation) LoginUser(email string, password string) (bool, *entity.User) {
	user := entity.User{}
	find, _ := models.DB.Query().
		From("users").
		Where("email", "=", email).
		Find(&user)
	userFind := find.(*entity.User)
	return userFind.Password == password, userFind
}

package handlers
import (
	"DigiPassAuthenticationApi/services"
)

type AccountUsersHandler struct {
	service *services.AccountUsersService
}

func NewAccountUsers(service *services.AccountUsersService) *AccountUsersHandler{
	return &AccountUsersHandler{service: service}
}
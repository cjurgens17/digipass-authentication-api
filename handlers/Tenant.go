package handlers

import(
	"DigiPassAuthenticationApi/services"
)

type TenantHandler struct {
	service *services.TenantService
}

func NewTenantHandler(service *services.TenantService) *TenantHandler{
	return &TenantHandler{service: service}
}
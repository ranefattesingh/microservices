package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranefattesingh/ecommerce-platform/user-service/handler/dto"
	"github.com/ranefattesingh/ecommerce-platform/user-service/service"
	"github.com/ranefattesingh/ecommerce-platform/user-service/service/models"
)

type AccountsHandler interface {
	CreateAccount(c *gin.Context)
}

type accountsHandler struct {
	accountsService service.AccountsService
}

func NewAccountsHandler(as service.AccountsService) *accountsHandler {
	return &accountsHandler{
		accountsService: as,
	}
}

func (ah *accountsHandler) CreateAccount(c *gin.Context) {
	var account dto.CreateAccount

	if err := c.BindJSON(&account); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	model := models.Account{
		Email:    account.Email,
		Password: account.Password,
	}

	publicID, err := ah.accountsService.CreateAccount(c, model)
	if err != nil {
		// IDENTIFY ERROR and RETURN

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": publicID})
}

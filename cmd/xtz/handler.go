package xtz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiln-mid/pkg/delegations"
	"github.com/kiln-mid/pkg/models"
)

// Handler represent the handler of delegationsRepository
type Handler struct {
	DelegationsClient *delegations.Client
}

// RegisterRouter expose all endpoint for the `xtz` group.
func (a *Handler) RegisterRouter(router *gin.Engine) {
	delegationsRouter := router.Group("/xtz")

	delegationsRouter.GET("/delegations", a.getLastDelegations)
}

// Response represent the response gived by to the client
type Response struct {
	Data []models.Delegations `json:"data"`
	Page int                  `json:"Page"`
}

// getLastDelegations return all last delegations found if no query params are found.
// If a `year` param is provided, it will search all delegations based on the year provided.
// page and limit param try to mitigate the volume of data returned to the client.
func (a *Handler) getLastDelegations(c *gin.Context) {
	var queryParams struct {
		Year  int `form:"year" binding:"omitempty,min=1000,max=9999"`
		Page  int `form:"page" binding:"omitempty,min=1"`
		Limit int `form:"limit" binding:"omitempty,min=1,max=5000"`
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Check Year field is valid and follow the following format `YYYY`",
		})
		return
	}

	if queryParams.Page == 0 {
		queryParams.Page = 1
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 100
	}

	delegations, err := a.DelegationsClient.GetDelegations(c.Request.Context(), queryParams.Year, queryParams.Page, queryParams.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := Response{
		Data: *delegations,
		Page: queryParams.Page,
	}

	c.JSON(http.StatusOK, response)
}

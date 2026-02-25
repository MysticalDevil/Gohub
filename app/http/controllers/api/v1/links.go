package v1

import (
	"gohub/app/models/link"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

type LinksController struct {
	BaseAPIController
}

func (ctrl *LinksController) Index(c *gin.Context) {
	links := link.AllCached(c.Request.Context())
	response.Data(c, links)
}

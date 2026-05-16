package v1

import (
	"github.com/gin-gonic/gin"
	"gohub/app/models/link"
	"gohub/pkg/response"
)

type LinksController struct {
	BaseAPIController
}

func (ctrl *LinksController) Index(c *gin.Context) {
	links, err := link.AllCached(c.Request.Context())
	if err != nil {
		response.Abort500(c, "Failed to get links")
		return
	}
	response.Data(c, links)
}

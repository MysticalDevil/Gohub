package v1

import (
	"github.com/gin-gonic/gin"
	"gohub/app/models/category"
	"gohub/app/requests"
	"gohub/pkg/response"
)

type CategoriesController struct {
	BaseAPIController
}

func (ctrl *CategoriesController) Store(c *gin.Context) {
	request := requests.CategoryRequest{}
	if ok := requests.Validate(c, &request, requests.CategorySave); !ok {
		return
	}

	categoryModel := category.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	categoryModel.Create()
	if categoryModel.ID > 0 {
		response.Created(c, categoryModel)
	} else {
		response.Abort500(c, "Failed to create, please try later ~")
	}
}

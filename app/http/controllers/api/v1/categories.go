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

	categoryModel.Create(c.Request.Context())
	if categoryModel.ID > 0 {
		response.Created(c, categoryModel)
	} else {
		response.Abort500(c, "Failed to create, please try later ~")
	}
}

func (ctrl *CategoriesController) Update(c *gin.Context) {
	// Verify that the url parameter id is correct
	categoryModel := category.Get(c.Request.Context(), c.Param("id"))
	if categoryModel.ID == 0 {
		response.Abort404(c)
		return
	}

	request := requests.CategoryRequest{}
	if ok := requests.Validate(c, &request, requests.CategorySave); !ok {
		return
	}

	categoryModel.Name = request.Name
	categoryModel.Description = request.Description
	rowsAffected := categoryModel.Save(c.Request.Context())

	if rowsAffected > 0 {
		response.Data(c, categoryModel)
	} else {
		response.Abort500(c)
	}
}

func (ctrl *CategoriesController) Index(c *gin.Context) {
	request := requests.PaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}

	data, pager := category.Paginate(c.Request.Context(), c, 10)
	response.Paginated(c, data, pager)
}

func (ctrl *CategoriesController) Delete(c *gin.Context) {
	categoryModel := category.Get(c.Request.Context(), c.Param("id"))
	if categoryModel.ID == 0 {
		response.Abort404(c)
		return
	}

	rowsAffected := categoryModel.Delete(c.Request.Context())
	if rowsAffected > 0 {
		response.Success(c)
		return
	}

	response.Abort500(c, "Deletion failed, please try later~")
}

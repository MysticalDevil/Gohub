package v1

import (
	"github.com/gin-gonic/gin"
	"gohub/app/models/category"
	"gohub/app/policies"
	"gohub/app/requests"
	"gohub/pkg/auth"
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
		UserID:      auth.CurrentUID(c),
	}

	if err := categoryModel.Create(c.Request.Context()); err != nil {
		response.Abort500(c, "Failed to create, please try later ~")
		return
	}
	response.Created(c, categoryModel)
}

func (ctrl *CategoriesController) Update(c *gin.Context) {
	categoryModel, err := category.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Abort500(c, "Failed to query category")
		return
	}
	if categoryModel.ID == 0 {
		response.Abort404(c)
		return
	}

	if ok := policies.CanModifyCategory(c, categoryModel); !ok {
		response.Abort403(c)
		return
	}

	request := requests.CategoryRequest{}
	if ok := requests.Validate(c, &request, requests.CategorySave); !ok {
		return
	}

	categoryModel.Name = request.Name
	categoryModel.Description = request.Description
	rowsAffected, err := categoryModel.UpdateFields(c.Request.Context(), request.Name, request.Description)
	if err != nil {
		response.Abort500(c)
		return
	}

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
	categoryModel, err := category.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Abort500(c, "Failed to query category")
		return
	}
	if categoryModel.ID == 0 {
		response.Abort404(c)
		return
	}

	if ok := policies.CanModifyCategory(c, categoryModel); !ok {
		response.Abort403(c)
		return
	}

	rowsAffected, err := categoryModel.Delete(c.Request.Context())
	if err != nil {
		response.Abort500(c, "Deletion failed, please try later~")
		return
	}
	if rowsAffected > 0 {
		response.Success(c)
		return
	}

	response.Abort500(c, "Deletion failed, please try later~")
}

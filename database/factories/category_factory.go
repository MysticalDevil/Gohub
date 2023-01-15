package factories

import (
	"github.com/go-faker/faker/v4"
	"gohub/app/models/category"
)

func MakeCategories(count int) []category.Category {
	var objs []category.Category

	faker.SetGenerateUniqueValues(true)

	for i := 0; i < count; i++ {
		categoryModel := category.Category{
			Name:        faker.Username(),
			Description: faker.Sentence(),
		}
		objs = append(objs, categoryModel)
	}

	return objs
}

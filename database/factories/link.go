package factories

import (
	"github.com/go-faker/faker/v4"
	"gohub/app/models/link"
)

func MakeLinks(count int) []link.Link {
	var objs []link.Link

	for range count {
		linkModel := link.Link{
			Name: faker.Username(),
			URL:  faker.URL(),
		}
		objs = append(objs, linkModel)
	}

	return objs
}

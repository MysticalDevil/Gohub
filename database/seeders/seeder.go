// Package seeders Storing data fill files
package seeders

import "gohub/pkg/seed"

func Initialize() {
	seed.SetRunOrder([]string{
		"SeedUsersTable",
	})
}

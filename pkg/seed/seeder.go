// Package seed Handling database populating logic
package seed

import "gorm.io/gorm"

// Store all Seeder
var seeders []Seeder

// Sequential execution of the Seeder array
// Support some seeders that must be executed in order,
// for example, topic creation must depend on user,
// so the TopicSeeder should be executed after the UserSeeder
var orderSeederNames []string

type SeederFunc func(*gorm.DB)

// Seeder Corresponds to the Seeder file in each database/seeders directory
type Seeder struct {
	Func SeederFunc
	Name string
}

// Add Register to the seeders array
func Add(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Name: name,
		Func: fn,
	})
}

// SetRunOrder Set [seeder array for sequential execution]
func SetRunOrder(names []string) {
	orderSeederNames = names
}

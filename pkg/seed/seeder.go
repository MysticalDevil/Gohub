// Package seed Handling database populating logic
package seed

import (
	"gohub/pkg/console"
	"gohub/pkg/database"
	"gorm.io/gorm"
)

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

// GetSeeder Get Seeder object by name
func GetSeeder(name string) Seeder {
	for _, sdr := range seeders {
		if name == sdr.Name {
			return sdr
		}
	}
	return Seeder{}
}

// RunAll Run all Seeder
func RunAll() {
	// Run the ordered first
	executed := make(map[string]string)
	for _, name := range orderSeederNames {
		sdr := GetSeeder(name)
		if len(sdr.Name) > 0 {
			console.Warning("Running Ordered Seeder: " + sdr.Name)
			sdr.Func(database.DB)
			executed[name] = name
		}
	}

	for _, sdr := range seeders {
		// Filtering already running
		if _, ok := executed[sdr.Name]; !ok {
			console.Warning("Running Seeder: " + sdr.Name)
			sdr.Func(database.DB)
		}
	}
}

// RunSeeder Run single Seeder
func RunSeeder(name string) {
	for _, sdr := range seeders {
		if name == sdr.Name {
			sdr.Func(database.DB)
			break
		}
	}
}

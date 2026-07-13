// internal/database/models/schedule.go
package models

const DBName = "config.db"

type TableRepresenter interface {
	TableName() string
}

type Schedule struct {
	ID         int    `db:"id" pk:"true"`
	Expression string `db:"expression"`
}

func (Schedule) TableName() string {
	return "schedules"
}

// GetSeedDataList now returns a slice of rows.
// You can add as many starting configurations here as you want.
func (Schedule) GetSeedDataList() []Schedule {
	return []Schedule{
		{Expression: "@every 60s"},  // Seed 1: Cache Sync
	}
}

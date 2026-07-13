package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fuadop/beacon/internal/database/ddl"
	"github.com/fuadop/beacon/internal/database/dml"
	"github.com/fuadop/beacon/internal/database/models"

	_ "modernc.org/sqlite"
)

func InitAndSeedDB() (*sql.DB, error) {
	// 1. Resolve OS Temp Directory + '.beacon-assets'
	assetsDir := filepath.Join(os.TempDir(), ".beacon-assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create asset storage path: %w", err)
	}

	dbPath := filepath.Join(assetsDir, models.DBName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database file: %w", err)
	}

	// 2. Structural table creation via DDL
	if err := ddl.Create[models.Schedule](db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to process DDL for Schedule model: %w", err)
	}

	// 3. MULTIPLE SEEDS ZONE (Cleaner using generic Count)
	count, err := dml.Count[models.Schedule](db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to verify table row integrity: %w", err)
	}

	// Only seed if the table is totally fresh and empty
	if count == 0 {
		fmt.Println("[DB Seed] Table is empty. Injecting dataset collection payload...")

		// Fetch our slice of seeds from the model definition layer
		seeds := models.Schedule{}.GetSeedDataList()

		// Loop over the dataset dynamically
		for _, seedRecord := range seeds {
			if err := dml.Insert(db, &seedRecord); err != nil {
				db.Close()
				return nil, fmt.Errorf("failed execution on dynamic seeder payload item: %w", err)
			}
			fmt.Printf("  ✔ Seeded row successfully -> ID: %d | Schedule: %s\n", seedRecord.ID, seedRecord.Expression)
		}
		fmt.Printf("✔ Seeding complete. Successfully injected %d rows.\n", len(seeds))
	} else {
		var sample models.Schedule
		fmt.Printf("ℹ Table '%s' already contains %d rows. Seeding skipped.\n", sample.TableName(), count)
	}

	return db, nil
}

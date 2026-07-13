// cmd/main.go
package main

import (
	"fmt"
	"log"

	"github.com/fuadop/beacon/internal/database"
	"github.com/fuadop/beacon/internal/database/dml"
	"github.com/fuadop/beacon/internal/database/models"
)

func main() {
	// 1. Resolve OS Temp path directories, run structural DDL, and apply seeds
	db, err := database.InitAndSeedDB()
	if err != nil {
		log.Fatalf("❌ Critical Engine Boot Failure: %v", err)
	}
	defer db.Close()

	fmt.Println("🚀 Core database asset environment setup successfully passed.")

	// 2. Verification Step: Use your generic GetAll method to print current records
	activeSchedules, err := dml.GetAll[models.Schedule](db)
	if err != nil {
		log.Fatalf("❌ Failed to fetch database records via DML layer: %v", err)
	}

	fmt.Printf("\n[Verification Check] Current Database Records Found (%d):\n", len(activeSchedules))
	for _, item := range activeSchedules {
		fmt.Printf(" 🔹 ID: %d | Schedule Expression: %s\n", item.ID, item.Expression)
	}
	fmt.Println("\n✔ Database step verification complete.")
}

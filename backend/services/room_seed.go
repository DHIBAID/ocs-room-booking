package services

import (
	"errors"
	"log"
	"strings"

	"ocs-room-booking/db"
	"ocs-room-booking/models"

	"gorm.io/gorm"
)

type seedRoom struct {
	Name     string
	Capacity int
}

var roomSeedData = []seedRoom{
	{Name: "A-Class Room 320", Capacity: 80},
	{Name: "A-AUDITORIUM", Capacity: 289},
	{Name: "A-Class Room 111", Capacity: 70},
	{Name: "A-Class Room 112", Capacity: 80},
	{Name: "A-Class Room 114", Capacity: 36},
	{Name: "A-Class Room 117", Capacity: 84},
	{Name: "A-Class Room 118", Capacity: 84},
	{Name: "A-Class Room 119", Capacity: 108},
	{Name: "A-Class Room 220", Capacity: 40},
	{Name: "A-Class Room 221", Capacity: 120},
	{Name: "A-LH-1", Capacity: 184},
	{Name: "A-LH-2", Capacity: 184},
	{Name: "BT/BM-009", Capacity: 24},
	{Name: "BT/BM-010", Capacity: 24},
	{Name: "BT/BM-118", Capacity: 60},
	{Name: "C-LH-10", Capacity: 68},
	{Name: "C-LH-2", Capacity: 138},
	{Name: "C-LH-3", Capacity: 100},
	{Name: "C-LH-4", Capacity: 60},
	{Name: "C-LH-5", Capacity: 60},
	{Name: "C-LH-6", Capacity: 60},
	{Name: "C-LH-7", Capacity: 70},
	{Name: "C-LH-9", Capacity: 66},
	{Name: "CSE-LH-01", Capacity: 70},
	{Name: "CSE-LH-02", Capacity: 70},
	{Name: "CSE-LH-03", Capacity: 70},
	{Name: "CY-LH-1", Capacity: 30},
	{Name: "CY-LH-2", Capacity: 40},
	{Name: "CY-LH-3", Capacity: 90},
	{Name: "EE-004(GF)", Capacity: 80},
	{Name: "EE-20 (SF)", Capacity: 60},
	{Name: "LHC-01", Capacity: 72},
	{Name: "LHC-02", Capacity: 72},
	{Name: "LHC-03", Capacity: 120},
	{Name: "LHC-04", Capacity: 200},
	{Name: "LHC-05", Capacity: 800},
	{Name: "LHC-06", Capacity: 320},
	{Name: "LHC-07", Capacity: 200},
	{Name: "LHC-08", Capacity: 120},
	{Name: "LHC-09", Capacity: 72},
	{Name: "LHC-10", Capacity: 72},
	{Name: "LHC-11", Capacity: 120},
	{Name: "LHC-12", Capacity: 200},
	{Name: "LHC-13", Capacity: 320},
	{Name: "LHC-14", Capacity: 200},
	{Name: "LHC-15", Capacity: 120},
	{Name: "MA-01", Capacity: 56},
	{Name: "MA-02", Capacity: 56},
	{Name: "MA-114", Capacity: 30},
	{Name: "MSME-LH-1", Capacity: 36},
	{Name: "MSME-LH-2", Capacity: 60},
	{Name: "MSME-LH-3", Capacity: 106},
	{Name: "PH-1", Capacity: 80},
	{Name: "PH-2", Capacity: 60},
	{Name: "PH-3", Capacity: 50},
}

func SeedRooms() {
	created := 0
	updated := 0
	for _, item := range roomSeedData {
		block := inferBlock(item.Name)
		var existing models.Room
		err := db.DB.Where("name = ?", item.Name).First(&existing).Error
		if err == nil {
			update := map[string]interface{}{
				"block":    block,
				"capacity": item.Capacity,
				"status":   "available",
			}
			if err := db.DB.Model(&existing).Updates(update).Error; err != nil {
				log.Println("Failed to update room:", item.Name, err)
				continue
			}
			updated++
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Failed to query room:", item.Name, err)
			continue
		}
		room := models.Room{
			Block:    block,
			Name:     item.Name,
			Capacity: item.Capacity,
			Status:   "available",
		}
		if err := db.DB.Create(&room).Error; err != nil {
			log.Println("Failed to create room:", item.Name, err)
			continue
		}
		created++
	}
	log.Printf("Room seed complete. Created: %d, Updated: %d", created, updated)
}

func inferBlock(name string) string {
	prefixes := []string{"CSE-", "CY-", "MSME-", "LHC-", "EE-", "MA-", "PH-", "BT/BM-", "A-", "C-"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return strings.TrimSuffix(prefix, "-")
		}
	}
	if idx := strings.Index(name, "-"); idx > 0 {
		return name[:idx]
	}
	if strings.HasPrefix(name, "A ") {
		return "A"
	}
	return "General"
}

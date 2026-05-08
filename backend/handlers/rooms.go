package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ocs-room-booking/db"
	"ocs-room-booking/models"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

type createRoomRequest struct {
	Block           string `json:"block"`
	Name            string `json:"name"`
	Capacity        int    `json:"capacity"`
	Status          string `json:"status"`
	AllowedPurposes string `json:"allowed_purposes"`
	Notes           string `json:"notes"`
}

type updateRoomRequest struct {
	Block           string `json:"block"`
	Name            string `json:"name"`
	Capacity        *int   `json:"capacity"`
	Status          string `json:"status"`
	AllowedPurposes string `json:"allowed_purposes"`
	Notes           string `json:"notes"`
}

func CreateRoom(c *gin.Context) {
	var req createRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Block) == "" || strings.TrimSpace(req.Name) == "" || req.Capacity <= 0 {
		utils.JSONError(c, http.StatusBadRequest, "block, name, and capacity are required")
		return
	}
	allowed, err := utils.NormalizeAllowedPurposes(req.AllowedPurposes)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "available"
	}
	room := models.Room{
		Block:           req.Block,
		Name:            req.Name,
		Capacity:        req.Capacity,
		Status:          status,
		AllowedPurposes: allowed,
		Notes:           req.Notes,
	}
	if err := db.DB.Create(&room).Error; err != nil {
		utils.JSONError(c, http.StatusConflict, "room name already exists")
		return
	}
	c.JSON(http.StatusCreated, room)
}

func ListRooms(c *gin.Context) {
	var rooms []models.Room
	if err := db.DB.Order("block asc, name asc").Find(&rooms).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch rooms")
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func UpdateRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid room id")
		return
	}
	var req updateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	var room models.Room
	if err := db.DB.First(&room, id).Error; err != nil {
		utils.JSONError(c, http.StatusNotFound, "room not found")
		return
	}
	updates := map[string]interface{}{}
	if strings.TrimSpace(req.Block) != "" {
		updates["block"] = req.Block
	}
	if strings.TrimSpace(req.Name) != "" {
		updates["name"] = req.Name
	}
	if req.Capacity != nil {
		if *req.Capacity <= 0 {
			utils.JSONError(c, http.StatusBadRequest, "capacity must be positive")
			return
		}
		updates["capacity"] = *req.Capacity
	}
	if strings.TrimSpace(req.Status) != "" {
		updates["status"] = req.Status
	}
	if req.AllowedPurposes != "" {
		allowed, err := utils.NormalizeAllowedPurposes(req.AllowedPurposes)
		if err != nil {
			utils.JSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		updates["allowed_purposes"] = allowed
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := db.DB.Model(&room).Updates(updates).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to update room")
		return
	}
	c.JSON(http.StatusOK, room)
}

func DeleteRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid room id")
		return
	}
	if err := db.DB.Delete(&models.Room{}, id).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to delete room")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func ListAvailableRooms(c *gin.Context) {
	block := strings.TrimSpace(c.Query("block"))
	purpose := strings.TrimSpace(c.Query("purpose"))
	minCapStr := strings.TrimSpace(c.Query("minCapacity"))
	dateStr := strings.TrimSpace(c.Query("date"))
	startStr := strings.TrimSpace(c.Query("start_time"))
	endStr := strings.TrimSpace(c.Query("end_time"))
	startRFC := strings.TrimSpace(c.Query("start"))
	endRFC := strings.TrimSpace(c.Query("end"))

	var minCap int
	if minCapStr != "" {
		value, err := strconv.Atoi(minCapStr)
		if err != nil {
			utils.JSONError(c, http.StatusBadRequest, "invalid minCapacity")
			return
		}
		minCap = value
	}

	var start time.Time
	var end time.Time
	var timeFilter bool
	if startRFC != "" || endRFC != "" || dateStr != "" || startStr != "" || endStr != "" {
		parsedStart, parsedEnd, err := utils.ParseBookingTimes(dateStr, chooseTimeInput(startRFC, startStr), chooseTimeInput(endRFC, endStr))
		if err != nil {
			utils.JSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		start = parsedStart
		end = parsedEnd
		timeFilter = true
	}

	query := db.DB.Model(&models.Room{}).Where("status = ?", "available")
	if block != "" {
		query = query.Where("block = ?", block)
	}
	if minCap > 0 {
		query = query.Where("capacity >= ?", minCap)
	}
	var rooms []models.Room
	if err := query.Order("block asc, name asc").Find(&rooms).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch rooms")
		return
	}

	filtered := make([]models.Room, 0, len(rooms))
	for _, room := range rooms {
		if purpose != "" && !utils.PurposeAllowedForRoom(purpose, room.AllowedPurposes) {
			continue
		}
		if timeFilter {
			conflict, err := hasBookingConflict(room.ID, start, end)
			if err != nil {
				utils.JSONError(c, http.StatusInternalServerError, "failed to check conflicts")
				return
			}
			if conflict {
				continue
			}
		}
		filtered = append(filtered, room)
	}

	c.JSON(http.StatusOK, filtered)
}

func chooseTimeInput(primary string, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}

func hasBookingConflict(roomID uint, start time.Time, end time.Time) (bool, error) {
	var count int64
	err := db.DB.Model(&models.Booking{}).
		Where("room_id = ?", roomID).
		Where("status = ?", "confirmed").
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

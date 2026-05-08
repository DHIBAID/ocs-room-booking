package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ocs-room-booking/db"
	"ocs-room-booking/middleware"
	"ocs-room-booking/models"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

type createBookingRequest struct {
	RoomID       uint   `json:"room_id"`
	Purpose      string `json:"purpose"`
	Participants int    `json:"participants"`
	Date         string `json:"date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Start        string `json:"start"`
	End          string `json:"end"`
}

type updateBookingRequest struct {
	Status string `json:"status"`
}

func CreateBooking(c *gin.Context) {
	user, ok := middleware.GetAuthUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	if user.Role == utils.RoleViewer {
		utils.JSONError(c, http.StatusForbidden, "viewer role cannot create bookings")
		return
	}
	var req createBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RoomID == 0 || req.Participants <= 0 {
		utils.JSONError(c, http.StatusBadRequest, "room_id and participants are required")
		return
	}
	purpose, ok := utils.CanonicalPurpose(req.Purpose)
	if !ok {
		utils.JSONError(c, http.StatusBadRequest, "invalid purpose")
		return
	}
	start, end, err := parseBookingTimes(req)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.EnsureNotPast(start); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	var room models.Room
	if err := db.DB.First(&room, req.RoomID).Error; err != nil {
		utils.JSONError(c, http.StatusNotFound, "room not found")
		return
	}
	if room.Status != "available" {
		utils.JSONError(c, http.StatusBadRequest, "room is not available")
		return
	}
	if req.Participants > room.Capacity {
		utils.JSONError(c, http.StatusBadRequest, "room capacity exceeded")
		return
	}
	if !utils.PurposeAllowedForRoom(purpose, room.AllowedPurposes) {
		utils.JSONError(c, http.StatusBadRequest, "purpose not allowed for room")
		return
	}
	conflict, err := hasBookingConflict(room.ID, start, end)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to check conflicts")
		return
	}
	if conflict {
		utils.JSONError(c, http.StatusConflict, "room already booked for this time")
		return
	}
	booking := models.Booking{
		RoomID:       room.ID,
		UserID:       user.ID,
		Purpose:      purpose,
		Participants: req.Participants,
		StartTime:    start,
		EndTime:      end,
		Status:       "confirmed",
	}
	if err := db.DB.Create(&booking).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to create booking")
		return
	}
	c.JSON(http.StatusCreated, booking)
}

func ListMyBookings(c *gin.Context) {
	user, ok := middleware.GetAuthUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var bookings []models.Booking
	if err := db.DB.Preload("Room").Where("user_id = ?", user.ID).Order("start_time desc").Find(&bookings).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch bookings")
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func ListAllBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := db.DB.Preload("Room").Preload("User").Order("start_time desc").Find(&bookings).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch bookings")
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func UpdateBookingStatus(c *gin.Context) {
	user, ok := middleware.GetAuthUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid booking id")
		return
	}
	var req updateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	status := strings.ToLower(req.Status)
	if status != "cancelled" && status != "rejected" {
		utils.JSONError(c, http.StatusBadRequest, "only cancellation or rejection is supported")
		return
	}
	var booking models.Booking
	if err := db.DB.First(&booking, id).Error; err != nil {
		utils.JSONError(c, http.StatusNotFound, "booking not found")
		return
	}
	if status == "rejected" && user.Role != utils.RoleAdmin {
		utils.JSONError(c, http.StatusForbidden, "admin access required to reject bookings")
		return
	}
	if status == "cancelled" && user.Role != utils.RoleAdmin && booking.UserID != user.ID {
		utils.JSONError(c, http.StatusForbidden, "cannot modify this booking")
		return
	}
	if err := db.DB.Model(&booking).Update("status", status).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to cancel booking")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func parseBookingTimes(req createBookingRequest) (time.Time, time.Time, error) {
	if strings.TrimSpace(req.Start) != "" || strings.TrimSpace(req.End) != "" {
		start, err := time.Parse(time.RFC3339, req.Start)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		end, err := time.Parse(time.RFC3339, req.End)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		if !end.After(start) {
			return time.Time{}, time.Time{}, fmt.Errorf("end time must be after start time")
		}
		return start, end, nil
	}
	return utils.ParseBookingTimes(req.Date, req.StartTime, req.EndTime)
}

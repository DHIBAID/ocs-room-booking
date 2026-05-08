package utils

import (
	"fmt"
	"strings"
	"time"
)

func ParseBookingTimes(dateStr string, startStr string, endStr string) (time.Time, time.Time, error) {
	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("start and end time are required")
	}
	if strings.Contains(startStr, "T") || strings.Contains(endStr, "T") {
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start time")
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end time")
		}
		if !end.After(start) {
			return time.Time{}, time.Time{}, fmt.Errorf("end time must be after start time")
		}
		return start, end, nil
	}
	if dateStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("date is required")
	}
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.Local
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", dateStr, startStr), loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", dateStr, endStr), loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time")
	}
	if !end.After(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end time must be after start time")
	}
	return start, end, nil
}

func EnsureNotPast(start time.Time) error {
	loc := start.Location()
	now := time.Now().In(loc)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	if start.Before(startOfToday) {
		return fmt.Errorf("booking date cannot be in the past")
	}
	return nil
}

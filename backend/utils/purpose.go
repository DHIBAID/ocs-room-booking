package utils

import (
	"fmt"
	"strings"
)

const (
	RoleAdmin  = "admin"
	RoleCore   = "core"
	RoleViewer = "viewer"
)

var AllowedPurposes = []string{"OA", "Interview", "PPT"}

func CanonicalPurpose(input string) (string, bool) {
	value := strings.TrimSpace(input)
	switch strings.ToLower(value) {
	case "oa":
		return "OA", true
	case "interview", "interviews":
		return "Interview", true
	case "ppt", "pre-placement talk", "pre placement talk", "pre-placement talks":
		return "PPT", true
	default:
		return "", false
	}
}

func NormalizeAllowedPurposes(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", nil
	}
	parts := strings.Split(raw, ",")
	seen := map[string]bool{}
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		canonical, ok := CanonicalPurpose(part)
		if !ok {
			return "", fmt.Errorf("invalid purpose: %s", strings.TrimSpace(part))
		}
		if !seen[canonical] {
			seen[canonical] = true
			result = append(result, canonical)
		}
	}
	return strings.Join(result, ","), nil
}

func PurposeAllowedForRoom(purpose string, allowedRaw string) bool {
	if strings.TrimSpace(allowedRaw) == "" {
		return true
	}
	canonical, ok := CanonicalPurpose(purpose)
	if !ok {
		return false
	}
	parts := strings.Split(allowedRaw, ",")
	for _, part := range parts {
		if strings.EqualFold(strings.TrimSpace(part), canonical) {
			return true
		}
	}
	return false
}

func CanonicalRole(input string) (string, bool) {
	value := strings.TrimSpace(input)
	switch strings.ToLower(value) {
	case RoleAdmin:
		return RoleAdmin, true
	case RoleCore:
		return RoleCore, true
	case RoleViewer:
		return RoleViewer, true
	default:
		return "", false
	}
}

package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TitleCase(s string) string {
	titleCaser := cases.Title(language.English)
	return titleCaser.String(s)
}

// GetSessionEntityType returns the appropriate entity_type based on session_type
func GetSessionEntityType(sessionType string) string {
	if sessionType == SessionTypeCoaching {
		return EntityTLCoaching
	}
	return EntityTLBriefing
}

func UnderscoreToDash(entityType string) string {
	return strings.ReplaceAll(entityType, "_", "-")
}

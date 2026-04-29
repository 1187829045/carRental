package service

import (
	"fmt"
	"strings"
	"time"
)

func parsePage(pageStr, limitStr string) (page, limit int, offset int, err error) {
	page = 1
	limit = 10
	if pageStr != "" {
		_, err = fmt.Sscanf(pageStr, "%d", &page)
		if err != nil || page <= 0 {
			return 0, 0, 0, fmt.Errorf("invalid page")
		}
	}
	if limitStr != "" {
		_, err = fmt.Sscanf(limitStr, "%d", &limit)
		if err != nil || limit <= 0 {
			return 0, 0, 0, fmt.Errorf("invalid limit")
		}
	}
	offset = (page - 1) * limit
	return page, limit, offset, nil
}

func likeArg(v string) string { return "%" + strings.TrimSpace(v) + "%" }

func parseDateTime(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
	}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, time.Local); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("invalid datetime")
}

func ParseTimeFlexible(s string) (*time.Time, error) {
	return parseDateTime(s)
}

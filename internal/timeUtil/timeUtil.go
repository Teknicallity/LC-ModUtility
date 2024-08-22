package timeUtil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseTimeString(timeStr string) time.Time {
	now := time.Now()

	//timeStr: '18 hours ago'
	parts := strings.Split(timeStr, " ")
	if len(parts) < 3 {
		fmt.Printf("Invalid time format: %v\n", parts)
		return now
	}

	var numStr string
	if strings.Contains(strings.ToLower(parts[0]), "a") {
		numStr = "1"
	} else {
		numStr = parts[0]
	}

	unit := parts[1]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		fmt.Println("Error parsing time string:", err)
		return now
	}

	switch {
	case strings.Contains(unit, "second"):
		return now.Add(-time.Duration(num) * time.Second)
	case strings.Contains(unit, "minute"):
		return now.Add(-time.Duration(num) * time.Minute)
	case strings.Contains(unit, "hour"):
		return now.Add(-time.Duration(num) * time.Hour)
	case strings.Contains(unit, "day"):
		return now.AddDate(0, 0, -num)
	case strings.Contains(unit, "week"):
		return now.AddDate(0, 0, -7*num)
	case strings.Contains(unit, "month"):
		return now.AddDate(0, -num, 0)
	case strings.Contains(unit, "year"):
		return now.AddDate(-num, 0, 0)
	default:
		fmt.Println("Unsupported time unit:", unit)
		return now
	}
}

package utils

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseTime(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case primitive.DateTime:
		return v.Time(), nil
	case time.Time:
		return v, nil
	case string:
		return time.Parse(time.RFC3339, v)
	default:
		return time.Time{}, fmt.Errorf("invalid time format")
	}
}

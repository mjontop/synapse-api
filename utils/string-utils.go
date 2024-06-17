package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func GenerateSlug(title string) string {
	re := regexp.MustCompile(`[^\w\-]+`)
	slug := re.ReplaceAllString(strings.ToLower(title), "-")

	// Truncate to 30 characters
	if len(slug) > 30 {
		slug = slug[:30]
	}

	// Append a random 5-digit suffix
	suffix := fmt.Sprintf("%05d", rand.Intn(99999))
	slug += "-" + suffix

	return slug
}

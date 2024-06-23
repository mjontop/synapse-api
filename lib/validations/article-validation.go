package validations

import "strings"

func ValidateTags(tagList []string) []string {
	var validatedTags []string

	for _, tag := range tagList {
		validTag := strings.ToLower(tag)
		fields := strings.Fields(validTag)
		validTag = strings.Join(fields, "")
		validatedTags = append(validatedTags, validTag)
	}
	return validatedTags
}

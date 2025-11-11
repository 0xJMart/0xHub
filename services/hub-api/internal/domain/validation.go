package domain

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	slugRegex        = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	maxTitleLength   = 160
	maxSummaryLength = 512
	maxTagLength     = 64
)

func appendError(errs *ValidationErrors, field, message string) {
	*errs = append(*errs, ValidationError{
		Field:   field,
		Message: message,
	})
}

func validateSlug(field, value string, errs *ValidationErrors) {
	if strings.TrimSpace(value) == "" {
		appendError(errs, field, "cannot be empty")
		return
	}
	if len(value) > 128 {
		appendError(errs, field, "must be 128 characters or fewer")
		return
	}
	if !slugRegex.MatchString(value) {
		appendError(errs, field, "must contain lowercase letters, numbers, or hyphens")
	}
}

func validateURL(field, raw string, errs *ValidationErrors) {
	if strings.TrimSpace(raw) == "" {
		appendError(errs, field, "cannot be empty")
		return
	}
	parsed, err := url.Parse(raw)
	if err != nil || !parsed.IsAbs() || parsed.Host == "" {
		appendError(errs, field, "must be an absolute URL")
	}
}

func prefixErrors(prefix string, errs ValidationErrors) ValidationErrors {
	if len(errs) == 0 {
		return nil
	}
	prefixed := make(ValidationErrors, 0, len(errs))
	for _, err := range errs {
		field := prefix
		if err.Field != "" {
			field = prefix + "." + err.Field
		}
		prefixed = append(prefixed, ValidationError{
			Field:   field,
			Message: err.Message,
		})
	}
	return prefixed
}

package http

import (
	"github.com/synapsis-test/internal/category"
)

// formatCategory formats the given category
// into the respective HTTP-format object.
func formatCategory(c category.Category) (categoryHTTP, error) {
	return categoryHTTP{
		ID:          &c.ID,
		Name:        &c.Name,
		Description: &c.Description,
	}, nil
}

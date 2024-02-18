package category

import "errors"

// Followings are the known errors returned from category.
var (
	// ErrDataNotFound is returned when the desired data is
	// not found.
	ErrDataNotFound = errors.New("data not found")

	// ErrInvalidCategoryID is returned when the given category ID is
	// invalid.
	ErrInvalidCategoryID = errors.New("invalid category id")
)

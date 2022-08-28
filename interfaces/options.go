package interfaces

import (
	"net/http"
)

// KeyContext describe key type for ngamux context
type KeyContext int

const (
	// KeyContextParams is key context for url params
	KeyContextParams KeyContext = 1 << iota
)

var (
	DefaultPattern = `[\w]+`
	IDPattern      = `[\d]+`
	IDKey          = `id`

	METHOD = map[string]struct{}{
		http.MethodGet:    {},
		http.MethodPost:   {},
		http.MethodPut:    {},
		http.MethodDelete: {},
		http.MethodPatch:  {},
	}
)

package binding

import "net/http"

// Binder describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binder interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

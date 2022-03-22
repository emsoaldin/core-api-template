package vm

// Response is JSON view model for all controller responses
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

// ResponseError view model
type ResponseError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Stack      string            `json:"stack"`
	Cause      *ResponseError    `json:"cause"`
	Validation map[string]string `json:"validation"`
}

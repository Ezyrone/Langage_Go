package api

type CheckRequest struct {
	URLs    []string       `json:"urls"`
	Options *CheckOptions  `json:"options,omitempty"`
}

type CheckOptions struct {
	Concurrency *int `json:"concurrency,omitempty"`
	TimeoutMs   *int `json:"timeout_ms,omitempty"`
}

type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

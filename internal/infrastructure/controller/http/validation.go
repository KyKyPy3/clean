package http

type ValidationError struct {
	Field  string      `json:"field"`
	Value  interface{} `json:"value"`
	Reason string      `json:"reason"`
}

type ResponseDTO struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    interface{}        `json:"data,omitempty"`
	Errors  []*ValidationError `json:"errors,omitempty"`
	Error   string             `json:"error,omitempty"`
}

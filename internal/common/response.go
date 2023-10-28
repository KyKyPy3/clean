package common

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const (
	MessageOk    = "ok"
	MessageError = "error"
)

func OK() Response {
	return Response{
		Message: MessageOk,
	}
}

func Error() Response {
	return Response{
		Message: MessageError,
	}
}

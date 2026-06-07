package response

type WebResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

func Success(data interface{}) WebResponse {
	return WebResponse{Code: 200, Status: "OK", Data: data}
}

func Created(data interface{}) WebResponse {
	return WebResponse{Code: 201, Status: "CREATED", Data: data}
}

func SuccessMessage(message string) WebResponse {
	return WebResponse{Code: 200, Status: "OK", Data: message}
}

func Error(code int, message string) WebResponse {
	return WebResponse{Code: code, Status: "ERROR", Data: nil, Message: message}
}

func BadRequest(message string) WebResponse {
	return WebResponse{Code: 400, Status: "BAD_REQUEST", Data: nil, Message: message}
}

func Unauthorized(message string) WebResponse {
	return WebResponse{Code: 401, Status: "UNAUTHORIZED", Data: nil, Message: message}
}

func Forbidden(message string) WebResponse {
	return WebResponse{Code: 403, Status: "FORBIDDEN", Data: nil, Message: message}
}

func NotFound(message string) WebResponse {
	return WebResponse{Code: 404, Status: "NOT_FOUND", Data: nil, Message: message}
}

func InternalServerError(message string) WebResponse {
	return WebResponse{Code: 500, Status: "INTERNAL_SERVER_ERROR", Data: nil, Message: message}
}

func Conflict(message string) WebResponse {
	return WebResponse{Code: 409, Status: "CONFLICT", Data: nil, Message: message}
}

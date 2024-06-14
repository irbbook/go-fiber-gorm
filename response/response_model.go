package MainResponse

type MainResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success returns a standard success response without data
func Success() *MainResponse {
	return &MainResponse{
		Status:  200,
		Message: "Success",
		Data:    nil,
	}
}

func Custom(status int, message string, data interface{}) *MainResponse {
	return &MainResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// SuccessWithData returns a success response with data
func SuccessWithData(data interface{}) *MainResponse {
	return &MainResponse{
		Status:  200,
		Message: "Success",
		Data:    data,
	}
}

// NotFound returns a standard not found response
func NotFound(message string) *MainResponse {
	return &MainResponse{
		Status:  404,
		Message: message,
		Data:    nil,
	}
}

// BadRequest returns a standard bad request response
func BadRequest(message string) *MainResponse {
	return &MainResponse{
		Status:  400,
		Message: message,
		Data:    nil,
	}
}

// InternalServerError returns a standard internal server error response
func InternalServerError(message string) *MainResponse {
	return &MainResponse{
		Status:  500,
		Message: message,
		Data:    nil,
	}
}

// CustomError returns a response with a custom status code and message
func CustomError(status int, message string) *MainResponse {
	return &MainResponse{
		Status:  status,
		Message: message,
		Data:    nil,
	}
}

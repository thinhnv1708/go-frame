package response

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Result  any    `json:"result,omitempty"`
}

func NewSuccessResponse(result any) ApiResponse {
	return ApiResponse{
		Code:   0,
		Result: result,
	}
}

func NewErrorResponse(code int, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
	}
}

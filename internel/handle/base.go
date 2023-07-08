package handle

import (
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func success(data interface{}) Response {
	return Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    data,
	}
}

func fail(err error, data interface{}) Response {
	return Response{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Data:    data,
	}
}

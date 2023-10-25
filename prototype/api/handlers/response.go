package handlers

import (
	"net/http"
)

type Responder interface {
	InvalidSession() ErrorResponse
	Unauthorized() ErrorResponse
	BadRequest() ErrorResponse
	InternalError(error) ErrorResponse
	OK(headers ...map[string]string) Response
}

type Response struct {
	Headers map[string]string `json:"-"`
	Code    int               `json:"-"`
}

type ErrorResponse struct {
	Reason string `json:"reason"`
	Response
}

type BaseResponder struct {
	GlobalHeaders map[string]string
}

type responderOpt func(r *BaseResponder)

func WithGlobalHeader(name, value string) responderOpt {
	return func(b *BaseResponder) {
		b.GlobalHeaders[name] = value
	}
}

func NewBaseResponder(opts ...responderOpt) *BaseResponder {

	r := &BaseResponder{
		GlobalHeaders: map[string]string{},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (b *BaseResponder) InvalidSession() ErrorResponse {
	return ErrorResponse{
		Response: Response{
			Code:    http.StatusBadRequest,
			Headers: b.GlobalHeaders,
		},
		Reason: "Invalid Session. Bad request from client, please reload page",
	}
}

func (b *BaseResponder) BadRequest() ErrorResponse {
	return ErrorResponse{
		Response: Response{
			Code:    http.StatusBadRequest,
			Headers: b.GlobalHeaders,
		},
		Reason: "Invalid json request",
	}

}

func (b *BaseResponder) Unauthorized() ErrorResponse {
	return ErrorResponse{
		Response: Response{
			Code:    http.StatusUnauthorized,
			Headers: b.GlobalHeaders,
		},
		Reason: "Unauthorized. Please login again",
	}

}

func (b *BaseResponder) InternalError(_ error) ErrorResponse {
	return ErrorResponse{
		Response: Response{
			Code:    http.StatusInternalServerError,
			Headers: b.GlobalHeaders,
		},
		Reason: "Something went wrong outside of your control. Please try again",
	}

}

func (b *BaseResponder) OK(headers ...map[string]string) Response {
	hs := make([]map[string]string, 1+len(headers))

	hs[0] = b.GlobalHeaders

	for i, v := range headers {
		hs[i+1] = v
	}

	return Response{
		Code:    http.StatusOK,
		Headers: mergeHeaders(hs...),
	}

}

func mergeHeaders(headerSets ...map[string]string) map[string]string {

	headers := make(map[string]string)

	for _, h := range headerSets {
		for k, v := range h {
			headers[k] = v
		}
	}

	return headers
}

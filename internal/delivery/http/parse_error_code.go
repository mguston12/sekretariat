package http

import (
	"sekretariat/pkg/response"
	"strings"
)

// ParseErrorCode ...
func ParseErrorCode(err string) response.Response {
	errResp := response.Error{
		Status: true,
		Msg:    err,
		Code:   200,
	}

	switch {
	case strings.Contains(err, "401"):
		errResp = response.Error{
			Status: true,
			Msg:    "Unauthorized",
			Code:   401,
		}
	}

	return response.Response{
		Error: errResp,
	}
}

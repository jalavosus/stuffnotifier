package errs

import (
	"github.com/pkg/errors"
)

const (
	buildHttpRequestErrMsg      string = "error building http request"
	httpResponseErrMsg          string = "error performing http request"
	httpReadBodyErrMsg          string = "error reading response body"
	httpUnmarshalResponseErrMsg string = "error unmarshalling response body"
)

func HttpBuildRequestError(cause error) error {
	return errors.WithMessage(cause, buildHttpRequestErrMsg)
}

func HttpResponseError(cause error) error {
	return errors.WithMessage(cause, httpResponseErrMsg)
}

func HttpReadBodyError(cause error) error {
	return errors.WithMessage(cause, httpReadBodyErrMsg)
}

func HttpUnmarshalResponseBodyError(cause error) error {
	return errors.WithMessage(cause, httpUnmarshalResponseErrMsg)
}

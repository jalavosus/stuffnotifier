package errs

import (
	"github.com/pkg/errors"
)

const (
	readFileErrMsg string = "error reading data from file %[1]s"
)

func ReadFileError(cause error, filePath string) error {
	return errors.WithMessagef(cause, readFileErrMsg, filePath)
}

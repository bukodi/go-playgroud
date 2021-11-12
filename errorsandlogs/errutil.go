package errorsandlogs

import "context"

type ErrorHandler func(ctx context.Context, err error, fields map[string]interface{}) error

var errorHandler ErrorHandler = DefaultErrorHandler

func SetErrorHandler(handler ErrorHandler) {
	errorHandler = handler
}

func HandleError(err error) error {
	if errorHandler == nil {
		return err
	} else {
		return errorHandler(nil, err, nil)
	}
}

func DefaultErrorHandler(ctx context.Context, err error, fields map[string]interface{}) error {
	return err
}

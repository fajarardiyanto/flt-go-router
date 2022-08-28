package main

import (
	"github.com/fajarardiyanto/flt-go-router/interfaces"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func MiddlewareLogger() interfaces.MiddlewareFunc {
	return func(next interfaces.Handler) interfaces.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			now := time.Now()
			GetLogger().Info("ttc(μs) : %d/ms", time.Since(now).Milliseconds())
			return next(w, r)
		}
	}
}

func MiddlewareError() interfaces.MiddlewareFunc {
	return func(next interfaces.Handler) interfaces.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := next(w, r); err != nil {
				response := formatResponse(err)

				if err := interfaces.JSON(w, response.Code, response); err != nil {
					return err
				}

				GetLogger().Error(err).Quit()
			}

			return nil
		}
	}
}

// formatResponse returns response error format based on error type.
func formatResponse(err error) *interfaces.APIResponseError {
	switch et := errors.Cause(err).(type) {
	case *interfaces.RequestError:
		return interfaces.NewAPIResponseError(et.Error(), et.Status)
	case interfaces.ValidationError:
		return interfaces.NewAPIValidationError(et, http.StatusBadRequest)
	default:
		return interfaces.NewAPIResponseError(
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}

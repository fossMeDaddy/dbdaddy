package middlewares

import (
	"github.com/fossmedaddy/dbdaddy/types"
)

func Apply(cmdRunFn types.CobraCmdFn, middlewares ...types.MiddlewareFunc) types.CobraCmdFn {
	result := cmdRunFn
	for _, middlewareFn := range middlewares {
		result = middlewareFn(result)
	}

	return result
}

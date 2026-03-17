package errcodes

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// GRPCToHTTP maps gRPC status codes to HTTP status codes.
// Based on standard gRPC to HTTP status code mapping conventions:
// https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
func GRPCToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK // 200
	case codes.Canceled:
		return 499 // Client Closed Request (non-standard but widely used)
	case codes.Unknown:
		return http.StatusInternalServerError // 500
	case codes.InvalidArgument:
		return http.StatusBadRequest // 400
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout // 504
	case codes.NotFound:
		return http.StatusNotFound // 404
	case codes.AlreadyExists:
		return http.StatusConflict // 409
	case codes.PermissionDenied:
		return http.StatusForbidden // 403
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests // 429
	case codes.FailedPrecondition:
		return http.StatusBadRequest // 400
	case codes.Aborted:
		return http.StatusConflict // 409
	case codes.OutOfRange:
		return http.StatusBadRequest // 400
	case codes.Unimplemented:
		return http.StatusNotImplemented // 501
	case codes.Internal:
		return http.StatusInternalServerError // 500
	case codes.Unavailable:
		return http.StatusServiceUnavailable // 503
	case codes.DataLoss:
		return http.StatusInternalServerError // 500
	case codes.Unauthenticated:
		return http.StatusUnauthorized // 401
	default:
		return http.StatusInternalServerError // 500
	}
}

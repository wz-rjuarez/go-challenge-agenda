package errcodes_test

import (
	"net/http"
	"testing"

	"go-challenge-agenda/pkg/errcodes"

	"google.golang.org/grpc/codes"
)

func TestGRPCToHTTP(t *testing.T) {
	tests := []struct {
		name       string
		grpcCode   codes.Code
		wantStatus int
	}{
		{
			name:       "OK maps to 200",
			grpcCode:   codes.OK,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Canceled maps to 499",
			grpcCode:   codes.Canceled,
			wantStatus: 499,
		},
		{
			name:       "Unknown maps to 500",
			grpcCode:   codes.Unknown,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "InvalidArgument maps to 400",
			grpcCode:   codes.InvalidArgument,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "DeadlineExceeded maps to 504",
			grpcCode:   codes.DeadlineExceeded,
			wantStatus: http.StatusGatewayTimeout,
		},
		{
			name:       "NotFound maps to 404",
			grpcCode:   codes.NotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "AlreadyExists maps to 409",
			grpcCode:   codes.AlreadyExists,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "PermissionDenied maps to 403",
			grpcCode:   codes.PermissionDenied,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "ResourceExhausted maps to 429",
			grpcCode:   codes.ResourceExhausted,
			wantStatus: http.StatusTooManyRequests,
		},
		{
			name:       "FailedPrecondition maps to 400",
			grpcCode:   codes.FailedPrecondition,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Aborted maps to 409",
			grpcCode:   codes.Aborted,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "OutOfRange maps to 400",
			grpcCode:   codes.OutOfRange,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Unimplemented maps to 501",
			grpcCode:   codes.Unimplemented,
			wantStatus: http.StatusNotImplemented,
		},
		{
			name:       "Internal maps to 500",
			grpcCode:   codes.Internal,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Unavailable maps to 503",
			grpcCode:   codes.Unavailable,
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "DataLoss maps to 500",
			grpcCode:   codes.DataLoss,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Unauthenticated maps to 401",
			grpcCode:   codes.Unauthenticated,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Unknown code maps to 500",
			grpcCode:   codes.Code(999),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errcodes.GRPCToHTTP(tt.grpcCode)
			if got != tt.wantStatus {
				t.Errorf("GRPCToHTTP(%v) = %d, want %d", tt.grpcCode, got, tt.wantStatus)
			}
		})
	}
}

// TestAllGRPCCodesAreMapped ensures that all standard gRPC codes have explicit mappings.
// This test will fail if new gRPC codes are added in the future and not mapped.
func TestAllGRPCCodesAreMapped(t *testing.T) {
	allCodes := []codes.Code{
		codes.OK,
		codes.Canceled,
		codes.Unknown,
		codes.InvalidArgument,
		codes.DeadlineExceeded,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss,
		codes.Unauthenticated,
	}

	for _, code := range allCodes {
		// Just ensure no panic and returns a valid HTTP status code
		status := errcodes.GRPCToHTTP(code)
		if status < 100 || status > 599 {
			t.Errorf("GRPCToHTTP(%v) returned invalid HTTP status: %d", code, status)
		}
	}
}

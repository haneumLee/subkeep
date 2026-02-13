package utils

import "testing"

func TestErrNotFound(t *testing.T) {
	err := ErrNotFound("user not found")
	if err.Code != 404 {
		t.Errorf("ErrNotFound code = %d, want 404", err.Code)
	}
	if err.Detail != "user not found" {
		t.Errorf("ErrNotFound detail = %q, want %q", err.Detail, "user not found")
	}
}

func TestErrUnauthorized(t *testing.T) {
	err := ErrUnauthorized("invalid token")
	if err.Code != 401 {
		t.Errorf("ErrUnauthorized code = %d, want 401", err.Code)
	}
}

func TestErrBadRequest(t *testing.T) {
	err := ErrBadRequest("missing field")
	if err.Code != 400 {
		t.Errorf("ErrBadRequest code = %d, want 400", err.Code)
	}
}

func TestErrForbidden(t *testing.T) {
	err := ErrForbidden("access denied")
	if err.Code != 403 {
		t.Errorf("ErrForbidden code = %d, want 403", err.Code)
	}
}

func TestErrValidation(t *testing.T) {
	err := ErrValidation("invalid input")
	if err.Code != 422 {
		t.Errorf("ErrValidation code = %d, want 422", err.Code)
	}
}

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name   string
		err    *AppError
		want   string
	}{
		{
			name: "with detail",
			err:  &AppError{Code: 404, Message: "Not Found", Detail: "user 123"},
			want: "[404] Not Found: user 123",
		},
		{
			name: "without detail",
			err:  &AppError{Code: 500, Message: "Internal Error"},
			want: "[500] Internal Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewAppError(t *testing.T) {
	err := NewAppError(503, "Service Unavailable", "db down")
	if err.Code != 503 {
		t.Errorf("code = %d, want 503", err.Code)
	}
	if err.Message != "Service Unavailable" {
		t.Errorf("message = %q, want %q", err.Message, "Service Unavailable")
	}
	if err.Detail != "db down" {
		t.Errorf("detail = %q, want %q", err.Detail, "db down")
	}
}

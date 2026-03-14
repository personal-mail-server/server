package auth

import "testing"

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name      string
		req       LoginRequest
		wantError string
	}{
		{
			name: "valid request",
			req: LoginRequest{
				LoginID:  "user-01",
				Password: "abc12345",
			},
			wantError: "",
		},
		{
			name: "missing login id",
			req: LoginRequest{
				Password: "abc12345",
			},
			wantError: CodeMissingRequired,
		},
		{
			name: "invalid login id uppercase",
			req: LoginRequest{
				LoginID:  "User-01",
				Password: "abc12345",
			},
			wantError: CodeInvalidLoginID,
		},
		{
			name: "invalid login id underscore",
			req: LoginRequest{
				LoginID:  "user_01",
				Password: "abc12345",
			},
			wantError: CodeInvalidLoginID,
		},
		{
			name: "password has no digit",
			req: LoginRequest{
				LoginID:  "user-01",
				Password: "abcdefgh",
			},
			wantError: CodeInvalidPassword,
		},
		{
			name: "password has no letter",
			req: LoginRequest{
				LoginID:  "user-01",
				Password: "12345678",
			},
			wantError: CodeInvalidPassword,
		},
		{
			name: "password has whitespace",
			req: LoginRequest{
				LoginID:  "user-01",
				Password: "abc1 2345",
			},
			wantError: CodeInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLoginRequest(tt.req)
			if tt.wantError == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.wantError != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantError)
				}
				if err.Code != tt.wantError {
					t.Fatalf("expected error code %q, got %q", tt.wantError, err.Code)
				}
			}
		})
	}
}

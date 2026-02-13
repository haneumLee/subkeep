package utils

import "testing"

func TestMonthlyAmount(t *testing.T) {
	tests := []struct {
		name         string
		amount       int64
		billingCycle string
		want         int64
	}{
		{
			name:         "monthly returns amount as-is",
			amount:       10000,
			billingCycle: "monthly",
			want:         10000,
		},
		{
			name:         "yearly divides by 12",
			amount:       120000,
			billingCycle: "yearly",
			want:         10000,
		},
		{
			name:         "weekly multiplies by 52 then divides by 12",
			amount:       5000,
			billingCycle: "weekly",
			want:         21667,
		},
		{
			name:         "invalid cycle defaults to monthly",
			amount:       10000,
			billingCycle: "daily",
			want:         10000,
		},
		{
			name:         "zero amount monthly",
			amount:       0,
			billingCycle: "monthly",
			want:         0,
		},
		{
			name:         "zero amount yearly",
			amount:       0,
			billingCycle: "yearly",
			want:         0,
		},
		{
			name:         "empty string defaults to monthly",
			amount:       7777,
			billingCycle: "",
			want:         7777,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MonthlyAmount(tt.amount, tt.billingCycle)
			if got != tt.want {
				t.Errorf("MonthlyAmount(%d, %q) = %d, want %d", tt.amount, tt.billingCycle, got, tt.want)
			}
		})
	}
}

type testValidateInput struct {
	Name  string `validate:"required,min=1,max=100"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0"`
}

func TestValidateStruct(t *testing.T) {
	t.Run("valid struct returns nil", func(t *testing.T) {
		input := testValidateInput{
			Name:  "홍길동",
			Email: "test@example.com",
			Age:   25,
		}
		err := ValidateStruct(input)
		if err != nil {
			t.Errorf("ValidateStruct() = %v, want nil", err)
		}
	})

	t.Run("missing required field returns AppError with 422", func(t *testing.T) {
		input := testValidateInput{
			Name:  "",
			Email: "",
			Age:   0,
		}
		err := ValidateStruct(input)
		if err == nil {
			t.Fatal("ValidateStruct() = nil, want error")
		}
		if err.Code != 422 {
			t.Errorf("error code = %d, want 422", err.Code)
		}
	})

	t.Run("invalid email returns AppError with 422", func(t *testing.T) {
		input := testValidateInput{
			Name:  "Test",
			Email: "not-an-email",
			Age:   10,
		}
		err := ValidateStruct(input)
		if err == nil {
			t.Fatal("ValidateStruct() = nil, want error")
		}
		if err.Code != 422 {
			t.Errorf("error code = %d, want 422", err.Code)
		}
	})
}

package models

import "testing"

func TestSubscription_MonthlyAmount(t *testing.T) {
	tests := []struct {
		name         string
		amount       int
		billingCycle BillingCycle
		want         int
	}{
		{
			name:         "monthly billing cycle returns amount as-is",
			amount:       10000,
			billingCycle: BillingCycleMonthly,
			want:         10000,
		},
		{
			name:         "yearly billing cycle divides by 12 evenly",
			amount:       120000,
			billingCycle: BillingCycleYearly,
			want:         10000,
		},
		{
			name:         "yearly billing cycle divides by 12 with rounding",
			amount:       10000,
			billingCycle: BillingCycleYearly,
			want:         833,
		},
		{
			name:         "weekly billing cycle multiplies by 52/12",
			amount:       5000,
			billingCycle: BillingCycleWeekly,
			want:         21667,
		},
		{
			name:         "zero amount monthly",
			amount:       0,
			billingCycle: BillingCycleMonthly,
			want:         0,
		},
		{
			name:         "zero amount yearly",
			amount:       0,
			billingCycle: BillingCycleYearly,
			want:         0,
		},
		{
			name:         "zero amount weekly",
			amount:       0,
			billingCycle: BillingCycleWeekly,
			want:         0,
		},
		{
			name:         "amount 1 yearly rounds to 0",
			amount:       1,
			billingCycle: BillingCycleYearly,
			want:         0,
		},
		{
			name:         "large amount yearly",
			amount:       1200000,
			billingCycle: BillingCycleYearly,
			want:         100000,
		},
		{
			name:         "unknown billing cycle defaults to amount as-is",
			amount:       9999,
			billingCycle: BillingCycle("daily"),
			want:         9999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				Amount:       tt.amount,
				BillingCycle: tt.billingCycle,
			}
			got := s.MonthlyAmount()
			if got != tt.want {
				t.Errorf("MonthlyAmount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSubscription_AnnualAmount(t *testing.T) {
	tests := []struct {
		name         string
		amount       int
		billingCycle BillingCycle
		want         int
	}{
		{
			name:         "monthly 10000 annual is 120000",
			amount:       10000,
			billingCycle: BillingCycleMonthly,
			want:         120000,
		},
		{
			name:         "yearly 120000 annual is 120000",
			amount:       120000,
			billingCycle: BillingCycleYearly,
			want:         120000,
		},
		{
			name:         "weekly 5000 annual is 260004",
			amount:       5000,
			billingCycle: BillingCycleWeekly,
			want:         260004,
		},
		{
			name:         "zero amount annual is 0",
			amount:       0,
			billingCycle: BillingCycleMonthly,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				Amount:       tt.amount,
				BillingCycle: tt.billingCycle,
			}
			got := s.AnnualAmount()
			if got != tt.want {
				t.Errorf("AnnualAmount() = %d, want %d", got, tt.want)
			}
		})
	}
}

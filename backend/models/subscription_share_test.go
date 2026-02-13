package models

import "testing"

func intPtr(v int) *int         { return &v }
func float64Ptr(v float64) *float64 { return &v }

func TestSubscriptionShare_PersonalAmount(t *testing.T) {
	tests := []struct {
		name          string
		share         SubscriptionShare
		monthlyAmount int
		want          int
	}{
		{
			name: "equal split 4 members amount 20000",
			share: SubscriptionShare{
				SplitType:            SplitTypeEqual,
				TotalMembersSnapshot: 4,
			},
			monthlyAmount: 20000,
			want:          5000,
		},
		{
			name: "equal split 3 members amount 10000 rounds",
			share: SubscriptionShare{
				SplitType:            SplitTypeEqual,
				TotalMembersSnapshot: 3,
			},
			monthlyAmount: 10000,
			want:          3333,
		},
		{
			name: "equal split 1 member returns full amount",
			share: SubscriptionShare{
				SplitType:            SplitTypeEqual,
				TotalMembersSnapshot: 1,
			},
			monthlyAmount: 10000,
			want:          10000,
		},
		{
			name: "equal split 0 members returns full amount (avoid division by zero)",
			share: SubscriptionShare{
				SplitType:            SplitTypeEqual,
				TotalMembersSnapshot: 0,
			},
			monthlyAmount: 10000,
			want:          10000,
		},
		{
			name: "custom amount with myShareAmount 5000",
			share: SubscriptionShare{
				SplitType:     SplitTypeCustomAmount,
				MyShareAmount: intPtr(5000),
			},
			monthlyAmount: 20000,
			want:          5000,
		},
		{
			name: "custom amount with nil myShareAmount returns 0",
			share: SubscriptionShare{
				SplitType:     SplitTypeCustomAmount,
				MyShareAmount: nil,
			},
			monthlyAmount: 20000,
			want:          0,
		},
		{
			name: "custom ratio 0.3 with amount 10000",
			share: SubscriptionShare{
				SplitType:    SplitTypeCustomRatio,
				MyShareRatio: float64Ptr(0.3),
			},
			monthlyAmount: 10000,
			want:          3000,
		},
		{
			name: "custom ratio nil returns 0",
			share: SubscriptionShare{
				SplitType:    SplitTypeCustomRatio,
				MyShareRatio: nil,
			},
			monthlyAmount: 10000,
			want:          0,
		},
		{
			name: "custom ratio 1.0 with amount 15000",
			share: SubscriptionShare{
				SplitType:    SplitTypeCustomRatio,
				MyShareRatio: float64Ptr(1.0),
			},
			monthlyAmount: 15000,
			want:          15000,
		},
		{
			name: "custom ratio 0.0 returns 0",
			share: SubscriptionShare{
				SplitType:    SplitTypeCustomRatio,
				MyShareRatio: float64Ptr(0.0),
			},
			monthlyAmount: 10000,
			want:          0,
		},
		{
			name: "unknown split type returns full amount",
			share: SubscriptionShare{
				SplitType: SplitType("unknown"),
			},
			monthlyAmount: 7000,
			want:          7000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.share.PersonalAmount(tt.monthlyAmount)
			if got != tt.want {
				t.Errorf("PersonalAmount(%d) = %d, want %d", tt.monthlyAmount, got, tt.want)
			}
		})
	}
}

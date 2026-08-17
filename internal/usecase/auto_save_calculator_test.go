package usecase

import (
	"testing"

	"backend_keping/internal/domain"
)

func TestAutoSaveCalculator_CalculateAllocation(t *testing.T) {
	calc := NewAutoSaveCalculator()

	tests := []struct {
		name        string
		income      float64
		wishlist    domain.Wishlist
		expectedAmt float64
	}{
		{
			name:   "Standard 10% allocation",
			income: 500000,
			wishlist: domain.Wishlist{
				TargetAmount:       1500000,
				SavedAmount:        0,
				AutoSavePercentage: 10,
				IsCompleted:        false,
			},
			expectedAmt: 50000,
		},
		{
			name:   "Cap at remaining target",
			income: 10000000,
			wishlist: domain.Wishlist{
				TargetAmount:       1000000,
				SavedAmount:        950000,
				AutoSavePercentage: 10, // 10% of 10M is 1M, but only 50k is needed
				IsCompleted:        false,
			},
			expectedAmt: 50000,
		},
		{
			name:   "Already completed wishlist",
			income: 1000000,
			wishlist: domain.Wishlist{
				TargetAmount:       1000000,
				SavedAmount:        1000000,
				AutoSavePercentage: 10,
				IsCompleted:        true,
			},
			expectedAmt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := calc.CalculateAllocation(tt.income, tt.wishlist)
			if res != tt.expectedAmt {
				t.Errorf("expected %v, got %v", tt.expectedAmt, res)
			}
		})
	}
}

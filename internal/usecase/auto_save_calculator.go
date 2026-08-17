package usecase

import (
	"math"

	"backend_keping/internal/domain"
)

type AutoSaveCalculator interface {
	CalculateAllocation(incomeAmount float64, wishlist domain.Wishlist) float64
}

type autoSaveCalculator struct{}

func NewAutoSaveCalculator() AutoSaveCalculator {
	return &autoSaveCalculator{}
}

func (c *autoSaveCalculator) CalculateAllocation(incomeAmount float64, wishlist domain.Wishlist) float64 {
	if incomeAmount <= 0 || wishlist.AutoSavePercentage <= 0 || wishlist.IsCompleted {
		return 0
	}

	// Formula: Allocation = Income * (Percentage / 100)
	rawAllocation := incomeAmount * (wishlist.AutoSavePercentage / 100.0)

	// Round to 2 decimal places to avoid floating point precision issues
	allocation := math.Round(rawAllocation*100) / 100

	// Cap allocation if it exceeds remaining target amount
	remainingNeeded := wishlist.TargetAmount - wishlist.SavedAmount
	if remainingNeeded <= 0 {
		return 0
	}

	if allocation > remainingNeeded {
		return remainingNeeded
	}

	return allocation
}

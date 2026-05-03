package service

import (
	"fmt"
	"math"
	"time"
)

// TradeOperation handles trading operations between entities.
type TradeOperation struct {
	cooldownMap    map[string]time.Time
	cooldownPeriod time.Duration
}

// NewTradeOperation creates a new TradeOperation.
func NewTradeOperation(cooldownPeriod time.Duration) *TradeOperation {
	return &TradeOperation{
		cooldownMap:    make(map[string]time.Time),
		cooldownPeriod: cooldownPeriod,
	}
}

// TradeInput holds inputs for a trade operation.
type TradeInput struct {
	BuyerID     string
	SellerID    string
	ItemName    string
	ItemType    string
	ItemQuality string
	Price       int64 // in low-grade spirit stones
}

// TradeResult holds the outcome of a trade operation.
type TradeResult struct {
	Success      bool
	HeavenlyTax  int64
	PlatformFee  int64
	SellerGets   int64
	BuyerPays    int64
	Message      string
}

// ExecuteTrade executes a trade between buyer and seller.
func (op *TradeOperation) ExecuteTrade(input TradeInput, buyerStones, sellerStones int64, now time.Time) (*TradeResult, error) {
	// Check buyer cooldown
	if lastTime, ok := op.cooldownMap[input.BuyerID]; ok {
		if now.Sub(lastTime) < op.cooldownPeriod {
			return nil, fmt.Errorf("trade cooldown: %v remaining", op.cooldownPeriod-now.Sub(lastTime))
		}
	}

	// Validate price
	if input.Price <= 0 {
		return nil, fmt.Errorf("价格必须大于0")
	}

	// Calculate heavenly tax (progressive rate)
	heavenlyTax := CalculateHeavenlyTax(input.Price)

	// Calculate platform fee (2% base)
	platformFee := int64(math.Ceil(float64(input.Price) * 0.02))

	// Total buyer pays
	buyerPays := input.Price + heavenlyTax

	// Seller receives (price minus platform fee)
	sellerGets := input.Price - platformFee

	// Check buyer can afford
	if buyerStones < buyerPays {
		return nil, fmt.Errorf("灵石不足：需要 %d，当前有 %d", buyerPays, buyerStones)
	}

	// Set cooldown
	op.cooldownMap[input.BuyerID] = now

	return &TradeResult{
		Success:     true,
		HeavenlyTax: heavenlyTax,
		PlatformFee: platformFee,
		SellerGets:  sellerGets,
		BuyerPays:   buyerPays,
		Message:     fmt.Sprintf("交易成功：%s 以 %d 灵石购买了 %s", input.BuyerID, buyerPays, input.ItemName),
	}, nil
}

// CalculateHeavenlyTax calculates the progressive heavenly tax rate.
//
// Progressive tax brackets:
//   - 0-100: 0%
//   - 101-1000: 1%
//   - 1001-10000: 3%
//   - 10001-50000: 5%
//   - 50001-100000: 8%
//   - 100001+: 12%
func CalculateHeavenlyTax(amount int64) int64 {
	var tax int64

	brackets := []struct {
		Limit int64
		Rate  float64
	}{
		{100, 0.00},
		{1000, 0.01},
		{10000, 0.03},
		{50000, 0.05},
		{100000, 0.08},
		{math.MaxInt64, 0.12},
	}

	remaining := amount
	prevLimit := int64(0)
	for _, bracket := range brackets {
		if remaining <= 0 {
			break
		}

		bracketSize := bracket.Limit - prevLimit
		taxableInBracket := minInt64(remaining, bracketSize)
		tax += int64(math.Ceil(float64(taxableInBracket) * bracket.Rate))
		remaining -= taxableInBracket
		prevLimit = bracket.Limit
	}

	return tax
}

// CalculateSpiritStoneExchange calculates the exchange rate and fees for
// converting between spirit stone grades.
//
// Exchange rates:
//   - 1 medium = 100 low
//   - 1 high = 100 medium = 10000 low
//   - 1 premium = 100 high = 10000 medium = 1000000 low
//
// Exchange fee: 1% of the converted amount
func CalculateSpiritStoneExchange(fromGrade, toGrade string, amount int64) (exchanged int64, fee int64, err error) {
	gradeValues := map[string]int64{
		"low":     1,
		"medium":  100,
		"high":    10000,
		"premium": 1000000,
	}

	fromValue, ok := gradeValues[fromGrade]
	if !ok {
		return 0, 0, fmt.Errorf("unknown grade: %s", fromGrade)
	}

	toValue, ok := gradeValues[toGrade]
	if !ok {
		return 0, 0, fmt.Errorf("unknown grade: %s", toGrade)
	}

	// Convert to low-grade equivalent
	lowEquivalent := amount * fromValue

	// Convert to target grade
	exchanged = lowEquivalent / toValue

	// Fee: 1%
	fee = int64(math.Ceil(float64(exchanged) * 0.01))
	exchanged -= fee

	if exchanged < 0 {
		exchanged = 0
	}

	return exchanged, fee, nil
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

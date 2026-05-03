package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateHeavenlyTax_Zero(t *testing.T) {
	assert.Equal(t, int64(0), CalculateHeavenlyTax(0))
}

func TestCalculateHeavenlyTax_FreeTier(t *testing.T) {
	// 0-100: 0%
	assert.Equal(t, int64(0), CalculateHeavenlyTax(50))
	assert.Equal(t, int64(0), CalculateHeavenlyTax(100))
}

func TestCalculateHeavenlyTax_LowTier(t *testing.T) {
	// 101-1000: 1%
	assert.Equal(t, int64(1), CalculateHeavenlyTax(101))   // 1 on the amount above 100
	assert.Equal(t, int64(9), CalculateHeavenlyTax(1000)) // 900 * 1% = 9
}

func TestCalculateHeavenlyTax_MidTier(t *testing.T) {
	// 1001-10000: 3%
	// First 1000 = 10 tax, next 1 = 3% = 1
	assert.True(t, CalculateHeavenlyTax(1001) >= 10)
	assert.True(t, CalculateHeavenlyTax(10000) > 10)
}

func TestCalculateHeavenlyTax_HighTier(t *testing.T) {
	// Progressive, should always be reasonable
	tax := CalculateHeavenlyTax(100000)
	assert.True(t, tax > 0)
	assert.True(t, tax < 50000) // less than 50%
}

func TestCalculateHeavenlyTax_Progressive(t *testing.T) {
	// Higher amounts should have higher effective rates
	tax1 := CalculateHeavenlyTax(1000)
	tax2 := CalculateHeavenlyTax(10000)
	tax3 := CalculateHeavenlyTax(100000)

	// Effective rate should increase
	rate1 := float64(tax1) / 1000.0
	rate2 := float64(tax2) / 10000.0
	rate3 := float64(tax3) / 100000.0

	assert.True(t, rate2 > rate1)
	assert.True(t, rate3 > rate2)
}

func TestCalculateSpiritStoneExchange_Upgrade(t *testing.T) {
	// 1000 low -> medium (minus 1% fee)
	exchanged, fee, err := CalculateSpiritStoneExchange("low", "medium", 1000)
	assert.NoError(t, err)
	assert.Equal(t, int64(9), exchanged)  // 10 - 1 fee
	assert.Equal(t, int64(1), fee)
}

func TestCalculateSpiritStoneExchange_Downgrade(t *testing.T) {
	// 1 high -> low (minus 1% fee)
	exchanged, fee, err := CalculateSpiritStoneExchange("high", "low", 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(9900), exchanged) // 10000 - 100 fee
	assert.Equal(t, int64(100), fee)
}

func TestCalculateSpiritStoneExchange_UnknownGrade(t *testing.T) {
	_, _, err := CalculateSpiritStoneExchange("unknown", "low", 100)
	assert.Error(t, err)

	_, _, err = CalculateSpiritStoneExchange("high", "unknown", 100)
	assert.Error(t, err)
}

func TestExecuteTrade_Basic(t *testing.T) {
	op := NewTradeOperation(time.Minute)

	input := TradeInput{
		BuyerID:  "buyer_1",
		SellerID: "seller_1",
		ItemName: "healing_pill",
		ItemType: "pill",
		Price:    500,
	}

	now := time.Now()
	result, err := op.ExecuteTrade(input, 1000, 0, now)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.True(t, result.BuyerPays >= 500)
	assert.True(t, result.SellerGets > 0)
	assert.True(t, result.SellerGets < 500) // minus platform fee
	assert.Contains(t, result.Message, "交易成功")
}

func TestExecuteTrade_InsufficientFunds(t *testing.T) {
	op := NewTradeOperation(time.Minute)

	input := TradeInput{
		BuyerID:  "buyer_1",
		SellerID: "seller_1",
		Price:    500,
	}

	now := time.Now()
	_, err := op.ExecuteTrade(input, 100, 0, now)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "灵石不足")
}

func TestExecuteTrade_InvalidPrice(t *testing.T) {
	op := NewTradeOperation(time.Minute)

	input := TradeInput{
		Price: -100,
	}

	now := time.Now()
	_, err := op.ExecuteTrade(input, 1000, 0, now)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "大于0")
}

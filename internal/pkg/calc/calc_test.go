package calc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-demobot/internal/pkg/calc"
)

func TestClamp(t *testing.T) {
	t.Run("below min", func(t *testing.T) {
		assert.Equal(t, 5.0, calc.Clamp(2.5, 5.0, 10.0))
		assert.Equal(t, 5, calc.Clamp(2, 5, 10))
	})

	t.Run("between min and max", func(t *testing.T) {
		assert.Equal(t, 7.5, calc.Clamp(7.5, 5.0, 10.0))
		assert.Equal(t, 7, calc.Clamp(7, 5, 10))
	})

	t.Run("above max", func(t *testing.T) {
		assert.Equal(t, 10.0, calc.Clamp(15.0, 5.0, 10.0))
		assert.Equal(t, 10, calc.Clamp(15, 5, 10))
	})
}

func TestRoundFloat64(t *testing.T) {
	assert.Equal(t, 3.33, calc.RoundFloat64(10.0/3, 2))
	assert.Equal(t, 3.0, calc.RoundFloat64(10.0/3, 0))
}

func TestStaticTextScale(t *testing.T) {
	assert.Equal(t, 1.2, calc.StaticTextScale(""))
	assert.Equal(t, 1.2, calc.StaticTextScale("getquad semi"))
	assert.Equal(t, 1.0, calc.StaticTextScale("getquad semigetquad semi"))
	assert.Equal(t, 0.8, calc.StaticTextScale("getquad semigetquad semigetquad semi"))
}

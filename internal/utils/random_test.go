package utils_test

import (
	"iron-stream/internal/utils"
	"testing"
)

func TestRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		code := utils.GenerateCode()
		if code < 100000 || code > 999999 {
			t.Errorf("Generated code %d is out of range", code)
		}
	}
}

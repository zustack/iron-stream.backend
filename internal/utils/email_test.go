package utils_test

import (
	"fmt"
	"iron-stream/internal/utils"
	"testing"
)

func TestEmail(t *testing.T) {
	subjet := "Verify your email on Iron Stream"
  code := utils.GenerateCode()
  email := "agustfricke@gmail.com"
  err := utils.SendEmail(code, email, subjet)
	if err != nil {
    t.Errorf("test failed because: %v", err)
	}
  fmt.Printf("=== Email to %s was sent. Check your inbox\n", email)
}

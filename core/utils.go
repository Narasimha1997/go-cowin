package core

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// SignOTP signs OTP by hasing with SHA-256 algorithm
func SignOTP(otp string) string {
	hashData := sha256.Sum256([]byte(otp))
	return fmt.Sprintf("%x", hashData)
}

// ToJSONString returns JSON string given the type
func ToJSONString(typedData *interface{}) (string, error) {
	bytes, err := json.Marshal(typedData)
	return string(bytes), err
}

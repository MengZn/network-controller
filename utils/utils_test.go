package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256String(t *testing.T) {
	o := sha256String("12345678")
	assert.Equal(t, "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f", o)
}

func TestGenerateVethName(t *testing.T) {
	o := GenerateVethName("12345678", "12345678")
	assert.Equal(t, "vtx33cdbc387", o)
}

func TestIsValidIP(t *testing.T) {
	o := IsValidIP("8.8.8.8")
	assert.Equal(t, true, o)
}

func TestIsInvalidIP(t *testing.T) {
	o := IsValidIP("abcd.efgh.ijkl.mnop")
	assert.Equal(t, false, o)
}

func TestIsValidCIDR(t *testing.T) {
	o := IsValidCIDR("10.0.0.1/24")
	assert.Equal(t, true, o)
}

func TestIsInvalidCIDR(t *testing.T) {
	o := IsValidCIDR("0.0.0.0")
	assert.Equal(t, false, o)
}

func TestIsValidVLAN(t *testing.T) {
	o := IsValidVLANTag(42)
	assert.Equal(t, true, o)
}

func TestIsInValidVLAN(t *testing.T) {
	o := IsValidVLANTag(4096)
	assert.Equal(t, false, o)
}

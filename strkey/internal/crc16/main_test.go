package crc16

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	result := Checksum([]byte{0x12, 0x34, 0x56, 0x78, 0x90})
	assert.Equal(t, uint16(0x48e6), result)
}

func TestValidate(t *testing.T) {
	err := Validate([]byte{0x12, 0x34, 0x56, 0x78, 0x90}, 0x48e6)
	assert.NoError(t, err)

	err = Validate([]byte{0x12, 0x34, 0x56, 0x78, 0x90}, 0x48e7)
	assert.ErrorIs(t, err, ErrInvalidChecksum)
}

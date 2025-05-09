package pkg

import (
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateUUID() (string, error) {
	u := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, u)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant (RFC 4122)
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10xxxxxx

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}

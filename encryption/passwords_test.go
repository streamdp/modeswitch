package encryption

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestBuildSecret(t *testing.T) {
	testCases := []struct {
		username string
		secret   []byte
	}{
		{
			username: "Alice",
			secret: []byte{0x41, 0x6c, 0x69, 0x63, 0x65, 0x35, 0x42, 0x55, 0x70, 0x74, 0x61, 0x4d, 0x70, 0x6b, 0x6f,
				0x70, 0x61, 0x6d, 0x66, 0x6a, 0x52, 0x65, 0x35, 0x6d, 0x62, 0x53, 0x53, 0x4e, 0x64, 0x73, 0x2b, 0x55},
		},
		{
			username: "Bob",
			secret: []byte{0x42, 0x6f, 0x62, 0x35, 0x42, 0x55, 0x70, 0x74, 0x61, 0x4d, 0x70, 0x6b, 0x6f, 0x70, 0x61,
				0x6d, 0x66, 0x6a, 0x52, 0x65, 0x35, 0x6d, 0x62, 0x53, 0x53, 0x4e, 0x64, 0x73, 0x2b, 0x55, 0x30, 0x57},
		},
	}

	for n, tt := range testCases {
		tt := tt
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			assert.Equal(t, tt.secret, buildSecret(tt.username))
		})
	}
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		text     string
		expected string
	}{
		{
			text:     "c29tZVRleHQK",
			expected: "someText\n",
		},
		{
			text:     "YW5vdGhlck9uZQo=",
			expected: "anotherOne\n",
		},
		{
			text:     "",
			expected: "",
		},
	}

	for n, tt := range testCases {
		tt := tt
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			decoded, err := decode(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(decoded))
		})
	}
}

func TestEncrypt(t *testing.T) {
	testCases := []struct {
		username string
		password string
	}{
		{
			username: "Alice",
			password: "dwuMSqRwXG+9jbLy",
		},
		{
			username: "Bob",
			password: "GjIQzRFNGmFrHecT",
		},
	}

	for n, tt := range testCases {
		tt := tt
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			encryptedPassword, err := Encrypt(tt.password, tt.username)
			assert.NoError(t, err)
			assert.NotEqual(t, tt.password, encryptedPassword)

			decryptedPassword, err := Decrypt(encryptedPassword, tt.username)
			assert.NoError(t, err)
			assert.Equal(t, tt.password, decryptedPassword)
		})
	}
}

func TestDecrypt(t *testing.T) {
	testCases := []struct {
		username          string
		decryptedPassword string
		encryptedPassword string
	}{
		{
			username:          "Alice",
			decryptedPassword: "E9GTbejOsm2Egud5",
			encryptedPassword: "tZkwWudzv5xm0DvOK3mfbxyMYQZaZdOIu+LCQWcvF1Lp1OvzJsNHaA==",
		},
		{
			username:          "Bob",
			decryptedPassword: "pA+pi9dO0dzr0+OJ",
			encryptedPassword: "0KBG2X10HHNxDzb9FBhs7LRIoQxbrdR/ln6RP+Afocd34ExX",
		},
		{
			username:          "Empty",
			decryptedPassword: "",
			encryptedPassword: "",
		},
	}

	for n, tt := range testCases {
		tt := tt
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			decryptedPassword, err := Decrypt(tt.encryptedPassword, tt.username)
			assert.NoError(t, err)
			assert.Equal(t, tt.decryptedPassword, decryptedPassword)
		})
	}
}

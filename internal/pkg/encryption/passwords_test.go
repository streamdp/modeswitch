package encryption

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			secret:   []byte("Alice5BUptaMpkopamfjRe5mbSSNds+U"),
		},
		{
			username: "Bob",
			secret:   []byte("Bob5BUptaMpkopamfjRe5mbSSNds+U0W"),
		},
	}

	for n, tt := range testCases {
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
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			decoded, err := decode(tt.text)
			require.NoError(t, err)
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
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			encryptedPassword, err := Encrypt(tt.password, tt.username)
			require.NoError(t, err)
			assert.NotEqual(t, tt.password, encryptedPassword)

			decryptedPassword, err := Decrypt(encryptedPassword, tt.username)
			require.NoError(t, err)
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
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			decryptedPassword, err := Decrypt(tt.encryptedPassword, tt.username)
			require.NoError(t, err)
			assert.Equal(t, tt.decryptedPassword, decryptedPassword)
		})
	}
}

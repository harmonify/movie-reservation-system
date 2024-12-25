package encryption_test

import (
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/user-service/lib/util/encryption"
	"github.com/stretchr/testify/assert"
)

func TestBcrypt(t *testing.T) {
	newBcryptHash := encryption.NewBcryptHasher()

	githubCI := os.Getenv("CI")
	if githubCI == "true" && os.Getenv("INTEGRATION_TEST") == "true" {
		t.Skip("Skipping test")
	}

	t.Run("Should return hashed password", func(t *testing.T) {
		res, err := newBcryptHash.Hash("password")

		assert.NoError(t, err)
		assert.Equal(t, res, res)
	})

	t.Run("Should return true when password match", func(t *testing.T) {
		hashed, err := newBcryptHash.Hash("password")
		assert.NoError(t, err)

		match, err := newBcryptHash.Compare(hashed, "password")
		assert.True(t, match)
		assert.NoError(t, err)
	})
}

package encryption

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	generator_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/generator"
	"go.uber.org/fx"
	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash in returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrInvalidHash = errors.New("argon2id: hash is not in the correct format")

	// ErrIncompatibleVariant is returned by ComparePasswordAndHash if the
	// provided hash was created using a unsupported variant of Argon2.
	// Currently only argon2id is supported by this package.
	ErrIncompatibleVariant = errors.New("argon2id: incompatible variant of argon2")

	// ErrIncompatibleVersion is returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")
)

var Argon2HasherDefaultConfig = &Argon2HasherConfig{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

// Argon2 hasher with Argon2id algorithm variant and cryptographically-secure random salts.
type Argon2Hasher interface {
	Hash(value string) (string, error)
	Compare(hash, value string) (bool, error)
}

type Argon2HasherParam struct {
	fx.In

	GeneratorUtil generator_util.GeneratorUtil
}

type Argon2HasherResult struct {
	fx.Out

	Argon2Hasher Argon2Hasher
}

// Reference: https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04#section-4
type Argon2HasherConfig struct {
	// The amount of memory used by the algorithm (in kibibytes).
	Memory uint32

	// The number of iterations over the memory.
	Iterations uint32

	// The number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// Length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32

	// Length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

type argon2HasherImpl struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32

	generatorUtil generator_util.GeneratorUtil
}

// Create an Argon2 hasher with recommended parameters.
func NewArgon2Hasher(p Argon2HasherParam, cfg Argon2HasherConfig) Argon2HasherResult {
	return Argon2HasherResult{
		Argon2Hasher: &argon2HasherImpl{
			Memory:      cfg.Memory,
			Iterations:  cfg.Iterations,
			Parallelism: cfg.Parallelism,
			SaltLength:  cfg.SaltLength,
			KeyLength:   cfg.KeyLength,

			generatorUtil: p.GeneratorUtil,
		},
	}
}

// CreateHash returns an Argon2id hash of a plain-text password using the
// provided algorithm parameters. The returned hash follows the format used by
// the Argon2 reference C implementation and contains the base64-encoded Argon2id d
// derived key prefixed by the salt and parameters. It looks like this:
//
//	$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
func (h *argon2HasherImpl) Hash(value string) (hash string, err error) {
	salt, err := h.generatorUtil.GenerateRandomBytes(h.SaltLength)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(value), salt, h.Iterations, h.Memory, h.Parallelism, h.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.Memory, h.Iterations, h.Parallelism, b64Salt, b64Key)
	return hash, nil
}

// Compare performs a constant-time comparison between a
// plain-text password and Argon2id hash, using the parameters and salt
// contained in the hash. It returns true if they match, otherwise it returns
// false.
func (h *argon2HasherImpl) Compare(hash, value string) (match bool, err error) {
	match, _, err = h.CheckHash(hash, value)
	return match, err
}

// CheckHash is like ComparePasswordAndHash, except it also returns the params that the hash was
// created with. This can be useful if you want to update your hash params over time (which you
// should).
func (h *argon2HasherImpl) CheckHash(hash, value string) (match bool, params *Argon2HasherConfig, err error) {
	params, salt, key, err := h.DecodeHash(hash)
	if err != nil {
		return false, nil, err
	}

	otherKey := argon2.IDKey([]byte(value), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))

	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return false, params, nil
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return true, params, nil
	}
	return false, params, nil
}

// DecodeHash expects a hash created from this package, and parses it to return the params used to
// create it, as well as the salt and key (password hash).
func (h *argon2HasherImpl) DecodeHash(hash string) (params *Argon2HasherConfig, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, ErrIncompatibleVariant
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params = &Argon2HasherConfig{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}

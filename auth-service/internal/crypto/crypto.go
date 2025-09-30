package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const argon2NumSections = 6

var (
	ErrInvalidHash         = errors.New("the encoded hash is in the wrong format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

const (
	defaultMemory      = 64 * 1024
	defaultIterations  = 3
	defaultParallelism = 2
	defaultSaltLength  = 16
	defaultKeyLength   = 32
)

type Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
decodeHash takes an encoded argon2id hash,
extracts the params, hash and salt from it
and then returns them.
@params
encodedHash - encoded argon2id hash.
@returns
ps - the pointer to argon2id params.
salt - the salt used to hash the password.
hash - the hashed version of the password.
err - for checking for the successful execution of the function.
*/
func decodeHash(encodedHash string) (ps *Params, salt, hash []byte, err error) {
	p := &Params{}

	sections := strings.Split(encodedHash, "$")
	if len(sections) != argon2NumSections {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(sections[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	_, err = fmt.Sscanf(sections[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}
	ps = p

	salt, err = base64.RawStdEncoding.Strict().DecodeString(sections[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(sections[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash)) //nolint:gosec // Conversion required by the IDKey function

	return ps, salt, hash, nil
}

func GetDefaultParams() Params {
	defaultParams := Params{
		memory:      defaultMemory,
		iterations:  defaultIterations,
		parallelism: defaultParallelism,
		saltLength:  defaultSaltLength,
		keyLength:   defaultKeyLength,
	}

	return defaultParams
}

/*
VerifyPassword takes the password to be verified against a hash
and the hash, and returns a bool that represents if the passwords
match.
@params
pass - the password to be verified.
encodedHash - the hashed password.
@returns
bool - whether or not the passwords match.
error - for checking for the successful execution of the function.
*/
func VerifyPassword(pass, encodedHash string) (bool, error) {
	// Extract the params, salt and hash from the encoded hash string.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	verifiedHash := argon2.IDKey([]byte(pass), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Using ConstantTimeCompare to prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, verifiedHash) == 1 {
		return true, nil
	}

	return false, nil
}

/*
HashPasswords takes the password to be hashed and a params struct,
hashes the password using the settings specified in the params,
and then returns the encoded version of the hash.
@params
pass - the password to be hashed.
params - parameters for argon2id.
@returns
string - the encoded password hash.
error - for checking for the successful execution of the function.
*/
func HashPassword(pass string, p *Params) (string, error) {
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(pass), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

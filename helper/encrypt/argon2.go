package encrypt

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"

	"github.com/magical/argon2" // https://godoc.org/github.com/magical/argon2
	"github.com/zhengxiaoyao0716/zmodule"
	"github.com/zhengxiaoyao0716/zmodule/config"
)

var (
	aSaltLen       = 16
	aN             = 3
	aPar           = 1
	aMem     int64 = 12
	aKeyLen        = 32
)

func init() {
	zmodule.Args["argon2-saltLen"] = zmodule.Argument{
		Type:    "int",
		Default: 16,
		Usage:   "Set the salt length of argon2.",
	}
	zmodule.Args["argon2-n"] = zmodule.Argument{
		Type:    "int",
		Default: 3,
		Usage:   "Set the iterations count of argon2.",
	}
	zmodule.Args["argon2-par"] = zmodule.Argument{
		Type:    "int",
		Default: 1,
		Usage:   "Set the parallelism threads amount of argon2.",
	}
	zmodule.Args["argon2-mem"] = zmodule.Argument{
		Type:    "int",
		Default: 12,
		Usage:   "Set the memory usage of argon2.",
	}
	zmodule.Args["argon2-keyLen"] = zmodule.Argument{
		Type:    "int",
		Default: 32,
		Usage:   "Set the key length of argon2.",
	}
	inits["argon2"] = func() {
		encrypt = Argon2
		compare = Argon2Compare
		aSaltLen = config.GetInt("argon2-saltLen")
		aN = config.GetInt("argon2-n")
		aPar = config.GetInt("argon2-par")
		aMem = config.GetI64("argon2-mem")
		aKeyLen = config.GetInt("argon2-keyLen")
	}
}

func randomSalt() ([]byte, error) {
	salt := make([]byte, aSaltLen)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// Argon2 .
func Argon2(plain string) ([]byte, error) {
	hashed, err := sha512bytes(plain)
	if err != nil {
		return nil, err
	}
	salt, err := randomSalt()
	if err != nil {
		return nil, err
	}
	key, err := argon2.Key(hashed, salt, aN, aPar, aMem, aKeyLen)
	if err != nil {
		return nil, err
	}
	secret := append(key, salt...)
	return secret, nil
}

// Argon2Compare .
func Argon2Compare(secret []byte, plain string) error {
	hashed, err := sha512bytes(plain)
	if err != nil {
		return err
	}
	salt := secret[aKeyLen:]
	key, err := argon2.Key(hashed, salt, aN, aPar, aMem, aKeyLen)
	if err != nil {
		return err
	}
	if bytes.Compare(secret[0:aKeyLen], key) != 0 {
		return errors.New("plain text not match with secret")
	}
	return nil
}

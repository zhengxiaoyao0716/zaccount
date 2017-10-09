package encrypt

import (
	"github.com/zhengxiaoyao0716/zmodule"
	"github.com/zhengxiaoyao0716/zmodule/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	bCost = bcrypt.DefaultCost
)

func init() {
	zmodule.Args["bcrypt-cost"] = zmodule.Argument{
		Type:    "int",
		Default: bcrypt.DefaultCost,
		Usage:   "Set the cost of bcrypt.",
	}
	inits["bcrypt"] = func() {
		encrypt = Bcrypt
		compare = BcryptCompare
		bCost = config.GetInt("bcrypt-cost")
	}
}

// Bcrypt .
func Bcrypt(plain string) ([]byte, error) {
	hashed, err := sha512bytes(plain)
	if err != nil {
		return nil, err
	}
	secret, err := bcrypt.GenerateFromPassword(hashed, bCost)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// BcryptCompare .
func BcryptCompare(secret []byte, plain string) error {
	hashed, err := sha512bytes(plain)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(secret, hashed)
}

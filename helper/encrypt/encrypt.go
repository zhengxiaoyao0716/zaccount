package encrypt

import (
	"crypto/sha512"
	"os"

	"github.com/zhengxiaoyao0716/util/console"
	"github.com/zhengxiaoyao0716/util/cout"
	"github.com/zhengxiaoyao0716/zmodule"

	"github.com/zhengxiaoyao0716/zmodule/config"
	"github.com/zhengxiaoyao0716/zmodule/event"
)

var (
	encrypt = Bcrypt
	compare = BcryptCompare
	inits   = map[string]func(){}
)

func init() {
	zmodule.Args["encrypt"] = zmodule.Argument{
		Type:    "string",
		Default: "bcrypt",
		Usage:   "-encrypt <bcrypt|argon2> : select the encrypt method.",
	}
	event.OnInit(func(e event.Event) error {
		event.On(event.KeyStart, func(event.Event) error {
			m := config.GetString("encrypt")
			init, ok := inits[m]
			if !ok {
				console.Log("unknown encrypt method: %s", cout.Err(m))
				os.Exit(1)
			}
			init()

			return nil
		})
		return nil
	})
}

// Encrypt .
func Encrypt(plain string) ([]byte, error) {
	secret, err := encrypt(plain)
	if err != nil {
		return nil, err
	}
	secret, err = AesEncrry(secret, pepper)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// Compare .
func Compare(secret []byte, plain string) error {
	secret, err := AesDecrypt(secret, pepper)
	if err != nil {
		return err
	}
	return compare(secret, plain)
}

func sha512bytes(plain string) ([]byte, error) {
	h := sha512.New()
	if _, err := h.Write([]byte(plain)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

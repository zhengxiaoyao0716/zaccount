package encrypt

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestBcrypt(t *testing.T) {
	hashed, err := Bcrypt("password")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", hashed)
	if err := BcryptCompare(hashed, "password"); err != nil {
		t.Error(err)
	}
}
func TestArgon2(t *testing.T) {
	hashed, err := Argon2("password")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", hashed)
	if err := Argon2Compare(hashed, "password"); err != nil {
		t.Error(err)
	}
}

func TestAes(t *testing.T) {
	plain := []byte("123456")
	secret, err := AesEncrry(plain, pepper)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", secret)
	restore, err := AesDecrypt(secret, pepper)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(plain, restore) != 0 {
		t.Error("plain text and the restore result after encrypt not match")
	}
}

func TestEncrypt(t *testing.T) { testEncrypt("password") }
func testEncrypt(plain string) {
	hashed, err := Encrypt(plain)
	if err != nil {
		log.Fatalln(err)
	}
	if err := Compare(hashed, plain); err != nil {
		log.Fatalln(err)
	}
}

func BenchmarkEcrypt(b *testing.B) {
	encrypt = Bcrypt
	compare = BcryptCompare
	// bCost = 4
	plain := "password"
	b.SetBytes(int64(len(plain)))
	for i := 0; i < b.N; i++ {
		testEncrypt(plain)
	}
}

func BenchmarkArgon2(b *testing.B) {
	encrypt = Argon2
	compare = Argon2Compare
	plain := "password"
	b.SetBytes(int64(len(plain)))
	for i := 0; i < b.N; i++ {
		testEncrypt(plain)
	}
}

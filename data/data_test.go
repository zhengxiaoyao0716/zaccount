package data

import (
	"testing"

	"github.com/zhengxiaoyao0716/zmodule/event"
)

func TestResetPapper(t *testing.T) {
	defer event.Emit(event.KeyStop, nil)
	id := "TEST_RESET_PAPPER"

	if err := db.Update(PutAccount(id, "password")); err != nil {
		t.Fatal(err)
	}

	if err := db.View(CheckAccount(id, "password")); err != nil {
		t.Error(err)
	}

	if err := ResetPepper("", "This is a test pepper"); err != nil {
		t.Error(err)
	}

	if err := db.View(CheckAccount(id, "password")); err != nil {
		t.Error(err)
	}

	if err := db.Update(DeleteAccount(id)); err != nil {
		t.Error(err)
	}

	if err := ResetPepper("This is a test pepper", ""); err != nil {
		t.Error(err)
	}
}

func init() {
	dbPath = ".data_test.db"
	event.Init(map[string]string{})
	Init()
}

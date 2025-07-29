package auth_test

import (
	"testing"

	"github.com/Kharitopolus/Myberries/internal/infrastructure/auth"
)

func TestPasswordManager(t *testing.T) {
	pm := auth.PasswordManager{}

	_, err := pm.GetHash("")
	if err == nil {
		t.Errorf("err == nil for GetHash from empty string")
	}

	hackmeHash, err := pm.GetHash("hackme")
	if err != nil {
		t.Errorf("got error: '%v' for password 'hackme' ", err)
	}

	err = pm.CheckHash(hackmeHash, "hackme")
	if err != nil {
		t.Errorf(
			"got error: '%v' for CheckHash from hackmeHash and 'hackme'",
			err,
		)
	}

	err = pm.CheckHash(hackmeHash, "lfkdjsalkfjdlsk")
	if err == nil {
		t.Errorf("err == nil for CheckHash for hackmeHash and random string")
	}
}

package auth_test

import (
	"testing"
	"time"

	"github.com/Kharitopolus/Myberries/internal/infrastructure/auth"
	"github.com/google/uuid"
)

func TestTokenManager(t *testing.T) {
	tm := auth.TokenManager{
		TokenSecret:         "7V22DsRkMeXRrv8V3bzzwma+++QFcHjr4l/e9nFMEkbiNYUAWa6Wq2lqCy0NcZa3ykWeS/16PHQA3+YXMjF+Lw==",
		AccessTokenExpTime:  time.Duration(1) * time.Second,
		RefreshTokenExpTime: time.Duration(1) * time.Second,
	}

	id := uuid.New()

	at, err := tm.MakeAccessToken(id)
	if err != nil {
		t.Errorf("err: '%v' on MakeAccessToken(id)", err)
	}

	rt, err := tm.MakeRefreshToken(id)
	if err != nil {
		t.Errorf("err: '%v' on MakeRefreshToken(id)", err)
	}

	atid, err := tm.ValidateAccessToken(at)
	if err != nil {
		t.Errorf("err: '%v' on ValidateAccessToken(at)", err)
	}
	if atid != id {
		t.Error("atid from ValidateAccessToken != id")
	}

	rtid, err := tm.ValidateRefreshToken(rt)
	if err != nil {
		t.Errorf("err: '%v' on ValidateRefreshToken(at)", err)
	}
	if rtid != id {
		t.Error("rtid from ValidateRefreshToken != id")
	}

	_, err = tm.ValidateAccessToken(rt)
	if err == nil {
		t.Error("don't have error on ValidateAccessToken(rt)")
	}
	_, err = tm.ValidateRefreshToken(at)
	if err == nil {
		t.Error("don't have error on ValidateRefreshToken(at)")
	}

	time.Sleep(time.Duration(1) * time.Second)

	_, err = tm.ValidateAccessToken(at)
	if err == nil {
		t.Error(
			"don't have error on ValidateAccessToken(at) when token expired",
		)
	}
	_, err = tm.ValidateRefreshToken(rt)
	if err == nil {
		t.Error(
			"don't have error on ValidateRefreshToken(rt) when token expired",
		)
	}
}

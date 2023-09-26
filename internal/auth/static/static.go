package static

import (
	"fmt"
	"time"
)

const authFalseDelay = time.Second * 10

type Static struct {
	StaticUser string
	StaticPass string
}

func NewStatic(user string, pass string) (*Static, error) {
	if user == "" || pass == "" {
		return &Static{}, fmt.Errorf("user or pass cannot be empty")
	}
	return &Static{
		StaticUser: user,
		StaticPass: pass,
	}, nil
}

func (s Static) Valid(user string, pass string) bool {
	if s.StaticUser == user && s.StaticPass == pass {
		return true
	} else {
		time.Sleep(authFalseDelay)
		return false
	}
}

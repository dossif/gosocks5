package static

type Static struct {
	StaticUser string
	StaticPass string
}

func NewStatic(user string, pass string) *Static {
	return &Static{
		StaticUser: user,
		StaticPass: pass,
	}
}

func (s Static) Valid(user string, pass string) bool {
	if s.StaticUser == user && s.StaticPass == pass {
		return true
	} else {
		return false
	}
}

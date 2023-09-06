package ldap

// Credentials enables using a map directly as a credential store
type Credentials map[string]string

func (s Credentials) Valid(user, password string) bool {
	pass, ok := s[user]
	if !ok {
		return false
	}
	return password == pass
}

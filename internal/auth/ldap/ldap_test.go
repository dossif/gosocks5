package ldap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStaticCredentials(t *testing.T) {
	cred := LdapCredentials{
		"user1": "pass1",
		"user2": "",
	}
	assert.True(t, cred.Valid("user1", "pass1"))
	assert.True(t, cred.Valid("user2", ""))
	assert.False(t, cred.Valid("user1", ""))
	assert.False(t, cred.Valid("", "pass1"))
}

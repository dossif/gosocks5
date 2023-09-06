package ldap

import (
	"fmt"
	"github.com/dossif/gosocks5/pkg/logger"
	"github.com/go-ldap/ldap"
)

type Ldap struct {
	Client ldap.Client
	BaseDn string
	Filter string
	Log    *logger.Logger
}

func NewLdap(log logger.Logger, url, bindUser, bindPass, baseDn, filter string) (*Ldap, error) {
	l, err := ldap.DialURL(url)
	if err != nil {
		return &Ldap{}, fmt.Errorf("failed to connect to ldap server %v: %v", url, err)
	}
	err = l.Bind(bindUser, bindPass)
	if err != nil {
		return &Ldap{}, fmt.Errorf("failed to bind ldap with user %v: %v", bindUser, err)
	}
	return &Ldap{
		Client: l,
		BaseDn: baseDn,
		Filter: filter,
		Log:    &log,
	}, nil
}

func (l *Ldap) Valid(user string, pass string) bool {
	filter := fmt.Sprintf(l.Filter, ldap.EscapeFilter(user))
	searchReq := ldap.NewSearchRequest(l.BaseDn, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{}, []ldap.Control{})
	result, err := l.Client.Search(searchReq)
	if err != nil {
		l.Log.Lg.Warn().Msgf("failed to search ldap user %v with filter %v: %v", user, filter, err)
		return false
	}
	if len(result.Entries) == 0 {
		l.Log.Lg.Warn().Msgf("ldap user %v not found", user)
		l.Log.Lg.Debug().Msgf("ldap filter %v returns empty result", filter)
		return false
	}
	var userDn string
	for _, res := range result.Entries {
		userDn = res.DN
	}
	err = l.Client.Bind(userDn, pass)
	if err != nil {
		l.Log.Lg.Warn().Msgf("failed to auth ldap user %v", user)
		l.Log.Lg.Debug().Msgf("failed to auth ldap user %v: %v", userDn, err)
		return false
	}
	return true
}

func (l *Ldap) Close() {
	l.Client.Close()
}

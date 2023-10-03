package ldap

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/dossif/gosocks5/pkg/logger"
	"github.com/go-ldap/ldap"
	"github.com/jellydator/ttlcache/v3"
	"time"
)

const (
	authTrueCache  = time.Second * 30
	authFalseCache = time.Second * 10
)

type Ldap struct {
	Url      string
	BindUser string
	BindPass string
	BaseDn   string
	Filter   string
	Log      *logger.Logger
	Cache    *ttlcache.Cache[string, string]
}

func NewLdap(log logger.Logger, url, bindUser, bindPass, baseDn, filter string) (*Ldap, error) {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](ttlcache.DefaultTTL), ttlcache.WithDisableTouchOnHit[string, string](),
	)
	return &Ldap{
		Url:      url,
		BindUser: bindUser,
		BindPass: bindPass,
		BaseDn:   baseDn,
		Filter:   filter,
		Log:      &log,
		Cache:    cache,
	}, nil
}

func (l *Ldap) Valid(user string, pass string) bool {
	l.Cache.DeleteExpired()
	switch l.checkCache(user, pass) {
	case 1:
		return true
	case 0:
		return false
	case -1: // expired
		switch l.checkLdap(user, pass) {
		case false:
			l.Cache.Set(user, "fail", authFalseCache)
			return false
		case true:
			l.Cache.Set(user, hashPass(pass), authTrueCache)
			return true
		}
	}
	return false
}

// checkCache
// return 1  = pass valid
// return 0  = pass invalid
// return -1 = pass expired
func (l *Ldap) checkCache(user string, pass string) int {
	if l.Cache.Has(user) != true {
		return -1 // expired
	}
	hashedPass := hashPass(pass)
	cachedPass := l.Cache.Get(user).Value()
	if hashedPass == cachedPass {
		return 1 // valid
	} else {
		return 0 // invalid
	}
}

func (l *Ldap) checkLdap(user string, pass string) bool {
	client, err := ldap.DialURL(l.Url)
	if err != nil {
		l.Log.Lg.Error().Msgf("failed to connect to ldap server %v: %v", l.Url, err)
		return false
	}
	err = client.Bind(l.BindUser, l.BindPass)
	if err != nil {
		l.Log.Lg.Error().Msgf("failed to bind ldap with user %v: %v", l.BindUser, err)
		return false
	}
	filter := fmt.Sprintf(l.Filter, ldap.EscapeFilter(user))
	searchReq := ldap.NewSearchRequest(l.BaseDn, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{}, []ldap.Control{})
	result, err := client.Search(searchReq)
	if err != nil {
		l.Log.Lg.Warn().Msgf("failed to search ldap user %v with filter %v: %v", user, filter, err)
		return false
	}
	if len(result.Entries) == 0 {
		l.Log.Lg.Debug().Msgf("ldap filter %v returns empty result", filter)
		return false
	}
	var userDn string
	for _, res := range result.Entries {
		userDn = res.DN
	}
	err = client.Bind(userDn, pass)
	if err != nil {
		l.Log.Lg.Debug().Msgf("failed to auth ldap user %v: %v", userDn, err)
		return false
	}
	defer client.Close()
	return true
}

func hashPass(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	bs := h.Sum(nil)
	b64 := base64.StdEncoding.EncodeToString(bs)
	return b64
}

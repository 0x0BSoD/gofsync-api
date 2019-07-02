package user

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/ldap.v3"
)

func LdapGet(username string, password string, ctx *GlobalCTX) (string, error) {
	// The username and password we want to chec
	bindUsername := ctx.Config.LDAP.BindUser
	bindPassword := ctx.Config.LDAP.BindPassword
	var l *ldap.Conn
	l, err := ldap.Dial("tcp",
		fmt.Sprintf("%s:%d", ctx.Config.LDAP.LdapServer[0], ctx.Config.LDAP.LdapServerPort))
	if err != nil {
		l, err = ldap.Dial("tcp",
			fmt.Sprintf("%s:%d", ctx.Config.LDAP.LdapServer[1], ctx.Config.LDAP.LdapServerPort))
		if err != nil {
			return "", err
		}
	}
	defer l.Close()

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return "", err
	}

	// First bind with a read only user
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		return "", err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		ctx.Config.LDAP.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(ctx.Config.LDAP.MatchStr, username),
		[]string{"*"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", fmt.Errorf("user does not exist or too many entries returned")
	}

	userDN := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userDN, password)
	if err != nil {
		return "", err
	}

	// Rebind as the read only user for any further queries
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		return "", err
	}

	return userDN, nil
}

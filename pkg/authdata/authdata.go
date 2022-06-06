package authdata

import (
	"fmt"
)

type AuthData interface {
	Account() string
	Key() string
	Secret() string
}

type ServiceAuthData interface {
	AuthData
	Host() string
	Port() int
	Username() string
	Password() string
}

type authData struct {
	// account is some sort of account identifier
	account string
	// key is an API key, auth token, etc.
	key string
	// secret is only defined for APIs which provide an API secret (such as Gemini)
	secret string
}

func NewAuthData(account, key, secret string) AuthData {
	return authData{
		account: account,
		key:     key,
		secret:  secret,
	}
}

func (a authData) Account() string {
	return a.account
}

func (a authData) Key() string {
	return a.key
}

func (a authData) Secret() string {
	return a.secret
}

type serviceAuthData struct {
	host     string
	username string
	password string
	port     int
}

func NewServiceAuthData(host string, port int, username, password string) ServiceAuthData {
	return serviceAuthData{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s serviceAuthData) Host() string {
	return s.host
}

func (s serviceAuthData) Port() int {
	return s.port
}

func (s serviceAuthData) Username() string {
	return s.username
}

func (s serviceAuthData) Password() string {
	return s.password
}

func (s serviceAuthData) Account() string {
	return fmt.Sprintf("%[1]s:%[2]d", s.Host(), s.Port())
}

func (s serviceAuthData) Key() string {
	return s.username
}

func (s serviceAuthData) Secret() string {
	return s.password
}

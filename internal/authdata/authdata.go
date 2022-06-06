package authdata

type AuthData interface {
	Account() string
	Key() string
	Secret() string
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

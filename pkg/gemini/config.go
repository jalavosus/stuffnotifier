package gemini

type Config struct {
	Auth *AuthConfig `json:"auth,omitempty" yaml:"auth,omitempty"`
}

type AuthConfig struct {
	ApiKey    string `json:"api_key" yaml:"api_key"`
	ApiSecret string `json:"api_secret" yaml:"api_secret"`
}

func (c AuthConfig) Account() string {
	return ""
}

func (c AuthConfig) Key() string {
	return c.ApiKey
}

func (c AuthConfig) Secret() string {
	return c.ApiSecret
}

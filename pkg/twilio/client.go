package twilio

import (
	"github.com/pkg/errors"
	"github.com/twilio/twilio-go"

	"github.com/stoicturtle/stuffnotifier/internal/authdata"
	"github.com/stoicturtle/stuffnotifier/internal/authdata/env"
	"github.com/stoicturtle/stuffnotifier/pkg/singleton"
)

type Client struct {
	auth       authdata.AuthData
	client     *twilio.RestClient
	sendNumber string
}

func NewClient() (*Client, error) {
	auth, authErr := authdata.TwilioAPITokenAuth()
	if authErr != nil {
		return nil, authErr
	}

	sendNumber, sendNumberErr := env.TwilioFromNumber()
	if sendNumberErr != nil {
		return nil, sendNumberErr
	}

	return newClient(auth, sendNumber), nil
}

func NewClientFromConfig(conf Config) (*Client, error) {
	switch {
	case conf.Auth.AuthToken != nil:
		return newClient(conf.Auth.AuthToken, conf.SenderNumber), nil
	case conf.Auth.ApiKey != nil:
		return newClient(conf.Auth.ApiKey, conf.SenderNumber), nil
	default:
		return nil, errors.New("auth_token or api_key must be configured in config")
	}
}

func newClient(authData authdata.AuthData, twilioNumber string) *Client {
	return &Client{
		auth:       authData,
		sendNumber: twilioNumber,
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			AccountSid: authData.Account(),
			Username:   authData.Key(),
			Password:   authData.Secret(),
		}),
	}
}

var (
	singletonInstance = newClientSingleton()
)

type clientSingleton struct {
	*singleton.BaseInstance[Client, Config]
	client *Client
}

// should be called exactly one time, and only by the `var singletonInstance = ...` declaration.
func newClientSingleton() *clientSingleton {
	c := new(clientSingleton)
	c.BaseInstance = singleton.NewBaseInstance[Client, Config](c.init, c.initFromConfig)

	return c
}

func ClientSingleton(conf *Config) (singleton.Singleton[Client, Config], error) {
	var initErr error

	if conf != nil {
		initErr = singletonInstance.InitFromConfig(*conf)
	} else {
		initErr = singletonInstance.Init()
	}

	return singletonInstance, initErr
}

func (c *clientSingleton) Client() *Client {
	return c.client
}

func (c *clientSingleton) init() error {
	client, initErr := NewClient()
	if initErr != nil {
		return initErr
	}

	c.client = client

	return nil
}

func (c *clientSingleton) initFromConfig(conf Config) error {
	client, initErr := NewClientFromConfig(conf)
	if initErr != nil {
		return initErr
	}

	c.client = client

	return nil
}

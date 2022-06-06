package gemini

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/stoicturtle/stuffnotifier/internal/nonceticker"
	"github.com/stoicturtle/stuffnotifier/pkg/authdata"
)

type NoncePayload struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
}

func (np NoncePayload) Serialize() []byte {
	marshalled, err := json.Marshal(np)
	if err != nil {
		logger.Panic("error marshalling nonce payload", zap.Error(err))
	}

	return marshalled
}

func (np NoncePayload) Encode() (encoded []byte) {
	serialized := np.Serialize()
	encoded = make([]byte, base64.StdEncoding.EncodedLen(len(serialized)))

	base64.StdEncoding.Encode(encoded, serialized)

	return
}

func (np NoncePayload) EncodeString() (encoded string) {
	serialized := np.Serialize()
	encoded = base64.StdEncoding.EncodeToString(serialized)

	return
}

func (np NoncePayload) HashHmac(key string) (sig string) {
	h := hmac.New(crypto.SHA384.New, []byte(key))
	s := h.Sum(np.Encode())

	sig = hex.EncodeToString(s)

	return
}

func BuildNonceWithPayload(authData authdata.AuthData, endpoint string) (sig, payload string) {
	nonceTicker := nonceticker.GetNonceTicker()

	noncePayload := NoncePayload{
		Request: endpoint,
		Nonce:   fmt.Sprintf("%d", nonceTicker.GetTick()),
	}

	sig = noncePayload.HashHmac(authData.Secret())
	payload = noncePayload.EncodeString()

	return
}

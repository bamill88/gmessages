package payload

import (
	"go.mau.fi/mautrix-gmessages/libgm/binary"
	"go.mau.fi/mautrix-gmessages/libgm/util"
)

func GetWebEncryptionKey(WebPairKey []byte) ([]byte, *binary.AuthenticationContainer, error) {
	id := util.RandomUUIDv4()
	payload := &binary.AuthenticationContainer{
		AuthMessage: &binary.AuthenticationMessage{
			RequestID:        id,
			TachyonAuthToken: WebPairKey,
			ConfigVersion:    ConfigMessage,
		},
	}
	encodedPayload, err2 := binary.EncodeProtoMessage(payload)
	if err2 != nil {
		return nil, payload, err2
	}
	return encodedPayload, payload, nil
}

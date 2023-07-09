package libgm

import (
	"encoding/json"
	"fmt"

	"go.mau.fi/mautrix-gmessages/libgm/pblite"

	"go.mau.fi/mautrix-gmessages/libgm/binary"
)

var skipCount int32

func (r *RPC) HandleRPCMsg(msgArr []interface{}) {
	response, decodeErr := pblite.DecodeAndDecryptInternalMessage(msgArr, r.client.authData.Cryptor)
	if decodeErr != nil {
		r.client.Logger.Error().Err(fmt.Errorf("failed to deserialize response %s", msgArr)).Msg("rpc deserialize msg err")
		return
	}
	//r.client.Logger.Debug().Any("byteLength", len(data)).Any("unmarshaled", response).Any("raw", string(data)).Msg("RPC Msg")
	if response == nil {
		r.client.Logger.Error().Err(fmt.Errorf("response data was nil %s", msgArr)).Msg("rpc msg data err")
		return
	}
	//r.client.Logger.Debug().Any("response", response).Msg("decrypted & decoded response")
	_, waitingForResponse := r.client.sessionHandler.requests[response.Data.RequestId]

	//r.client.Logger.Info().Any("raw", msgArr).Msg("Got msg")
	//r.client.Logger.Debug().Any("waiting", waitingForResponse).Msg("got request! waiting?")
	r.client.sessionHandler.addResponseAck(response.ResponseId)
	if waitingForResponse {
		r.client.sessionHandler.respondToRequestChannel(response)
	} else {
		switch response.BugleRoute {
		case binary.BugleRoute_PairEvent:
			r.client.handlePairingEvent(response)
		case binary.BugleRoute_DataEvent:
			if skipCount > 0 {
				skipCount--
				r.client.Logger.Info().Any("action", response.Data.Action).Any("toSkip", skipCount).Msg("Skipped DataEvent")
				return
			}
			r.client.handleUpdatesEvent(response)
		default:
			r.client.Logger.Debug().Any("res", response).Msg("Got unknown bugleroute")
		}
	}

}

func (r *RPC) tryUnmarshalJSON(jsonData []byte, msgArr *[]interface{}) error {
	err := json.Unmarshal(jsonData, &msgArr)
	return err
}

func (r *RPC) HandleByLength(data []byte) {
	r.client.Logger.Debug().Any("byteLength", len(data)).Any("corrupt raw", string(data)).Msg("RPC Corrupt json")
}

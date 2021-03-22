package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ RequestSpec           = &OracleRequestPacketData{}
	_ RequestWithSenderSpec = &MsgRequestData{}
)

// RequestSpec captures the essence of what it means to be a request-making object.
type RequestSpec interface {
	GetOracleScriptID() OracleScriptID
	GetCalldata() []byte
	GetAskCount() uint64
	GetMinCount() uint64
	GetClientID() string
}

type RequestWithSenderSpec interface {
	RequestSpec
	GetSender() sdk.AccAddress
}

package nameservice

import (
	"github.com/yxuco/cosmos-nameservice/x/nameservice/keeper"
	"github.com/yxuco/cosmos-nameservice/x/nameservice/types"
)

// alias used by handler.go

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey // nameservice store key
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper        = keeper.NewKeeper  // function to create new nameservice KV store
	NewQuerier       = keeper.NewQuerier // function to create querier handler
	NewMsgBuyName    = types.NewMsgBuyName
	NewMsgSetName    = types.NewMsgSetName
	NewMsgDeleteName = types.NewMsgDeleteName
	NewWhois         = types.NewWhois
	ModuleCdc        = types.ModuleCdc
	RegisterCodec    = types.RegisterCodec
)

type (
	Keeper          = keeper.Keeper // nameservice KV store
	MsgSetName      = types.MsgSetName
	MsgBuyName      = types.MsgBuyName
	MsgDeleteName   = types.MsgDeleteName
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	Whois           = types.Whois
)

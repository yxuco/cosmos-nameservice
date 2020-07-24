package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/yxuco/cosmos-nameservice/x/nameservice/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
// called on chain start to import genesis state into the keeper
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	for _, record := range data.WhoisRecords {
		k.SetWhois(ctx, record.Value, record)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data types.GenesisState) {
	var records []types.Whois
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		name := string(iterator.Key())
		whois := k.GetWhois(ctx, name)
		records = append(records, whois)

	}
	return types.NewGenesisState(records)
}

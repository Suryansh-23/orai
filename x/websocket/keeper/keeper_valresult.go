package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/exported"
	"github.com/oraichain/orai/x/websocket/types"
)

// DefaultValResultI returns the default ai data source object
func (k Keeper) DefaultValResultI() exported.ValResultI {
	return k.DefaultValResult()
}

// DefaultValResult is a default constructor for the validator result
func (k Keeper) DefaultValResult() types.ValResult {
	return types.ValResult{
		Validator: nil,
		Result:    []byte{},
	}
}

// GetKeyResultSuccess is a getter to collect the result success key for validator result verification using by other modules.
func (k Keeper) GetKeyResultSuccess() string {
	return types.ResultSuccess
}

// NewValResult is a wrapper function of the websocket module that allow others to initiate a new valresult entity through the keeper
func (k Keeper) NewValResult(val sdk.ValAddress, result []byte) exported.ValResultI {
	return types.NewValResult(val, result)
}

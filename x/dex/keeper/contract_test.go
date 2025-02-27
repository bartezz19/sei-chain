package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/sei-protocol/sei-chain/testutil/keeper"
	"github.com/sei-protocol/sei-chain/x/dex/types"
	"github.com/stretchr/testify/require"
)

func TestChargeRentForGas(t *testing.T) {
	keeper, ctx := keepertest.DexKeeper(t)
	keeper.SetParams(ctx, types.Params{SudoCallGasPrice: sdk.NewDecWithPrec(1, 1), PriceSnapshotRetention: 1})
	err := keeper.SetContract(ctx, &types.ContractInfoV2{
		Creator:      keepertest.TestAccount,
		ContractAddr: keepertest.TestContract,
		CodeId:       1,
		RentBalance:  1000000,
	})
	require.Nil(t, err)
	err = keeper.ChargeRentForGas(ctx, keepertest.TestContract, 5000000, 0)
	require.Nil(t, err)
	contract, err := keeper.GetContract(ctx, keepertest.TestContract)
	require.Nil(t, err)
	require.Equal(t, uint64(500000), contract.RentBalance)
	err = keeper.ChargeRentForGas(ctx, keepertest.TestContract, 6000000, 0)
	require.NotNil(t, err)
	contract, err = keeper.GetContract(ctx, keepertest.TestContract)
	require.Nil(t, err)
	require.Equal(t, uint64(0), contract.RentBalance)
	err = keeper.SetContract(ctx, &types.ContractInfoV2{
		Creator:      keepertest.TestAccount,
		ContractAddr: keepertest.TestContract,
		CodeId:       1,
		RentBalance:  1000000,
	})
	require.Nil(t, err)
	err = keeper.ChargeRentForGas(ctx, keepertest.TestContract, 5000000, 4000000)
	require.Nil(t, err)
	contract, err = keeper.GetContract(ctx, keepertest.TestContract)
	require.Nil(t, err)
	require.Equal(t, uint64(900000), contract.RentBalance)
	err = keeper.ChargeRentForGas(ctx, keepertest.TestContract, 5000000, 6000000)
	require.Nil(t, err)
	contract, err = keeper.GetContract(ctx, keepertest.TestContract)
	require.Nil(t, err)
	require.Equal(t, uint64(900000), contract.RentBalance)

	// delete contract
	keeper.DeleteContract(ctx, keepertest.TestContract)
	_, err = keeper.GetContract(ctx, keepertest.TestContract)
	require.NotNil(t, err)
}

func TestGetAllContractInfo(t *testing.T) {
	keeper, ctx := keepertest.DexKeeper(t)
	keeper.SetParams(ctx, types.Params{SudoCallGasPrice: sdk.NewDecWithPrec(1, 1), PriceSnapshotRetention: 1})
	keeper.SetContract(ctx, &types.ContractInfoV2{
		Creator:      keepertest.TestAccount,
		ContractAddr: keepertest.TestContract,
		CodeId:       1,
		RentBalance:  1000000,
	})
	keeper.SetContract(ctx, &types.ContractInfoV2{
		Creator:      "ta2",
		ContractAddr: "tc2",
		CodeId:       2,
		RentBalance:  1000000,
	})
	contracts := keeper.GetAllContractInfo(ctx)
	require.Equal(t, uint64(1000000), contracts[0].RentBalance)
	require.Equal(t, uint64(1000000), contracts[1].RentBalance)
	require.Equal(t, uint64(1), contracts[0].CodeId)
	require.Equal(t, uint64(2), contracts[1].CodeId)
	require.Equal(t, keepertest.TestAccount, contracts[0].Creator)
	require.Equal(t, keepertest.TestContract, contracts[0].ContractAddr)
	require.Equal(t, "ta2", contracts[1].Creator)
	require.Equal(t, "tc2", contracts[1].ContractAddr)

}

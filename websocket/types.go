package websocket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

type KYCRequest struct {
	ImageHash string    `json:"image_hash"`
	ImageName string    `json:"image_name"`
	AIRequest AIRequest `json:"ai_request"`
}

func NewKYCRequest(
	imageHash string,
	imageName string,
	aiRequest AIRequest,
) KYCRequest {
	return KYCRequest{
		ImageHash: imageHash,
		ImageName: imageName,
		AIRequest: aiRequest,
	}
}

type PriceRequest struct {
	AIRequest AIRequest `json:"ai_request"`
}

func NewPriceRequest(
	aiRequest AIRequest,
) PriceRequest {
	return PriceRequest{
		AIRequest: aiRequest,
	}
}

type AIRequest struct {
	RequestID        string         `json:"request_id"`
	OracleScriptName string         `json:"oscript_name"`
	Creator          sdk.AccAddress `json:"creator"`
	ValidatorCount   int64          `json:"validator_count"`
	Input            string         `json:"request_input"`
	ExpectedOutput   string         `json:"expected_output"`
}

func NewAIRequest(
	requestID string,
	oscriptName string,
	creator sdk.AccAddress,
	validatorCount int64,
	input string,
	expectedOutput string,
) AIRequest {
	return AIRequest{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		ValidatorCount:   validatorCount,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}

type Report struct {
	RequestID string         `json:"request_id"`
	Validator sdk.ValAddress `json:"validator"`
	// DataSourceResults are the actual results, not from the test cases
	DataSourceResults []types.DataSourceResult `json:"data_source_results"`
	TestCaseResults   []types.TestCaseResult   `json:"test_case_results"`
	Reporter          sdk.AccAddress           `json:"reporter"`
	Fees              sdk.Coins                `json:"report_fee"`
	AggregatedResult  []byte                   `json:"aggregated_result"`
}

// NewReport is a constructor function for MsgCreateReport
func NewReport(
	requestID string,
	validator sdk.ValAddress,
	dataSourceResults []types.DataSourceResult,
	testCaseResults []types.TestCaseResult,
	reporter sdk.AccAddress,
	fees sdk.Coins,
	aggregatedResult []byte,
) Report {
	return Report{
		RequestID:         requestID,
		Validator:         validator,
		DataSourceResults: dataSourceResults,
		TestCaseResults:   testCaseResults,
		Reporter:          reporter,
		Fees:              fees,
		AggregatedResult:  aggregatedResult,
	}
}

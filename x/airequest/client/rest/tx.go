package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (

	// "bytes"
	// "net/http"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/segmentio/ksuid"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	// "github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type setPriceRequestReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	OracleScriptName string       `json:"oracle_script_name"`
	Input            string       `json:"input"`
	ExpectedOutput   string       `json:"expected_output"`
	Fees             string       `json:"fees"`
	ValidatorCount   int          `json:"validator_count"`
}

type setKYCRequestReq struct {
	OracleScriptName string       `json:"oracle_script_name"`
	From             string       `json:"from"`
	Memo             string       `json:"memo"`
	ChainID          string       `json:"chain_id"`
	AccountNumber    uint64       `json:"account_number"`
	Sequence         uint64       `json:"sequence"`
	Fees             sdk.Coins    `json:"fees"`
	GasPrices        sdk.DecCoins `json:"gas_prices"`
	Gas              string       `json:"gas"`
	GasAdjustment    string       `json:"gas_adjustment"`
	Simulate         bool         `json:"simulate"`
	Input            string       `json:"input"`
	ExpectedOutput   string       `json:"expected_output"`
	ValidatorCount   int          `json:"validator_count"`
}

type setClassificationRequestReq struct {
	OracleScriptName string       `json:"oracle_script_name"`
	From             string       `json:"from"`
	Memo             string       `json:"memo"`
	ChainID          string       `json:"chain_id"`
	AccountNumber    uint64       `json:"account_number"`
	Sequence         uint64       `json:"sequence"`
	Fees             sdk.Coins    `json:"fees"`
	GasPrices        sdk.DecCoins `json:"gas_prices"`
	Gas              string       `json:"gas"`
	GasAdjustment    string       `json:"gas_adjustment"`
	Simulate         bool         `json:"simulate"`
	Input            string       `json:"input"`
	ExpectedOutput   string       `json:"expected_output"`
	ValidatorCount   int          `json:"validator_count"`
}

type setIPFSImage struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/kycreq", storeName),
		setKYCRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/clreq", storeName),
		setClassificationRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/pricereq", storeName),
		setPriceRequestHandlerFn(cliCtx),
	).Methods("POST")
}

func setAIRequestHandlerFn(cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request) setKYCRequestReq {
	// convert from form-data string to correct data type of baseReq
	accNum, err := strconv.Atoi(r.FormValue("account_number"))
	if len(r.FormValue("account_number")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for account number")
		return setKYCRequestReq{}
	}

	// convert from form-data string to correct data type of baseReq
	sequence, err := strconv.Atoi(r.FormValue("sequence"))
	if len(r.FormValue("sequence")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for sequence")
		return setKYCRequestReq{}
	}

	// convert from form-data string to correct data type of baseReq
	simulate, err := strconv.ParseBool(r.FormValue("simulate"))
	if len(r.FormValue("simulate")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for simulate")
		return setKYCRequestReq{}
	}

	// convert from form-data string to correct data type of baseReq
	feeStr := r.FormValue("fees")
	var fees sdk.Coins

	if len(feeStr) != 0 {
		fees, err = sdk.ParseCoins(feeStr)
		if len(r.FormValue("fees")) != 0 && err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for fees")
			return setKYCRequestReq{}
		}
	}

	// convert from form-data string to correct data type of baseReq
	gasPricesStr := r.FormValue("gas_prices")
	var gasPrices sdk.DecCoins

	if len(gasPricesStr) != 0 {
		gasPrices, err = sdk.ParseDecCoins(gasPricesStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for gas price")
			return setKYCRequestReq{}
		}
	}

	// convert from form-data string to correct data type of baseReq
	valCount, err := strconv.Atoi(r.FormValue("validator_count"))
	if len(r.FormValue("validator_count")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for validator count")
		return setKYCRequestReq{}
	}

	// This may be redundant, but it does not affect much so leave this for now.
	req := setKYCRequestReq{
		OracleScriptName: r.FormValue("oracle_script_name"),
		From:             r.FormValue("from"),
		Memo:             r.FormValue("memo"),
		ChainID:          r.FormValue("chain_id"),
		AccountNumber:    uint64(accNum),
		Sequence:         uint64(sequence),
		Fees:             fees,
		GasPrices:        gasPrices,
		Gas:              r.FormValue("gas"),
		GasAdjustment:    r.FormValue("gas_adjustment"),
		Simulate:         simulate,
		Input:            r.FormValue("input"),
		ExpectedOutput:   r.FormValue("expected_output"),
		ValidatorCount:   valCount,
	}

	return req
}

func setKYCRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		imageHash := r.FormValue("image_hash")
		imageName := r.FormValue("image_name")

		// // collect image file from user
		// file, handler, err := r.FormFile("image")
		// if err != nil {
		// 	fmt.Println("Error Retrieving the File")
		// 	fmt.Println(err)
		// 	return
		// }
		// defer file.Close()

		// fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")

		// // Create a temp file in local storage for IPFS http request
		// tempFile, err := ioutil.TempFile("./", "upload-*.png")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// defer tempFile.Close()
		// fileBytes, err := ioutil.ReadAll(file)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// tempFile.Write(fileBytes)

		// fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

		// // Prepare to send the image onto IPFS
		// b, writer, err := filehandling.CreateMultipartFormData("image", tempFile.Name())

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to create multiform data image")
		// 	return
		// }

		// httpReq, err := http.NewRequest("POST", types.IPFSUrl+types.IPFSAdd, &b)
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to create new request for sending image to IPFS")
		// 	return
		// }
		// // Don't forget to set the content type, this will contain the boundary.
		// httpReq.Header.Set("Content-Type", writer.FormDataContentType())

		// client := &http.Client{}
		// resp, err := client.Do(httpReq)

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to execute request for sending image to IPFS")
		// 	return
		// }

		// defer resp.Body.Close()

		// result := setIPFSImage{}

		// // Collect the result in json form. Remember that we need to create a corresponding struct to do this
		// json.NewDecoder(resp.Body).Decode(&result)

		// // After collecting the hash image, we need to clear the image file stored temporary
		// err = os.Remove(tempFile.Name())
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to remove the temporary image file: %s", err.Error()))
		// 	return
		// }

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to push image onto IPFS: %s", err.Error()))
		// 	return
		// }

		req := setAIRequestHandlerFn(cliCtx, w, r)

		// Need to create a baseReq to write tx response. We cannot use baseReq in the AIRequest struct because AIRequest needs to be in form data to be able to send images
		baseReq := rest.BaseReq{
			From:          req.From,
			Memo:          req.Memo,
			ChainID:       req.ChainID,
			AccountNumber: req.AccountNumber,
			Sequence:      req.Sequence,
			Fees:          req.Fees,
			GasPrices:     req.GasPrices,
			Gas:           req.Gas,
			GasAdjustment: req.GasAdjustment,
			Simulate:      req.Simulate,
		}

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(req.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "AVXSD")
			return
		}

		// create the message
		// msg := types.NewMsgSetKYCRequest(result.Hash, handler.Filename, types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees.String(), req.ValidatorCount, req.Input, req.ExpectedOutput))
		msg := types.NewMsgSetKYCRequest(imageHash, imageName, types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees.String(), req.ValidatorCount, req.Input, req.ExpectedOutput))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GHYK")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func setPriceRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req setPriceRequestReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "AVXSD")
			return
		}

		// create the message
		msg := types.NewMsgSetPriceRequest(types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees, req.ValidatorCount, req.Input, req.ExpectedOutput))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GHYK")
			return
		}
		// Collect fees in Coins type. Bug: cannot set fee through json using REST API => This is the workaround
		fees, _ := sdk.ParseCoins(req.Fees)
		baseReq.Fees = fees
		if !baseReq.ValidateBasic(w) {
			return
		}
		fmt.Println("base req: ", baseReq)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func setClassificationRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// // collect image file from user
		// file, handler, err := r.FormFile("image")
		// if err != nil {
		// 	fmt.Println("Error Retrieving the File")
		// 	fmt.Println(err)
		// 	return
		// }
		// defer file.Close()

		// fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")

		// // Create a temp file in local storage for IPFS http request
		// tempFile, err := ioutil.TempFile("./", "upload-*.png")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// defer tempFile.Close()
		// fileBytes, err := ioutil.ReadAll(file)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// tempFile.Write(fileBytes)

		// fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

		// // Prepare to send the image onto IPFS
		// b, writer, err := filehandling.CreateMultipartFormData("image", tempFile.Name())

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to create multiform data image")
		// 	return
		// }

		// httpReq, err := http.NewRequest("POST", types.IPFSUrl+types.IPFSAdd, &b)
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to create new request for sending image to IPFS")
		// 	return
		// }
		// // Don't forget to set the content type, this will contain the boundary.
		// httpReq.Header.Set("Content-Type", writer.FormDataContentType())

		// client := &http.Client{}
		// resp, err := client.Do(httpReq)

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to execute request for sending image to IPFS")
		// 	return
		// }

		// defer resp.Body.Close()

		// result := setIPFSImage{}

		// // Collect the result in json form. Remember that we need to create a corresponding struct to do this
		// json.NewDecoder(resp.Body).Decode(&result)

		// // After collecting the hash image, we need to clear the image file stored temporary
		// err = os.Remove(tempFile.Name())
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to remove the temporary image file: %s", err.Error()))
		// 	return
		// }

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to push image onto IPFS: %s", err.Error()))
		// 	return
		// }

		imageHash := r.FormValue("image_hash")
		imageName := r.FormValue("image_name")

		req := setAIRequestHandlerFn(cliCtx, w, r)

		// Need to create a baseReq to write tx response. We cannot use baseReq in the AIRequest struct because AIRequest needs to be in form data to be able to send images
		baseReq := rest.BaseReq{
			From:          req.From,
			Memo:          req.Memo,
			ChainID:       req.ChainID,
			AccountNumber: req.AccountNumber,
			Sequence:      req.Sequence,
			Fees:          req.Fees,
			GasPrices:     req.GasPrices,
			Gas:           req.Gas,
			GasAdjustment: req.GasAdjustment,
			Simulate:      req.Simulate,
		}

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(req.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "AVXSD")
			return
		}

		// create the message
		msg := types.NewMsgSetClassificationRequest(imageHash, imageName, types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees.String(), req.ValidatorCount, req.Input, req.ExpectedOutput))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GHYK")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

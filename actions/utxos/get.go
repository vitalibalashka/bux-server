package utxos

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will fetch a given utxo according to conditions
// Get UTXO godoc
//	@Summary		Get UTXO
//	@Description	Get UTXO
//	@Tags			UTXO
//	@Produce		json
//	@Param			tx_id			query	string	true	"tx_id"
//	@Param			output_index	query	int		true	"output_index"
//	@Success		200
//	@Router			/v1/utxo [get]
//	@Security		bux-auth-xpub
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	txID := params.GetString("tx_id")
	outputIndex := uint32(params.GetUint64("output_index"))

	// Get a utxo using a xPub
	utxo, err := a.Services.Bux.GetUtxo(
		req.Context(),
		reqXPubID,
		txID,
		outputIndex,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(utxo))
}

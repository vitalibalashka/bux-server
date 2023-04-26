package transactions

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// count will fetch a count of transactions filtered on conditions and metadata
// Count of transactions godoc
//	@Summary		Count of transactions
//	@Description	Count of transactions
//	@Tags			Transactions
//	@Produce		json
//	@Param			metadata	query	string	false	"metadata"
//	@Param			conditions	query	string	false	"conditions"
//	@Success		200
//	@Router			/v1/transaction/count [post]
//	@Security		bux-auth-xpub
func (a *Action) count(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	_, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)a
	var count int64
	if count, err = a.Services.Bux.GetTransactionsByXpubIDCount(
		req.Context(),
		reqXPubID,
		metadata,
		conditions,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, count)
}

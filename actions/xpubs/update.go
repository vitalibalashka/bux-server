package xpubs

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// update will update an existing model
// Update xPub godoc
// @Summary     	Update xPub
// @Description 	Update xPub
// @Tags			xPub
// @Produce     	json
// @Param			metadata query string false "metadata"
// @Success     	200
// @Router      	/v1/xpub [patch]
// @Security 		bux-auth-xpub
func (a *Action) update(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPub, _ := bux.GetXpubFromRequest(req)
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	metadata := params.GetJSON(actions.MetadataField)

	// Get an xPub
	var xPub *bux.Xpub
	var err error
	xPub, err = a.Services.Bux.UpdateXpubMetadata(
		req.Context(), reqXPubID, metadata,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	signed := req.Context().Value("auth_signed")
	if signed == nil || !signed.(bool) || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(xPub))
}

package message

import (
	"encoding/json"
	"net/http"

	"github.com/xRebellion/MessageDistributor/recipient"
	"github.com/xRebellion/MessageDistributor/utils"
)

type Controller struct {
	service *Service
}

type sendMesssageParams struct {
	Subject    string              `json:"subject"`
	Message    string              `json:"message"`
	Recipients []recipient.Contact `json:"recipients"`
}

func (c *Controller) SendMessage(w http.ResponseWriter, r *http.Request) {
	var params *sendMesssageParams
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)

	errs := c.service.SendMessage(ctx, params.Subject, params.Message, params.Recipients)
	if len(errs) > 0 {
		utils.ResponseErrors(w, errs)
		return
	}

	utils.ResponseJSON(w, 200, nil)
}

func NewController(messageService *Service) *Controller {
	return &Controller{service: messageService}
}

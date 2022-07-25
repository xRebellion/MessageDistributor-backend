package recipient

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xRebellion/MessageDistributor/utils"
)

type Controller struct {
	service *Service
}

type createRecipientParams struct {
	Name        string    `json:"name"`
	ContactInfo []Contact `json:"contactInfo"`
}

func (c *Controller) CreateRecipient(w http.ResponseWriter, r *http.Request) {
	var params *createRecipientParams

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	recipient, err := c.service.CreateRecipient(ctx, params.Name, params.ContactInfo)
	if err != nil {
		utils.ResponseError(w, err)
		return
	}
	utils.ResponseJSON(w, 201, recipient)

}
func (c *Controller) GetRecipientByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	recipientID := chi.URLParam(r, "recipientID")

	recipient, err := c.service.GetRecipientByID(ctx, recipientID)
	if err != nil {
		utils.ResponseError(w, err)
		return
	}
	utils.ResponseJSON(w, 200, recipient)
}

type listRecipientsResponse struct {
	Data []Recipient `json:"data"`
}

func (c *Controller) ListRecipients(w http.ResponseWriter, r *http.Request) {
	var recipients listRecipientsResponse

	ctx := r.Context()
	res, err := c.service.ListRecipients(ctx)
	if err != nil {
		utils.ResponseError(w, err)
		return
	}
	recipients.Data = res

	utils.ResponseJSON(w, 200, recipients)
}

func NewController(recipientService *Service) *Controller {
	return &Controller{
		service: recipientService,
	}
}

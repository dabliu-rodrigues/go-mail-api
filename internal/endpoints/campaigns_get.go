package endpoints

import (
	"net/http"
)

func (h *Handler) ListCampaigns(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	//campaigns, err := h.CampaignService.Repository.List()
	return nil, http.StatusOK, nil
}

package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id := chi.URLParam(r, "id")
	campaign, err := h.CampaignService.GetByID(id)
	return campaign, http.StatusOK, err
}

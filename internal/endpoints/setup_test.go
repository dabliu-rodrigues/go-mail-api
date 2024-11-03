package endpoints

import (
	"bytes"
	"context"
	internalmock "emailn/internal/test/internal-mock"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

var (
	service *internalmock.CampaingServiceMock
	handler = Handler{}
)

func setup() {
	service = new(internalmock.CampaingServiceMock)
	handler.CampaignService = service
}

func newHttpTest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}

	req, _ := http.NewRequest(method, url, &buf)
	rr := httptest.NewRecorder()
	return req, rr
}

func addParameter(req *http.Request, key, value string) *http.Request {
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add(key, value)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiContext))
	return req
}

func addContext(req *http.Request, key, value string) *http.Request {
	ctx := context.WithValue(req.Context(), key, value)
	return req.WithContext(ctx)
}

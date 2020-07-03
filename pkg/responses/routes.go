package responses

import (
	"encoding/json"

	f "github.com/bmaynard/apimock/pkg/filesystem"
	l "github.com/bmaynard/apimock/pkg/utils/logger"

	m "github.com/bmaynard/apimock/pkg/mocks"
	"github.com/gorilla/mux"
)

type MockResponseList struct {
	Domain        string
	Path          string
	RequestMethod string
}

func BuildRoutes(r *mux.Router) {
	lfo := f.GetAdapter()

	response := make(map[MockResponseList][]*MockResponseItem)

	for _, mock := range lfo.GetMocks() {
		data := m.MockResponse{}
		err := json.Unmarshal([]byte(mock.Contents), &data)

		if err != nil || data.Response == nil || data.Meta.RequestMethod == "" || data.Meta.RequestPath == "" {
			l.Log.Errorf("Unable to parse file: %s", mock.FilePath)
			continue
		}

		var statusCode int

		if data.Meta.StatusCode >= 100 && data.Meta.StatusCode < 600 {
			statusCode = data.Meta.StatusCode
		} else {
			statusCode = 200
		}

		mockData := &MockResponseItem{
			Response:   data.Response,
			StatusCode: statusCode,
		}

		response[MockResponseList{mock.Domain, data.Meta.RequestPath, data.Meta.RequestMethod}] = append(response[MockResponseList{mock.Domain, data.Meta.RequestPath, data.Meta.RequestMethod}], mockData)
	}

	for record, mockResponses := range response {
		if record.Domain == "_all_" {
			r.NewRoute().Path(record.Path).Handler(Handler{MockResponses: mockResponses, H: JsonResponse}).Methods(record.RequestMethod)
			continue
		}

		s := r.Host(record.Domain).Subrouter()
		s.NewRoute().Path(record.Path).Handler(Handler{MockResponses: mockResponses, H: JsonResponse}).Methods(record.RequestMethod)
	}
}

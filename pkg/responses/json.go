package responses

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func JsonResponse(mockResponse []*MockResponseItem, w http.ResponseWriter, r *http.Request) error {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(mockResponse) // pick a random mock

	fmt.Printf("Endpoint Hit: %s\n", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mockResponse[n].StatusCode)
	json.NewEncoder(w).Encode(mockResponse[n].Response)
	return nil
}

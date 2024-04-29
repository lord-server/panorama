package web

import "net/http"

type GetMetadataResponse struct {
}

func (s *server) GetMetadata(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, map[string]any{
		"center": []int{700, 200},
	})
}

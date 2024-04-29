package web

import "net/http"

type View struct {
	Type     string `json:"type"`
	Rotation string `json:"rotation"`
}

type GetViewsResponse struct {
	Views map[string]View `json:"views"`
}

func (s *server) GetViews(w http.ResponseWriter, r *http.Request) {
	response := GetViewsResponse{
		Views: make(map[string]View),
	}

	for name, view := range s.config.Views {
		response.Views[name] = View{
			Type:     view.Type,
			Rotation: view.Rotation,
		}
	}

	sendJSON(w, response)
}

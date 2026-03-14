package pushtrigger

import "encoding/json"

type prView struct {
	Number      int    `json:"number"`
	URL         string `json:"url"`
	BaseRefName string `json:"baseRefName"`
	State       string `json:"state"`
}

func parsePRView(s string) (*prView, error) {
	var view prView
	if err := json.Unmarshal([]byte(s), &view); err != nil {
		return nil, err
	}
	return &view, nil
}

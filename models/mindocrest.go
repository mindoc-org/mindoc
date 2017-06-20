package models

// MinDocRest ...
type MinDocRest struct {
	Token    string `json:"token"`
	Folder   string `json:"folder"`
	Title    string `json:"title"`
	Identify string `json:"identify"`
	TextMD   string `json:"textmd"`
	TextHTML string `json:"texthtml"`
}

// NewMinDocRest ...
func NewMinDocRest() *MinDocRest {
	return &MinDocRest{}
}

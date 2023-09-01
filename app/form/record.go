package form

type RecordForm struct {
	Account     string `json:"account"`
	GroupName   string `json:"group_name"`
	AccessToken string `json:"access_token"`
}

type CsvFileForm struct {
	GroupName string `json:"group_name"`
}

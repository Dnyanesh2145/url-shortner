package helpers

import "time"

type ModelURL struct {
	Id        string    `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
}

type LogModel struct {
	Userid        string    `json:"user_id" `
	Requestid     string    `json:"request_id" `
	DeviceInfo    string    `json:"deviceinfo" `
	Endpoint      string    `json:"apiendpoint"  `
	IP            string    `json:"ip_address" `
	Location      string    `json:"geolocation" `
	Method        string    `json:"httpmethod" `
	Referrer      string    `json:"referrer" `
	ResponseSize  int       `json:"responsesize" `
	ResponseTime  int       `json:"responsetime" `
	StatusCode    int       `json:"statuscode" `
	StatusMessage string    `json:"statusmessage"`
	UserAgent     string    `json:"useragent" `
	CreatedAt     time.Time `json:"created_at" `
	//UpdatedAt     time.Time `json:"updated_at"`
}

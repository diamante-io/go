package health

// Response implements the most basic required fields for the health response
// based on the format defined in the draft IETF network working group
// standard, Health Check Response Format for HTTP APIs.
//
// https://tools.ietf.org/id/draft-inadarei-api-health-check-01.html
type Response struct {
	Status Status `json:"status"`
}

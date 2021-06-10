package alert

type CreateTeamRequest struct {
	Team       Team         `json:"team"`
	Developers []Developers `json:"developers"`
}

type Team struct {
	Name string `json:"name"`
}

type Developers struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

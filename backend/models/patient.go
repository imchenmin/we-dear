package models

type Patient struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Age       int       `json:"age"`
	Phone     string    `json:"phone"`
	Diagnosis string    `json:"diagnosis"`
	Doctor    string    `json:"doctor"`
	Avatar    string    `json:"avatar"`
	Messages  []Message `json:"messages"`
}

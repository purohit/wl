package wundergo

type User struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Revision   uint   `json:"revision"`
	TypeString string `json:"type"`
}

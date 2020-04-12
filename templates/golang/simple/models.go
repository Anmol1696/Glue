package main

// Identity struct used for user identity
type UserIdentity struct {
	Subject string   `json:"subject,omitempty"`
	Email   string   `json:"email,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}

func (ui *UserIdentity) HasRole(role string) bool {
	for _, userRole := range ui.Roles {
		if role == userRole {
			return true
		}
	}
	return false
}

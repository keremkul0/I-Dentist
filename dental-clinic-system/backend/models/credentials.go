package models

// Credentials represents the structure for user login credentials
type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

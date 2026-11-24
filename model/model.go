package model

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var Users = []User{
	{ID: 1, Name: "Rahim", Email: "rahim@example.com", Age: 25},
	{ID: 2, Name: "Karim", Email: "karim@example.com", Age: 24},
	{ID: 3, Name: "Kobirul", Email: "kobirul@example.com", Age: 26},
}

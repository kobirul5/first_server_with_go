package main

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "Rahim", Email: "rahim@example.com"},
	{ID: 2, Name: "Karim", Email: "karim@example.com"},
	{ID: 3, Name: "Kobirul", Email: "kobirul@example.com"},
}

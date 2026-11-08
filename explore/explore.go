package main

var users = []string{"kobirul", "rakib", "habib", "mahi", "nabil", "rabiul"}

func deleteUser(id int) {

	idx:= 2

	users = append(users[:idx], users[idx+1:]...)


	
}


// user[:2] => ["kobirul", "rakib"]
//users[2+1:] => ["mahi", "nabil", "rabiul"]
//append(["kobirul","rakib"] ,["mahi", "nabil", "rabiul"]) => ["kobirul","rakib", ["mahi", "nabil", "rabiul"]]
//append(["kobirul","rakib"], users[2+1:]...) => ["kobirul","rakib", "mahi", "nabil", "rabiul"]

package main

type Item struct {
	Name string `json:"name"`
}

type Collection struct {
	Items []Item `ginja:"hasMany",json:"items"`
}

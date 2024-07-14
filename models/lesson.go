package models

type Lesson struct {
	Id       int    `json:"id"`
	Code     string `json:"code"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Video    string `json:"video"`
	Duration string `json:"duration"`
}

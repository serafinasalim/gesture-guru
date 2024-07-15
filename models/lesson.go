package models

type Lesson struct {
	Id       int    `json:"id"`
	Code     string `json:"code"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Video    string `json:"video"`
	Duration string `json:"duration"`
	Saved    bool   `json:"saved"`
	Status   string `json:"status"`
}

type LessonBrowse struct {
	UserId   int `json:"userId" validate:"required"`
	LessonId int `json:"lessonId"`
}

type LessonAttempt struct {
	AttemptNumber int     `json:"attemptNumber"`
	Score         float64 `json:"score" validate:"required,min=1,max=100"`
}

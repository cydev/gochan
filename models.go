package gochan

import (
	"fmt"
	"time"
)

type Post struct {
	Id    uint
	Date  time.Time
	Text  string
	User  string
	Media Media
}

type Media interface {
	Render() string
}

type MediaYouTube struct {
	Id string
}

func (media MediaYouTube) Render() string {
	return fmt.Sprintf("<iframe type=\"text/html\" src=\"%s\" frameborder=\"0\"/>", media.Id)
}

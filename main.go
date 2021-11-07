package main

import (
	"fmt"

	movie "github.com/oohyun15/naver-movie-crawler/model"
)

func main() {
	m := movie.New("57723", "타짜")
	fmt.Println(m)
}

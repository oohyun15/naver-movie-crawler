package main

import (
	"naver-movie-crawler/models/movie"
)

func main() {
	batchSize := 100
	movie := movie.New("57723", "타짜")
	movie.Scrape(batchSize)
}

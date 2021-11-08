package main

import (
	"fmt"
	"naver-movie-crawler/models/movie"
	"time"
)

func main() {
	startTime := time.Now()
	fmt.Println("start:", startTime)

	batchSize := 200
	fmt.Println("batch size: ", batchSize)

	movies := []movie.Movie{
		movie.New("57723", "타짜"),
		movie.New("100931", "겨울왕국"),
		movie.New("167651", "극한직업"),
		movie.New("161967", "기생충"),
		movie.New("184517", "소울"),
	}

	for _, m := range movies {
		m.Scrape(batchSize)
	}

	endTime := time.Now()
	fmt.Println("end:", endTime)
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("elapsed:", elapsedTime)
}

package movie

import (
	"fmt"
	"naver-movie-crawler/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Movie struct {
	identifier string
	title      string
	size       int
	page       int
	score      map[int]int
}

type Review struct {
	name        string
	score       int
	date        string
	description string
}

func New(identifier string, title string) Movie {
	size, page, score := initialize(identifier)
	movie := Movie{identifier, title, size, page, score}
	fmt.Println(fmt.Sprintf("title: %v, identifier: %v, size: %v, page: %v", title, identifier, size, page))
	return movie
}

func initialize(identifier string) (int, int, map[int]int) {
	score := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0}

	pageURL := "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=" + identifier + "&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=lowest"
	res, err := http.Get(pageURL)
	utils.CheckErr(err)
	utils.CheckCode(res, pageURL)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	utils.CheckErr(err)
	size, _ := strconv.Atoi(strings.Replace(doc.Find("strong.total em").Text(), ",", "", -1))
	sizePerPage := doc.Find("div.score_result li").Size()
	page := (size / sizePerPage) + 1

	return size, page, score
}

func (movie *Movie) Scrape(batchSize int) {
	startTime := time.Now()
	fmt.Println("start:", startTime)
	fmt.Print("batch size: ", batchSize)

	var reviews []Review
	var baseURL string = "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=" + movie.identifier + "&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=lowest"
	pages := movie.page
	c := make(chan []Review)

	for idx := 0; idx < pages/batchSize+1; idx++ {
		start := idx*batchSize + 1
		end := (idx + 1) * batchSize
		if end > pages {
			end = pages
		}

		for i := start; i <= end; i++ {
			go movie.getPage(i, baseURL, c)
		}

		for i := start; i < end; i++ {
			extractedReviews := <-c
			reviews = append(reviews, extractedReviews...)
		}
	}

	fmt.Println(reviews)
	// writeReviews(reviews)

	fmt.Println("Done, extracted")
	endTime := time.Now()
	fmt.Println("end: ", endTime)
}

func (movie *Movie) getPage(page int, url string, mainC chan<- []Review) {
	var reviews []Review
	c := make(chan Review)
	pageURL := url + "&page=" + strconv.Itoa(page)
	// fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	utils.CheckErr(err)
	utils.CheckCode(res, pageURL)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	utils.CheckErr(err)
	reviewLists := doc.Find("div.score_result li")
	reviewLists.Each(func(i int, list *goquery.Selection) {
		go movie.extractReview(list, i, c)
	})

	for i := 0; i < reviewLists.Length(); i++ {
		review := <-c
		reviews = append(reviews, review)
	}
	mainC <- reviews
}

func (movie *Movie) extractReview(list *goquery.Selection, num int, c chan<- Review) {
	name := list.Find("div.score_reple dl dt em a span").Text()
	score, _ := strconv.Atoi(list.Find("div.star_score em").Text())
	description := utils.CleanString(list.Find("div.score_reple p span#_filtered_ment_" + strconv.Itoa(num)).Text())
	date := utils.CleanString(list.Find("div.score_reple dl dt em").Last().Text())
	// movie.score[score] += 1
	c <- Review{
		name:        name,
		score:       score,
		date:        date,
		description: description,
	}
}

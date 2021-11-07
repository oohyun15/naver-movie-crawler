package movie

import (
	"encoding/csv"
	"fmt"
	"naver-movie-crawler/utils"
	"net/http"
	"os"
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
	// mutex      *sync.Mutex
}

type Review struct {
	name        string
	score       int
	description string
	date        string
}

func New(identifier string, title string) Movie {
	size, page, score := initialize(identifier)
	movie := Movie{identifier, title, size, page, score}
	// mutex := &sync.Mutex{}
	// movie := Movie{identifier, title, size, page, score, mutex}
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
	fmt.Println("Scraping " + movie.title + "...")
	var reviews []Review
	var baseURL string = "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=" + movie.identifier
	pages := movie.page
	c := make(chan []Review)

	for idx := 0; idx < pages/batchSize+1; idx++ {
		start := idx*batchSize + 1
		end := (idx + 1) * batchSize
		if end > pages {
			end = pages
		}

		fmt.Println("start: ", start, "end:", end, "reviews:", len(reviews))
		for i := start; i <= end; i++ {
			go movie.getPage(i, baseURL, c)
		}

		for i := start; i <= end; i++ {
			extractedReviews := <-c
			reviews = append(reviews, extractedReviews...)
		}
	}

	movie.setScore(reviews)
	fmt.Println("Result:", movie.score)
	movie.writeReviews(reviews)
}

func (movie *Movie) getPage(page int, url string, mainC chan<- []Review) {
	var reviews []Review
	c := make(chan Review)
	pageURL := url + "&page=" + strconv.Itoa(page)
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

	// movie.mutex.Lock()
	// movie.score[score] += 1
	// movie.mutex.Unlock()

	c <- Review{
		name:        name,
		score:       score,
		description: description,
		date:        date,
	}
}

func (movie *Movie) setScore(reviews []Review) {
	for _, review := range reviews {
		movie.score[review.score] += 1
	}
}

func (movie *Movie) writeReviews(reviews []Review) {
	file, err := os.Create(movie.title + "_reviews(" + time.Now().Format("2006-01-02") + ").csv")
	utils.CheckErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"name", "score", "description", "date"}
	wErr := w.Write(headers)
	utils.CheckErr(wErr)
	c := make(chan error)

	for _, review := range reviews {
		go movie.writeReview(review, w, c)
		utils.CheckErr(<-c)
	}
}

func (movie *Movie) writeReview(review Review, w *csv.Writer, c chan<- error) {
	row := []string{review.name, strconv.Itoa(review.score), review.description, review.date}
	err := w.Write(row)
	c <- err
}

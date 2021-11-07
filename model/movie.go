package movie

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	errhd "github.com/oohyun15/naver-movie-crawler/helper"
)

type Movie struct {
	identifier string
	title      string
	size       int
	page       int
	score      map[int]int
}

func New(identifier string, title string) Movie {
	size, page, score := initialize(identifier)
	movie := Movie{identifier, title, size, page, score}
	return movie
}

func initialize(identifier string) (int, int, map[int]int) {
	score := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0}

	pageURL := "https://movie.naver.com/movie/bi/mi/pointWriteFormList.nhn?code=" + identifier + "&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=lowest"
	res, err := http.Get(pageURL)
	errhd.CheckErr(err)
	errhd.CheckCode(res, pageURL)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	errhd.CheckErr(err)
	size, _ := strconv.Atoi(strings.Replace(doc.Find("strong.total em").Text(), ",", "", -1))
	sizePerPage := doc.Find("div.score_result li").Size()
	page := (size / sizePerPage) + 1

	return size, page, score
}

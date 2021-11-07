package movie

type Movie struct {
	identifier string
	title      string
	size       int
	page       int
	score      map[int]int
}

func New(identifier string, title string) Movie {
	size := 0
	page := getTotalPage(identifier)
	score := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0}

	movie := Movie{identifier, title, size, page, score}
	return movie
}

func getTotalPage(identifier string) int {
	return 0
}

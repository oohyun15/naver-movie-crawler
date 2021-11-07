package utils

import (
	"fmt"
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err := recover(); err != nil {
		panic(err)
	}
}

func CheckCode(res *http.Response, url string) {
	defer func() {
		if c := recover(); c != nil {
			fmt.Println("recover", url)
		}
	}()
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

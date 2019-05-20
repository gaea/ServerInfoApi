package helpers

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AnalizeScraping(host string) (string, string) {
	title := ""
	logo := ""

	responseScrap, err := http.Get("http://" + host)

	if err != nil {
		log.Println(err)
	} else {
		defer responseScrap.Body.Close()
		document, err := goquery.NewDocumentFromReader(responseScrap.Body)

		if err != nil {
			log.Println(err)
		} else {
			document.Find("head > title").Each(func(index int, element *goquery.Selection) {
				title = element.Text()
			})

			document.Find("head > link").EachWithBreak(func(index int, element *goquery.Selection) bool {
				imgLogo, exists := element.Attr("href")

				if exists {
					imgPath := strings.ToLower(imgLogo)

					if strings.HasSuffix(imgPath, "png") || strings.HasSuffix(imgPath, "jpeg") || strings.HasSuffix(imgPath, "jpg") {
						logo = imgLogo

						return false
					}
				}

				return true
			})
		}
	}

	return title, logo
}

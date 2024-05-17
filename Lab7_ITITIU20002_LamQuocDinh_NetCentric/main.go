package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"golang.org/x/net/html"
)

type Manga struct {
	genre string
	link  string
}

type MangaJSON struct {
	Genre  string `json:"Genre"`
    Titles []string `json:"Titles"`
}

func getLinks(n *html.Node, className string) []Manga {
	var mangas []Manga
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "ul" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == className {
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && c.Data == "li" {
							var manga Manga
							for _, attr := range c.Attr {
								if attr.Key == "data-genre" {
									manga.genre = attr.Val
								}
							}
							for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
								if cc.Type == html.ElementNode && cc.Data == "a" {
									for _, aa := range cc.Attr {
										if aa.Key == "href" {
											manga.link = aa.Val
										}
									}
								}
							}
							mangas = append(mangas, manga)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return mangas
}

func getAllTitles(n *html.Node, className string) []string {
	var titles []string
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "p" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == className {
					if n.FirstChild != nil {
						titles = append(titles, n.FirstChild.Data)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result := f(c)
			if result != "" {
				return result
			}
		}
		return ""
	}
	f(n)
	return titles
}

func getDoc(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func main() {
	doc, err := getDoc("https://www.webtoons.com/en/genres/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mangas := getLinks(doc, "snb _genre")
	mangas = mangas[:len(mangas)-1]

	var mangasJSON []MangaJSON
	for _, manga := range mangas {
		doc, err := getDoc(manga.link)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		titles := getAllTitles(doc, "subj")
		mangasJSON = append(mangasJSON, MangaJSON{Genre: manga.genre, Titles: titles})
	}

	jsonData, err := json.Marshal(mangasJSON)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile("mangas.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		os.Exit(1)
	}

}

package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/RealJK/rss-parser-go"
)

type rssO struct {
	Header Heading
	Items  []Item
}

type Heading struct {
	Title         string
	TitleLen      bool
	Generator     string
	LastBuildDate string
	Description   string
	LenItems      int
}

type Item struct {
	Number      int
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	TitleLen    bool
	Link        string `xml:"link"`
	LinkLen     bool
	Description string `xml:"description"`
	Author      string `xml:"author,omitempty"`
	Category    string `xml:"category,omitempty"`
	CategotyLen bool
	Comments    string `xml:"comments,omitempty"`
	Guid        Guid
	PubDate     string `xml:"pubDate"`
	PubDateLen  bool
	Source      Source
}

type Source struct {
	XMLName xml.Name `xml:"source"`
	Url     string   `xml:"url,attr"`
	Value   string   `xml:",innerxml"`
}

type Guid struct {
	XMLName   xml.Name `xml:"guid"`
	PermaLink string   `xml:"isPermaLink,attr"`
	Value     string   `xml:",innerxml"`
}

type News []Item

func (this News) Len() int {
	return len(this)
}

func (this News) Less(i, j int) bool {
	return this[i].Number > this[j].Number
}
func (this News) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func makeItem(items []Item, rssObejct rss.Channel) []Item {
	for _, item := range rssObejct.Items {
		date, err := item.PubDate.Parse()
		if err != nil {
			date, err = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", string(item.PubDate))
			if err != nil {
				date, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", string(item.PubDate))
				if err != nil {
					date, err = time.Parse("Mon, 02 Jan 2006 15:04:05 MST", string(item.PubDate))
					if err != nil {
						date, err = time.Parse("02 Jan 2006 15:04:05 GMT", string(item.PubDate))
						if err != nil {
							log.Fatal("Date: ", err)
						}
					}
				}
			}
		}
		number := date.Year()*10000000000 + int(date.Month())*100000000 + date.Day()*1000000 + date.Hour()*100*100 + date.Minute()*100 + date.Second()
		item0 := Item{number, item.XMLName, item.Title, len(item.Title) != 0, item.Link,
			len(item.Link) != 0, item.Description, item.Author, item.Category.Value,
			len(item.Category.Value) != 0, item.Comments, Guid{item.Guid.XMLName, item.Guid.PermaLink,
				item.Guid.Value}, string(item.PubDate),
			len(item.PubDate) != 0, Source{item.Source.XMLName, item.Source.Url, item.Source.Value}}
		items = append(items, item0)
	}
	return items
}

func makeOneNews(w http.ResponseWriter, rssObject rss.Channel) {
	head := Heading{rssObject.Title, len(rssObject.Title) != 0, rssObject.Generator,
		rssObject.LastBuildDate, rssObject.Description, len(rssObject.Items)}

	items := make([]Item, 0, 1)
	items = makeItem(items, rssObject)

	tmpl, err := template.ParseFiles("oneNews.html", "header.html", "footer.html")
	if err != nil {
		log.Fatal("Error news.html :", err)
	}

	if err := tmpl.ExecuteTemplate(w, "oneNews", rssO{head, items}); err != nil {
		log.Fatal("Error news.html :", err)
	}

}

func makeNews(w http.ResponseWriter, r *http.Request) {
	rssObject, err := rss.ParseRSS("https://rg.ru/xml/index.xml")
	rssObject1, err1 := rss.ParseRSS("https://lenta.ru/rss")
	rssObject2, err2 := rss.ParseRSS("https://news.mail.ru/rss/90/")

	if err != nil && err1 != nil && err2 != nil {

		items := make([]Item, 0, 1)
		items = makeItem(items, rssObject.Channel)
		items = makeItem(items, rssObject1.Channel)
		items = makeItem(items, rssObject2.Channel)

		sort.Sort(News(items))

		tmpl, err := template.ParseFiles("news.html", "header.html", "footer.html")
		if err != nil {
			log.Fatal("Error news.html :", err)
		}

		if err := tmpl.ExecuteTemplate(w, "news", items); err != nil {
			log.Fatal("Error news.html :", err)
		}

	}
}

func makeNewsLenta(w http.ResponseWriter, r *http.Request) {
	rssObject, err := rss.ParseRSS("https://lenta.ru/rss")

	if err != nil {
		makeOneNews(w, rssObject.Channel)

	}
}

func makeNewsRG(w http.ResponseWriter, r *http.Request) {
	rssObject, err := rss.ParseRSS("https://rg.ru/xml/index.xml")

	if err != nil {
		makeOneNews(w, rssObject.Channel)
	}
}

func makeNewsMail(w http.ResponseWriter, r *http.Request) {
	rssObject, err := rss.ParseRSS("https://news.mail.ru/rss/90/")

	if err != nil {
		makeOneNews(w, rssObject.Channel)
	}
}

func makeNewsServer(w http.ResponseWriter, r *http.Request) {
	rssObject, err := rss.ParseRSS("https://news.mail.ru/rss/90/")

	if err != nil {
		items := make([]Item, 0, 1)
		items = makeItem(items, rssObject.Channel)
		if len(rssObject.Channel.Title) != 0 {
			fmt.Printf("Title : %s\n", rssObject.Channel.Title)
		}
		if len(rssObject.Channel.Generator) != 0 {
			fmt.Printf("Generator : %s\n", rssObject.Channel.Generator)
		}
		if len(rssObject.Channel.LastBuildDate) != 0 {
			fmt.Printf("LastBuildDate : %s\n", rssObject.Channel.LastBuildDate)
		}
		if len(rssObject.Channel.Description) != 0 {
			fmt.Printf("Description : %s\n", rssObject.Channel.Description)
		}
		fmt.Printf("Number of Items : %d\n", len(rssObject.Channel.Items))

		for v, item := range rssObject.Channel.Items {
			fmt.Println()
			fmt.Printf("Item Number : %d\n", v)
			if len(item.Title) != 0 {
				fmt.Printf("Title : %s\n", item.Title)
			}
			if len(item.Link) != 0 {
				fmt.Printf("Link : %s\n", item.Link)
			}
			if len(item.Description) != 0 {
				fmt.Printf("Description : %s\n", item.Description)
			}
			if len(item.Guid.Value) != 0 {
				fmt.Printf("Guid : %s\n", item.Guid.Value)
			}
		}
	}
}

func main() {
	println("Server listen on port 9005")
	http.HandleFunc("/news", makeNews)
	http.HandleFunc("/newsLenta", makeNewsLenta)
	http.HandleFunc("/newsRG", makeNewsRG)
	http.HandleFunc("/newsMail", makeNewsMail)
	http.HandleFunc("/server", makeNewsServer)
	err := http.ListenAndServe(":9005", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

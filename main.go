package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	//create your file with desired read/write permissions
	f, err := os.OpenFile("dalog.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	//test case
	log.Println("check to make sure it works")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "The id query parameter is missing", http.StatusBadRequest)
			return
		}

		title := getTitle(url)
		fmt.Println(title)
		log.Println(title + " " + url)
		fmt.Fprintf(w, "Downloading, %q", title)
		go runCommand(url)

	})

	log.Fatal(http.ListenAndServe(":17945", nil))

}

func runCommand(url string) {
	cmd := exec.Command("/usr/local/bin/yt-dlp", "-f", "bv*+ba/b", "-o", "vids/%(title)s-%(id)s.%(ext)s", url)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

// func formatTitle(title string) string {
// title = strings.ReplaceAll(title, " ", "")
// return title
// }

func getTitle(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	title := doc.Find("title").Text()
	return title
}

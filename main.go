package main

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/PuerkitoBio/goquery"
)

var downloadsPathName = "downloads"

func main() {
	fmt.Println("Main running")

	initLog()

	// createPath(downloadsPathName)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		audioOnly := r.URL.Query().Get("audioonly")
		downloadAudio := false

		if url == "" {
			http.Error(w, "The URL query parameter is missing", http.StatusBadRequest)
			return
		}

		if audioOnly != "" && audioOnly == "true" {
			downloadAudio = true
		}

		title, err := getTitle(url)

		if err != nil {
			fmt.Println("ERROR:", err)
			fmt.Fprintf(w, "ERROR, %q", err)
			// log.Fatal(err)
		} else {
			fmt.Println(title)
			log.Println(title + " " + url)
			fmt.Fprintf(w, "Downloading, %q", title)

			go runDownload(url, downloadAudio)
		}

	})

	log.Println(http.ListenAndServe(":17945", nil))

}

func runDownload(url string, downloadAudio bool) {
	fmt.Println("Download Merged Video And Audio")
	downloadsPath := getPwd() + "/" + downloadsPathName

	var cmd *exec.Cmd
	if downloadAudio {
		filePath := downloadsPath + "/%(title)s-%(id)s-audio.%(ext)s"
		cmd = exec.Command("/usr/local/bin/yt-dlp", "-f", "ba", "-x", "--audio-format", "mp3", "-o", filePath, url)
	} else {
		filePath := downloadsPath + "/%(title)s-%(id)s.%(ext)s"
		// cmd = exec.Command("/usr/local/bin/yt-dlp", "-f", "bv+ba/b", "-o", filePath, url)
		cmd = exec.Command("/usr/local/bin/yt-dlp", "-f", "bv+mergeall[vcodec=none]", "-o", filePath, url)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
	// if err := cmd.Run(); err != nil {
	// 	fmt.Println("ERROR!", err)
	// 	log.Println(err)
	// }
}

func getTitle(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
		err := errors.New("Couldn't get Title from " + url)
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	// if err != nil {
	// log.Fatal(err)
	// }

	title := doc.Find("title").Text()
	return title, err
}

func createPath(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func getPwd() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("Current path:", path)
	return path
}

func initLog() {
	//create your file with desired read/write permissions
	f, err := os.OpenFile("downloads.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	//test case
	log.Println("check to make sure it works")
}

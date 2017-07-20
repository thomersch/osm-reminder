package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var url = "https://www.openstreetmap.org/api/0.6/changesets?time=%v&closed=true"

type osmChangesets struct {
	Changeset []struct {
		ID   string `xml:"id,attr"`
		User string `xml:"user,attr"`
		Tags []struct {
			Key   string `xml:"k,attr"`
			Value string `xml:"v,attr"`
		} `xml:"tag"`
	} `xml:"changeset"`
}

type changeset struct {
	ID       int
	Username string
	Comment  string
}

type notification struct {
}

func main() {
	runEngine(context.Background(), time.Minute)
}

func runEngine(ctx context.Context, crawlingFrequency time.Duration) {
	var eventCh = make(chan notification, 100)
	go func() {
		for {
			select {
			case noti := <-eventCh:
				log.Println(noti)
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			chs, err := downloadChangesets()
			if err != nil {
				log.Println(err)
				time.Sleep(crawlingFrequency)
				continue
			}
			for range relevantChangesets(chs) {
				eventCh <- notification{}
			}
		}
		time.Sleep(crawlingFrequency)
	}
}

func downloadChangesets() ([]changeset, error) {
	resp, err := http.Get(fmt.Sprintf(url, time.Now().Format(time.RFC3339)))
	if err != nil {
		return nil, err
	}

	chs := osmChangesets{}
	err = xml.NewDecoder(resp.Body).Decode(&chs)
	if err != nil {
		return nil, err
	}

	var changesets []changeset
	for _, ch := range chs.Changeset {
		var c changeset
		for _, tag := range ch.Tags {
			if tag.Key == "comment" {
				c.Comment = tag.Value
			}
		}
		c.ID, err = strconv.Atoi(ch.ID)
		if err != nil {
			return nil, err
		}
		c.Username = ch.User
		changesets = append(changesets, c)
	}

	return changesets, nil
}

func relevantChangesets(chs []changeset) []changeset {
	var rel []changeset
	for _, ch := range chs {
		if strings.Contains(ch.Comment, "!remindme") {
			rel = append(rel, ch)
		}
	}
	return rel
}

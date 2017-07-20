package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReminder(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<osm>
			<changeset id="123123" user="foobar">
				<tag k="comment" v="nobody cares" />
			</changeset>
			<changeset id="123123" user="foobar">
				<tag k="comment" v="i care !remindme"/>
			</changeset>
		</osm>`))
	}))
	defer srv.Close()
	url = srv.URL

	ctx, cancel := context.WithCancel(context.Background())
	go runEngine(ctx, 100*time.Millisecond)
	time.Sleep(time.Second)
	cancel()
}

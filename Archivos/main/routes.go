package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	//"github.com/tomiok/course-phones-review/gadgets/smartphones/web"
	//reviews "github.com/tomiok/course-phones-review/reviews/web"
)

func Routes() *chi.Mux {
	mux := chi.NewMux()
	mux.Use(
		middleware.Logger,    //log every http request
		middleware.Recoverer, // recover if a panic occurs
	)

	mux.Get("/getEmails", getEmailHandler)

	return mux
}

func getEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//q := r.URL.Query()

	//from := q.Get("id")
	query := `{
		        "search_type": "alldocuments",
				"from": 5000,
				"max_results": 10,
		        "_source": []

		    }`
	req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "0208Mavl")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(body)
	}
	fmt.Println(string(body))
	//res := map[string]interface{}{"index": string(body)}
	_ = json.NewEncoder(w).Encode(string(body))
}

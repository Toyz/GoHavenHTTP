package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Toyz/GoHaven"
	"github.com/julienschmidt/httprouter"
)

var (
	Haven *GoHaven.WallHaven
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var options []GoHaven.Option

	sorting := r.URL.Query().Get("sort")
	if len(sorting) <= 0 {
		options = append(options, GoHaven.SortRelevance)
	} else {
		var p GoHaven.Sorting
		p.Set(sorting)

		options = append(options, p)
	}

	purity := r.URL.Query().Get("purity")
	if len(purity) <= 0 {
		options = append(options, GoHaven.PuritySFW)
	} else {
		var p GoHaven.Purity
		for _, pure := range strings.Split(purity, ",") {
			p.Set(pure)
		}
		options = append(options, p)
	}

	cata := r.URL.Query().Get("category")
	if len(cata) <= 0 {
		options = append(options, GoHaven.CatGeneral)
	} else {
		var p GoHaven.Categories
		for _, cat := range strings.Split(cata, ",") {
			p.Set(cat)
		}

		options = append(options, p)
	}

	order := r.URL.Query().Get("order")
	if len(order) <= 0 {
		options = append(options, GoHaven.OrderDesc)
	} else {
		var p GoHaven.Order
		p.Set(order)

		options = append(options, p)
	}

	page := r.URL.Query().Get("page")
	if len(page) > 0 {
		var p GoHaven.Page
		p.Set(page)

		options = append(options, p)
	}

	ratios := r.URL.Query().Get("ratio")
	if len(ratios) > 0 {
		var p GoHaven.Ratios

		for _, ratio := range strings.Split(ratios, ",") {
			p.Set(ratio)
		}

		options = append(options, p)
	}

	search := r.URL.Query().Get("q")

	var data []byte
	results, err := Haven.Search(search, options...)
	if err != nil {
		w.WriteHeader(500)
		data, _ = json.Marshal(err)
	}

	data, _ = json.Marshal(results)
	w.Write(data)
}

func Info(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var data []byte

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(500)
		data, _ = json.Marshal(err)
	}
	imageID := GoHaven.ID(id)
	details, _ := imageID.Details()
	data, _ = json.Marshal(details)

	w.Write(data)
}

func main() {
	Haven = GoHaven.New()

	router := httprouter.New()
	router.GET("/search", Index)
	router.GET("/detail/:id", Info)

	log.Fatal(http.ListenAndServe(getEnv("LISTEN", ":8080"), router))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

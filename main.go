package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

type Point struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	Distance int `json:"distance"`
}

func (p1 *Point) SetDistance(p2 *Point) {
	p1.Distance = int(math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))) // Manhattan distance
}

type SortablePoints []Point

func (i SortablePoints) InRange(origin *Point, maxDistance int) SortablePoints {
	log.Println("Searching points within", maxDistance, "units of (", origin.X, origin.Y, ")")
	filteredPoints := make(SortablePoints, 0)
	for _, p := range points {
		p.SetDistance(origin)
		if p.Distance <= maxDistance {
			filteredPoints = append(filteredPoints, p)
		}
	}
	log.Println("Found", len(filteredPoints), "points matching criteria")
	return filteredPoints
}

func (i SortablePoints) Len() int {
	return len(i)
}

func (i SortablePoints) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i SortablePoints) Less(a, b int) bool {
	return i[a].Distance < i[b].Distance
}

var points SortablePoints

func init() {
	f, err := os.Open("data/points.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(b, &points)
	log.Println("Loaded", len(points), "points from 'data/points.json'")
}

func main() {
	http.HandleFunc("/api/points", GetPoints)
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	qs := r.URL.Query()
	param := make(map[string]int)
	for _, name := range []string{"x", "y", "distance"} {
		val, err := ValidateParameter(qs, name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		param[name] = val
	}

	filteredPoints := points.InRange(&Point{param["x"], param["y"], 0}, param["distance"])
	sort.Sort(filteredPoints)

	b, err := json.MarshalIndent(filteredPoints, "", "  ")
	if err != nil {
		log.Println("Error:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Write(b)
}

func ValidateParameter(qs url.Values, name string) (int, error) {
	s := qs.Get(name)
	if s == "" {
		return 0, errors.New(fmt.Sprintf("Missing required parameter: %s", name))
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Must be a valid integer: %s", name))
	}
	return i, nil
}

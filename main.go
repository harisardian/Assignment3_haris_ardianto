package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/go-co-op/gocron"
)

type Data struct {
	Status StsData
}

type StsData struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := gocron.NewScheduler(time.UTC)

		s.Every(15).Seconds().Do(func() { createFile() })

		water, wind := readFile()
		statusWater, statusWind := checkStatus(water, wind)

		var filepath = path.Join("views", "index.html")
		var tmpl, err = template.ParseFiles(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data = map[string]interface{}{
			"title":       "Assignment 3",
			"water":       water,
			"wind":        wind,
			"statusWater": statusWater,
			"statusWind":  statusWind,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		s.StartAsync()
	})

	http.HandleFunc("/tes", func(w http.ResponseWriter, r *http.Request) {
		s := gocron.NewScheduler(time.UTC)

		s.Every(5).Seconds().Do(func() { fmt.Println("cron") })
		s.StartAsync()
	})

	http.ListenAndServe(":8080", nil)
}

func createFile() {
	rand.Seed(time.Now().Unix())
	water := 0 + rand.Intn(100-0)
	wind := 0 + rand.Intn(100-0)
	file, err := os.Create("file.json")
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	a := fmt.Sprintf(`{
	"status": {
		"water": %d,
		"wind": %d
	}
}`, water, wind)
	file.WriteString(a)

	file.Sync()
}

func readFile() (int, int) {
	var list_data Data
	jsonFile, err := os.Open("file.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &list_data)

	return list_data.Status.Water, list_data.Status.Wind
}

func checkStatus(water int, wind int) (stsWater string, stsWind string) {
	if water <= 5 {
		stsWater = "Aman"
	} else if water >= 6 && water <= 8 {
		stsWater = "Siaga"
	} else {
		stsWater = "Bahaya"
	}

	if wind <= 6 {
		stsWind = "Aman"
	} else if wind >= 7 && wind <= 15 {
		stsWind = "Siaga"
	} else {
		stsWind = "Bahaya"
	}

	return
}

package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"quotes/pkg/uid"
)

func main() {
	file, err := os.OpenFile("posts.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	id := "8c257b3d-5409-4a1b-93ae-8009c71ed441"
	for i := 0; i < 100*1000000; i++ {
		if err := writer.Write([]string{uid.New().String(), id, "auto", time.Now().Format("2006-01-02 15:04:05"), "", ""}); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("done")
}

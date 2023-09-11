package main

import (
	"BlueSoftRecruitmentTask/pkg/currency"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	xDefault = 10
	yDefault = 5

	minRange = 4.5
	maxRange = 4.7
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	var x, y int
	flag.IntVar(&x, "x", xDefault, "Number of checks to be performed")
	flag.IntVar(&y, "y", yDefault, "Interval in seconds between checks")
	flag.Parse()

	client := currency.NewClient()
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", y))
	if err != nil {
		log.Fatal(err)
		return
	}

	for i := 0; i < x; i++ {
		waitForChannel := time.After(duration)
		last100, err := client.GetLast100()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Printf("Last 100: %v", last100)
		for _, rate := range last100.Rates {
			if !(minRange < rate.Mid && rate.Mid < maxRange) {
				log.Println(rate.EffectiveDate)
			}
		}
		<-waitForChannel
	}
}

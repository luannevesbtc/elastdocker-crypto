package main

import (
	"encoding/json"
	"fmt"
	"github.com/luannevesb/workers-go-routine/model"
	"log"
	"net/http"
	"time"
)

func main() {
	ini := time.Now()
	pages := make(chan int, 50)
	passengers := make(chan model.Passengers)
	for i := 1; i <= 50; i++ {
		pages <- i
	}

	close(pages)

	go initGetPassengersWorkers(pages, passengers)

	for ps := range passengers {
		for i := range ps.Data {
			log.Printf("PASSENGERS: %x", ps.Data[i].ID)
		}
	}

	log.Println("(Took ", time.Since(ini).Seconds(), "secs)")
}

func initGetPassengersWorkers(pages <-chan int, passengers chan<- model.Passengers) {
	log.Print("[INIT] initGetPassengersWorkers running!")
	go workerGetPassengers(pages, passengers)
}

func workerGetPassengers(page <-chan int, result chan<- model.Passengers) {
	for pg := range page {
		passengers, err := getPassagersRequest(pg)
		if err != nil {
			log.Fatalf("ooopsss an error occurred, please try again: Err:%s", err)
		}
		result <- passengers
	}
	close(result)
}

func getPassagersRequest(page int) (model.Passengers, error) {
	//We make HTTP request using the Get function
	resp, err := http.Get(fmt.Sprintf("https://api.instantwebtools.net/v1/passenger?page=%d&size=10", page))
	if err != nil {
		return model.Passengers{}, err
	}
	defer resp.Body.Close()

	//Create a variable of the same type as our model
	var cResp model.Passengers

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return model.Passengers{}, err
	}

	return cResp, nil
}

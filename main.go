package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/luannevesb/workers-go-routine/model"
)

var numberOfWorkers = 50

func main() {
	//Initiation time
	ini := time.Now()

	//Make a buffered channel with status open for page numbers
	pages := make(chan int, 300)

	//Make a unbuffered channel with status open for passengers infos
	//It is in this channel that we send info for the go routines
	passengers := make(chan model.Passengers)

	//Make a unbuffered channel to main function know when its all done
	done := make(chan bool)

	//Load the buffered channel with all pages to the requests
	for i := 1; i <= 300; i++ {
		pages <- i
	}

	//And i can close because i already send the infos
	close(pages)

	//Call the func init for workers
	go initGetPassengersWorkers(pages, passengers)
	go initSendPassengersWorkers(passengers, done)

	//Waits for the info that comes through the passengers channel
	<-done
	//Logging the time it takes for completion
	log.Print("(Took ", time.Since(ini).Seconds(), "secs)")
}

func initGetPassengersWorkers(pages <-chan int, passengers chan model.Passengers) {
	log.Print("[INIT] initGetPassengersWorkers running!")

	//Creating a wait group to wait for all go routines finish
	var wg sync.WaitGroup

	//In Here you can increase the number of workers to make it faster
	for i := 0; i < numberOfWorkers; i++ {
		log.Println("Main: Starting worker", i)
		//Add one worker to the wait group
		wg.Add(1)

		//Create the worker using a go routine
		go workerGetPassengers(&wg, pages, passengers)
	}

	//Here we wait until all go routines finish
	wg.Wait()
	//Then we close the chan because we finish on writing all the infos
	close(passengers)

}

func workerGetPassengers(wg *sync.WaitGroup, page <-chan int, passengers chan model.Passengers) {
	defer wg.Done()
	//For every page that are in the buffered channel make the request
	for pg := range page {
		resp, err := getPassagersRequest(pg)
		if err != nil {
			log.Fatalf("ooopsss an error occurred, please try again: Err:%s", err)
		}
		//Send data to the passengers unbuffered channel
		passengers <- resp
	}
}

func initSendPassengersWorkers(passengers <-chan model.Passengers, done chan<- bool) {
	var wg sync.WaitGroup

	//In Here you can increase the number of workers to make it faster
	for i := 0; i < numberOfWorkers; i++ {
		fmt.Println("Main: Starting worker", i)
		//Add one worker to the wait group
		wg.Add(1)

		//Create the worker using a go routine
		go workerSendPassengers(&wg, passengers, done)
	}

	//Here we wait until all go routines finish
	wg.Wait()
	//And send done to the main func to tell the process of get and send is done
	done <- true
}

func workerSendPassengers(wg *sync.WaitGroup, passengers <-chan model.Passengers, done chan<- bool) {
	defer wg.Done()
	for ps := range passengers {
		for i := range ps.Data {
			//In Here you have total control over the data and can do whatever you want
			log.Printf("PASSENGERS: %s", ps.Data[i].ID)
		}
	}
}

func getPassagersRequest(page int) (model.Passengers, error) {
	//We make HTTP request using the Get function
	resp, err := http.Get(fmt.Sprintf("https://api.instantwebtools.net/v1/passenger?page=%d&size=100", page))
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

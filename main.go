package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const Usage =
	`
./fizzbuzz <command> [<args>]" 

2 commands are supported:

client			Runs a few concurrent clients that poll the fizzbuzz endpoint and prints the output
server			Runs a FizzBuzzServer
`


type FizzBuzzResponse struct {
	Response string `json:"response"`
}


type FizzBuzzHandler struct {
	counter uint8
	mutex sync.Mutex
}

func (fbh *FizzBuzzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("{}"))
		return
	}
	var resp FizzBuzzResponse
	fbh.mutex.Lock()
	defer fbh.mutex.Unlock()

	fbh.counter++
	if fbh.counter % 3 == 0 && fbh.counter % 5 == 0 {
		resp.Response = "fizzbuzz"
	} else if fbh.counter % 5 == 0 {
		resp.Response = "buzz"
	} else if fbh.counter % 3 == 0 {
		resp.Response = "fizz"
	} else {
		resp.Response = strconv.Itoa(int(fbh.counter))
	}

	if fbh.counter >= 100 {
		fbh.counter = 0
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	_, err = w.Write(respBytes)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

type FizzBuzzClient struct {
	http.Client
}

func (fbc *FizzBuzzClient) GetFizzBuzz(fizzBuzzAddr string, exit chan bool) {
	for {
		resp, err := fbc.Post(fizzBuzzAddr + "/fizzbuzz", "application/json", &bytes.Buffer{})
		if err != nil {
			fmt.Println(err)
			exit <- true
			return
		}
		rb , err  := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			exit <- true
			return
		}
		fmt.Println(string(rb))
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)) + 1000)
	}
}

func getEnv(name, defaultValue string) string {
	val, found := os.LookupEnv(name)
	if !found {
		return defaultValue
	}
	return val
}


func main() {
	if len(os.Args) != 2 {
		log.Fatalf(Usage)
	}

	port := getEnv("FIZZBUZZ_PORT", "4343")
	remoteAddr := getEnv("FIZZBUZZ_REMOTE_ADDR", "http://localhost:4343")

	switch os.Args[1] {
	case "server":
		{
			fbh := &FizzBuzzHandler{}
			mux := http.NewServeMux()
			mux.Handle("/fizzbuzz", fbh)
			mux.Handle("/", http.NotFoundHandler())
			fmt.Printf("listening on port: %v\n", port)
			err := http.ListenAndServe(":" + port, mux)
			log.Fatal(err)
		}
	case "client":
		{
			exit := make(chan bool)
			for i := 0; i < 5; i++ {
				fbc := &FizzBuzzClient{http.Client{Timeout: time.Second * 10}}
				go func() { fbc.GetFizzBuzz(remoteAddr, exit) }()
			}
			// Wait until one of the clients determines that it's time to finish
			<-exit
		}
	default:
		{
			log.Fatalf(Usage)
		}
	}
}

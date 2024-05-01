package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	memStorage "github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
)

func main() {
	type Student struct {
		Fname  string
		Lname  string
		City   string
		Mobile int64
	}
	// storage
	_, err := memStorage.New()
	if err != nil {
		panic("Failed to init storage")
	}
	pollCount := 0
	m := new(runtime.MemStats)
	for {
		runtime.ReadMemStats(m)
		pollCount++
		if pollCount%5 == 0 {
			// отправить метрики на сервер
			buf := bytes.NewReader([]byte{})

			v := reflect.ValueOf(*m)
			typeOfS := v.Type()

			for i := 0; i < v.NumField(); i++ {
				resp, err := http.Post(fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", typeOfS.Field(i).Name, v.Field(i).Interface()), "text/plain", buf)
				if err != nil {
					panic(err)
				}
				fmt.Println(resp)
			}

			resp, err := http.Post(fmt.Sprintf("http://localhost:8080/update/counter/%v/%v", "PollCount", pollCount), "text/plain", buf)
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)

			randomValue := randFloats(1.10, 101.98, 5)

			resp, err = http.Post(fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", "RandomValue", randomValue), "text/plain", buf)
			if err != nil {
				panic(err)
			}

			fmt.Println(resp)
		}
		time.Sleep(2 * time.Second)
	}
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

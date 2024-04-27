package update

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MetricsSaver interface {
	SaveGauge(name string, value float64) (err error)
	SaveCount(name string, value int64) (err error)
}

const (
	gauge   string = "gauge"
	counter string = "counter"
)

func New(storage MetricsSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
			return
		}

		p := strings.Split(r.URL.Path, "/")
		if len(p) < 4 {
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}

		//    1          2           3                4
		// update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		mType := p[2]
		mName := p[3]
		mVal := p[4]

		if mType == "" || (mType != gauge && mType != counter) {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if mName == "" {
			http.Error(w, "Invalid metric name", http.StatusNotFound)
			return
		}

		fmt.Println(mType)
		if mType == gauge {
			if val, err := strconv.ParseFloat(mVal, 64); err == nil {
				err := storage.SaveGauge(mName, val)
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
			}
		}

		if mType == counter {
			if val, err := strconv.ParseInt(mVal, 0, 64); err == nil {
				err := storage.SaveCount(mName, val)
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
			}
		}
	}
}

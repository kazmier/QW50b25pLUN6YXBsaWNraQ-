package url

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type model struct {
	ID       int            `json:"id"`
	Url      string         `json:"url"`
	Interval float32        `json:"interval"`
	history  []historyEntry `json:"-"`
}

type historyEntry struct {
	Response  interface{} `json:"response"`
	Duration  float64     `json:"duration"`
	CreatedAt int64       `json:"created_at"`
}

type jsonResponse struct {
	ID int `json:"id"`
}

const maxFileSize = 1 << (10 * 2)

var urls = make([]model, 0)

func SaveUrl(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxFileSize {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	input := model{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// if ID has zero value, create new record
	if input.ID == 0 {
		input.ID = 1
		if len(urls) != 0 {
			input.ID = urls[len(urls)-1].ID + 1
		}
		urls = append(urls, input)
		if err := respond(w, jsonResponse{ID: input.ID}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// update existing record
	if _, m := findUrl(input.ID); m != nil {
		m.Url = input.Url
		m.Interval = input.Interval
		if err := respond(w, jsonResponse{ID: m.ID}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, m := findUrl(id); m != nil {
		if err := respond(w, m); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		end := time.Now()
		m.history = append(m.history, historyEntry{
			Response:  m,
			Duration:  end.Sub(start).Seconds(),
			CreatedAt: end.Unix(),
		})
		return
	}
	w.WriteHeader(http.StatusNotFound)
	return
}

func DeleteUrl(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if k, m := findUrl(id); m != nil {
		copy(urls[k:], urls[k+1:])
		urls[len(urls)-1] = model{}
		urls = urls[:len(urls)-1]
		if err := respond(w, jsonResponse{ID: id}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
	return
}

func GetAllUrls(w http.ResponseWriter, r *http.Request) {
	if err := respond(w, urls); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func GetUrlHistory(w http.ResponseWriter, r *http.Request) {
	idVar := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idVar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, m := findUrl(id); m != nil {
		if err = respond(w, m.history); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	return
}

func findUrl(id int) (int, *model) {
	for k := range urls {
		if urls[k].ID == id {
			return k, &urls[k]
		}
	}
	return 0, nil
}

func respond(w http.ResponseWriter, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(payload)
}

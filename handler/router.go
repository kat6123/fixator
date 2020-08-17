package handler

import (
	"encoding/json"
	"fixator/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const layout = "02.01.2006"

// API provides methods to work with handler package.
type (
	Fixator interface {
		Fix(fixation *model.Fixation) error
		Select(date model.FixationTime, velocity model.FixationFloat) ([]*model.Fixation, error)
		SelectRange(date model.FixationTime) ([2]*model.Fixation, error)
	}

	Config struct {
		SelectStartHour TimePeriod `yaml:"select_start"`
		SelectEndHour   TimePeriod `yaml:"select_end"`
	}

	API struct {
		s Fixator
		r *mux.Router
	}
)

// New returns new API instance and initializes api router.
func New(s Fixator, config Config) *API {
	api := &API{
		s: s,
		r: mux.NewRouter().StrictSlash(true),
	}

	api.r.HandleFunc("/fixation", api.Fix).Methods("POST")

	inPeriodMiddleware := inPeriod(config.SelectStartHour, config.SelectEndHour)
	api.r.HandleFunc("/fixation/select", inPeriodMiddleware(api.SelectFixations)).
		Queries("date", "{date}", "start", "{start}").Methods("GET")
	api.r.HandleFunc("/fixation/select/range", inPeriodMiddleware(api.SelectRange)).
		Queries("date", "{date}").Methods("GET")

	return api
}

// Router returns Router of API instance. It will be initialized after New method is called.
func (a API) Router() http.Handler {
	return a.r
}

func (a API) Fix(w http.ResponseWriter, r *http.Request) {
	var f model.Fixation

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.s.Fix(&f); err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a API) SelectFixations(w http.ResponseWriter, r *http.Request) {
	parsedTime, err := time.Parse(layout, r.FormValue("date"))
	if err != nil {
		err = fmt.Errorf("parse %s as fixation time has failed: %v", "date", err)
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var start model.FixationFloat
	if err := start.UnmarshalJSON([]byte(r.FormValue("start"))); err != nil {
		err = fmt.Errorf("parse %s as fixation velocity has failed: %v", "start", err)
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fixations, err := a.s.Select(model.FixationTime(parsedTime), start)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, fixations)
}

func (a API) SelectRange(w http.ResponseWriter, r *http.Request) {
	parsedTime, err := time.Parse(layout, r.FormValue("date"))
	if err != nil {
		err = fmt.Errorf("parse %s as fixation time has failed: %v", "date", err)
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	velocityRange, err := a.s.SelectRange(model.FixationTime(parsedTime))
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, velocityRange)
}

func jsonResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		writeJSON(w, fmt.Sprintf("encode %s as json: %v", v, err), http.StatusInternalServerError)
		return
	}
}
func writeJSON(w http.ResponseWriter, errMsg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	message := struct {
		Content string `json:"error"`
	}{errMsg}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		panic(fmt.Sprintf("encode error to json has failed: %v", err))
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"microsomes.com/stgo/utils"
)

type AddJob struct {
	Callback func(jobDetails utils.SteamJob)
}

func (a *AddJob) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jobDetails := utils.SteamJob{}

	err = json.Unmarshal(body, &jobDetails)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Job: ", jobDetails.JobID)

	a.Callback(jobDetails)

	w.Write([]byte("OK"))
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"microsomes.com/stgo/utils"
)

type JobStatus struct {
	GetStatusesByJobID func(id string) utils.JobStatus
}

func (j *JobStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	jobID := r.URL.Query().Get("jobid")

	fmt.Println("JobID: ", jobID)

	if jobID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := j.GetStatusesByJobID(jobID)

	if status.State == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(status)
}

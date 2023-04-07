package handlers

import "net/http"

type WorkerStatus struct {
	GetWorkerStatus func() string
}

func (w *WorkerStatus) ServeHTTP(http.ResponseWriter, *http.Request) {
}

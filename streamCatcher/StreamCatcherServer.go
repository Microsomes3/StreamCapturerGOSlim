package streamcatcher

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"syscall"

	"microsomes.com/stgo/streamCatcher/handlers"
)

type StreamCatcherServer struct {
	client        http.Server
	StreamCatcher *StreamCatcher
}

func NewStreamCatcherServer() *StreamCatcherServer {

	return &StreamCatcherServer{
		client: http.Server{
			Addr: ":9005",
		},
		StreamCatcher: NewStreamCatcher(),
	}
}

func Drain() {}

func (s *StreamCatcherServer) StartAndServe(wg *sync.WaitGroup, ctx context.Context) {

	queueCtx, cancelQueue := context.WithCancel(ctx)
	go s.StreamCatcher.StartAllWorkers(wg)
	go s.StreamCatcher.StartQueues(queueCtx)

	files := handlers.Files{}
	addJob := handlers.AddJob{
		Callback: s.StreamCatcher.AddJob,
	}

	jobStatus := handlers.JobStatus{
		GetStatusesByJobID:    s.StreamCatcher.GetJobStatus,
		GetAllStatusesByJobID: s.StreamCatcher.GetAllStatusesByJobID,
	}

	workerQueue := handlers.WorkerStatus{
		GetWorkerStatus: s.StreamCatcher.GetWorkerStatus,
	}

	drain := handlers.Drain{
		Callback: func() {
			fmt.Println("Draining")

			//send siganal sigint
			syscall.Kill(os.Getpid(), syscall.SIGINT)
		},
	}

	http.HandleFunc("/", files.ServeHTTP)
	http.HandleFunc("/addJob", addJob.ServeHTTP)
	http.HandleFunc("/drain", drain.ServeHTTP)
	http.HandleFunc("/jobStatus", jobStatus.ServeHTTP)
	http.HandleFunc("/workerStatus", workerQueue.ServeHTTP)
	go s.client.ListenAndServe()

	select {
	case <-ctx.Done():
		cancelQueue()
		s.client.Shutdown(ctx)
		fmt.Println("Server shutdown")
	}
}

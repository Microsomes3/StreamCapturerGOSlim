package streamcatcher

import (
	"context"
	"fmt"
	"sync"

	"microsomes.com/stgo/streamCatcher/streamutil"
	"microsomes.com/stgo/utils"
)

type StreamCatcher struct {
	JobQueue        chan utils.SteamJob
	WorkQueue       chan utils.SteamJob
	ConcurrentLimit int
	JobStatuses     map[string]utils.JobStatus
	JobStatusEvents map[string]utils.JobStatusEvents
}

func (s *StreamCatcher) GetJobStatus(id string) utils.JobStatus {
	return s.JobStatuses[id]
}

func NewStreamCatcher() *StreamCatcher {
	return &StreamCatcher{
		JobQueue:        make(chan utils.SteamJob, 1),
		WorkQueue:       make(chan utils.SteamJob, 100),
		ConcurrentLimit: 10,
		JobStatuses:     make(map[string]utils.JobStatus),
		JobStatusEvents: make(map[string]utils.JobStatusEvents),
	}
}

func (s *StreamCatcher) AddStatusEvent(job *utils.SteamJob, status string, result []string) {

	nstatus := utils.JobStatus{
		State:  status,
		Result: result,
	}

	s.JobStatusEvents[job.JobID] = append(s.JobStatusEvents[job.JobID], nstatus)

	s.JobStatuses[job.JobID] = nstatus

}

func (s *StreamCatcher) AddJob(job utils.SteamJob) {
	fmt.Println("Adding job to job queue: ", job.JobID)
	s.JobQueue <- job
}

func (s *StreamCatcher) StartWork(wg *sync.WaitGroup) {
	defer wg.Done()

	for Job := range s.WorkQueue {

		data, err := streamutil.ProcessDownload(Job.YoutubeLink, Job.TimeoutSeconds, Job.JobID, Job.IsStart)
		if err != nil {
			fmt.Println("Error: ", err)
			s.AddStatusEvent(&Job, "error", []string{err.Error()})
		}

		s.AddStatusEvent(&Job, "done", data.Paths)

		fmt.Println("Processed: ", Job.JobID, " done")
	}

	fmt.Println("interrupted")

}

func (s *StreamCatcher) StartQueues(ctx context.Context) {
	for {
		select {
		case job := <-s.JobQueue:
			s.WorkQueue <- job
		case <-ctx.Done():
			close(s.WorkQueue)
			close(s.JobQueue)
			return
		}
	}
}

func (s *StreamCatcher) StartAllWorkers(diewg *sync.WaitGroup) {
	wg1 := &sync.WaitGroup{}
	wg1.Add(s.ConcurrentLimit)
	for i := 0; i < s.ConcurrentLimit; i++ {
		go s.StartWork(wg1)
	}
	wg1.Wait()
	diewg.Done()
	fmt.Println("work done die")
}

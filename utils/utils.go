package utils

type SteamJob struct {
	JobID          string `json:"jobId"`
	YoutubeLink    string `json:"youtubeLink"`
	TimeoutSeconds int    `json:"timeout"`
	IsStart        bool   `json:"isStart"`
	UpdateHook     string `json:"updateHook"`
}

type JobStatus struct {
	State  string   `json:"state"`
	Result []string `json:"result"`
}

type JobStatusEvents []JobStatus

type JobResponse struct {
	Status   string   `json:"status"`
	Reason   string   `json:"reason"`
	Paths    []string `json:"paths"`
	Comments []string `json:"comments"`
}

type WorkerStatus struct {
	TotalQueue     int `json:"totalQueue"`
	TotalRecording int `json:"totalRecording"`
	TotalDone      int `json:"totalDone"`
}

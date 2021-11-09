package godisApi

type SetResult struct {
	Success bool `json:"success"`
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

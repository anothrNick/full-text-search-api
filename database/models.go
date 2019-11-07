package database

// ProjectData represents a `project_data` database record
type ProjectData struct {
	ID      string      `json:"id"`
	Project string      `json:"project"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

package web

import (
	"encoding/json"
	"net/http"

	"github.com/anothrNick/full-text-search-api/database"
	"github.com/gin-gonic/gin"
)

// ProjectIndex represents a searchable record, `Data` and `Meta` can be anything, so Unmarshal the raw data
type ProjectIndex struct {
	Data json.RawMessage `json:"data"`
	Meta json.RawMessage `json:"meta"`
}

// Database is the required interface to the DB layer from the HTTP handlers
type Database interface {
	TranslateError(err error) *database.TranslatedError

	AddProjectData(projectName string, data []byte, meta []byte) error
	SearchProjectData(projectName, query string) ([]*database.ProjectData, error)
}

// Handlers contains all handler functions
type Handlers struct {
	db Database
}

// NewHandlers creates and returns a new instance of `Handlers` with the datastore
func NewHandlers(datastore Database) *Handlers {
	return &Handlers{
		db: datastore,
	}
}

// CreateRecord creates a new searchable record for a project
func (h *Handlers) CreateRecord(c *gin.Context) {
	project := c.Param("project")
	projectIndex := &ProjectIndex{}
	c.BindJSON(projectIndex)

	err := h.db.AddProjectData(project, projectIndex.Data, projectIndex.Meta)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// SearchRecords searches a project's records
func (h *Handlers) SearchRecords(c *gin.Context) {
	project := c.Param("project")
	search := c.Query("_search")

	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'_search' query parameter is required"})
	}

	records, err := h.db.SearchProjectData(project, search)

	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": records})
}

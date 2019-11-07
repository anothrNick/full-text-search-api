package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

// Postgres wraps the db connection to a postgres instance
type Postgres struct {
	db *sql.DB
}

// NewPostgres returns a new instance of `Postgres` with a connection to the postgres db based on provided parameters.
func NewPostgres(user, password, host, database string) (*Postgres, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		db: db,
	}, nil
}

// TranslateError translates the database specific error to a simple `error` to return to the user.
func (p *Postgres) TranslateError(err error) *TranslatedError {
	originalError := err.Error()

	if err, ok := err.(*pq.Error); ok {
		// postgres specific errors
		switch err.Code {
		case "23505":
			return NewTranslatedError(http.StatusBadRequest, errors.New("key already exists"))
		case "22023":
			return NewTranslatedError(http.StatusBadRequest, errors.New("key already exists"))
		default:
			fmt.Println(err)
			fmt.Println(err.Code)
			return NewTranslatedError(http.StatusInternalServerError, errors.New(strings.TrimPrefix(err.Error(), "pq: ")))
		}
	} else {
		// generic sql package errors
		switch originalError {
		case sql.ErrConnDone.Error():
			return NewTranslatedError(http.StatusInternalServerError, errors.New("internal server error"))
		case sql.ErrNoRows.Error():
			return NewTranslatedError(http.StatusNotFound, errors.New("not found"))
		case sql.ErrTxDone.Error():
			return NewTranslatedError(http.StatusInternalServerError, errors.New("internal server error"))
		}
	}

	fmt.Println(err)
	return NewTranslatedError(http.StatusInternalServerError, errors.New("internal server error"))
}

// AddProjectData creates a new searchable record for a project
func (p *Postgres) AddProjectData(projectName string, data []byte, meta []byte) error {
	_, err := p.db.Exec("INSERT INTO project_data (project, data, meta) VALUES ($1, $2, $3)", projectName, data, meta)
	return err
}

// SearchProjectData performs a full text search on a project
func (p *Postgres) SearchProjectData(projectName, query string) ([]*ProjectData, error) {
	rows, err := p.db.Query(
		"SELECT id, project, data, meta FROM project_data WHERE project=$1 AND to_tsvector('English', data::text) @@ plainto_tsquery('English', $2)",
		projectName,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]*ProjectData, 0)
	for rows.Next() {
		rec := ProjectData{}
		data := make([]byte, 0)
		meta := make([]byte, 0)
		err = rows.Scan(
			&rec.ID,
			&rec.Project,
			&data,
			&meta,
		)

		json.Unmarshal(data, &rec.Data)
		json.Unmarshal(meta, &rec.Meta)
		if err != nil {
			return nil, err
		}

		records = append(records, &rec)
	}

	return records, nil
}

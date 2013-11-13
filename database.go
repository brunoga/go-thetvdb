// Copyright 2012 Bruno Albuquerque.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package thetvdb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// LocalSeriesDatabase represents a local sqlite-based series database.
type LocalSeriesDatabase struct {
	db *sql.DB
}

// NewLocalSeriesDatabase creates and returns a new LocalSeriesDatabase instance.
func NewLocalSeriesDatabase(path string) (*LocalSeriesDatabase, error) {
	// Open database trying to create it if it does not exist.
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database %q : %v", path,
			err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Series (Id INTEGER " +
		"PRIMARY KEY, Name TEXT, Genre TEXT, Status TEXT, FetchDate " +
		"DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return nil, fmt.Errorf("error creating Series table : %v", err)
	}

	return &LocalSeriesDatabase{db: db}, nil
}

// Lookup searches for the series with the given seriesId in the local series
// database and, if found, returns a Series instance representing it.
func (db *LocalSeriesDatabase) Lookup(seriesId int) (*Series, error) {
	stmt, err := db.db.Prepare("SELECT * FROM Series WHERE Id = ?")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement : %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(seriesId)

	var id int64
	var name, status, genre string
	var fetchDate time.Time
	err = row.Scan(&id, &name, &genre, &status, &fetchDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	elapsedTimeHours := time.Now().Sub(fetchDate).Hours()
	if elapsedTimeHours > 48 {
		// Cache expired for this entry. Fetch again.
		db.Remove(seriesId)
		return nil, nil
	}

	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	return &Series{
		Id: seriesId,
		Name: name,
		Status: status,
		Genre: genre,
		FetchDate: fetchDate.Format(layout),
	}, nil
}

// Insert inserts the given series to the local series database.
func (db *LocalSeriesDatabase) Insert(series Series) error {
	stmt, err := db.db.Prepare("INSERT INTO Series (Id, Name, Genre, Status) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Errorf("error preparing statement : %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(series.Id, series.Name, series.Genre, series.Status)
	if err != nil {
		return fmt.Errorf("error executing statement : %v", err)
	}

	return nil
}


// Remove removes the series with the given seriesId from the local series
// database.
func (db *LocalSeriesDatabase) Remove(seriesId int) error {
	stmt, err := db.db.Prepare("DELETE FROM Series WHERE Id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement : %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(seriesId)
	if err != nil {
		return fmt.Errorf("error executing statement : %v", err)
	}

	return nil
}

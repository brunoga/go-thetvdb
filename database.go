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
	"fmt"
	"time"

	"github.com/feyeleanor/gosqlite3"
)

type LocalSeriesDatabase struct {
	db *sqlite3.Database
}

func NewLocalSeriesDatabase(path string) (*LocalSeriesDatabase, error) {
	// Open database trying to create it if it does not exist.
	db, err := sqlite3.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening database %q : %v", path,
			err)
	}

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS Series " +
		"(Id INTEGER PRIMARY KEY, Name TEXT, Status TEXT, FetchDate DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return nil, fmt.Errorf("error creating Series table : %v", err)
	}

	return &LocalSeriesDatabase{db: db}, nil
}

func (db *LocalSeriesDatabase) Lookup(seriesId int) (*Series, error) {
	gotSeries := false
	series := Series{}
	_, err := db.db.Execute(fmt.Sprintf(
		"SELECT * FROM Series WHERE Id = %d", seriesId),
		func(st *sqlite3.Statement, values ...interface{}) {
			series.Id = int(st.Column(0).(int64))
			series.Name = st.Column(1).(string)
			series.Status = st.Column(2).(string)
			series.FetchDate = st.Column(3).(string)
			gotSeries = true
		})
	if err != nil {
		return nil, err
	}

	if !gotSeries {
		return nil, nil
	}

	fetchDate, err := time.Parse("2006-01-02 15:04:05", series.FetchDate)
	if err != nil {
		// If we got an error, just force refetching the series entry.
		db.Remove(series.Id)
		return nil, nil
	} else {
		elapsedTimeHours := time.Now().Sub(fetchDate).Hours()
		if elapsedTimeHours > 48 {
			// Cache expired for this entry. Fetch again.
			db.Remove(series.Id)
			return nil, nil
		}
	}

	return &series, nil
}

func (db *LocalSeriesDatabase) Insert(series Series) error {
	sql := fmt.Sprintf(
		"INSERT INTO Series (Id, Name, Status) VALUES (%d, %q, %q)",
		series.Id, series.Name, series.Status)
	_, err := db.db.Execute(sql)
	return err
}

func (db *LocalSeriesDatabase) Remove(seriesId int) error {
	_, err := db.db.Execute(fmt.Sprintf(
		"DELETE FROM Series WHERE Id = %d", seriesId))
	return err
}

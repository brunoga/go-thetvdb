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
	"strconv"
)

// Methods related to series information.
//
// Requirements:
//      TODO(bga): Fill in requirements.

// A Series object contains information about a specific TV series.
type Series struct {
	Id               int `xml:"id"`
	Actors           string
	AirDay           string `xml:"Airs_DayOfWeek"`
	AirTime          string `xml:"Airs_Time"`
	ContentRating    string
	FirstAired       string
	Genre            string
	IMDBId           string `xml:"IMDB_ID"`
	Language         string
	Network          string
	Name             string `xml:"SeriesName"`
	BannerPathSuffix string `xml:"banner"`
	Overview         string
	Rating           string
	RatingCount      int
	Runtime          int
	Status           string
	FetchDate        string
}

// GetSeriesById fetches information for the series identified by seriesId and
// returns a Series struct filled with this information.
func (t *TheTVDB) GetSeriesById(seriesId int) (*Series, error) {
	series, err := t.db.Lookup(seriesId)
	if err != nil {
		return nil, err
	} else {
		if series != nil {
			// Found series on local database. Return it.
			return series, nil
		}
	}

	// Did not find series in local database. Fetch it and save.
	err = ValidateAPIKey(t.apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid API key : %v", err)
	}

	type Data struct {
		Series []Series
	}
	v := Data{}

	err = doRequest(t.apiKey+"/series/"+strconv.Itoa(seriesId)+
		"/en.xml", nil, &v)
	if err != nil {
		return nil, err
	}

	if len(v.Series) == 0 {
		return nil, fmt.Errorf("empty series data")
	}

	err = t.db.Insert(v.Series[0])
	if err != nil {
		fmt.Println(err)
	}

	return &v.Series[0], nil
}

// GetSeries fetches series that contains the seriesNameSubstring as part of
// its name and returns a slice of Series objects with data about those series,
// including their name, id and language.
func (t *TheTVDB) GetSeries(seriesNameSubstring string) ([]Series, error) {
	type Data struct {
		Series []Series
	}

	parameters := parameterMap{
		"seriesname": seriesNameSubstring,
		"language":   "en",
	}

	v := Data{}

	err := doRequest("GetSeries.php", parameters, &v)
	if err != nil {
		return nil, err
	}

	return v.Series, nil
}

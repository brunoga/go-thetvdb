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
	FetchDate	 string
}

func (t *TheTVDB) GetSeriesById(seriesId int) (*Series, error) {
	series, err := t.db.Lookup(seriesId)
	if err != nil {
		fmt.Println(err)
	} else {
		if series != nil {
			return series, nil
		}
	}

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

	err = t.db.Insert(v.Series[0])
	if err != nil {
		fmt.Println(err)
	}

	return &v.Series[0], nil
}

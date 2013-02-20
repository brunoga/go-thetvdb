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

package main

import (
	"flag"
	"fmt"

	thetvdb "github.com/brunoga/go-thetvdb"
)

// Flags.
var (
	theTVDBAccountId = flag.String("accountid", "", "TheTVDB account id")
        theTVDBApiKey = flag.String("apikey", "", "TheTVDB API key")
)

func main() {
	flag.Parse()

	tvdb, err := thetvdb.New(*theTVDBApiKey, *theTVDBAccountId)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	fmt.Printf("Fetching user favorites ... ");
	userFavorites, err := tvdb.GetUserFavorites()
	if err != nil {
		fmt.Println("Error retrieving favorites :", err)
		return
	} else {
		fmt.Println("Ok")
		fmt.Println("Number of favorites :", len(userFavorites))
	}

	var foundEnded = false
	for _, seriesId := range userFavorites {
		series, err := tvdb.GetSeriesById(seriesId)
		if err != nil {
			fmt.Println("Error getting series data :", err)
		} else {
			if series.Status == "Ended" {
				foundEnded = true
				fmt.Printf("Removing series name : %q (%s)\n", series.Name, series.Status)
				tvdb.RemoveUserFavorite(series.Id)
			}
		}
	}

	if !foundEnded {
		fmt.Println("No ended series found.")
	}
}


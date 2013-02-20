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

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage : list_series substring")
		return
	}

	// GetSeries does not require either an API key or an account id.
	tvdb, err := thetvdb.New("", "")
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	fmt.Printf("Fetching series ... ")
	series, err := tvdb.GetSeries(flag.Args()[0])
	if err != nil {
		fmt.Println("Error retrieving series :", err)
		return
	} else {
		fmt.Println("Ok")
		for _, serie := range series {
			fmt.Println(serie.Id, serie.Name)
		}
	}
}

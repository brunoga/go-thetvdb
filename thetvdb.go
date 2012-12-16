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

// Package thetvdb implements an interface to the TheTVDB API that allows
// querying their dabatase for information about TV series and episodes in
// general, including support for user favorites.
package thetvdb

// References:
//    TheTVDB API: http://www.thetvdb.com/wiki/index.php/Programmers_API

// TODO(bga):
//	Finish implementing all methods.
//	Add client side caching.
//      write tests.

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// A TheTVDB object represents an entry point to the TheTVDB API.
type TheTVDB struct {
	apiKey    string
	accountId string
	configDir string
	db        *LocalSeriesDatabase
}

// New creates a new TheTVDB instance associated with the given API key and/or
// account id. Note both parameters are optional. Calling methods that require
// any of the 2 that was not set will result in errors.
func New(apiKey, accountId string) (*TheTVDB, error) {
	if len(apiKey) != 0 {
		err := ValidateAPIKey(apiKey)
		if err != nil {
			return nil, fmt.Errorf(
				"error validating API key : %v", err)
		}
	}

	if len(accountId) != 0 {
		err := ValidateAccountId(accountId)
		if err != nil {
			return nil, fmt.Errorf(
				"error validating account id : %v", err)
		}
	}

	flag.Parse()

	err := os.Mkdir(*configDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Can not create config dir :", err)
		os.Exit(1)
	}

	db, err := NewLocalSeriesDatabase(*configDir + "/go-thetvdb.db")
	if err != nil {
		fmt.Println("Can not open local series database :", err)
		os.Exit(1)
	}

	return &TheTVDB{
		apiKey:    apiKey,
		accountId: accountId,
		configDir: *configDir,
		db:        db,
	}, nil
}

// ValidateAPIKey checks if the given API key is valid (it only checks for its
// format, so it may pass the check and not exist, for example). It will return
// nil if it is valid or a proper error if it is not.
func ValidateAPIKey(apiKey string) error {
	matched, err := regexp.Match("^[0-9A-F]{16}$", []byte(apiKey))
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("invalid API key : %q", apiKey)
	}

	return nil
}

// ValidateAccountId checks if the given account id is valid (it only checks
// for its format, so it may pass the check and not exist, for example). It will
// return nil if it is valid or a proper error if it is not.
func ValidateAccountId(accountId string) error {
	matched, err := regexp.Match("^[0-9A-F]{16}$", []byte(accountId))
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("invalid account id : %q", accountId)
	}

	return nil
}

const (
	apiUrlPrefix = "http://www.thetvdb.com/api/"
)

// Flags.
var (
	configDir = flag.String(
		"configdir", "/home/bga/.go-thetvdb", "config directory")
)

type parameterMap map[string]string

func doRequest(path string, parameters parameterMap,
	container interface{}) error {
	urlParameters := ""
	first := true
	for key, value := range parameters {
		if !first {
			urlParameters = urlParameters + "&"
		} else {
			first = false
		}
		urlParameters = urlParameters + key + "=" + value
	}

	finalUrl := apiUrlPrefix + path
	if urlParameters != "" {
		finalUrl = finalUrl + "?" + urlParameters
	}

	resp, err := http.Get(finalUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, container)
	if err != nil {
		return err
	}

	return nil
}

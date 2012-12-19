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
	"sort"
	"strconv"
)

// Methods related to user favorites list manipulation.
//
// Requirements:
// 	All methods in this file require a valid account id being set.

// GetUserFavorites returns an slice of Series associated with the account id in
// ascending order of their ids. Each Series object will only have their id
// field set and all other fields will be at their default values, so they are
// not meaningful.
func (t *TheTVDB) GetUserFavorites() ([]int, error) {
	err := ValidateAccountId(t.accountId)
	if err != nil {
		return nil, fmt.Errorf("invalid account id : %v", err)
	}

	parameters := parameterMap{
		"accountid": t.accountId,
	}
	seriesIds, err := doUserFavoritesRequest(parameters)
	if err != nil {
		return nil, err
	}

	return seriesIds, nil
}

// AddUserFavorites adds the series identified by the given series id to the
// list of series associated with the account id.
func (t *TheTVDB) AddUserFavorite(seriesId int) error {
	err := ValidateAccountId(t.accountId)
	if err != nil {
                return fmt.Errorf("invalid account id : %v", err)
        }

	if !isSeriesIdValid(seriesId) {
		return fmt.Errorf("invalid series id given : %v", seriesId)
	}

	parameters := parameterMap{
		"accountid": t.accountId,
		"type":      "add",
		"seriesid":  strconv.Itoa(seriesId),
	}
	seriesIds, err := doUserFavoritesRequest(parameters)
	if err != nil {
		return nil
	}

	if !findSeriesId(seriesIds, seriesId) {
		return fmt.Errorf("failed adding series to favorites")
	}

	return nil
}

// RemoveUserFavorite removes the series identified by the given series id from
// the list of series associated with the account id.
func (t *TheTVDB) RemoveUserFavorite(seriesId int) error {
	err := ValidateAccountId(t.accountId)
        if err != nil {
                return fmt.Errorf("invalid account id : %v", err)
        }

	if !isSeriesIdValid(seriesId) {
		return fmt.Errorf("invalid series id given : %v", seriesId)
	}

	parameters := parameterMap{
		"accountid": t.accountId,
		"type":      "remove",
		"seriesid":  strconv.Itoa(seriesId),
	}
	seriesIds, err := doUserFavoritesRequest(parameters)
	if err != nil {
		return nil
	}

	if findSeriesId(seriesIds, seriesId) {
		return fmt.Errorf("failed removing series from favorites")
	}

	err = t.db.Remove(seriesId)
	if err != nil {
		fmt.Println("Warning. Could not remove series from local database.")
	}

	return nil
}

// doUserfavoritesRequest sends a User Favorites request and and sorts the
// resulting list of series ids, returning it.
func doUserFavoritesRequest(parameters parameterMap) ([]int, error) {
	type Favorites struct {
		Series []int
	}
	v := Favorites{}

	err := doRequest("User_Favorites.php", parameters, &v)
	if err != nil {
		return nil, err
	}

	sort.Ints(v.Series)

	return v.Series, nil
}

// findSeriesId checks if the given series id is present in the given (sorted)
// seriesIds slice.
func findSeriesId(seriesIds []int, seriesId int) bool {
	i := sort.Search(len(seriesIds), func(i int) bool {
		return seriesIds[i] >= seriesId
	})
	if seriesIds[i] != seriesId {
		return false
	}

	return true
}

// isSeriesIdValid returns if a given series id is valid. it only does a basic
// check and does not handle cases like a given id simply not existing in the
// TheTVDB.com database.
func isSeriesIdValid(seriesId int) bool {
	if seriesId < 0 {
		return false
	}

	return true
}

// Copyright 2020 Eryx <evorui аt gmail dοt com>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hcapi

import (
	"sort"
	"strings"
	"sync"
)

var (
	dataMu sync.RWMutex
)

type DataPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	T string  `json:"t,omitempty"`
}

type DataItem struct {
	Name   string       `json:"name,omitempty"`
	Attrs  []string     `json:"attrs,omitempty"`
	Points []*DataPoint `json:"points,omitempty"`
}

type DataList struct {
	Items []*DataItem `json:"items,omitempty"`
}

func NewDataItem(name string) *DataItem {
	return &DataItem{
		Name: name,
	}
}

func (it *DataList) Has(ds *DataItem) bool {
	id := ds.id()
	for _, v := range it.Items {
		if v.id() == id {
			return true
		}
	}
	return false
}

func (it *DataList) Set(ds *DataItem) {

	dataMu.Lock()
	defer dataMu.Unlock()

	id := ds.id()
	for i, v := range it.Items {
		if v.id() == id {
			it.Items[i] = ds
			return
		}
	}
	it.Items = append(it.Items, ds)
}

func (it *DataItem) id() string {
	if len(it.Attrs) > 0 {
		sort.Strings(it.Attrs)
		if it.Name == "" {
			return strings.Join(it.Attrs, "/")
		}
		return it.Name + "/" + strings.Join(it.Attrs, "/")
	}
	return it.Name
}

func (it *DataItem) AttrSet(v string) {
	v = strings.ToLower(strings.TrimSpace(v))
	for _, v2 := range it.Attrs {
		if v == v2 {
			return
		}
	}
	it.Attrs = append(it.Attrs, v)
	sort.Strings(it.Attrs)
}

func (it *DataItem) PointSet(x, y float64) {
	it.Point(x).Y = y
}

func (it *DataItem) Point(x float64) *DataPoint {

	for _, v := range it.Points {
		if v.X == x {
			return v
		}
	}

	p := &DataPoint{
		X: x,
	}

	it.Points = append(it.Points, p)
	return p
}

func (it *DataItem) PointXs() []float64 {
	ar := []float64{}
	for _, v := range it.Points {
		ar = append(ar, v.X)
	}
	return ar
}

func (it *DataItem) PointYs() []float64 {
	ar := []float64{}
	for _, v := range it.Points {
		ar = append(ar, v.Y)
	}
	return ar
}

func (it *DataItem) PointTs() []string {
	ar := []string{}
	for _, v := range it.Points {
		ar = append(ar, v.T)
	}
	return ar
}

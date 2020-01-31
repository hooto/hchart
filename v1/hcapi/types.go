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
	"sync"
)

var (
	listMu    sync.RWMutex
	datasetMu sync.RWMutex
)

const (
	ChartTypeBar           = "bar"
	ChartTypeBarHorizontal = "bar-h"
	ChartTypeLine          = "line"
	ChartTypePie           = "pie"
)

type ChartEntry struct {
	Type    string       `json:"type"`
	Options ChartOptions `json:"options"`
	Data    ChartData    `json:"data"`
}

type ChartOptions struct {
	Title  string `json:"title,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

type ChartData struct {
	Labels   []string        `json:"labels,omitempty"`
	Datasets []*ChartDataset `json:"datasets,omitempty"`
}

type ChartDataset struct {
	Label string  `json:"label,omitempty"`
	Data  []int64 `json:"data,omitempty"`
}

func (it *ChartEntry) Valid() error {
	return nil
}

func (it *ChartData) Sync(legendLabel, dsLabel string, dsData int64) {

	datasetMu.Lock()
	defer datasetMu.Unlock()

	for k, v := range it.Datasets {
		if v.Label == legendLabel {
			it.Datasets[k].Data = append(v.Data, dsData)
			if len(it.Datasets[k].Data) > len(it.Labels) {
				it.Labels = append(it.Labels, dsLabel)
			}
			return
		}
	}

	it.Datasets = append(it.Datasets, &ChartDataset{
		Label: legendLabel,
		Data:  []int64{dsData},
	})
	if len(it.Labels) < 1 {
		it.Labels = append(it.Labels, legendLabel)
	}
}

type ChartList struct {
	Items []ChartEntry `json:"items"`
}

func (it *ChartList) Sync(c_type, c_title, legendLabel, dsLabel string, dsData int64) {

	listMu.Lock()
	defer listMu.Unlock()

	for k, v := range it.Items {
		if v.Type == c_type && v.Options.Title == c_title {
			it.Items[k].Data.Sync(legendLabel, dsLabel, dsData)
			return
		}
	}

	it.Items = append(it.Items, ChartEntry{
		Type: c_type,
		Options: ChartOptions{
			Title: c_title,
		},
		Data: ChartData{
			Labels: []string{legendLabel},
			Datasets: []*ChartDataset{
				{
					Label: dsLabel,
					Data:  []int64{dsData},
				},
			},
		},
	})
}

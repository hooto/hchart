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

const (
	ChartTypeLine      = "line"
	ChartTypeBar       = "bar"
	ChartTypeHistogram = "histogram"
)

var (
	chatOptionWidthDefault  float64 = 800
	chatOptionHeightDefault float64 = 400
)

type ChartItem struct {
	Type     string       `json:"type"`
	Options  ChartOptions `json:"options"`
	Labels   []string     `json:"labels,omitempty"`
	Datasets []*DataItem  `json:"datasets,omitempty"`
}

func (it *ChartItem) Valid() error {
	return nil
}

type ChartOptions struct {
	Title  string      `json:"title,omitempty"`
	Width  float64     `json:"width,omitempty"`
	Height float64     `json:"height,omitempty"`
	X      AxisOptions `json:"x,omitempty"`
	Y      AxisOptions `json:"y,omitempty"`
}

type AxisOptions struct {
	Title string `json:"title,omitempty"`
}

func (it *ChartOptions) WidthLength() float64 {
	if it.Width <= 0 {
		it.Width = chatOptionWidthDefault
	} else if it.Width < 100 {
		it.Width = 100
	} else if it.Width > 4000 {
		it.Width = 4000
	}
	return it.Width
}

func (it *ChartOptions) HeightLength() float64 {
	if it.Height <= 0 {
		it.Height = chatOptionHeightDefault
	} else if it.Height < 100 {
		it.Height = 100
	} else if it.Height > 4000 {
		it.Height = 4000
	}
	return it.Height
}

type ChartRenderOptions struct {
	Name      string `json:"name"`
	SvgEnable bool   `json:"svg_enable"`
	PngEnable bool   `json:"png_enable"`
}

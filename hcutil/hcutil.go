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

package hcutil

import (
	"errors"
	"image/color"

	"github.com/hooto/hchart/hcapi"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	// [0, 5]
	colorThemeG = []color.RGBA{
		{0x00, 0x57, 0xe7, 0xff},
		{0xd6, 0x2d, 0x20, 0xff},
		{0xff, 0xa7, 0x00, 0xff},
		{0x00, 0x87, 0x44, 0xff},
		{0x9c, 0x27, 0xb0, 0xff},
		{0x33, 0x33, 0x33, 0xff},
	}
	// [0 ~ 8]
	colorGrays = []color.RGBA{
		{0xf8, 0xf9, 0xfa, 0xff},
		{0xe9, 0xec, 0xef, 0xff},
		{0xde, 0xe2, 0xe6, 0xff},
		{0xce, 0xd4, 0xda, 0xff},
		{0xad, 0xb5, 0xbd, 0xff},
		{0x6c, 0x75, 0x7d, 0xff},
		{0x49, 0x50, 0x57, 0xff},
		{0x34, 0x3a, 0x40, 0xff},
		{0x21, 0x25, 0x29, 0xff},
	}
)

func ColorTheme(n int) color.RGBA {
	return colorThemeG[n%len(colorThemeG)]
}

func ColorGray(n int) color.RGBA {
	return colorGrays[n%len(colorGrays)]
}

func Render(item *hcapi.ChartEntry, opts *hcapi.ChartRenderOptions) error {

	if err := item.Valid(); err != nil {
		return err
	}

	p, _ := plot.New()

	p.Title.Text = item.Options.Title

	switch item.Type {

	case hcapi.ChartTypeLine:
		lineRender(p, item)

	case hcapi.ChartTypeBar:
		barRender(p, item)

	case hcapi.ChartTypeHistogram:
		histRender(p, item)

	default:
		return errors.New("invalid type")
	}

	if err := p.Save(
		vg.Length(item.Options.WidthLength()),
		vg.Length(item.Options.HeightLength()), opts.Name+".png"); err != nil {
		return err
	}

	return nil
}

func lineRender(p *plot.Plot, item *hcapi.ChartEntry) error {

	{
		p.X.Label.Text = item.Options.X.Title
		p.X.Label.TextStyle.Color = ColorGray(8)

		p.X.LineStyle.Color = ColorGray(8)

		p.X.Tick.Label.Color = ColorGray(8)
		p.X.Tick.LineStyle.Width = 0
	}

	{
		p.Y.Label.Text = item.Options.Y.Title
		p.Y.LineStyle.Width = 0
		p.Y.LineStyle.Color = ColorGray(8)

		p.Y.Tick.LineStyle.Width = 0
		p.Y.Tick.Label.Color = ColorGray(8)
	}

	/**
	{
		ticks := []plot.Tick{}
		for k, v := range item.Data.Labels {
			ticks = append(ticks, plot.Tick{
				Value: float64(k),
				Label: v,
			})
		}
		p.X.Tick.Marker = plot.ConstantTicks(ticks)
	}
	*/

	for k, v := range item.Data.Datasets {

		var data plotter.XYs
		for _, v2 := range v.Points {
			data = append(data, plotter.XY{X: v2.X, Y: v2.Y})
		}

		l, s, _ := plotter.NewLinePoints(data)

		l.Color = ColorTheme(k)
		l.Width = 2

		s.Color = ColorTheme(k)
		s.Shape = plotutil.Shape(k)

		p.Add(l, s)
		p.Legend.Add(v.Label, l, s)
	}

	{
		grid := plotter.NewGrid()
		grid.Vertical.Width = 0
		grid.Horizontal.Color = ColorGray(6)
		p.Add(grid)
	}

	{
		p.Legend.Position = "bottom"
	}

	return nil
}

func barRender(p *plot.Plot, item *hcapi.ChartEntry) error {

	legendN := len(item.Data.Labels)

	if legendN == 0 ||
		len(item.Data.Datasets) < 1 ||
		legendN != len(item.Data.Datasets[0].Values) {
		return errors.New("invalid datasets")
	}

	labels := []string{}
	for _, v := range item.Data.Labels {
		labels = append(labels, v)
	}
	p.NominalX(labels...)

	w := vg.Length(15)
	wMax := vg.Length(item.Options.WidthLength())
	wMax /= vg.Length(len(labels) * (len(item.Data.Datasets) + 1))
	if wMax < w {
		w = wMax
		if w < 1 {
			w = 1
		}
	}

	if item.Options.X.Title != "" {
		p.X.Label.Text = item.Options.X.Title
	}
	if item.Options.Y.Title != "" {
		//
		p.Y.Label.Text = item.Options.Y.Title
	}

	offsetK := vg.Length(0)
	if n := len(item.Data.Datasets); n > 1 {
		offsetK = (w * vg.Length(n)) / -2
		offsetK += w / 2
	}

	for k, v := range item.Data.Datasets {

		data := plotter.Values(v.Values)

		bcPlot, err := plotter.NewBarChart(data, w)
		if err != nil {
			continue
		}
		bcPlot.LineStyle.Width = vg.Length(0)
		bcPlot.Color = ColorTheme(k)
		bcPlot.Offset = offsetK + (vg.Length(k) * w)

		p.Add(bcPlot)
		p.Legend.Add(v.Label, bcPlot)
	}

	{
		grid := plotter.NewGrid()
		grid.Vertical.Width = 0
		grid.Horizontal.Color = ColorGray(6)
		p.Add(grid)
	}

	p.Legend.Position = "bottom"

	return nil
}

func histRender(p *plot.Plot, item *hcapi.ChartEntry) error {

	if len(item.Data.Datasets) < 1 || len(item.Data.Datasets[0].Values) < 1 {
		return errors.New("invalid datasets")
	}

	for k, v := range item.Data.Datasets {

		data := plotter.Values(v.Values)

		hPlot, err := plotter.NewHist(data, 100)
		if err != nil {
			return err
		}

		hPlot.FillColor = ColorTheme(k)
		hPlot.LineStyle.Width = 0

		p.Add(hPlot)
	}

	{
		grid := plotter.NewGrid()
		grid.Vertical.Width = 0
		grid.Horizontal.Color = ColorGray(6)
		p.Add(grid)
	}

	p.Legend.Position = "bottom"

	return nil
}

/*
func dumpPng(name string, item hcapi.ChartEntry, c chart.Chart) error {

	fp, err := os.OpenFile(name+".png", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(0)

	img := image.NewRGBA(image.Rect(0, 0, 800, 400))
	bg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff})
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)

	igr := imgg.AddTo(img, 0, 0, 800, 400, color.RGBA{0xff, 0xff, 0xff, 0xff}, nil, nil)
	c.Plot(igr)

	png.Encode(fp, img)

	return nil
}
*/

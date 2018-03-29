package visualization

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/sharath/lambo/database"
	"gopkg.in/mgo.v2"
	"image/color"
	"path"
	"strconv"
	"strings"
	"time"
)

func NewGraph(name, kind string, entries *mgo.Collection) string {
	var all []*database.MongoEntry
	entries.Find(nil).All(&all)
	var p *plot.Plot
	switch kind {
	case "timeusd":
		p = timeusd(name, all)
		break
	case "timebtc":
		p = timebtc(name, all)
		break
	case "timeeth":
		p = timeeth(name, all)
		break
	case "change1hr":
		p = change1hr(name, all)
		break
	}
	seed := make([]byte, 20)
	rand.Read(seed)
	seed = append(seed, []byte(time.Now().String())...)
	filename := path.Join("static", "graph", strings.Replace(base64.StdEncoding.EncodeToString(seed)+".png", "/", "<", -1))
	p.Save(vg.Inch*7, vg.Inch*7, filename)
	return filename
}

func timeusd(name string, all []*database.MongoEntry) *plot.Plot {
	pts := make(plotter.XYs, len(all))
	for i, e := range all {
		for _, t := range e.Tokens {
			if t.Name == name {
				x, _ := strconv.ParseFloat(t.LastUpdated, 63)
				y, _ := strconv.ParseFloat(t.PriceUsd, 63)
				pts[i].X = x
				pts[i].Y = y
			}
		}
	}
	p, _ := plot.New()
	p.Title.Text = name
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "PriceUSD"
	p.Add(plotter.NewGrid())

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
	p.Add(l)
	return p
}

func timebtc(name string, all []*database.MongoEntry) *plot.Plot {
	pts := make(plotter.XYs, len(all))
	for i, e := range all {
		for _, t := range e.Tokens {
			if t.Name == name {
				x, _ := strconv.ParseFloat(t.LastUpdated, 63)
				y, _ := strconv.ParseFloat(t.PriceBtc, 63)
				pts[i].X = x
				pts[i].Y = y
			}
		}
	}
	p, _ := plot.New()
	p.Title.Text = name
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "PriceBTC"
	p.Add(plotter.NewGrid())

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
	p.Add(l)
	return p
}

func timeeth(name string, all []*database.MongoEntry) *plot.Plot {
	pts := make(plotter.XYs, len(all))
	var eth = 0.1
	for i, e := range all {
		for _, t := range e.Tokens {
			if t.Name == "Ethereum" {
				eth, _ = strconv.ParseFloat(t.PriceUsd, 63)
			}
		}
		for _, t := range e.Tokens {
			if t.Name == name {
				x, _ := strconv.ParseFloat(t.LastUpdated, 63)
				y, _ := strconv.ParseFloat(t.PriceUsd, 63)
				pts[i].X = x
				pts[i].Y = y / eth
			}
		}
	}
	p, _ := plot.New()
	p.Title.Text = name
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "PriceETH"
	p.Add(plotter.NewGrid())

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
	p.Add(l)
	return p
}

func change1hr(name string, all []*database.MongoEntry) *plot.Plot {
	pts := make(plotter.XYs, len(all))
	for i, e := range all {
		for _, t := range e.Tokens {
			if t.Name == name {
				x, _ := strconv.ParseFloat(t.LastUpdated, 63)
				y, _ := strconv.ParseFloat(t.PercentChange1H, 63)
				pts[i].X = x
				pts[i].Y = y
			}
		}
	}
	p, _ := plot.New()
	p.Title.Text = name
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Hourly Change"
	p.Add(plotter.NewGrid())

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
	p.Add(l)
	return p
}

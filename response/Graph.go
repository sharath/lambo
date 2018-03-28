package response

import (
	"encoding/base64"
	"github.com/sharath/lambo/database"
	"gopkg.in/mgo.v2"
	"image/color"
	"path"
	"strconv"

	"crypto/rand"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"strings"
	"time"
)

// Graph is a wrapper for the graph url
type Graph struct {
	Data string `json:"data" bson:"data"`
}

// NewGraph makes a new graph and returns the url wrapper
func NewGraph(name string, entries *mgo.Collection) *Graph {
	var all []*database.MongoEntry
	entries.Find(nil).All(&all)
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
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l)
	seed := make([]byte, 20)
	rand.Read(seed)
	seed = append(seed, []byte(time.Now().String())...)
	filen := path.Join("static", "graph", strings.Replace(base64.StdEncoding.EncodeToString(seed)+".png", "/", "<", -1))
	p.Save(7*vg.Inch, 7*vg.Inch, filen)
	g := new(Graph)
	g.Data = filen
	return g
}

package response

import (
	"gopkg.in/mgo.v2"
	"github.com/sharath/lambo/database"
	"strconv"
	"image/color"
	"os"
	"path"
	"encoding/base64"

	"crypto/rand"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"io/ioutil"
)

type Graph struct {
	Data string `json:"data" bson:"data"`
}

func clear() {
	os.RemoveAll("gen")
}

func NewGraph(name string, entries *mgo.Collection) *Graph {
	var all []*database.MongoEntry
	entries.Find(nil).All(&all)
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
	os.Mkdir("gen", 0777)
	seed := make([]byte, 20)
	rand.Read(seed)
	filen := path.Join("gen", base64.StdEncoding.EncodeToString(seed)+".png")
	p.Save(4*vg.Inch, 4*vg.Inch, filen)
	c, _ := ioutil.ReadFile(filen)
	g := new(Graph)
	g.Data = base64.StdEncoding.EncodeToString(c)
	clear()
	return g
}

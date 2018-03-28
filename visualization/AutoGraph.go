package visualization

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/sharath/lambo/database"
	"github.com/sharath/lambo/poller"
	"gopkg.in/mgo.v2"
	"image/color"
	"path"
	"strconv"
	"time"
)

func AutoGraph(updater *poller.MongoUpdater, entries *mgo.Collection) {
	auto := func() {
		time.Sleep(time.Second)
		for range updater.P.Update {
			var all []*database.MongoEntry
			entries.Find(nil).All(&all)
			for i := 0; i < len(all); i++ {
				name := all[0].Tokens[i].Name
				pts := make(plotter.XYs, len(all))
				for _, e := range all {
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
				l, _ := plotter.NewLine(pts)
				l.LineStyle.Width = vg.Points(1)
				l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
				l.LineStyle.Color = color.RGBA{B: 255, A: 255}

				p.Add(l)
				filen := path.Join("static", "graph", name+".png")
				p.Save(7*vg.Inch, 7*vg.Inch, filen)
			}
		}
	}
	go auto()
}

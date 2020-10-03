package index

import (
	"html/template"
	"net/http"

	"github.com/diamondburned/l4d2lb/internal/durafmt"
	"github.com/diamondburned/l4d2lb/pages"
	"github.com/diamondburned/l4d2lb/pages/components/footer"
	"github.com/diamondburned/l4d2lb/pages/components/header"
	"github.com/diamondburned/l4d2lb/pages/components/loading"
	"github.com/diamondburned/l4d2lb/pages/index/leaderboard"
	"github.com/diamondburned/l4d2lb/pages/internal/flushw"
	"github.com/diamondburned/l4d2lb/pages/internal/pararender"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi"
)

func Mount(state *pages.RenderState) http.Handler {
	tpl := state.Template("index", map[string]interface{}{
		"human": humanize.Comma,
		"dural": durafmt.Long,
		"duras": durafmt.Short,

		"maps":    func(tpl ParallelRender) template.HTML { return tpl.TopMaps.Render() },
		"players": func(tpl ParallelRender) template.HTML { return tpl.TopPlayers.Render() },
		"stats":   func(tpl ParallelRender) template.HTML { return tpl.Statistics.Render() },
	})

	tpl.AddComponent(header.Path)
	tpl.AddComponent(footer.Path)
	tpl.AddComponent(loading.Header) // loading-header
	tpl.AddComponent(loading.Footer) // loading-footer
	tpl.AddComponent("index/index-maps.html")
	tpl.AddComponent("index/index-stats.html")
	tpl.AddComponent("index/index-players.html")

	r := chi.NewMux()
	r.Use(flushw.Middleware)
	r.Get("/", renderPage(tpl))
	r.Mount("/leaderboard", leaderboard.Mount(state))

	return r
}

type ParallelRender struct {
	*pages.Template
	TopMaps    pararender.Task
	TopPlayers pararender.Task
	Statistics pararender.Task
}

func renderPage(tpl *pages.Template) http.HandlerFunc {
	var (
		topMaps    = tpl.HTMLComponentRenderer("index-maps")
		topPlayers = tpl.HTMLComponentRenderer("index-players")
		statistics = tpl.HTMLComponentRenderer("index-stats")
	)

	return func(w http.ResponseWriter, r *http.Request) {
		var render = ParallelRender{
			Template:   tpl,
			TopMaps:    pararender.EmptyTask(topMaps),
			TopPlayers: pararender.EmptyTask(topPlayers),
			Statistics: pararender.EmptyTask(statistics),
		}

		rv := tpl.WithContext(r.Context())

		go func() { render.TopMaps.Send(rv.TopMaps(5)) }()
		go func() { render.TopPlayers.Send(rv.TopPlayers("playtime", 5)) }()
		go func() { render.Statistics.Send(rv.Statistics()) }()

		tpl.Execute(w, render)
	}
}

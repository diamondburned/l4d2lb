package index

import (
	"html/template"
	"net/http"

	"github.com/diamondburned/l4d2lb/internal/flushw"
	"github.com/diamondburned/l4d2lb/pages"
	"github.com/diamondburned/l4d2lb/pages/components/errbox"
	"github.com/diamondburned/l4d2lb/pages/components/header"
	"github.com/diamondburned/l4d2lb/pages/components/loading"
	"github.com/diamondburned/l4d2lb/pages/index/leaderboard"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func Mount(state *pages.RenderState) http.Handler {
	tpl := state.Template("index", map[string]interface{}{
		"maps": func(parent *pages.Template) template.HTML {
			maps, err := parent.Top10Maps()
			if err != nil {
				return errbox.RenderHTML(errors.Wrap(err, "failed to get maps"))
			}

			return parent.RenderHTMLComponent("index-maps", maps)
		},

		"stats": func(parent *pages.Template) template.HTML {
			stats, err := parent.Statistics()
			if err != nil {
				return errbox.RenderHTML(errors.Wrap(err, "failed to get stats"))
			}

			return parent.RenderHTMLComponent("index-stats", stats)
		},

		"players": func(parent *pages.Template) template.HTML {
			players, err := parent.Leaderboard(10, 0)
			if err != nil {
				return errbox.RenderHTML(errors.Wrap(err, "failed to get leaderboard"))
			}

			return parent.RenderHTMLComponent("index-players", players)
		},
	})

	tpl.AddComponent(header.Path)
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

func renderPage(tpl *pages.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, tpl)
	}
}

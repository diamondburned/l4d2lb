package leaderboard

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/diamondburned/l4d2lb/pages"
	"github.com/diamondburned/l4d2lb/pages/components/errbox"
	"github.com/diamondburned/l4d2lb/pages/components/header"
	"github.com/diamondburned/l4d2lb/pages/components/loading"
	"github.com/diamondburned/l4d2lb/stats"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func Mount(state *pages.RenderState) http.Handler {
	tpl := state.Template("index/leaderboard", map[string]interface{}{
		"dec": func(i int) int { return i - 1 },
		"inc": func(i int) int { return i + 1 },

		"leaderboard": renderLeaderboard,
	})
	tpl.AddComponent(header.Path)
	tpl.AddComponent(loading.Header) // loading-header
	tpl.AddComponent(loading.Footer) // loading-footer
	tpl.AddComponent("index/leaderboard/leaderboard-table.html")

	r := chi.NewMux()
	r.Get("/", renderPage(tpl))

	return r
}

func getPageNumber(r *http.Request) int {
	r.ParseForm()

	v := r.FormValue("p")
	i, err := strconv.Atoi(v)
	if err != nil {
		return 1
	}

	if i < 1 {
		return 1
	}

	return i
}

const PlayerPerPage = 100

type Info struct {
	*pages.Template
	Page  int
	Query string
}

func renderPage(tpl *pages.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pageNum = getPageNumber(r)
		var query = r.FormValue("q")

		tpl.Execute(w, Info{
			Template: tpl,
			Page:     pageNum,
			Query:    query,
		})
	}
}

type Leaderboard struct {
	Info
	stats.PlayerResults
}

func renderLeaderboard(info Info) template.HTML {
	var results stats.PlayerResults
	var err error

	if info.Query != "" {
		results, err = info.SearchLeaderboard(info.Query, PlayerPerPage, info.Page-1)
	} else {
		results, err = info.Leaderboard(PlayerPerPage, info.Page-1)
	}

	if err != nil {
		return errbox.RenderHTML(errors.Wrap(err, "failed to get leaderboard"))
	}

	return info.RenderHTMLComponent("leaderboard-table", Leaderboard{
		Info:          info,
		PlayerResults: results,
	})
}

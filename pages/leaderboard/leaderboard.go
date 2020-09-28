package leaderboard

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/diamondburned/l4d2lb/pages"
	"github.com/diamondburned/l4d2lb/pages/errpage"
	"github.com/diamondburned/l4d2lb/stats"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var tmpl *template.Template

func init() {
	tmpl = pages.Template("leaderboard")
}

func Mount(statsDB *stats.Database) http.Handler {
	r := chi.NewMux()
	r.Get("/", renderPage(statsDB))

	return r
}

func getPageNumber(r *http.Request) int {
	v := r.FormValue("p")
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	if i < 1 {
		return 0
	}

	return i - 1
}

const PlayerPerPage = 100

func renderPage(statsDB *stats.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		players, err := statsDB.Leaderboard(PlayerPerPage, getPageNumber(r))
		if err != nil {
			errpage.RenderError(w, 400, errors.Wrap(err, "failed to get leaderboard"))
			return
		}

		pages.Execute(tmpl, w, players)
	}
}

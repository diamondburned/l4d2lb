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
	tmpl = pages.Template("leaderboard", map[string]interface{}{
		"dec": func(i int) int { return i - 1 },
		"inc": func(i int) int { return i + 1 },
	})
}

func Mount(statsDB *stats.Database) http.Handler {
	r := chi.NewMux()
	r.Get("/", renderPage(statsDB))

	return r
}

func getPageNumber(r *http.Request) int {
	r.ParseForm()

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

type PageInfo struct {
	stats.PlayerResults
	Page int
}

func renderPage(statsDB *stats.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pageNum = getPageNumber(r)

		results, err := statsDB.Leaderboard(PlayerPerPage, pageNum)
		if err != nil {
			errpage.RenderError(w, 400, errors.Wrap(err, "failed to get leaderboard"))
			return
		}

		pages.Execute(tmpl, w, PageInfo{
			PlayerResults: results,
			Page:          pageNum + 1,
		})
	}
}

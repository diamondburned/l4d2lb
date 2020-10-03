package leaderboard

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/diamondburned/l4d2lb/pages"
	"github.com/diamondburned/l4d2lb/pages/components/footer"
	"github.com/diamondburned/l4d2lb/pages/components/header"
	"github.com/diamondburned/l4d2lb/pages/components/loading"
	"github.com/diamondburned/l4d2lb/pages/internal/pararender"
	"github.com/diamondburned/l4d2lb/stats"
	"github.com/go-chi/chi"
)

func Mount(state *pages.RenderState) http.Handler {
	tpl := state.Template("index/leaderboard", map[string]interface{}{
		"dec": func(i int) int { return i - 1 },
		"inc": func(i int) int { return i + 1 },

		"leaderboard": func(info Info) template.HTML { return info.Leaderboard.Render() },
	})
	tpl.AddComponent(header.Path)
	tpl.AddComponent(footer.Path)
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
	Leaderboard pararender.Task
}

type Leaderboard struct {
	stats.PlayerResults
	Page  int
	Sort  string
	Query string
}

func renderPage(tpl *pages.Template) http.HandlerFunc {
	var leaderboard = tpl.HTMLComponentRenderer("leaderboard-table")

	return func(w http.ResponseWriter, r *http.Request) {
		var info = Info{
			Template:    tpl,
			Leaderboard: pararender.EmptyTask(leaderboard),
		}

		go func() {
			rv := tpl.WithContext(r.Context())
			lb := Leaderboard{
				Page:  getPageNumber(r),
				Sort:  r.FormValue("s"),
				Query: r.FormValue("q"),
			}

			v, err := rv.Leaderboard(lb.Query, lb.Sort, PlayerPerPage, lb.Page-1)
			if err != nil {
				info.Leaderboard.Send(nil, err)
				return
			}
			lb.PlayerResults = v
			info.Leaderboard.Send(lb, nil)
		}()

		tpl.Execute(w, info)
	}
}

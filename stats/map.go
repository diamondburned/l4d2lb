package stats

import (
	"time"

	"github.com/pkg/errors"
)

type Map struct {
	Name     string `db:"name"`
	Kills    int64  `db:"kills"`
	Points   int64  `db:"points"`
	Playtime int64  `db:"playtime"` // minutes
}

func (m Map) ConvertPlaytime() time.Duration {
	return time.Duration(m.Playtime) * time.Minute
}

func (rv *ReadView) TopMaps(limit int) (maps []Map, err error) {
	const query = `
		SELECT
			name,
	    	kills_nor + kills_adv + kills_exp AS kills,
	    	points_nor + points_adv + points_exp AS points,
	    	playtime_nor + playtime_adv + playtime_exp AS playtime
		FROM maps ORDER BY playtime DESC LIMIT ?
	`

	if err := rv.SelectContext(rv.ctx, &maps, query, limit); err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	return
}

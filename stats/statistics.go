package stats

import (
	"time"

	"github.com/pkg/errors"
)

// Statistics contains general statistics.
type Statistics struct {
	PlayersServed  int64 `db:"players_served"`
	TotalKills     int64 `db:"total_kills"`
	TotalPoints    int64 `db:"total_points"`
	TotalPlaytime  int64 `db:"total_playtime"`
	TotalHeadshots int64 `db:"total_headshots"`
}

func (s Statistics) ConvertTotalPlaytime() time.Duration {
	return time.Duration(s.TotalPlaytime) * time.Minute
}

func (rv *ReadView) Statistics() (*Statistics, error) {
	const query = `
		SELECT
			COUNT(*) AS players_served,
			SUM(kills) AS total_kills,
			SUM(points) AS total_points,
			SUM(playtime) AS total_playtime,
			SUM(headshots) AS total_headshots
		FROM players
	`

	var stats Statistics

	if err := rv.GetContext(rv.ctx, &stats, query); err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	return &stats, nil
}

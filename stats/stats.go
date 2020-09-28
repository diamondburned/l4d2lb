package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/MrWaggel/gosteamconv"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type SteamID string

// ToSteamID64 converts the steam ID to its ID64 variant.
func (id SteamID) ToSteamID64() int64 {
	i, _ := gosteamconv.SteamStringToInt64(string(id))
	return i
}

// URL returns the link to the user's profile.
func (id SteamID) URL() string {
	return fmt.Sprintf("http://steamcommunity.com/profiles/%d", id.ToSteamID64())
}

// Database represents a database.
type Database struct {
	*sqlx.DB
}

func Connect(addr string) (*Database, error) {
	d, err := sqlx.Open("mysql", addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db")
	}

	d.SetConnMaxIdleTime(3 * time.Minute)
	d.SetMaxOpenConns(10)
	d.SetMaxIdleConns(10)

	return &Database{d}, nil
}

var playerQuery = fmt.Sprintf(
	"SELECT %s FROM players ORDER BY points DESC LIMIT ? OFFSET ?",
	strings.Join(playerColumns(), ","),
)

// Leaderboard fetches the leaderboard. The page count is 0-indexed.
func (d *Database) Leaderboard(count, page int) ([]Player, error) {
	page *= count
	r, err := d.DB.Queryx(playerQuery, count, page)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	defer r.Close()

	var players []Player

	for r.Next() {
		var player Player

		if err := r.StructScan(&player); err != nil {
			return nil, errors.Wrap(err, "failed to scan to player")
		}

		players = append(players, player)
	}

	if err := r.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to finish scanning")
	}

	return players, nil
}

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

type PlayerResults struct {
	Players []Player
	HasMore bool
}

var noResults = PlayerResults{}

var playerQuery = fmt.Sprintf(
	"SELECT %s FROM players ORDER BY points DESC LIMIT ? OFFSET ?",
	strings.Join(playerColumns(), ","),
)

// Leaderboard fetches the leaderboard. The page count is 0-indexed.
func (d *Database) Leaderboard(count, page int) (PlayerResults, error) {
	return d.queryPlayers(playerQuery, count, count+1, page*count)
}

var searchQuery = fmt.Sprintf(
	"SELECT %s FROM players WHERE name LIKE ? ORDER BY points DESC LIMIT ? OFFSET ?",
	strings.Join(playerColumns(), ","),
)

var queryEscaper = strings.NewReplacer(
	"%", `\%`,
	`\`, `\\`,
)

// Search searches the leaderboard for players. The page count is 0-indexed.
func (d *Database) Search(queryString string, count, page int) (PlayerResults, error) {
	queryString = queryEscaper.Replace(queryString)
	return d.queryPlayers(searchQuery, count, queryString, count+1, page*count)
}

func (d *Database) queryPlayers(query string, lim int, v ...interface{}) (PlayerResults, error) {
	r, err := d.DB.Queryx(query, v...)
	if err != nil {
		return noResults, errors.Wrap(err, "failed to query")
	}

	defer r.Close()

	var results PlayerResults

	for i := 0; r.Next(); i++ {
		var player Player

		if err := r.StructScan(&player); err != nil {
			return noResults, errors.Wrap(err, "failed to scan to player")
		}

		if i == lim {
			results.HasMore = true
		} else {
			results.Players = append(results.Players, player)
		}
	}

	if err := r.Err(); err != nil {
		return noResults, errors.Wrap(err, "failed to finish scanning")
	}

	return results, nil
}

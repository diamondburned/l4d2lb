package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/MrWaggel/gosteamconv"
	"github.com/diamondburned/l4d2lb/internal/rfsql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

// Player contains some statistics of each player.
type Player struct {
	SteamID           SteamID `db:"steamid"`
	Name              string  `db:"name"`
	Kills             int64   `db:"kills"`
	Points            int64   `db:"points"`
	MeleeKills        int64   `db:"melee_kills"`
	Headshots         int64   `db:"headshots"`
	AwardFriendlyfire int64   `db:"award_friendlyfire"`
	AwardFincap       int64   `db:"award_fincap"`
}

// URL returns the link to the player's profile.
func (p Player) URL() string {
	return p.SteamID.URL()
}

type PlayerResults struct {
	Players []Player
	HasMore bool
}

var noResults = PlayerResults{}

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

	d.SetConnMaxIdleTime(30 * time.Second)
	d.SetMaxOpenConns(4)
	d.SetMaxIdleConns(4)

	return &Database{d}, nil
}

var playerQuery = fmt.Sprintf(
	"SELECT %s FROM players ORDER BY points DESC LIMIT ? OFFSET ?",
	strings.Join(rfsql.Columns(Player{}), ","),
)

// Leaderboard fetches the leaderboard. The page count is 0-indexed.
func (d *Database) Leaderboard(count, page int) (PlayerResults, error) {
	return d.queryPlayers(playerQuery, count, count+1, page*count)
}

var searchQuery = fmt.Sprintf(
	"SELECT %s FROM players WHERE name LIKE ? ORDER BY points DESC LIMIT ? OFFSET ?",
	strings.Join(rfsql.Columns(Player{}), ","),
)

var queryEscaper = strings.NewReplacer(
	"%", `\%`,
	`\`, `\\`,
)

// SearchLeaderboard searches the leaderboard for players. The page count is
// 0-indexed.
func (d *Database) SearchLeaderboard(queryString string, count, page int) (PlayerResults, error) {
	queryString = "%" + queryEscaper.Replace(queryString) + "%"
	return d.queryPlayers(searchQuery, count, queryString, count+1, page*count)
}

func (d *Database) queryPlayers(query string, lim int, v ...interface{}) (PlayerResults, error) {
	if lim > 100 {
		return noResults, errors.New("limit too high")
	}

	var results PlayerResults

	if err := d.DB.Select(&results.Players, query, v...); err != nil {
		return noResults, errors.Wrap(err, "failed to query")
	}

	if len(results.Players) > lim {
		results.Players = results.Players[:lim]
		results.HasMore = true
	}

	return results, nil
}

type Map struct {
	Name     string `db:"name"`
	Kills    int    `db:"kills"`
	Points   int    `db:"points"`
	Playtime int    `db:"playtime"` // minutes
}

func (d *Database) Top10Maps() (maps []Map, err error) {
	const query = `
		SELECT
			name,
	    	kills_nor + kills_adv + kills_exp AS kills,
	    	points_nor + points_adv + points_exp AS points,
	    	playtime_nor + playtime_adv + playtime_exp AS playtime
		FROM maps ORDER BY playtime DESC LIMIT 10
	`

	if err := d.DB.Select(&maps, query); err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	return
}

// Statistics contains general statistics.
type Statistics struct {
	PlayersServed  int `db:"players_served"`
	TotalKills     int `db:"total_kills"`
	TotalPoints    int `db:"total_points"`
	TotalPlaytime  int `db:"total_playtime"`
	TotalHeadshots int `db:"total_headshots"`
}

func (d *Database) Statistics() (*Statistics, error) {
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

	if err := d.DB.Get(&stats, query); err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	return &stats, nil
}

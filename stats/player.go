package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/MrWaggel/gosteamconv"
	"github.com/diamondburned/l4d2lb/internal/rfsql"
	"github.com/pkg/errors"
)

// Player contains some statistics of each player.
type Player struct {
	Index             int     `db:"ix"`
	SteamID           SteamID `db:"steamid"`
	Name              string  `db:"name"`
	Kills             int64   `db:"kills"`
	Points            int64   `db:"points"`
	Playtime          int64   `db:"playtime"`
	MeleeKills        int64   `db:"melee_kills"`
	Headshots         int64   `db:"headshots"`
	AwardFriendlyfire int64   `db:"award_friendlyfire"`
	AwardFincap       int64   `db:"award_fincap"`
}

func (p Player) ConvertPlaytime() time.Duration {
	return time.Duration(p.Playtime) * time.Minute
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

// TopPlayers returns the top players from the leaderboard. It does not
// paginate.
func (rv *ReadView) TopPlayers(sort string, count int) ([]Player, error) {
	r, err := rv.Leaderboard("", sort, count, 0)
	if err != nil {
		return nil, err
	}
	return r.Players, nil
}

const (
	playerSuffix = "LIMIT ? OFFSET ?"
	// UTF8_GENERAL_CI for case-insensitivity.
	searchSuffix = "WHERE CONVERT(name USING utf8mb4) COLLATE utf8mb4_general_ci LIKE ? " + playerSuffix
)

var queryEscaper = strings.NewReplacer(
	"%", `\%`,
	`\`, `\\`,
)

var playerColumnsNoIndex = strings.Join(rfsql.Columns(Player{})[1:], ",")

func buildPlayerQuery(isSearch bool, sort string) (query string) {
	// Constants used in Sprintf.
	const sflimit = "LIMIT ? OFFSET ?"
	const fstring = "SELECT ROW_NUMBER() OVER(ORDER BY %s DESC) ix, %s FROM players ORDER BY ix"
	const fstrlim = fstring + " " + sflimit
	const fsearch = "SELECT * FROM (" + fstring + ") AS sorted " +
		"WHERE CONVERT(sorted.name USING utf8mb4) COLLATE utf8mb4_general_ci LIKE ? " + sflimit

	if isSearch {
		query = fmt.Sprintf(fsearch, sort, playerColumnsNoIndex)
	} else {
		query = fmt.Sprintf(fstrlim, sort, playerColumnsNoIndex)
	}

	return
}

// Leaderboard searches the leaderboard for players. The page count is
// zero-indexed.
func (rv *ReadView) Leaderboard(search, sort string, count, page int) (PlayerResults, error) {
	page *= count

	// Validate the sorted column.
	if sort == "" {
		sort = "points"
	} else {
		if !rfsql.IsValidColumn(Player{}, sort) {
			return noResults, fmt.Errorf("unknown column %q", sort)
		}
	}

	if search != "" {
		search = "%" + queryEscaper.Replace(search) + "%"
		return rv.queryPlayers(buildPlayerQuery(true, sort), count, search, count+1, page)
	}

	return rv.queryPlayers(buildPlayerQuery(false, sort), count, count+1, page)
}

func (rv *ReadView) queryPlayers(query string, lim int, v ...interface{}) (PlayerResults, error) {
	if lim > 100 {
		return noResults, errors.New("limit too high")
	}

	var results PlayerResults

	if err := rv.SelectContext(rv.ctx, &results.Players, query, v...); err != nil {
		return noResults, errors.Wrap(err, "failed to query")
	}

	if len(results.Players) > lim {
		results.Players = results.Players[:lim]
		results.HasMore = true
	}

	return results, nil
}

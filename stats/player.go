package stats

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

func playerColumns() []string {
	return []string{
		"steamid",
		"name",
		"kills",
		"points",
		"melee_kills",
		"headshots",
		"award_friendlyfire",
		"award_fincap",
	}
}

// URL returns the link to the player's profile.
func (p Player) URL() string {
	return p.SteamID.URL()
}

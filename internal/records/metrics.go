package records

import (
	"fmt"
	"sort"

	"example.com/othello/internal/store"
)

type OutcomeStats struct {
	Total     int
	Published int
	Archived  int
	BlackWins int
	WhiteWins int
	Draws     int
	ByRegion  map[string]int
}

func (r *Repository) Stats(filter Filter) (OutcomeStats, error) {
	values, err := r.Search(filter)
	if err != nil {
		return OutcomeStats{}, err
	}
	stats := OutcomeStats{Total: len(values), ByRegion: make(map[string]int)}
	for _, record := range values {
		stats.ByRegion[record.Region]++
		if record.Published {
			stats.Published++
		}
		if record.Archived {
			stats.Archived++
		}
		switch record.Winner {
		case "B":
			stats.BlackWins++
		case "W":
			stats.WhiteWins++
		default:
			stats.Draws++
		}
	}
	return stats, nil
}

func (stats OutcomeStats) RegionNames() []string {
	values := make([]string, 0, len(stats.ByRegion))
	for region := range stats.ByRegion {
		values = append(values, region)
	}
	sort.Strings(values)
	return values
}

func (stats OutcomeStats) String() string {
	return fmt.Sprintf("total=%d published=%d archived=%d black=%d white=%d draws=%d", stats.Total, stats.Published, stats.Archived, stats.BlackWins, stats.WhiteWins, stats.Draws)
}

func RecordScore(record store.Record) int {
	return record.BlackScore + record.WhiteScore
}

func (stats OutcomeStats) PublishedRecords() int {
	return stats.Published
}

func (stats OutcomeStats) IsEmpty() bool {
	return stats.Total == 0
}

func (stats OutcomeStats) RegionCount(region string) int {
	return stats.ByRegion[region]
}

func (stats OutcomeStats) WinnerBreakdown() map[string]int {
	return map[string]int{"black": stats.BlackWins, "white": stats.WhiteWins, "draw": stats.Draws}
}

func (stats OutcomeStats) Decided() int {
	return stats.BlackWins + stats.WhiteWins
}

func (stats OutcomeStats) StatusSummary() string {
	return fmt.Sprintf("published=%d archived=%d", stats.Published, stats.Archived)
}

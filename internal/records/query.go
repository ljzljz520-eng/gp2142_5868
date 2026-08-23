package records

import (
	"fmt"
	"sort"
	"strings"

	"example.com/othello/internal/store"
)

type History struct {
	Records []store.Record
	Total   int
	Regions []string
}

func (r *Repository) History(filter Filter) (History, error) {
	records, err := r.Search(filter)
	if err != nil {
		return History{}, err
	}
	regions := make(map[string]struct{})
	for _, record := range records {
		regions[record.Region] = struct{}{}
	}
	values := make([]string, 0, len(regions))
	for region := range regions {
		values = append(values, region)
	}
	sort.Strings(values)
	return History{Records: records, Total: len(records), Regions: values}, nil
}

func FormatRecord(record store.Record) string {
	state := record.Status
	if record.Archived {
		state = "archived"
	}
	return fmt.Sprintf("%s [%s] %s B:%d W:%d notes:%s", record.ID, state, strings.TrimSpace(record.Summary), record.BlackScore, record.WhiteScore, strings.TrimSpace(record.Notes))
}

func SummarizeHistory(history History) string {
	if history.Total == 0 {
		return "no matching records"
	}
	lines := make([]string, 0, history.Total+1)
	lines = append(lines, fmt.Sprintf("%d record(s)", history.Total))
	for _, record := range history.Records {
		lines = append(lines, FormatRecord(record))
	}
	return strings.Join(lines, "\n")
}

func (r *Repository) FindFirst(filter Filter) (store.Record, error) {
	records, err := r.Search(filter)
	if err != nil {
		return store.Record{}, err
	}
	if len(records) == 0 {
		return store.Record{}, fmt.Errorf("no record matches region %q", filter.Region)
	}
	return records[0], nil
}

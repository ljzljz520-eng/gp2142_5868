package service

import (
	"fmt"
	"strings"

	"example.com/othello/internal/records"
	"example.com/othello/internal/store"
)

func DescribeRecord(record store.Record) string {
	parts := []string{record.ID, record.Status, record.Region, record.Summary}
	if record.Published {
		parts = append(parts, "published")
	}
	if record.Archived {
		parts = append(parts, "archived")
	}
	return strings.Join(parts, " | ")
}

func DescribeSearch(history records.History) string {
	if history.Total == 0 {
		return "no records"
	}
	return fmt.Sprintf("%d result(s): %s", history.Total, records.SummarizeHistory(history))
}

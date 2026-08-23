package cli

import (
	"errors"
	"fmt"
	"strings"

	"example.com/othello/internal/game"
	"example.com/othello/internal/importer"
	"example.com/othello/internal/records"
	"example.com/othello/internal/service"
	"example.com/othello/internal/store"
)

func runStart(args []string) error {
	if len(args) != 2 {
		return errors.New("start needs game-id and mode")
	}
	mode, err := game.ParseMode(args[1])
	if err != nil {
		return err
	}
	state, err := game.New(args[0], mode)
	if err != nil {
		return err
	}
	fmt.Println(renderBoard(state))
	fmt.Println("legal moves:", coordinateList(state.LegalMoves()))
	return nil
}

func runPlay(args []string) error {
	if len(args) < 3 {
		return errors.New("play needs game-id, mode, and at least one coordinate")
	}
	mode, err := game.ParseMode(args[1])
	if err != nil {
		return err
	}
	state, err := game.ParseMoves(args[0], mode, args[2:])
	if err != nil {
		return err
	}
	fmt.Println(renderBoard(state))
	fmt.Println("transcript:", state.Transcript())
	return nil
}

func runHistory(args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errors.New("history needs database and optional region")
	}
	database, application, closeStore, err := openService(args[0])
	if err != nil {
		return err
	}
	defer closeStore()
	filter := records.Filter{}
	if len(args) == 2 {
		filter.Region = args[1]
	}
	history, err := application.Search(filter)
	if err != nil {
		return err
	}
	for _, record := range history.Records {
		fmt.Println(renderRecord(record))
	}
	stats, err := application.OutcomeStats(filter)
	if err != nil {
		return err
	}
	fmt.Println(stats.String())
	_ = database
	return nil
}

func runImport(args []string) error {
	if len(args) < 5 {
		return errors.New("import needs database, game-id, mode, region, and moves")
	}
	mode, err := game.ParseMode(args[2])
	if err != nil {
		return err
	}
	_, application, closeStore, err := openService(args[0])
	if err != nil {
		return err
	}
	defer closeStore()
	request := importer.Request{ID: args[1], Mode: mode, Region: args[3], Source: "command-line", Content: strings.Join(args[4:], " ")}
	result, err := application.Import(request, "cli")
	if err != nil {
		return err
	}
	fmt.Println(result.Report)
	return nil
}

func runDemo(args []string) error {
	if len(args) != 1 {
		return errors.New("demo needs database")
	}
	_, application, closeStore, err := openService(args[0])
	if err != nil {
		return err
	}
	defer closeStore()
	state, err := service.FinishedFixture("demo-1")
	if err != nil {
		return err
	}
	lifecycle, err := application.CompleteLifecycle(state, "beginner", "cli", "coach")
	if err != nil {
		return err
	}
	record := lifecycle.Record
	if _, err := application.ArchiveRecord(record.ID, "coach", "demo complete"); err != nil {
		return err
	}
	fmt.Println(renderRecord(record))
	fmt.Println("archived", record.ID)
	return nil
}

func openService(path string) (*store.Store, *service.Service, func(), error) {
	database, err := store.Open(path)
	if err != nil {
		return nil, nil, func() {}, err
	}
	application, err := service.New(database)
	if err != nil {
		database.Close()
		return nil, nil, func() {}, err
	}
	return database, application, func() { _ = database.Close() }, nil
}

package importer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"example.com/othello/internal/game"
	"example.com/othello/internal/store"
)

type Request struct {
	ID      string
	Mode    game.Mode
	Region  string
	Source  string
	Content string
}

type Result struct {
	Game       game.Game
	Attachment store.Attachment
	Moves      int
	Report     string
}

func Parse(request Request) (game.Game, string, error) {
	if strings.TrimSpace(request.ID) == "" {
		return game.Game{}, "", errors.New("import id is required")
	}
	if strings.TrimSpace(request.Content) == "" {
		return game.Game{}, "", errors.New("棋谱内容不能为空")
	}
	tokens := Tokenize(request.Content)
	if err := ValidateTokens(tokens); err != nil {
		return game.Game{}, "", err
	}
	moves := make([]string, 0, len(tokens))
	for _, token := range tokens {
		value, err := moveToken(token)
		if err != nil {
			return game.Game{}, "", err
		}
		moves = append(moves, value)
	}
	result, err := game.ParseMoves(request.ID, request.Mode, moves)
	if err != nil {
		return game.Game{}, "", fmt.Errorf("棋谱校验失败: %w", err)
	}
	if err := ValidateGame(result); err != nil {
		return game.Game{}, "", err
	}
	return result, fmt.Sprintf("validated %d move(s), status=%s, score=%d:%d", len(result.Moves), result.Status, result.BlackScore, result.WhiteScore), nil
}

func moveToken(token string) (string, error) {
	clean := strings.TrimSpace(token)
	if strings.EqualFold(clean, "pass") {
		return "pass", nil
	}
	dot := strings.IndexByte(clean, '.')
	if dot >= 0 {
		number, err := strconv.Atoi(clean[:dot])
		if err != nil || number < 1 {
			return "", fmt.Errorf("棋谱编号 %q 无效", clean)
		}
		clean = clean[dot+1:]
	}
	if clean == "" {
		return "", errors.New("棋谱包含空落子")
	}
	return clean, nil
}

func SaveAttachment(database *store.Store, request Request, parsed game.Game, report string) (store.Attachment, error) {
	if database == nil {
		return store.Attachment{}, errors.New("store is required")
	}
	attachment := store.Attachment{ID: "attachment-" + request.ID, RecordID: "record-" + request.ID, Name: request.Source, Kind: "othello-transcript", Content: request.Content}
	if strings.TrimSpace(attachment.Name) == "" {
		attachment.Name = "manual-import"
	}
	if parsed.ID == "" || report == "" {
		return store.Attachment{}, errors.New("parsed game and report are required")
	}
	if err := database.SaveAttachment(attachment); err != nil {
		return store.Attachment{}, err
	}
	return attachment, nil
}

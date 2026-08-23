package importer

import (
	"fmt"
	"strings"

	"example.com/othello/internal/game"
)

func Tokenize(content string) []string {
	return strings.FieldsFunc(content, func(r rune) bool { return r == ',' || r == ';' || r == '\n' || r == '\t' || r == ' ' })
}

func ValidateTokens(tokens []string) error {
	if len(tokens) == 0 {
		return fmt.Errorf("棋谱至少需要一个落子")
	}
	for index, token := range tokens {
		if strings.TrimSpace(token) == "" {
			return fmt.Errorf("棋谱第 %d 项为空", index+1)
		}
		if _, err := moveToken(token); err != nil {
			return err
		}
	}
	return nil
}

func ValidateGame(parsed game.Game) error {
	if parsed.ID == "" {
		return fmt.Errorf("game id is required")
	}
	if parsed.BlackScore < 0 || parsed.WhiteScore < 0 {
		return fmt.Errorf("score cannot be negative")
	}
	if parsed.Status != game.StatusActive && parsed.Status != game.StatusEnded {
		return fmt.Errorf("unknown game status %q", parsed.Status)
	}
	return nil
}

func NormalizedContent(content string) string {
	return strings.Join(Tokenize(content), " ")
}

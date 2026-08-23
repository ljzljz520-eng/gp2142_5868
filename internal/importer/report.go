package importer

import (
	"fmt"
	"strings"

	"example.com/othello/internal/game"
	"example.com/othello/internal/store"
)

type Report struct {
	AttachmentID string
	GameID       string
	Valid        bool
	MoveCount    int
	BlackScore   int
	WhiteScore   int
	Status       string
	Message      string
}

func BuildReport(attachment store.Attachment, parsed game.Game, validationMessage string) Report {
	message := strings.TrimSpace(validationMessage)
	if message == "" {
		message = "validated"
	}
	return Report{AttachmentID: attachment.ID, GameID: parsed.ID, Valid: true, MoveCount: len(parsed.Moves), BlackScore: parsed.BlackScore, WhiteScore: parsed.WhiteScore, Status: string(parsed.Status), Message: message}
}

func (r Report) String() string {
	return fmt.Sprintf("%s: %s, moves=%d, score=%d:%d, attachment=%s", r.GameID, r.Message, r.MoveCount, r.BlackScore, r.WhiteScore, r.AttachmentID)
}

func ValidateAttachment(attachment store.Attachment) error {
	if strings.TrimSpace(attachment.ID) == "" || strings.TrimSpace(attachment.RecordID) == "" {
		return fmt.Errorf("attachment identity is incomplete")
	}
	if strings.TrimSpace(attachment.Content) == "" {
		return fmt.Errorf("attachment %s is empty", attachment.ID)
	}
	return nil
}

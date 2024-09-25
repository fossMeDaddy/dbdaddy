package migrationsLib

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/types"
)

func GetInfoTextFromDiff(actions []types.MigAction) string {
	infoText := "## Changes summary:"
	infoText += fmt.Sprintln()

	for _, action := range actions {
		entityIdMaxI := min(len(action.Entity.Id), 3)

		infoText += fmt.Sprintf(
			"- %s %s %v",
			action.ActionType,
			action.Entity.Type,
			action.Entity.Id[:entityIdMaxI],
		)

		infoText += fmt.Sprintln()
	}

	return infoText
}

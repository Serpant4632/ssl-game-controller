package statemachine

import "github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"

// CreateCommandChange creates a change with a new command
func CreateCommandChange(command *state.Command) *Change {
	return &Change{
		Change: &Change_NewCommandChange{
			NewCommandChange: &Change_NewCommand{
				Command: command,
			},
		},
	}
}

// CreateGameEventChange creates a change with a new game event
func CreateGameEventChange(eventType state.GameEvent_Type, event *state.GameEvent) *Change {
	event.Type = &eventType
	event.Origin = []string{changeOriginStateMachine}
	return &Change{
		Origin: &changeOriginStateMachine,
		Change: &Change_AddGameEventChange{
			AddGameEventChange: &Change_AddGameEvent{
				GameEvent: event,
			},
		},
	}
}

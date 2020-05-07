package engine

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/tracker"
	"log"
	"time"
)

func (e *Engine) ProcessTrackerFrame(wrapperFrame *tracker.TrackerWrapperPacket) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if wrapperFrame.TrackedFrame == nil {
		return
	}

	now := e.timeProvider()
	state := GcStateTracker{
		SourceName: wrapperFrame.SourceName,
		Uuid:       wrapperFrame.Uuid,
		Ball:       convertBalls(wrapperFrame.TrackedFrame.Balls),
		Robots:     convertRobots(wrapperFrame.TrackedFrame.Robots),
	}
	e.trackerLastUpdate[*wrapperFrame.Uuid] = now
	e.gcState.TrackerState[*wrapperFrame.Uuid] = &state

	if e.config.ActiveTrackerSource == nil {
		e.config.ActiveTrackerSource = state.Uuid
		log.Printf("Switched tracker source to %v (%v)", *state.Uuid, *state.SourceName)
	}

	if *e.config.ActiveTrackerSource == *wrapperFrame.Uuid {
		e.gcState.TrackerStateGc = &state
	}
}

func (e *Engine) processTrackerSources() {
	now := e.timeProvider()

	if e.config.ActiveTrackerSource != nil && now.Sub(e.trackerLastUpdate[*e.config.ActiveTrackerSource]) > time.Second {
		log.Printf("Tracker source %v timed out", *e.config.ActiveTrackerSource)
		e.config.ActiveTrackerSource = nil
		e.gcState.TrackerStateGc = &GcStateTracker{}
	}

	// remove old states
	for sourceId, state := range e.gcState.TrackerState {
		if now.Sub(e.trackerLastUpdate[*state.Uuid]) > time.Second {
			delete(e.gcState.TrackerState, sourceId)
		}
	}
}

func convertRobots(robots []*tracker.TrackedRobot) (rs []*Robot) {
	for _, robot := range robots {
		rs = append(rs, &Robot{
			Id:  robot.RobotId,
			Pos: robot.Pos,
		})
	}
	return
}

func convertBalls(balls []*tracker.TrackedBall) *Ball {
	if len(balls) == 0 {
		return nil
	}
	ball := Ball{
		Pos: balls[0].Pos,
		Vel: balls[0].Vel,
	}
	return &ball
}

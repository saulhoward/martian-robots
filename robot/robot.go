package robot

import (
	"container/ring"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type (
	coordinate struct {
		x, y int
	}

	world struct {
		upperRight                 *coordinate   // world boundary
		fallenPositions            []*coordinate // store positions of robots that go oob
		finalPositionAndDirections []string      // store final position/direction of robots as strings
	}

	robot struct {
		position  *coordinate
		direction *direction
		isFallen  bool // robot has "fallen" oob, final position inclues "LOST", and it stops moving
	}

	direction struct {
		ring *ring.Ring // store compass direction as ring
	}
)

var compassPoints = []string{"N", "E", "S", "W"}

// RunRobots takes the complete instructions for the world and robot
// movement, and returns a slice of strings for the robot final
// positions.
func RunRobots(instructions string) ([]string, error) {
	lines := strings.FieldsFunc(
		string(instructions), func(c rune) bool { return c == '\n' || c == '\r' },
	)
	if len(lines) < 3 {
		return nil, errors.New("invalid input")
	}

	world := &world{}
	var curRobot *robot

	for i, l := range lines {
		if i == 0 {
			// first line -- initalise the world
			coords := strings.Fields(l)
			if len(coords) != 2 {
				return nil, fmt.Errorf("unable to parse world size: %s", l)
			}
			upperRight, err := parseCoords(coords[0], coords[1])
			if err != nil {
				return nil, err
			}
			err = world.setUpperRight(upperRight)
			if err != nil {
				return nil, err
			}
			continue
		}
		if i%2 != 0 {
			// odd -- robot position/direction

			// add the previous robot's position
			if curRobot != nil {
				world.addFinalPositionAndDirection(curRobot)
			}

			coords := strings.Fields(l)
			if len(coords) != 3 {
				return nil, fmt.Errorf("unable to parse robot position: %s", l)
			}
			robotCoords, err := parseCoords(coords[0], coords[1])
			if err != nil {
				return nil, err
			}
			robotDir, err := newDirection(coords[2])
			if err != nil {
				return nil, err
			}
			newRobot, err := newRobot(robotCoords, robotDir, world)
			if err != nil {
				return nil, err
			}
			curRobot = newRobot

		} else {
			// even -- robot movement instructions
			for _, instr := range l {
				if !curRobot.isFallen {
					err := world.moveRobot(curRobot, string(instr))
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	// add the final robot's position
	if curRobot != nil {
		world.addFinalPositionAndDirection(curRobot)
	}

	return world.finalPositionAndDirections, nil
}

// NB: Sample data doesn't clearly indicate spaces between the coords,
// but I assume they're there, or we wouldn't have coords > 9
func parseCoords(xStr, yStr string) (*coordinate, error) {
	var x, y int
	var err error
	if x, err = strconv.Atoi(xStr); err != nil {
		return nil, err
	}
	if y, err = strconv.Atoi(yStr); err != nil {
		return nil, err
	}
	return &coordinate{x: x, y: y}, nil
}

// initialise a collections/ring for use as the compass points
func newDirection(start string) (*direction, error) {
	index, err := indexOfCompass(start)
	if err != nil {
		return nil, err
	}
	r := ring.New(4)
	for _, c := range compassPoints {
		r.Value = c
		r = r.Next()
	}
	r = r.Move(index)
	return &direction{ring: r}, nil
}

func newRobot(pos *coordinate, dir *direction, world *world) (*robot, error) {
	if pos.x > world.upperRight.x || pos.y > world.upperRight.y {
		return nil, errors.New("illegal starting position for robot")
	}
	return &robot{position: pos, direction: dir}, nil
}

// return the index of the provided compass point in the compassPoints
// slice
func indexOfCompass(c string) (int, error) {
	for k, v := range compassPoints {
		if c == v {
			return k, nil
		}
	}
	return -1, fmt.Errorf("compass point %s unknown", c)
}

func (w *world) setUpperRight(c *coordinate) error {
	if c.x > 50 || c.y > 50 || c.x < 0 || c.y < 0 {
		return errors.New("illegal world size")
	}
	w.upperRight = c
	return nil
}

func (w *world) addFallenPosition(c *coordinate) {
	w.fallenPositions = append(w.fallenPositions, c)
}

func (w *world) addFinalPositionAndDirection(r *robot) {
	str := fmt.Sprintf("%d %d %s", r.position.x, r.position.y, r.direction.ring.Value)
	if r.isFallen {
		str += " LOST"
	}
	w.finalPositionAndDirections = append(w.finalPositionAndDirections, str)
}

// Update the provided robot's direction/position according to a single
// command, (currently either L, R or F).
// Update the world if the robot goes out of bounds (falls).
func (w *world) moveRobot(robot *robot, instr string) error {
	prevPos := &coordinate{x: robot.position.x, y: robot.position.y}
	switch instr {
	case "L":
		robot.direction.ring = robot.direction.ring.Prev()
	case "R":
		robot.direction.ring = robot.direction.ring.Next()
	case "F":
		switch robot.direction.ring.Value {
		case "N":
			robot.position.y += 1
		case "E":
			robot.position.x += 1
		case "S":
			robot.position.y -= 1
		case "W":
			robot.position.x -= 1
		}
	default:
		return fmt.Errorf("illegal command %s", instr)
	}
	if robot.position.x < 0 ||
		robot.position.y < 0 ||
		robot.position.x > w.upperRight.x ||
		robot.position.y > w.upperRight.y {

		isFallenPos := false
		for _, c := range w.fallenPositions {
			if prevPos.x == c.x && prevPos.y == c.y {
				// if this is a position where other robots have fallen,
				// don't set this robot isFallen.
				isFallenPos = true
				break
			}
		}
		robot.position = prevPos
		if !isFallenPos {
			w.addFallenPosition(prevPos)
			robot.isFallen = true
		}
	}
	return nil
}

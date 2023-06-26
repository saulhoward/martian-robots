package main

import (
	"container/ring"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	Coordinate struct {
		X, Y int
	}

	World struct {
		UpperRight                 *Coordinate
		FallenPositions            []*Coordinate
		FinalPositionAndDirections []string
	}

	Robot struct {
		Position  *Coordinate
		Direction *Direction
		IsFallen  bool
	}

	Direction struct {
		Ring *ring.Ring
	}
)

var compassPoints = []string{"N", "E", "S", "W"}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "i", "", "path of file containing robot instructions")
	flag.Parse()

	instructions, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	lines := strings.FieldsFunc(
		string(instructions), func(c rune) bool { return c == '\n' || c == '\r' },
	)

	world := &World{}
	var curRobot *Robot

	for i, l := range lines {
		if i == 0 {
			// first line -- initalise the world
			coords := strings.Fields(l)
			if len(coords) != 2 {
				log.Fatalf("unable to parse world size: %s", l)
			}
			upperRight, err := ParseCoords(coords[0], coords[1])
			if err != nil {
				log.Fatalf("unable to parse upper right coords: %v", err)
			}
			world.UpperRight = upperRight
			continue
		}
		if i%2 != 0 {
			// odd -- robot position/direction

			// add the previous robot's position
			if curRobot != nil {
				world.AddFinalPositionAndDirection(curRobot)
			}

			coords := strings.Fields(l)
			if len(coords) != 3 {
				log.Fatalf("unable to parse robot position: %s", l)
			}
			robotCoords, err := ParseCoords(coords[0], coords[1])
			if err != nil {
				log.Fatalf("unable to parse robot coords: %v", err)
			}
			robotDir, err := NewDirection(coords[2])
			if err != nil {
				log.Fatalf("unable to parse robot direction: %v", err)
			}
			curRobot = &Robot{Position: robotCoords, Direction: robotDir}

		} else {
			// even -- robot movement instructions
			for _, instr := range l {
				if !curRobot.IsFallen {
					world.MoveRobot(curRobot, string(instr))
				}
			}
		}
	}

	// add the final robot's position
	if curRobot != nil {
		world.AddFinalPositionAndDirection(curRobot)
	}

	for _, p := range world.FinalPositionAndDirections {
		fmt.Println(p)
	}
}

// NB: Sample data doesn't clearly indicate spaces between the coords,
// but I assume they're there, or we wouldn't have coords > 9
func ParseCoords(xStr, yStr string) (*Coordinate, error) {
	var x, y int
	var err error
	if x, err = strconv.Atoi(xStr); err != nil {
		return nil, err
	}
	if y, err = strconv.Atoi(yStr); err != nil {
		return nil, err
	}
	return &Coordinate{X: x, Y: y}, nil
}

func NewDirection(start string) (*Direction, error) {
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
	return &Direction{Ring: r}, nil
}

func indexOfCompass(c string) (int, error) {
	for k, v := range compassPoints {
		if c == v {
			return k, nil
		}
	}
	return -1, fmt.Errorf("compass point %s unknown", c)
}

func (w *World) AddFallenPosition(c *Coordinate) {
	w.FallenPositions = append(w.FallenPositions, c)
}

func (w *World) AddFinalPositionAndDirection(r *Robot) {
	str := fmt.Sprintf("%d %d %s", r.Position.X, r.Position.Y, r.Direction.Ring.Value)
	if r.IsFallen {
		str += " LOST"
	}
	w.FinalPositionAndDirections = append(w.FinalPositionAndDirections, str)
}

func (w *World) MoveRobot(robot *Robot, instr string) {
	prevPos := &Coordinate{X: robot.Position.X, Y: robot.Position.Y}
	switch instr {
	case "L":
		robot.Direction.Ring = robot.Direction.Ring.Prev()
	case "R":
		robot.Direction.Ring = robot.Direction.Ring.Next()
	case "F":
		switch robot.Direction.Ring.Value {
		case "N":
			robot.Position.Y += 1
		case "E":
			robot.Position.X += 1
		case "S":
			robot.Position.Y -= 1
		case "W":
			robot.Position.X -= 1
		}
		// case "other command types":
	}
	if robot.Position.X < 0 ||
		robot.Position.Y < 0 ||
		robot.Position.X > w.UpperRight.X ||
		robot.Position.Y > w.UpperRight.Y {

		isFallenPos := false
		for _, c := range w.FallenPositions {
			if prevPos.X == c.X && prevPos.Y == c.Y {
				isFallenPos = true
				break
			}
		}
		robot.Position = prevPos
		if !isFallenPos {
			w.AddFallenPosition(prevPos)
			robot.IsFallen = true
		}
	}
}

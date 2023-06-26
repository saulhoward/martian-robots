package robot_test

import (
	"strings"
	"testing"

	"github.com/saulhoward/martian-robots/robot"
)

func TestRobots(t *testing.T) {

	fix := []struct {
		input  string
		result string
	}{
		{
			input: `5 3
1 1 E
RFRFRFRF

3 2 N
FRRFLLFFRRFLL

0 3 W
LLFFFLFLFL`,
			result: `1 1 E
3 3 N LOST
2 3 S`,
		},
		{
			input: `4 4
0 0 N
FFFFF

0 0 N
FFFFF`,
			result: `0 4 N LOST
0 4 N`,
		},
		{
			input: `4 4
0 0 E
FFFFF

0 0 E
FFFFF`,
			result: `4 0 E LOST
4 0 E`,
		},
		{
			input: `4 4
0 0 N
FFLLLLFF`,
			result: `0 4 N`,
		},
	}

	for _, f := range fix {
		resultArr, err := robot.RunRobots(f.input)
		if err != nil {
			t.Fatalf("failed to run with err %v", err)
		}
		result := strings.Join(resultArr, "\n")
		if result != f.result {
			t.Fatalf("failed to match expected result\n\n%s\n\n%s", result, f.result)
		}
	}

}

func TestRobotsIllegalInput(t *testing.T) {

	fix := []struct {
		input string
	}{
		{
			input: `foobar`,
		},
		{
			input: `100 100
0 0 N
FFFFF`,
		},
		{
			input: `4 4
6 6 E
FFFFF`,
		},
		{
			input: `4 4
0 0 X
RFLEF`,
		},
		{
			input: `4 4
0 0 E
FX`,
		},
	}

	for _, f := range fix {
		_, err := robot.RunRobots(f.input)
		if err == nil {
			t.Fatalf("failed to raise error with input %s", f.input)
		}
	}

}

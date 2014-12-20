package playground

import "testing"

func TestCalculateDistance(t *testing.T) {
	var tests = []*struct {
		from, to *Dot
		result   uint16
	}{
		{NewDot(0, 0), NewDot(0, 0), 0},
		{NewDot(10, 10), NewDot(10, 10), 0},
		{NewDot(10, 0), NewDot(10, 10), 10},
		{NewDot(10, 10), NewDot(10, 0), 10},
		{NewDot(10, 0), NewDot(0, 0), 10},
		{NewDot(0, 0), NewDot(11, 0), 11},
		{NewDot(10, 5), NewDot(0, 0), 15},
		{NewDot(0, 0), NewDot(11, 22), 33},
	}

	for i, test := range tests {
		if CalculateDistance(test.from, test.to) != test.result {
			t.Fatalf("Calculating distance error: test #%i", i+1)
		}
	}
}

func TestCalculateDirection(t *testing.T) {
	var tests = []*struct {
		from, to *Dot
		result   Direction
	}{
		{NewDot(10, 0), NewDot(10, 10), DIR_SOUTH},
		{NewDot(10, 10), NewDot(10, 0), DIR_NORTH},
		{NewDot(10, 0), NewDot(0, 0), DIR_WEST},
		{NewDot(0, 0), NewDot(11, 0), DIR_EAST},
		{NewDot(10, 5), NewDot(0, 0), DIR_WEST},
		{NewDot(0, 0), NewDot(11, 22), DIR_SOUTH},
	}

	for i, test := range tests {
		if CalculateDirection(test.from, test.to) != test.result {
			t.Fatalf("Calculating direction error: test #%i", i+1)
		}
	}
}

func TestNavigationOnPlayground(t *testing.T) {
	var pg, _ = NewPlayground(30, 30)

	var tests = []*struct {
		dot *Dot
		dir Direction
		dis uint8
		res *Dot
	}{
		{NewDot(10, 0), DIR_SOUTH, 10, NewDot(10, 10)},
		{NewDot(10, 10), DIR_NORTH, 10, NewDot(10, 0)},
		{NewDot(10, 0), DIR_WEST, 10, NewDot(0, 0)},
		{NewDot(0, 0), DIR_EAST, 11, NewDot(11, 0)},
		{NewDot(0, 0), DIR_WEST, 3, NewDot(27, 0)},
		{NewDot(28, 0), DIR_EAST, 2, NewDot(0, 0)},
		{NewDot(28, 0), DIR_EAST, 92, NewDot(0, 0)},
		{NewDot(10, 0), DIR_NORTH, 10, NewDot(10, 20)},
		{NewDot(10, 10), DIR_SOUTH, 30, NewDot(10, 10)},
	}

	for i, test := range tests {
		if dot, _ := pg.Navigate(test.dot, test.dir,
			test.dis); !dot.Equals(test.res) {
			t.Fatalf("Navigation error: test #%i %s", i+1, dot)
		}
	}
}

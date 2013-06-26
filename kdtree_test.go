package kdtree

import (
	"fmt"
	"testing"
)

type Point struct {
	x, y float64
}

func (p Point) X() float64 {
	return p.x
}

func (p Point) Y() float64 {
	return p.y
}

func (p Point) DistTo(pt XY) float64 {
	xdiff := p.X() - pt.X()
	ydiff := p.Y() - pt.Y()
	return xdiff*xdiff + ydiff*ydiff
}

func (p Point) String() string {
	return fmt.Sprintf("(%f,%f)", p.x, p.y)
}

type Compare struct {
}

func (c Compare) CompareX(a, b interface{}) float64 {
	return a.(Point).x - b.(Point).x
}

func (c Compare) CompareY(a, b interface{}) float64 {
	return a.(Point).y - b.(Point).y
}

func TestNeighbourSearch(t *testing.T) {
	var c Compare
	points := []XY{
		Point{1, 2},
		Point{2, 1},
		Point{3, 2},
		Point{2, 0},
		Point{1.5, 0},
		Point{2.5, 3},
		Point{4, 5},
	}

	tree := New(c)
	tree.Build(points)

	p := Point{1.5, 0.3}
	xy, _, err := tree.FindNearest(p)
	if err != nil || xy.X() != 1.5 || xy.Y() != 0 {
		t.Errorf("Expected to find (1.5,0) - instead found %v", xy)
	}

	p = Point{2.8, 2}
	xy, _, err = tree.FindNearest(p)
	if err != nil || xy.X() != 3 || xy.Y() != 2 {
		t.Errorf("Expected to find (3,2) - instead found %v", xy)
	}

	p = Point{1.1, 1.7}
	xy, _, err = tree.FindNearest(p)
	if err != nil || xy.X() != 1 || xy.Y() != 2 {
		t.Errorf("Expected to find (1,2) - instead %v", xy)
	}

}

// This file impements 2d tree and nearest neighbour search

/*
 *	 Copyright 2013 Sarath Lakshman
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package kdtree

import (
	"errors"
	"math"
)

var (
	ErrNotFound = errors.New("Neighbour not found")
)

// Point interface
type XY interface {
	X() float64
	Y() float64
	DistTo(xy XY) float64
}

// Axis wise comparator
type XYComparator interface {
	CompareX(xy1, xy2 interface{}) float64
	CompareY(xy1, xy2 interface{}) float64
}

// Generic helper to abstract Point slices
type XYSlice []XY

func (s XYSlice) At(i int) interface{} {
	return s[i]
}

func (s XYSlice) Len() int {
	return len(s)
}

func (s XYSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s XYSlice) Slice(i, j int) Data {
	return XYSlice(s[i:j])
}

// 2d point node
type kdNode struct {
	val         XY
	cmp         func(xy1, xy2 interface{}) float64
	left, right *kdNode
}

// Best node sofar
type nearNode struct {
	dist float64
	node *kdNode
}

type KDTree struct {
	root    *kdNode
	compare XYComparator
}

func New(cmpr XYComparator) *KDTree {
	return &KDTree{
		compare: cmpr,
	}
}

// Helper - recursively build balanced 2d tree
func (t *KDTree) build(l []XY, level int) *kdNode {
	if len(l) == 0 {
		return nil
	}

	node := new(kdNode)
	if level%2 == 0 {
		node.cmp = t.compare.CompareX
	} else {
		node.cmp = t.compare.CompareY
	}

	if len(l) == 1 {
		node.val = l[0]
		return node
	}

	// Find median of current list, set that to val of current node
	// Recursively apply the same on left and right point set
	sl := XYSlice(l)
	m := len(l) / 2
	SelectK(sl, m, node.cmp)

	for m < sl.Len()-1 && node.cmp(sl.At(m-1), sl.At(m)) == 0 {
		m += 1
	}

	node.val = sl.At(m - 1).(XY)
	node.left = t.build(l[:m-1], level+1)
	node.right = t.build(l[m:], level+1)

	return node
}

// Build a balanced 2d tree from given set of points
func (t *KDTree) Build(nodeList []XY) {
	t.root = t.build(nodeList, 0)
}

// Helper
func (t KDTree) findNearest(xy XY, n *kdNode, near *nearNode) {
	if n == nil {
		return
	}

	// Set current node to best node if its the best dist sofar
	dist := xy.DistTo(n.val)
	if near.node == nil || near.dist > dist {
		near.node = n
		near.dist = dist
	}

	// Draw a circle with center as xy and passing through n.
	// If xy lies left or bottom of axis (y or x), large sector formed by current axis split
	// lies in left or bottom side of axis. The max best dist that can be formed using other section
	// is equal to axisDist. Hence if we find a dist < axisDist from left or bottom section, we do not
	// need to search in the other section.
	axisDist := n.cmp(xy, n.val)
	var nearest, farthest *kdNode

	if axisDist <= 0 {
		nearest = n.left
		farthest = n.right
	} else {
		nearest = n.right
		farthest = n.left
	}

	t.findNearest(xy, nearest, near)

	if axisDist*axisDist >= near.dist {
		return
	}

	t.findNearest(xy, farthest, near)
}

// Find nearest point from the given point
func (t KDTree) FindNearest(xy XY) (r XY, dist float64, err error) {
	near := new(nearNode)

	t.findNearest(xy, t.root, near)
	if near.node != nil {
		r = near.node.val
		dist = math.Sqrt(near.dist)
		return
	}
	err = ErrNotFound
	return
}


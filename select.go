// This file implements selection algorithm

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

type Data interface {
	At(i int) interface{}
	Len() int
	Swap(i, j int)
	Slice(i, j int) Data
}

func partition(l Data, cmp func(interface{}, interface{}) float64) int {
	var i int = 0
	var j int = l.Len()
	var p int = 0

	for {
		for i += 1; i < l.Len() && cmp(l.At(i), l.At(p)) < 0; i += 1 {
		}

		for j -= 1; j > 0 && cmp(l.At(j), l.At(p)) > 0; j -= 1 {
		}

		if i >= j {
			break
		}

		l.Swap(i, j)
	}

	l.Swap(j, p)

	return j
}

// Find kth minimum object from the given list in
// Complexity: O(n)
func SelectK(l Data, k int, cmp func(interface{}, interface{}) float64) {
	if l.Len() <= 1 {
		return
	}

	p := partition(l, cmp)
	if k <= p {
		SelectK(l.Slice(0, p), k, cmp)
	} else {
		SelectK(l.Slice(p+1, l.Len()), k-p-1, cmp)
	}
}


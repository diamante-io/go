// Copyright 2016 diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"fmt"
	"path"
)

type DirPrefix [3]uint8

func (d DirPrefix) Path() string {
	return d.pathPrefix(len(d))
}

func (d DirPrefix) pathPrefix(n int) string {
	tmp := []string{}
	for i, b := range d {
		if i > n {
			break
		}
		tmp = append(tmp, fmt.Sprintf("%02x", b))
	}
	return path.Join(tmp...)
}

func CheckpointPrefix(seq uint32) DirPrefix {
	return DirPrefix{
		uint8(seq >> 24),
		uint8(seq >> 16),
		uint8(seq >> 8),
	}
}

func HashPrefix(h Hash) DirPrefix {
	return DirPrefix{h[0], h[1], h[2]}
}

// Returns an array of path prefixes to walk to enumerate all the
// objects in the provided range.
func RangePaths(r Range) []string {
	res := []string{}
	lowpre := CheckpointPrefix(r.Low)
	highpre := CheckpointPrefix(r.High)
	diff := 0
	for i, e := range lowpre {
		diff = i
		if highpre[i] != e {
			break
		}
	}
	// log.Printf("prefix %s and %s differ at point %d",
	//            lowpre.Path(), highpre.Path(), diff)
	tmp := lowpre
	for i := int(lowpre[diff]); i <= int(highpre[diff]); i++ {
		tmp[diff] = uint8(i)
		res = append(res, tmp.pathPrefix(diff))
	}
	return res
}

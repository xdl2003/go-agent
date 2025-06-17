// Package util
/**
@author: xdl2003
@desc:
@date: 2025/6/11
**/
package util

func ContainArrStr(elems []string, target string) bool {
	for _, elem := range elems {
		if elem == target {
			return true
		}
	}
	return false
}

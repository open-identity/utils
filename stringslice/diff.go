package stringslice

import "sort"

func Difference(slice1 []string, slice2 []string) (added []string, removed []string) {

	var i, j int
	added = []string{}
	removed = []string{}

	sort.StringSlice(slice1).Sort()
	sort.StringSlice(slice2).Sort()

	for i < len(slice1) && j < len(slice2) {
		if slice1[i] < slice2[j] {
			removed = append(removed, slice1[i])
			i++
		} else if slice1[i] > slice2[j] {
			added = append(added, slice2[j])
			j++
		} else {
			i++
			j++
		}
	}
	for i < len(slice1) {
		removed = append(removed, slice1[i])
		i++
	}
	for j < len(slice2) {
		added = append(added, slice2[j])
		j++
	}
	return
}

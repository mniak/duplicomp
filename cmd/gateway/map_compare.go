package main

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type Difference struct {
	Key            int
	Left           any
	Right          any
	Message        string
	SubDifferences []Difference
}

func CompareMaps(mapLeft, mapRight map[int]any) []Difference {
	allkeysMap := make(map[int]struct{})
	for k := range mapLeft {
		allkeysMap[k] = struct{}{}
	}
	for k := range mapRight {
		allkeysMap[k] = struct{}{}
	}

	allKeys := maps.Keys(allkeysMap)
	slices.Sort(allKeys)

	var result []Difference
	for _, key := range allKeys {
		diff, hasDiff := CompareMapValues(key, mapLeft, mapRight)
		if hasDiff {
			result = append(result, diff)
		}
	}
	return result
}

func CompareMapValues(key int, mapLeft, mapRight map[int]any) (Difference, bool) {
	diff := Difference{
		Key: key,
	}

	left, hasValue1 := mapLeft[key]
	right, hasValue2 := mapRight[key]
	if hasValue1 {
		diff.Left = left
	}
	if hasValue2 {
		diff.Right = right
	}

	if !hasValue1 && !hasValue2 {
		return diff, false
	}

	switch {
	case hasValue1 && !hasValue2:
		diff.Message = fmt.Sprintf("value present on the left but missing on the right")
		return diff, true
	case !hasValue1 && hasValue2:
		diff.Message = fmt.Sprintf("value present on the right but missing on the left")
		return diff, true
	}

	objLeft, leftIsObject := left.(map[int]any)
	objRight, rightIsObject := left.(map[int]any)

	switch {
	case leftIsObject && !rightIsObject:
		diff.Message = fmt.Sprintf("value on the left is object but on the right not")
		return diff, true
	case !leftIsObject && rightIsObject:
		diff.Message = fmt.Sprintf("value on the right is object but on the left not")
		return diff, true
	}
	diff.SubDifferences = CompareMaps(objLeft, objRight)
	if len(diff.SubDifferences) > 0 {
		diff.Message = "there are differences on the subfields"
		return diff, true
	}

	if !reflect.DeepEqual(left, right) {
		diff.Message = "the values differ"
		return diff, true
	}
	return Difference{}, false
}

type KeyPath []int

func (kp KeyPath) String() string {
	var sb strings.Builder
	for i, p := range kp {
		if i != 0 {
			sb.WriteRune('.')
		}
		fmt.Fprintf(&sb, "%d", p)
	}
	return sb.String()
}

type FlatDifference struct {
	FieldPath KeyPath
	Message   string
}

func FlattenDifferences(keyPath KeyPath, diffs []Difference) []FlatDifference {
	var result []FlatDifference
	for _, diff := range diffs {
		keyPath = append(keyPath, diff.Key)
		if len(diff.SubDifferences) > 0 {
			FlattenDifferences(keyPath, diff.SubDifferences)
		} else {
			result = append(result, FlatDifference{
				FieldPath: append(keyPath, diff.Key),
				Message:   diff.Message,
			})
		}
	}
	return result
}

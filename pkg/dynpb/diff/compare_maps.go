package diff

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

func CompareMaps(mapLeft, mapRight map[Key]Value) Differences {
	allkeysMap := make(map[Key]Unit)
	for k := range mapLeft {
		allkeysMap[k] = Unit{}
	}
	for k := range mapRight {
		allkeysMap[k] = Unit{}
	}

	allKeys := maps.Keys(allkeysMap)
	slices.Sort(allKeys)

	var result Differences
	for _, key := range allKeys {
		diff, hasDiff := compareMapValues(key, mapLeft, mapRight)
		if hasDiff {
			result = append(result, diff)
		}
	}
	return result
}

func compareMapValues(key Key, mapLeft, mapRight map[Key]Value) (Difference, bool) {
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
		diff.Difference = RightIsMissing
		return diff, true
	case !hasValue1 && hasValue2:
		diff.Difference = LeftIsMissing
		return diff, true
	}

	objLeft, leftIsObject := left.(map[Key]Value)
	objRight, rightIsObject := left.(map[Key]Value)

	switch {
	case leftIsObject && !rightIsObject:
		diff.Difference = LeftIsObject
		return diff, true
	case !leftIsObject && rightIsObject:
		diff.Difference = RightIsObject
		return diff, true
	}
	diff.SubDifferences = CompareMaps(objLeft, objRight)
	if len(diff.SubDifferences) > 0 {
		diff.Difference = ValuesAreDifferent
		return diff, true
	}

	if !reflect.DeepEqual(left, right) {
		diff.Difference = SubfieldsAreDifferent
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
	Path       KeyPath
	Difference DiffKind
}

func FlattenDifferences(keyPath KeyPath, diffs Differences) []FlatDifference {
	return flattenDifferences(KeyPath{}, diffs)
}

func flattenDifferences(keyPath KeyPath, diffs Differences) []FlatDifference {
	var result []FlatDifference
	for _, diff := range diffs {
		keyPath = append(keyPath, diff.Key)
		if len(diff.SubDifferences) > 0 {
			flattenDifferences(keyPath, diff.SubDifferences)
		} else {
			result = append(result, FlatDifference{
				Path:       append(keyPath, diff.Key),
				Difference: diff.Difference,
			})
		}
	}
	return result
}

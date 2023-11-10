package diff

import (
	"reflect"
	"slices"

	"golang.org/x/exp/maps"
)

func CompareMaps(mapLeft, mapRight map[Key]Value) Differences {
	allKeysMap := make(map[Key]Unit)
	for k := range mapLeft {
		allKeysMap[k] = Unit{}
	}
	for k := range mapRight {
		allKeysMap[k] = Unit{}
	}

	allKeys := maps.Keys(allKeysMap)
	slices.Sort(allKeys)

	var result Differences
	for _, key := range allKeys {
		diff, areEqual := compareMapValues(key, mapLeft, mapRight)
		if !areEqual {
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
		return diff, true
	}

	switch {
	case hasValue1 && !hasValue2:
		diff.Difference = RightIsMissing
		return diff, false
	case !hasValue1 && hasValue2:
		diff.Difference = LeftIsMissing
		return diff, false
	}

	objLeft, leftIsObject := left.(map[Key]Value)
	objRight, rightIsObject := right.(map[Key]Value)

	switch {
	case leftIsObject && rightIsObject:
		diff.SubDifferences = CompareMaps(objLeft, objRight)
		if len(diff.SubDifferences) > 0 {
			diff.Difference = SubfieldsAreDifferent
			return diff, false
		}

	case leftIsObject && !rightIsObject:
		diff.Difference = LeftIsObject
		return diff, false
	case !leftIsObject && rightIsObject:
		diff.Difference = RightIsObject
		return diff, false

	}

	if !reflect.DeepEqual(left, right) {
		diff.Difference = ValuesAreDifferent
		return diff, false
	}
	return Difference{}, true
}

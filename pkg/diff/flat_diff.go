package diff

import (
	"fmt"
	"strings"
)

type KeyPath []int

type FlatDifference struct {
	Path       KeyPath
	Difference Kind
}

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

func (diffs Differences) Flatten() []FlatDifference {
	return diffs.flatten(KeyPath{})
}

func (diffs Differences) flatten(parentPath KeyPath) []FlatDifference {
	var result []FlatDifference
	for _, diff := range diffs {
		keyPath := append(parentPath, diff.Key)
		if len(diff.SubDifferences) > 0 {
			result = append(result, diff.SubDifferences.flatten(keyPath)...)
		} else {
			result = append(result, FlatDifference{
				Path:       keyPath,
				Difference: diff.Difference,
			})
		}
	}
	return result
}

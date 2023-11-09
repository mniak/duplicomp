package diff

type DiffKind string

const (
	ValuesAreDifferent    DiffKind = "ValuesAreDifferent"
	SubfieldsAreDifferent DiffKind = "SubfieldsAreDifferent"

	RightIsMissing DiffKind = "RightIsMissing"
	RightIsObject  DiffKind = "RightIsObject"

	LeftIsMissing DiffKind = "LeftIsMissing"
	LeftIsObject  DiffKind = "LeftIsObject"
)

type (
	Value       = any
	Key         = int
	Unit        = struct{}
	Differences = []ValueDiff
)

type ValueDiff struct {
	Key            Key
	Left           Value
	Right          Value
	Difference     DiffKind
	SubDifferences Differences
}

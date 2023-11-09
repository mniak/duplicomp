package diff

type DiffKind string

const (
	ValuesAreDifferent    DiffKind = "values differ"
	SubfieldsAreDifferent DiffKind = "subfields differ"

	RightIsMissing DiffKind = "right is missing"
	RightIsObject  DiffKind = "right is object, left is not"

	LeftIsMissing DiffKind = "left is missing"
	LeftIsObject  DiffKind = "left is object, right is not"
)

type (
	Value       = any
	Key         = int
	Unit        = struct{}
	Differences = []Difference
)

type Difference struct {
	Key            Key         `json:"key"`
	Left           Value       `json:"left"`
	Right          Value       `json:"right"`
	Difference     DiffKind    `json:"diff"`
	SubDifferences Differences `json:"sub,omitempty"`
}

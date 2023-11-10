package diff

type Kind string

const (
	ValuesAreDifferent    Kind = "values differ"
	SubfieldsAreDifferent Kind = "subfields differ"

	RightIsMissing Kind = "right is missing"
	RightIsObject  Kind = "right is object, left is not"

	LeftIsMissing Kind = "left is missing"
	LeftIsObject  Kind = "left is object, right is not"
)

type (
	Value       = any
	Key         = int
	Unit        = struct{}
	Differences []Difference
)

type Difference struct {
	Key            Key         `json:"key"`
	Left           Value       `json:"left"`
	Right          Value       `json:"right"`
	Difference     Kind        `json:"diff"`
	SubDifferences Differences `json:"sub,omitempty"`
}

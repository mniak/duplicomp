package main

import (
	"fmt"
	"strings"

	"github.com/mniak/duplicomp/log2"
	"github.com/mniak/duplicomp/pkg/dynpb"
	"github.com/pkg/errors"
)

type LogComparator struct {
	logger log2.Logger
}

func (lc LogComparator) Compare(
	primaryMsg []byte, primaryError error,
	shadowMsg []byte, shadowError error,
) error {
	if primaryError != nil {
		if shadowError != primaryError {
		}
	}
	if shadowError != primaryError {
		lc.logger.Printf("Errors not match",
			"primary_error", primaryError,
			"shadow_error", shadowError,
		)
		return nil
	}

	hints := dynpb.HintMap{}

	primaryData, err := dynpb.ParseWithHints(primaryMsg, hints)
	if err != nil {
		return errors.WithMessage(err, "failed to parse primary message")
	}
	shadowData, err := dynpb.ParseWithHints(shadowMsg, hints)
	if err != nil {
		return errors.WithMessage(err, "failed to parse shadow message")
	}

	diffs := CompareMaps(primaryData, shadowData)
	var attrs []any
	flatDiffs := FlattenDifferences(nil, diffs)
	var sb strings.Builder
	for _, diff := range flatDiffs {
		fmt.Fprintf(&sb, "%s=%v, ", diff.KeyPath.String(), diff.Message)
	}
	if len(attrs) > 0 {
		lc.logger.Print("the two messages are different", sb.String())
	}

	return nil
}

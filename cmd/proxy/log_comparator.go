package main

import (
	"fmt"
	"io"

	"github.com/mniak/ps121/pkg/diff"
	"github.com/mniak/ps121/pkg/dynpb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

type LogComparator struct {
	Logger           zerolog.Logger
	AliasesPerMethod map[string]AliasTree
	HintsPerMethod   map[string]dynpb.HintMap
}

func (lc LogComparator) Compare(
	methodName string,
	primaryMsg []byte, primaryError error,
	shadowMsg []byte, shadowError error,
) error {
	if primaryError == io.EOF && shadowError == io.EOF {
		return nil
	}

	if shadowError != primaryError {
		e := lc.Logger.Info().
			AnErr("primary_error", primaryError).
			AnErr("shadow_error", shadowError)
		switch {
		case primaryError == io.EOF:
			e.Msg("message received from shadow after the primary connection was closed")
		case shadowError == io.EOF:
			e.Msg("shadow connection was closed before the primary")
		default:
			e.Msg("errors dont match")
		}
		return nil
	}

	hints := lc.HintsPerMethod[methodName]

	primaryData, err := dynpb.ParseWithHints(primaryMsg, hints)
	if err != nil {
		return errors.WithMessage(err, "failed to parse primary message")
	}
	shadowData, err := dynpb.ParseWithHints(shadowMsg, hints)
	if err != nil {
		return errors.WithMessage(err, "failed to parse shadow message")
	}

	diffs := diff.CompareMaps(primaryData, shadowData)
	flatDiffs := diffs.Flatten()
	diffMap := DiffMap(flatDiffs, lc.AliasesPerMethod[methodName])

	var evt *zerolog.Event
	if len(diffMap) > 0 {
		evt = lc.Logger.Info()
	} else {
		evt = lc.Logger.Debug()
	}

	evt.Int("DifferenceCount", len(diffMap))
	evt.Str("Method", methodName)
	evt.Any("Diferences", diffMap)

	if len(diffMap) > 0 {
		evt.Msg("the two messages are different")
	} else {
		evt.Msg("the two messages are equal")
	}

	return nil
}

func DiffMap(diffs []diff.FlatDifference, aliases AliasTree) map[string]diff.Kind {
	result := lo.Associate[diff.FlatDifference, string, diff.Kind](diffs, func(item diff.FlatDifference) (string, diff.Kind) {
		path := item.Path.String()
		alias, _ := aliases.GetAlias(item.Path...)
		result := fmt.Sprintf("[%s] %s", path, alias)
		return result, item.Difference
	})

	return result
}

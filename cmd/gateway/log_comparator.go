package main

import (
	"io"

	"github.com/mniak/duplicomp/pkg/dynpb"
	"github.com/mniak/duplicomp/pkg/dynpb/diff"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

type LogComparator struct {
	logger zerolog.Logger
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
		e := lc.logger.Info().
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

	hints := dynpb.HintMap{}

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
	diffMap := DiffMap(flatDiffs, AliasTree{})

	var evt *zerolog.Event
	if len(diffMap) > 0 {
		evt = lc.logger.Info()
	} else {
		evt = lc.logger.Debug()
	}

	evt.Bool("differ", len(diffMap) > 0)
	evt.Str("method", methodName)
	evt.Any("diffs", diffs)
	evt.Any("flat_diffs", flatDiffs)
	evt.Any("diff_map", diffMap)

	if len(diffMap) > 0 {
		evt.Msg("the two messages are different")
	} else {
		evt.Msg("the two messages are equal")
	}

	return nil
}

type AliasTree map[int]AliasNode

type AliasNode struct {
	Alias string
	Nodes AliasTree
}

func (node AliasNode) Find(indexes ...int) (string, bool) {
	if len(indexes) == 0 {
		return node.Alias, node.Alias != ""
	}
	return node.Nodes.Find(indexes[1:]...)
}

func (tree AliasTree) Find(indexes ...int) (string, bool) {
	if tree == nil || len(indexes) == 0 {
		return "", false
	}
	node, ok := tree[indexes[0]]
	if !ok {
		return "", false
	}

	return node.Find(indexes[1:]...)
}

func DiffMap(diffs []diff.FlatDifference, aliases AliasTree) map[string]diff.Kind {
	result := lo.Associate[diff.FlatDifference, string, diff.Kind](diffs, func(item diff.FlatDifference) (string, diff.Kind) {
		alias, ok := aliases.Find(item.Path...)
		if !ok {
			alias = item.Path.String()
		}
		return alias, item.Difference
	})

	return result
}

package main

import (
	"fmt"
	"io"

	"github.com/mniak/duplicomp/pkg/dynpb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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

	diffs := CompareMaps(primaryData, shadowData)
	flatDiffs := FlattenDifferences(nil, diffs)

	evtInfo := lc.logger.Info()
	for _, diff := range flatDiffs {
		evtInfo.Str(fmt.Sprintf("key_%s", diff.KeyPath.String()), diff.Message)
	}
	if len(flatDiffs) > 0 {
		evtInfo.Bool("has_differences", true).
			Msg("the two messages are different")
	} else {
		lc.logger.Debug().
			Bool("has_differences", false).
			Msg("the two messages are equal")
	}

	return nil
}

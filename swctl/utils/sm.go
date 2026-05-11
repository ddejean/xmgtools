// Copyright (C) 2026 - Damien Dejean <dam.dejean@gmail.com>

package utils

import "context"

// State type defines a state function that runs with arguments in args.
type State[T any] func(ctx context.Context, args T) (T, State[T], error)

// Run runs the state machine starting with state <start> and <args> arguments.
func Run[T any](ctx context.Context, args T, start State[T]) (T, error) {
	var err error
	current := start
	for {
		if ctx.Err() != nil {
			return args, ctx.Err()
		}
		args, current, err = current(ctx, args)
		if err != nil {
			return args, err
		}
		if current == nil {
			return args, nil
		}
	}
}

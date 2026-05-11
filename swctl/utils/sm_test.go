package utils

import (
	"context"
	"errors"
	"testing"
	"time"
)

type TestArgs struct {
	Count int
}

func testStateEnd(ctx context.Context, args *TestArgs) (*TestArgs, State[*TestArgs], error) {
	args.Count++
	return args, nil, nil
}

func testState2(ctx context.Context, args *TestArgs) (*TestArgs, State[*TestArgs], error) {
	args.Count++
	return args, testStateEnd, nil
}

func testState1(ctx context.Context, args *TestArgs) (*TestArgs, State[*TestArgs], error) {
	args.Count++
	return args, testState2, nil
}

func testStateErr(ctx context.Context, args *TestArgs) (*TestArgs, State[*TestArgs], error) {
	args.Count++
	return args, nil, errors.New("state error")
}

func testStateLoop(ctx context.Context, args *TestArgs) (*TestArgs, State[*TestArgs], error) {
	args.Count++
	return args, testStateLoop, nil
}

func TestRun_Success(t *testing.T) {
	ctx := context.Background()
	args := &TestArgs{Count: 0}

	finalArgs, err := Run(ctx, args, testState1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if finalArgs.Count != 3 {
		t.Errorf("expected count to be 3, got %d", finalArgs.Count)
	}
}

func TestRun_StateError(t *testing.T) {
	ctx := context.Background()
	args := &TestArgs{Count: 0}

	finalArgs, err := Run(ctx, args, testStateErr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "state error" {
		t.Errorf("expected 'state error', got %v", err)
	}
	if finalArgs.Count != 1 {
		t.Errorf("expected count to be 1, got %d", finalArgs.Count)
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	args := &TestArgs{Count: 0}

	// Cancel the context immediately
	cancel()

	finalArgs, err := Run(ctx, args, testState1)
	if err == nil {
		t.Fatal("expected error due to cancelled context, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
	if finalArgs.Count != 0 {
		t.Errorf("expected count to be 0, got %d", finalArgs.Count)
	}
}

func TestRun_ContextCancelledDuringRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	args := &TestArgs{Count: 0}

	// Use an infinite loop state that will be interrupted by the context timeout
	_, err := Run(ctx, args, testStateLoop)
	if err == nil {
		t.Fatal("expected error due to cancelled context, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded error, got %v", err)
	}
	// We don't check Count precisely because it depends on scheduling,
	// but it should be greater than 0.
	if args.Count == 0 {
		t.Errorf("expected count > 0, got %d", args.Count)
	}
}

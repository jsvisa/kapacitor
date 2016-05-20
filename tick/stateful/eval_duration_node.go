package stateful

import (
	"regexp"
	"time"
)

type EvalDurationNode struct {
	Duration time.Duration
}

func (n *EvalDurationNode) Type(scope ReadOnlyScope, executionState ExecutionState) (ValueType, error) {
	return TDuration, nil
}

func (n *EvalDurationNode) EvalFloat(scope *Scope, executionState ExecutionState) (float64, error) {
	return float64(0), ErrTypeGuardFailed{RequestedType: TFloat64, ActualType: TDuration}
}

func (n *EvalDurationNode) EvalInt(scope *Scope, executionState ExecutionState) (int64, error) {
	return 0, ErrTypeGuardFailed{RequestedType: TInt64, ActualType: TDuration}
}

func (n *EvalDurationNode) EvalString(scope *Scope, executionState ExecutionState) (string, error) {
	return "", ErrTypeGuardFailed{RequestedType: TString, ActualType: TDuration}
}

func (n *EvalDurationNode) EvalBool(scope *Scope, executionState ExecutionState) (bool, error) {
	return false, ErrTypeGuardFailed{RequestedType: TBool, ActualType: TDuration}
}
func (n *EvalDurationNode) EvalRegex(scope *Scope, executionState ExecutionState) (*regexp.Regexp, error) {
	return nil, ErrTypeGuardFailed{RequestedType: TRegex, ActualType: TDuration}
}
func (n *EvalDurationNode) EvalTime(scope *Scope, executionState ExecutionState) (time.Time, error) {
	return time.Time{}, ErrTypeGuardFailed{RequestedType: TTime, ActualType: TDuration}
}

func (n *EvalDurationNode) EvalDuration(scope *Scope, executionState ExecutionState) (time.Duration, error) {
	return n.Duration, nil
}

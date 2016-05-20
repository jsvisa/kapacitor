package stateful

import (
	"regexp"
	"time"
)

type EvalFloatNode struct {
	Float64 float64
}

func (n *EvalFloatNode) Type(scope ReadOnlyScope, executionState ExecutionState) (ValueType, error) {
	return TFloat64, nil
}

func (n *EvalFloatNode) EvalFloat(scope *Scope, executionState ExecutionState) (float64, error) {
	return n.Float64, nil
}

func (n *EvalFloatNode) EvalInt(scope *Scope, executionState ExecutionState) (int64, error) {
	return int64(0), ErrTypeGuardFailed{RequestedType: TFloat64, ActualType: TFloat64}
}

func (n *EvalFloatNode) EvalString(scope *Scope, executionState ExecutionState) (string, error) {
	return "", ErrTypeGuardFailed{RequestedType: TString, ActualType: TFloat64}
}

func (n *EvalFloatNode) EvalBool(scope *Scope, executionState ExecutionState) (bool, error) {
	return false, ErrTypeGuardFailed{RequestedType: TBool, ActualType: TFloat64}
}
func (n *EvalFloatNode) EvalRegex(scope *Scope, executionState ExecutionState) (*regexp.Regexp, error) {
	return nil, ErrTypeGuardFailed{RequestedType: TRegex, ActualType: TFloat64}
}
func (n *EvalFloatNode) EvalTime(scope *Scope, executionState ExecutionState) (time.Time, error) {
	return time.Time{}, ErrTypeGuardFailed{RequestedType: TTime, ActualType: TFloat64}
}
func (n *EvalFloatNode) EvalDuration(scope *Scope, executionState ExecutionState) (time.Duration, error) {
	return 0, ErrTypeGuardFailed{RequestedType: TDuration, ActualType: TFloat64}
}

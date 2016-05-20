package stateful

import (
	"regexp"
	"time"
)

type EvalIntNode struct {
	Int64 int64
}

func (n *EvalIntNode) Type(scope ReadOnlyScope, executionState ExecutionState) (ValueType, error) {
	return TInt64, nil
}

func (n *EvalIntNode) EvalFloat(scope *Scope, executionState ExecutionState) (float64, error) {
	return float64(0), ErrTypeGuardFailed{RequestedType: TFloat64, ActualType: TInt64}
}

func (n *EvalIntNode) EvalInt(scope *Scope, executionState ExecutionState) (int64, error) {
	return n.Int64, nil
}

func (n *EvalIntNode) EvalString(scope *Scope, executionState ExecutionState) (string, error) {
	return "", ErrTypeGuardFailed{RequestedType: TString, ActualType: TInt64}
}

func (n *EvalIntNode) EvalBool(scope *Scope, executionState ExecutionState) (bool, error) {
	return false, ErrTypeGuardFailed{RequestedType: TBool, ActualType: TInt64}
}
func (n *EvalIntNode) EvalRegex(scope *Scope, executionState ExecutionState) (*regexp.Regexp, error) {
	return nil, ErrTypeGuardFailed{RequestedType: TRegex, ActualType: TInt64}
}
func (n *EvalIntNode) EvalTime(scope *Scope, executionState ExecutionState) (time.Time, error) {
	return time.Time{}, ErrTypeGuardFailed{RequestedType: TTime, ActualType: TInt64}
}
func (n *EvalIntNode) EvalDuration(scope *Scope, executionState ExecutionState) (time.Duration, error) {
	return 0, ErrTypeGuardFailed{RequestedType: TDuration, ActualType: TInt64}
}

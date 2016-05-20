package stateful

import (
	"regexp"
	"time"

	"github.com/influxdata/kapacitor/tick/ast"
)

type EvalRegexNode struct {
	Node *ast.RegexNode
}

func (n *EvalRegexNode) Type(scope ReadOnlyScope, executionState ExecutionState) (ValueType, error) {
	return TRegex, nil
}

func (n *EvalRegexNode) EvalRegex(scope *Scope, executionState ExecutionState) (*regexp.Regexp, error) {
	return n.Node.Regex, nil
}

func (n *EvalRegexNode) EvalString(scope *Scope, executionState ExecutionState) (string, error) {
	return "", ErrTypeGuardFailed{RequestedType: TString, ActualType: TRegex}
}

func (n *EvalRegexNode) EvalFloat(scope *Scope, executionState ExecutionState) (float64, error) {
	return float64(0), ErrTypeGuardFailed{RequestedType: TFloat64, ActualType: TRegex}
}

func (n *EvalRegexNode) EvalInt(scope *Scope, executionState ExecutionState) (int64, error) {
	return int64(0), ErrTypeGuardFailed{RequestedType: TInt64, ActualType: TRegex}
}

func (n *EvalRegexNode) EvalBool(scope *Scope, executionState ExecutionState) (bool, error) {
	return false, ErrTypeGuardFailed{RequestedType: TBool, ActualType: TRegex}
}
func (n *EvalRegexNode) EvalTime(scope *Scope, executionState ExecutionState) (time.Time, error) {
	return time.Time{}, ErrTypeGuardFailed{RequestedType: TTime, ActualType: TRegex}
}
func (n *EvalRegexNode) EvalDuration(scope *Scope, executionState ExecutionState) (time.Duration, error) {
	return 0, ErrTypeGuardFailed{RequestedType: TDuration, ActualType: TRegex}
}

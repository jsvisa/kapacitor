package stateful

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/influxdata/kapacitor/tick/ast"
)

// ErrTypeGuardFailed is returned when a speicifc value type is requested thorugh NodeEvaluator (for example: "Float64Value")
// when the node doesn't support the given type, for example "Float64Value" is called on BoolNode
type ErrTypeGuardFailed struct {
	RequestedType ValueType
	ActualType    ValueType
}

func (e ErrTypeGuardFailed) Error() string {
	return fmt.Sprintf("TypeGuard: expression returned unexpected type %s, expected %s", e.ActualType, e.RequestedType)
}

type ReadOnlyScope interface {
	Get(name string) (interface{}, error)
}

// NodeEvaluator provides a generic way for trying to fetch
// node value, if a speicifc type is requested (so Value isn't called, the *Value is called) ErrTypeGuardFailed must be returned
type NodeEvaluator interface {
	EvalFloat(scope *Scope, executionState ExecutionState) (float64, error)
	EvalInt(scope *Scope, executionState ExecutionState) (int64, error)
	EvalString(scope *Scope, executionState ExecutionState) (string, error)
	EvalBool(scope *Scope, executionState ExecutionState) (bool, error)
	EvalRegex(scope *Scope, executionState ExecutionState) (*regexp.Regexp, error)
	EvalTime(scope *Scope, executionState ExecutionState) (time.Time, error)
	EvalDuration(scope *Scope, executionState ExecutionState) (time.Duration, error)

	// Type returns the type of ValueType
	Type(scope ReadOnlyScope, executionState ExecutionState) (ValueType, error)
}

func createNodeEvaluator(n ast.Node) (NodeEvaluator, error) {
	switch node := n.(type) {

	case *ast.BoolNode:
		return &EvalBoolNode{Node: node}, nil

	case *ast.NumberNode:
		switch {
		case node.IsFloat:
			return &EvalFloatNode{Float64: node.Float64}, nil

		case node.IsInt:
			return &EvalIntNode{Int64: node.Int64}, nil

		default:
			// We wouldn't reach ever, unless there is bug in tick parsing ;)
			return nil, errors.New("Invalid NumberNode: Not float or int")
		}

	case *ast.DurationNode:
		return &EvalDurationNode{Duration: node.Dur}, nil
	case *ast.StringNode:
		return &EvalStringNode{Node: node}, nil

	case *ast.RegexNode:
		return &EvalRegexNode{Node: node}, nil

	case *ast.BinaryNode:
		return NewEvalBinaryNode(node)

	case *ast.ReferenceNode:
		return &EvalReferenceNode{Node: node}, nil

	case *ast.FunctionNode:
		return NewEvalFunctionNode(node)

	case *ast.UnaryNode:
		return NewEvalUnaryNode(node)
	}

	return nil, fmt.Errorf("Given node type is not valid evaluation node: %T", n)
}

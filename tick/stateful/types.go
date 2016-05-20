package stateful

import (
	"regexp"
	"time"

	"github.com/influxdata/kapacitor/tick/ast"
)

type ValueType uint8

const (
	InvalidType ValueType = iota << 1
	TFloat64
	TInt64
	TString
	TBool
	TRegex
	TTime
	TDuration
)

func (v ValueType) String() string {
	switch v {
	case TFloat64:
		return "float64"
	case TInt64:
		return "int64"
	case TString:
		return "string"
	case TBool:
		return "boolean"
	case TRegex:
		return "regex"
	case TTime:
		return "time"
	case TDuration:
		return "duration"
	}

	return "invalid type"
}

func valueTypeOf(v interface{}) ValueType {
	if v == nil {
		return InvalidType
	}
	switch v.(type) {
	case float64:
		return TFloat64
	case int64:
		return TInt64
	case string:
		return TString
	case bool:
		return TBool
	case *regexp.Regexp:
		return TRegex
	case time.Time:
		return TTime
	case time.Duration:
		return TDuration
	default:
		return InvalidType
	}
}

// getCostantNodeType - Given a ast.Node we want to know it's return type
// this method does exactly this, few examples:
// *) StringNode -> TString
// *) UnaryNode -> we base the type by the node type
func getConstantNodeType(n ast.Node) ValueType {
	switch node := n.(type) {
	case *ast.NumberNode:
		if node.IsInt {
			return TInt64
		}

		if node.IsFloat {
			return TFloat64
		}
	case *ast.DurationNode:
		return TDuration
	case *ast.StringNode:
		return TString
	case *ast.BoolNode:
		return TBool
	case *ast.RegexNode:
		return TRegex

	case *ast.UnaryNode:
		// If this is comparison operator we know for sure the output must be boolean
		if node.Operator == ast.TokenNot {
			return TBool
		}

		// Could be float int or duration
		if node.Operator == ast.TokenMinus {
			return getConstantNodeType(node.Node)
		}

	case *ast.BinaryNode:
		leftType := getConstantNodeType(node.Left)
		rightType := getConstantNodeType(node.Right)
		return binaryConstantTypes[operationKey{operator: node.Operator, leftType: leftType, rightType: rightType}]
	}

	return InvalidType
}

func isDynamicNode(n ast.Node) bool {
	switch node := n.(type) {
	case *ast.ReferenceNode:
		return true
	case *ast.FunctionNode:
		return true
	case *ast.UnaryNode:
		return isDynamicNode(node.Node)
	case *ast.BinaryNode:
		return isDynamicNode(node.Left) || isDynamicNode(node.Right)
	default:
		return false
	}
}

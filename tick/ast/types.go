package ast

import (
	"regexp"
	"time"
)

type Type int

const (
	InvalidType Type = iota
	TFloat
	TInt
	TString
	TBool
	TRegex
	TDuration
	TLambda
	TNode
)

func (v Type) String() string {
	switch v {
	case TFloat:
		return "float"
	case TInt:
		return "int"
	case TString:
		return "string"
	case TBool:
		return "bool"
	case TRegex:
		return "regex"
	case TDuration:
		return "duration"
	case TLambda:
		return "lambda"
	case TNode:
		return "node"
	}

	return "invalid type"
}

func TypeOf(v interface{}) Type {
	if v == nil {
		return InvalidType
	}
	switch v.(type) {
	case float64:
		return TFloat
	case int64:
		return TInt
	case string:
		return TString
	case bool:
		return TBool
	case *regexp.Regexp:
		return TRegex
	case time.Duration:
		return TDuration
	case *LambdaNode:
		return TLambda
	case Node:
		return TNode
	default:
		return InvalidType
	}
}

package ast

import (
	"errors"
	"regexp"
	"time"
)

type ValueType uint8

const (
	InvalidType ValueType = iota << 1
	TFloat
	TInt
	TString
	TBool
	TRegex
	TTime
	TDuration
)

func (v ValueType) String() string {
	switch v {
	case TFloat:
		return "float"
	case TInt:
		return "int"
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

func TypeOf(v interface{}) ValueType {
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
	case time.Time:
		return TTime
	case time.Duration:
		return TDuration
	default:
		return InvalidType
	}
}

func ZeroValue(t ValueType) interface{} {
	switch t {
	case TFloat:
		return float64(0)
	case TInt:
		return int64(0)
	case TString:
		return ""
	case TBool:
		return false
	case TRegex:
		return nil
	case TTime:
		return time.Time{}
	case TDuration:
		return time.Duration(0)
	default:
		return errors.New("invalid type")
	}
}

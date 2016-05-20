package stateful

import (
	"regexp"
	"time"

	"github.com/influxdata/kapacitor/tick/ast"
)

type operationKey struct {
	operator  ast.TokenType
	leftType  ValueType
	rightType ValueType
}

var boolTrueResultContainer = resultContainer{BoolValue: true, IsBoolValue: true}
var boolFalseResultContainer = resultContainer{BoolValue: false, IsBoolValue: true}
var emptyResultContainer = resultContainer{}

type evaluationFnInfo struct {
	f          evaluationFn
	returnType ValueType
}

// Constant return types of all binary operations
var binaryConstantTypes map[operationKey]ValueType

func init() {
	// Populate binaryConstantTypes from the evaluationFuncs map
	binaryConstantTypes = make(map[operationKey]ValueType, len(evaluationFuncs))
	for opKey, info := range evaluationFuncs {
		binaryConstantTypes[opKey] = info.returnType
	}
}

var evaluationFuncs = map[operationKey]*evaluationFnInfo{
	// -----------------------------------------
	//	Comparison evaluation funcs

	operationKey{operator: ast.TokenAnd, leftType: TBool, rightType: TBool}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left bool
			var right bool
			var err error

			if left, err = leftNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			// Short circuit evaluation
			if !left {
				return boolFalseResultContainer, nil
			}

			if right, err = rightNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left && right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenOr, leftType: TBool, rightType: TBool}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left bool
			var right bool
			var err error

			if left, err = leftNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			// Short circuit evaluation
			if left {
				return boolTrueResultContainer, nil
			}

			if right, err = rightNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left || right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenEqual, leftType: TBool, rightType: TBool}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left bool
			var right bool
			var err error

			if left, err = leftNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left == right, IsBoolValue: true}, nil
		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenNotEqual, leftType: TBool, rightType: TBool}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left bool
			var right bool
			var err error

			if left, err = leftNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalBool(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left != right, IsBoolValue: true}, nil
		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLess, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left < right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLessEqual, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left <= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenNotEqual, leftType: TInt64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right float64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: float64(left) != right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreaterEqual, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left >= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenEqual, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left == right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenNotEqual, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left != right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenNotEqual, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left != right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLessEqual, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left <= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenEqual, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left == right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreater, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left > right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreater, leftType: TFloat64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right int64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left > float64(right), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreaterEqual, leftType: TFloat64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right int64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left >= float64(right), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenEqual, leftType: TFloat64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right int64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left == float64(right), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLessEqual, leftType: TInt64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right float64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: float64(left) <= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenEqual, leftType: TInt64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right float64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: float64(left) == right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenNotEqual, leftType: TFloat64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right int64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left != float64(right), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLess, leftType: TFloat64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right int64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left < float64(right), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLess, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left < right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreaterEqual, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left >= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreater, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left > right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLessEqual, leftType: TFloat64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right int64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left <= float64(right), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreaterEqual, leftType: TInt64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right float64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: float64(left) >= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreater, leftType: TInt64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right float64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: float64(left) > right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLess, leftType: TInt64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right float64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: float64(left) < right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreater, leftType: TString, rightType: TString}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right string
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left > right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreaterEqual, leftType: TString, rightType: TString}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right string
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left >= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLess, leftType: TString, rightType: TString}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right string
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left < right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLessEqual, leftType: TString, rightType: TString}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right string
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left <= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenEqual, leftType: TString, rightType: TString}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right string
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left == right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenNotEqual, leftType: TString, rightType: TString}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right string
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left != right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenRegexNotEqual, leftType: TString, rightType: TRegex}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right *regexp.Regexp
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalRegex(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: !right.MatchString(left), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenRegexEqual, leftType: TString, rightType: TRegex}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right *regexp.Regexp
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalRegex(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: right.MatchString(left), IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenEqual, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left == right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenNotEqual, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left != right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},
	operationKey{operator: ast.TokenGreater, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left > right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenGreaterEqual, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left >= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLess, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left < right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	operationKey{operator: ast.TokenLessEqual, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return boolFalseResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{BoolValue: left <= right, IsBoolValue: true}, nil

		},
		returnType: TBool,
	},

	// -----------------------------------------
	//	Math evaluation funcs

	operationKey{operator: ast.TokenPlus, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Float64Value: left + right, IsFloat64Value: true}, nil
		},
		returnType: TFloat64,
	},

	operationKey{operator: ast.TokenMinus, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Float64Value: left - right, IsFloat64Value: true}, nil
		},
		returnType: TFloat64,
	},

	operationKey{operator: ast.TokenMult, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Float64Value: left * right, IsFloat64Value: true}, nil
		},
		returnType: TFloat64,
	},

	operationKey{operator: ast.TokenDiv, leftType: TFloat64, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right float64
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Float64Value: left / right, IsFloat64Value: true}, nil
		},
		returnType: TFloat64,
	},

	operationKey{operator: ast.TokenPlus, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Int64Value: left + right, IsInt64Value: true}, nil
		},
		returnType: TInt64,
	},

	operationKey{operator: ast.TokenMinus, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Int64Value: left - right, IsInt64Value: true}, nil
		},
		returnType: TInt64,
	},

	operationKey{operator: ast.TokenMult, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Int64Value: left * right, IsInt64Value: true}, nil
		},
		returnType: TInt64,
	},

	operationKey{operator: ast.TokenDiv, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Int64Value: left / right, IsInt64Value: true}, nil
		},
		returnType: TInt64,
	},

	operationKey{operator: ast.TokenMod, leftType: TInt64, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right int64
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{Int64Value: left % right, IsInt64Value: true}, nil
		},
		returnType: TInt64,
	},

	operationKey{operator: ast.TokenPlus, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: left + right, IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},

	operationKey{operator: ast.TokenMinus, leftType: TDuration, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right time.Duration
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: left - right, IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},

	operationKey{operator: ast.TokenMult, leftType: TDuration, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right int64
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: left * time.Duration(right), IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},

	operationKey{operator: ast.TokenMult, leftType: TInt64, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left int64
			var right time.Duration
			var err error

			if left, err = leftNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: time.Duration(left) * right, IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},

	operationKey{operator: ast.TokenMult, leftType: TDuration, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right float64
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: time.Duration(float64(left) * right), IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},

	operationKey{operator: ast.TokenMult, leftType: TFloat64, rightType: TDuration}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left float64
			var right time.Duration
			var err error

			if left, err = leftNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: time.Duration(left * float64(right)), IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},

	operationKey{operator: ast.TokenDiv, leftType: TDuration, rightType: TInt64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right int64
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalInt(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: left / time.Duration(right), IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},
	operationKey{operator: ast.TokenDiv, leftType: TDuration, rightType: TFloat64}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left time.Duration
			var right float64
			var err error

			if left, err = leftNode.EvalDuration(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalFloat(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{DurationValue: time.Duration(float64(left) / right), IsDurationValue: true}, nil
		},
		returnType: TDuration,
	},

	// -----------------------------------------
	//	String concatenation func

	operationKey{operator: ast.TokenPlus, leftType: TString, rightType: TString}: {
		f: func(scope *Scope, executionState ExecutionState, leftNode, rightNode NodeEvaluator) (resultContainer, *ErrSide) {
			var left string
			var right string
			var err error

			if left, err = leftNode.EvalString(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsLeft: true}
			}

			if right, err = rightNode.EvalString(scope, executionState); err != nil {
				return emptyResultContainer, &ErrSide{error: err, IsRight: true}
			}

			return resultContainer{StringValue: left + right, IsStringValue: true}, nil
		},
		returnType: TString,
	},
}

package pipeline

import "github.com/influxdata/kapacitor/tick/ast"

// The WhereNode filters the data stream by a given expression.
//
// Example:
// var sums = stream
//     |from()
//         .groupBy('service', 'host')
//     |sum('value')
// //Watch particular host for issues.
// sums
//    |where(lambda: "host" == 'h001.example.com')
//    |alert()
//        .crit(lambda: TRUE)
//        .email().to('user@example.com')
//
type WhereNode struct {
	chainnode
	// The expression predicate.
	// tick:ignore
	Expression ast.Node
}

func newWhereNode(wants EdgeType, predicate ast.Node) *WhereNode {
	return &WhereNode{
		chainnode:  newBasicChainNode("where", wants, wants),
		Expression: predicate,
	}
}

// And another expression onto the existing expression.
func (w *WhereNode) Where(expression ast.Node) *WhereNode {
	w.Expression = &ast.BinaryNode{
		Operator: ast.TokenAnd,
		Left:     w.Expression,
		Right:    expression,
	}
	return w
}

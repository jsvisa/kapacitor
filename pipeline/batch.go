package pipeline

import (
	"bytes"
	"reflect"
	"time"
)

// A node that handles creating several child QueryNodes.
// Each call to `query` creates a child batch node that
// can further be configured. See QueryNode
// The `batch` variable in batch tasks is an instance of
// a BatchNode.
//
// Example:
//     var errors = batch
//                      |query('SELECT value from errors')
//                      ...
//     var views = batch
//                      |query('SELECT value from views')
//                      ...
//
// Available Statistics:
//
//    * query_errors -- number of errors when querying
//    * connect_errors -- number of errors connecting to InfluxDB
//    * batches_queried -- number of batches returned from queries
//    * points_queried -- total number of points in batches
//
type BatchNode struct {
	node
}

func newBatchNode() *BatchNode {
	return &BatchNode{
		node: node{
			desc:     "batch",
			wants:    NoEdge,
			provides: BatchEdge,
		},
	}
}

// The query to execute. Must not contain a time condition
// in the `WHERE` clause or contain a `GROUP BY` clause.
// The time conditions are added dynamically according to the period, offset and schedule.
// The `GROUP BY` clause is added dynamically according to the dimensions
// passed to the `groupBy` method.
func (b *BatchNode) Query(q string) *QueryNode {
	n := newQueryNode()
	n.QueryStr = q
	b.linkChild(n)
	return n
}

// Do not add the source batch node to the dot output
// since its not really an edge.
// tick:ignore
func (b *BatchNode) dot(buf *bytes.Buffer) {
}

// A QueryNode defines a source and a schedule for
// processing batch data. The data is queried from
// an InfluxDB database and then passed into the data pipeline.
//
// Example:
// batch
//     |query('''
//         SELECT mean("value")
//         FROM "telegraf"."default".cpu_usage_idle
//         WHERE "host" = 'serverA'
//     ''')
//         .period(1m)
//         .every(20s)
//         .groupBy(time(10s), 'cpu')
//     ...
//
// In the above example InfluxDB is queried every 20 seconds; the window of time returned
// spans 1 minute and is grouped into 10 second buckets.
type QueryNode struct {
	chainnode

	// The query text
	//tick:ignore
	QueryStr string

	// The period or length of time that will be queried from InfluxDB
	Period time.Duration

	// How often to query InfluxDB.
	//
	// The Every property is mutually exclusive with the Cron property.
	Every time.Duration

	// Align start and end times with the Every value
	// Does not apply if Cron is used.
	// tick:ignore
	AlignFlag bool `tick:"Align"`

	// Define a schedule using a cron syntax.
	//
	// The specific cron implementation is documented here:
	// https://github.com/gorhill/cronexpr#implementation
	//
	// The Cron property is mutually exclusive with the Every property.
	Cron string

	// How far back in time to query from the current time
	//
	// For example an Offest of 2 hours and an Every of 5m,
	// Kapacitor will query InfluxDB every 5 minutes for the window of data 2 hours ago.
	//
	// This applies to Cron schedules as well. If the cron specifies to run every Sunday at
	// 1 AM and the Offset is 1 hour. Then at 1 AM on Sunday the data from 12 AM will be queried.
	Offset time.Duration

	// The list of dimensions for the group-by clause.
	//tick:ignore
	Dimensions []interface{} `tick:"GroupBy"`

	// Fill the data.
	// Options are:
	//
	//   - Any numerical value
	//   - null - exhibits the same behavior as the default
	//   - previous - reports the value of the previous window
	//   - none - suppresses timestamps and values where the value is null
	Fill interface{}

	// The name of a configured InfluxDB cluster.
	// If empty the default cluster will be used.
	Cluster string
}

func newQueryNode() *QueryNode {
	b := &QueryNode{
		chainnode: newBasicChainNode("query", BatchEdge, BatchEdge),
	}
	return b
}

//tick:ignore
func (n *QueryNode) ChainMethods() map[string]reflect.Value {
	return map[string]reflect.Value{
		"GroupBy": reflect.ValueOf(n.chainnode.GroupBy),
	}
}

// Group the data by a set of dimensions.
// Can specify one time dimension.
//
// This property adds a `GROUP BY` clause to the query
// so all the normal behaviors when quering InfluxDB with a `GROUP BY` apply.
// More details: https://influxdb.com/docs/v0.9/query_language/data_exploration.html#the-group-by-clause
//
// Example:
//    batch
//        |query(...)
//            .groupBy(time(10s), 'tag1', 'tag2'))
//
// tick:property
func (b *QueryNode) GroupBy(d ...interface{}) *QueryNode {
	b.Dimensions = d
	return b
}

// Align start and stop times for quiries with even boundaries of the QueryNode.Every property.
// Does not apply if using the QueryNode.Cron property.
// tick:property
func (b *QueryNode) Align() *QueryNode {
	b.AlignFlag = true
	return b
}

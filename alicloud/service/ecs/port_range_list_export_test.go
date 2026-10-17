package ecs

// DescribePortRangeListForTest exposes the unexported Read-side describe so the
// acceptance-test CheckDestroy can reuse the resource's real call chain — the
// same strict response-shape guards Read relies on — instead of a parallel
// query. It is declared in a _test.go file, so it is absent from non-test
// builds and does not add a production public API.
var DescribePortRangeListForTest = describePortRangeList

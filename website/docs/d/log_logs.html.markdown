---
layout: "alicloud"
page_title: "Alicloud: alicloud_log_logs"
sidebar_current: "docs-alicloud-datasource-log-logs"
description: |-
    Provides a list of LogService's logs to the user.
---

# alicloud\_log_logs

This data source query logs in a Logstore of a specific project. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29029.htm).

## Example Usage

```
data "alicloud_log_logs" "all" {
    project = "tf-log-project"
    logstore = "tf-log-logstore"
    from = 1544154502
    to = 1544155502
    query = "* and error"
    output_file = "./logs.json"
}
```

## Argument Reference

The following arguments are supported:

* `project` - Name of the project where the log to be queried belongs.
* `logstore` - The name of the Logstore where the log to be queried belongs.
* `from` - The query start time (the number of seconds since 1970-1-1 00:00:00 UTC).
* `to` - The query end time (the number of seconds since 1970-1-1 00:00:00 UTC).
* `query` - The query expression. For more information about the query expression syntax, see  Query syntax.
* `lines` - The maximum number of log lines returned by the request. The maximum number of logs returned from the request. The value range is 0–100 and the default value is 100.
* `offset` - The returned log start point of the request. The value can be 0 or a positive integer. The default value is 0.
* `reverse` - Whether or not logs are returned in reverse order according to the log timestamp.  true indicates reverse order and false indicates sequent order. The default value is false.
* `output_file` - The output file path to save the query's result, file content is a JSON array.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `__time__` - The log timestamp (the number of seconds since 1970-1-1 00:00:00 UTC).
* `__source__` - The log source, which is specified when logs are written.
* `__tag__` - The log tag, which is specified in log group.
* `__topic__` - The log topic, which is specified in log group.
* `[content]` - The original content of the log, which is organized in key-value pairs.

Example output:

```JSON
[
	{
		"__source__": "10.1.2.3",
		"__tag__:__client_ip__": "42.120.75.128",
		"__tag__:__receive_time__": "1544154692",
		"__tag__:tag1": "value1",
		"__tag__:tag2": "value2",
		"__time__": "1544154692",
		"__topic__": "test_topic",
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"provider": "terraform"
	}
]
```

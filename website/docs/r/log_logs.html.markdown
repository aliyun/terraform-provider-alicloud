---
layout: "alicloud"
page_title: "Alicloud: alicloud_log_logs"
sidebar_current: "docs-alicloud-resource-log-logs"
description: |-
  Write log data to a specified Logstore in the following modes.
---

# alicloud\_log\_logs

Write log data to a specified Logstore. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29026.htm).

There are two modes in writing:

* Load balancing mode: Automatically write logs to all writable shards in a Logstore in the load balancing mode. This mode is highly available for writing (SLA: 99.95%), applicable to scenarios in which data writing and consumption are independent of shards , for example, scenarios that do not preserve the order.
* KeyHash mode: A key is required when writing data. Log Service automatically writes data to the shard that meets the key range. For example, hash a producer  (for example, an instance) to a fixed shard based on the name to make sure the data writing and consumption in this shard are strictly ordered (when merging or splitting shards, a key can only appear in one shard at a time point). For more information, see [Shard](https://www.alibabacloud.com/help/doc-detail/28976.htm) .

## Example Usage

Basic Usage

```
resource "alicloud_log_logs" "example" {
      project  = "tf-log-project"
      logstore = "tf-log-store"
      source = "10.1.2.3"
      topic = "test_topic"
      logs = [
        {
          time = 1544154502
          contents = {
            key1 = "value1"
            key2 = "value2"
            key3 = "value3"
           }
        }
      ]
      tags = {
        tag1 = "value1"
        tag2 = "value2"
      }
}
```
## Argument Reference

The following arguments are supported:

* `project` - The name of the Project where logs are to be written.
* `logstore` - The name of the Logstore where logs are to be written.
* `source` - The log source, which is specified when logs are written. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29055.htm).
* `tags` - The log tag, which is specified in log group. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29055.htm).
* `topic` - The log topic, which is specified in log group. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29055.htm).
* `logs` - The logs, [Refer to details](https://www.alibabacloud.com/help/doc-detail/29055.htm): 
  * `time` - The log timestamp (the number of seconds since 1970-1-1 00:00:00 UTC). Default value is the timestamp when logs been send.
  * `contents` - The original content of the log, which is organized in key-value pairs.
* `retry_seconds` - The max retry seconds when post logs fails. Default value is 0, which means no retry.

## Attributes Reference

The following attributes are exported:

* `result` - The send result.



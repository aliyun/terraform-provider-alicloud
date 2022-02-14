---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_saved_search"
sidebar_current: "docs-alicloud-resource-log-saved-search"
description: |-
  Provides a Alicloud log saved search resource.
---

# alicloud\_log\_saved\_search
Log service data saved search management, this service provides the saved search feature to save the required data query and analysis operations. You can use a saved search to quickly perform query and analysis operations.
[Refer to details](https://www.alibabacloud.com/help/en/doc-detail/88985.htm).

-> **NOTE:** Available in 1.157.0+

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "tf-log-project"
  description = "created by terraform"
  tags        = { "test" : "test" }
}
resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "tf-log-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "alicloud_log_saved_search" "example" {
  project_name    = alicloud_log_project.example.name
  logstore_name   = alicloud_log_store.example.name
  search_name    = "tf-test-saved-search"
  search_query      = "* | select count(*) as c,__time__ as t group by t order by t DESC"
  topic      = "tf-test"
  display_name = "test-sls-saved-search"
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `logstore_name` - (Required，ForceNew) The name of the log logstore.
* `search_name` - (Required，ForceNew) The name of the Log Saved Search.
* `search_query` - (Required，ForceNew) Query statement.
* `topic` - (Optional) The topic of the saved search.
* `display_name` - (Optional) Saved Search alias.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log saved search.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when Creating LogSavedSearch instance. 
* `update` - (Defaults to 1 mins) Used when Updating LogSavedSearch instance. 
* `delete` - (Defaults to 1 mins) Used when terminating the LogSavedSearch instance.

## Import

Log saved search can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_saved_search.example tf-log-project:tf-log-logstore:tf-saved-search
```

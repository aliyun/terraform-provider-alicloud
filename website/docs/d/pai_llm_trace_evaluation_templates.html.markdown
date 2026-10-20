---
subcategory: "PAI LLM Trace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_llm_trace_evaluation_templates"
sidebar_current: "docs-alicloud-datasource-pai-llm-trace-evaluation-templates"
description: |-
    Provides the list of PAI LLM Trace evaluation templates.
---

# alicloud_pai_llm_trace_evaluation_templates

This data source provides the list of PAI LLM Trace (`PaiLLMTrace`) evaluation templates returned by
the `GetEvaluationTemplates` API. The API exposes public information, requires no input parameters and
carries no RAM authorization, so it can be read with any valid Alibaba Cloud account. Each template is
passed through as its JSON encoding, in the order returned by the API.

-> **NOTE:** Available since v1.247.0.

## Example Usage

```terraform
data "alicloud_pai_llm_trace_evaluation_templates" "default" {
}

output "templates_count" {
  value = length(data.alicloud_pai_llm_trace_evaluation_templates.default.templates)
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

The following attributes are exported:

* `id` - The provider region, used as the data source identifier.
* `templates` - The evaluation templates, each encoded as a JSON string in the order returned by the
  API. The element shape is free-form public information, so the raw payload is preserved verbatim.

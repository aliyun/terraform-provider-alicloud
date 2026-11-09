# Pull request titles

The title checker resolves every changed resource and data source from the diff
and its registrations. New entities use `New Resource: <full-name>` or
`New Data Source: <full-name>`. Resource and data source names are checked
separately, including when they share the same name.

Code MR titles have a 128-character limit. When several new entity names share a
prefix, one explicit suffix group can preserve all names within that limit:

```text
New Resource: alicloud_sandbox_template{,_cache}; New Data Source: alicloud_sandbox_template{s,_caches}
```

This expands to the resources `alicloud_sandbox_template` and
`alicloud_sandbox_template_cache`, and the data sources
`alicloud_sandbox_templates` and `alicloud_sandbox_template_caches`. An empty
choice includes the prefix itself. A nonempty choice appends that exact suffix.
The checker still requires every new entity in the diff to appear in the
expanded declarations. A longer name does not also declare a shorter prefix.
Declarations are separated by whitespace or `;`; the existing full-name form
continues to work, including comma-separated declarations.

Groups must have at least two distinct choices. Spaces inside groups, nesting,
multiple groups per name, ranges and wildcards are rejected. This syntax does
not waive the rules for modified entities, docs, tests or whole-product
migrations. Ellipses remain invalid.

Both Code and GitHub CI invoke `title_check.go` and restore the checker from the
exact PR base. Merge this policy change first, then update the original MR to
include the new base before using a compact title. A PR cannot grant itself a
title exemption by changing its own checker.

Run the regression suite with `go test ./scripts/title`.

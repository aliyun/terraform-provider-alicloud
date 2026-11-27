# Pull request titles

The title checker resolves every changed resource and data source from the diff
and its registrations. New entities use `New Resource: <full-name>` or
`New Data Source: <full-name>`. Resource and data source names are checked
separately, including when they share the same name.

The `testcase: ` title rule counts Go test files under `alicloud/`. Tests for
CI tools such as `scripts/ci/accscope` and `scripts/title` can use a title
describing the tooling change. When a change also touches provider tests or
documentation, their existing title rules still apply.

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

For a new resource and its companion data source, use two ordered suffixes:

```text
New Resource/Data Source: alicloud_express_connect_router_route_map_association{,s}
```

The first suffix names the resource (`…_association`), the second the data
source (`…_associations`). An empty suffix uses the prefix unchanged.
Use `{,}` when both kinds have exactly the same registered name. Additional
entities still need their own declarations.

When two or more distinct existing entities require more than 128 characters
for their mandatory title markers alone, use `provider: <change summary>`.
For example:

```text
provider: migrate legacy ECS, NAS, CMS and SLS RPC calls to shared request signing
```

CI determines eligibility from the diff, counts resources and data sources
separately, and lists their services and registered names in the job summary.
Repeated changes to one entity count once. At 128 characters or fewer the
existing rules apply. Mixed changes must still explicitly declare new entities
after the summary, e.g. `provider: fix shared signing; New Resource: alicloud_example`.

A new Code MR Title run uses the current target branch checker against the MR
commit. After a rule merges, edit the title and start a new run on the existing
commit. The GitHub workflow uses the PR base checker.

Run the regression suite with `go test ./scripts/title`.

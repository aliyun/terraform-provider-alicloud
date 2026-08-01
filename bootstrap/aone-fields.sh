#!/usr/bin/env bash
# Query and fill required custom fields that block updates to legacy Aone workitems.
# `missing` returns legal candidates, `fill` writes explicit assignments, `auto-fill`
# makes only deterministic choices, and `preflight` provides the fail-closed dispatch
# gate used by the bridge.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"
pools_cfg="$jarvis_root/config/pools.json"
A1="${JARVIS_A1:-${A1_BIN:-$jarvis_root/bin/a1id --}}"

usage() {
    cat >&2 <<'EOF'
Usage:
  aone-fields.sh missing <workitem-id>
  aone-fields.sh fill <workitem-id> <field-id>=<value> [<field-id>=<value> ...]
  aone-fields.sh auto-fill <workitem-id>
  aone-fields.sh preflight <workitem-id> [expected-project-id]
  aone-fields.sh inspect <workitem-id> [expected-project-id]
  aone-fields.sh apply <workitem-id> <expected-project-id> <expected-type-id> \
    <expected-revision> <candidate-digest> <field-id>=<value> [...]
EOF
}

valid_id() {
    case "$1" in
        ""|*[![:alnum:]_.-]*) return 1 ;;
        *) return 0 ;;
    esac
}

sha256_24() {
    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 | awk '{print substr($1,1,24)}'
    else
        sha256sum | awk '{print substr($1,1,24)}'
    fi
}

inspection_digest() {
    local canonical
    canonical="$(jq -S -c '{
      missing:(.missing // []),
      assignments:(.assignments // []),
      unresolved:(.unresolved // [])
    }')" || return 1
    printf '%s' "$canonical" | sha256_24
}

normalize_options() {
    jq -c '
        if type == "array" then .
        elif type == "object" and (.options | type) == "array" then .options
        elif type == "object" and (.Options | type) == "array" then .Options
        elif type == "object" and (.items | type) == "array" then .items
        elif type == "object" and (.Items | type) == "array" then .Items
        elif type == "object" and (.data | type) == "array" then .data
        elif type == "object" and (.Data | type) == "array" then .Data
        elif type == "object" and
             (has("identifier") or has("Identifier") or has("value") or
              has("Value") or has("displayValue") or has("DisplayValue") or
              has("name") or has("Name") or has("path") or has("Path")) then [.]
        else [] end
    '
}

workitem_value() {
    local wi_json="$1" identifier="$2"
    printf '%s' "$wi_json" | jq -r --arg id "$identifier" '
        (.fields // [])
        | map(select((.fieldIdentifier // .identifier // "") == $id))
        | .[0]
        | (.value // .displayValue // empty)
    ' 2>/dev/null
}

context_matches() {
    local wi_json="$1" expected_project="$2" expected_type="$3"
    [ "$(workitem_value "$wi_json" space)" = "$expected_project" ] \
        && [ "$(workitem_value "$wi_json" workitemType)" = "$expected_type" ]
}

assignment_mismatches() {
    local wi_json="$1" assignments="$2"
    jq -c -n --argjson wi "$wi_json" --argjson assignments "$assignments" '
        def ident: (.fieldIdentifier // .identifier // "" | tostring);
        def actual_value:
            if (.value != null and .value != "" and .value != [] and .value != {})
            then (.value | tostring)
            else (.displayValue // "" | tostring)
            end;
        (($wi.fields // []) |
         map({key:ident, value:actual_value}) | from_entries) as $actual |
        [$assignments[] |
         (.value | tostring) as $expected |
         select(($actual[.id] // "") != $expected) |
         {id:.id, expected:$expected, actual:($actual[.id] // "")}]
    '
}

pool_policy() {
    local project_id="$1"
    jq -c --arg p "$project_id" '
        [.pools[]? | select((.project | tostring) == $p)] | .[0] // {} |
        {
          defaults: (.claim_required_field_defaults // {}),
          fallbacks: (.claim_required_field_fallbacks // {}),
          placeholders: (.claim_required_field_placeholders // {})
        }
    ' "$pools_cfg" 2>/dev/null || printf '{"defaults":{},"fallbacks":{},"placeholders":{}}'
}

resolve_missing() {
    local missing_json="$1" policy_json="$2" title="$3"
    jq -c -n --argjson missing "$missing_json" \
        --argjson policy "$policy_json" --arg title "$title" '
        def first_nonempty($values):
            [$values[] | select(. != null) | tostring | select(length > 0)][0] // "";
        def option_value:
            first_nonempty([
              .value, .Value, .identifier, .Identifier, .id, .Id,
              .displayValue, .DisplayValue, .name, .Name, .path, .Path
            ]);
        def option_label:
            first_nonempty([
              .displayValue, .DisplayValue, .name, .Name, .label, .Label,
              .path, .Path, .value, .Value, .identifier, .Identifier
            ]);
        def option_aliases:
            [
              .value, .Value, .identifier, .Identifier, .id, .Id,
              .displayValue, .DisplayValue, .name, .Name, .label, .Label,
              .path, .Path
              | select(. != null) | tostring | select(length > 0) | ascii_downcase
            ] | unique;
        def configured_terms:
            if . == null then []
            elif type == "array" then map(tostring)
            else tostring | split("/")
            end
            | map(select(length > 0) | ascii_downcase) | unique;
        def configured_matches($options; $configured):
            ($configured | configured_terms) as $terms |
            [$options[] | . as $option |
             select(any($terms[]; . as $term |
                    any(($option | option_aliases)[]; . == $term)))]
            | unique_by(option_value);
        def title_tokens:
            ascii_downcase | [scan("[a-z0-9]+") as $token |
            $token | select(length >= 3) |
            select(["alicloud","terraform","resource","data","service","instance",
                    "company","information","group","whitelist","named","sharding",
                    "provider"] | index($token) == null) |
            select(test("[a-z]"))] | unique;
        ($title | title_tokens) as $tokens |
        ($policy.defaults // {}) as $defaults |
        ($policy.fallbacks // {}) as $fallbacks |
        ($policy.placeholders // {}) as $placeholders |
        [$missing[] |
            . as $field |
            (($defaults[$field.id] // $defaults[$field.name] // null)) as $default |
            (($fallbacks[$field.id] // $fallbacks[$field.name] // null)) as $fallback |
            (($placeholders[$field.id] // $placeholders[$field.name] // null)) as $placeholder |
            (($field.options // []) | map(select(option_value != "")) |
             unique_by(option_value)) as $options |
            ($options | map(. as $option | . + {
                _score: ([$tokens[] as $token |
                    select((((($option | option_label) + " " + ($option | option_value)) |
                             ascii_downcase) | contains($token)))] | length)
            })) as $scored |
            (($scored | map(._score) | max) // 0) as $best |
            ($scored | map(select(._score == $best and $best > 0))) as $title_matches |
            (configured_matches($options; $default)) as $default_matches |
            (configured_matches($options; $fallback)) as $fallback_matches |
            # A placeholder is not an answer — it is the pool saying "nobody can
            # decide this one, stamp the field'"'"'s own neutral option so the ticket
            # keeps moving and a human corrects it later". It is therefore only
            # offered on the genuinely-undecidable branch below, never on a
            # broken configured default/fallback (that is a config bug to fix,
            # not something to paper over), and only when it resolves to exactly
            # one legal option — otherwise no key is emitted at all.
            ((configured_matches($options; $placeholder)) as $m |
             if $placeholder != null and ($m | length) == 1
             then {placeholder: {value: ($m[0] | option_value),
                                 displayValue: ($m[0] | option_label)}}
             else {} end) as $placeholder_hint |
            if ($field.optionsLookupError // false) then
                {kind:"unresolved", id:$field.id, name:$field.name,
                 reason:"options_lookup_error", options:[]}
            elif $default != null then
                # An operator-configured default outranks title_match on purpose.
                # title_match is substring token scoring over the whole candidate
                # set; on a large org-wide tree it produces confident nonsense
                # (a Terraform provider ticket scored "Cloud Control API" from the
                # tokens "cloud"+"api"). Once a pool pins the answer in
                # pools.json, the heuristic must not be able to override it, and
                # a pinned-but-unmatched default stays fail-closed rather than
                # silently falling through to guessing.
                if ($default_matches | length) == 1 then
                    {kind:"assignment", id:$field.id, name:$field.name,
                     value:($default_matches[0] | option_value),
                     source:"configured_default"}
                else
                    {kind:"unresolved", id:$field.id, name:$field.name,
                     reason:"configured_default_not_unique",
                     configured:$default,
                     options:($options | map({value:option_value, displayValue:option_label}))}
                end
            elif ($title_matches | length) == 1 then
                {kind:"assignment", id:$field.id, name:$field.name,
                 value:($title_matches[0] | option_value), source:"title_match"}
            elif ($options | length) == 1 then
                {kind:"assignment", id:$field.id, name:$field.name,
                 value:($options[0] | option_value), source:"single_option"}
            elif $fallback != null then
                if ($fallback_matches | length) == 1 then
                    {kind:"assignment", id:$field.id, name:$field.name,
                     value:($fallback_matches[0] | option_value),
                     source:"validated_fallback"}
                else
                    {kind:"unresolved", id:$field.id, name:$field.name,
                     reason:"configured_fallback_not_unique",
                     configured:$fallback,
                     options:($options | map({value:option_value, displayValue:option_label}))}
                end
            else
                {kind:"unresolved", id:$field.id, name:$field.name,
                 reason:(if ($options | length) == 0 then "no_legal_options"
                         elif $best == 0 then "no_title_match"
                         else "ambiguous_title_match" end),
                 options:($options | map({value:option_value, displayValue:option_label}))}
                + $placeholder_hint
            end
        ] as $rows |
        {assignments:[$rows[] | select(.kind=="assignment") | del(.kind)],
         unresolved:[$rows[] | select(.kind=="unresolved") | del(.kind)]}
    '
}

# Field metadata (`field list` / `field options`) is per project*type and turns
# over on the order of weeks, but one dispatch burst re-reads it on every
# inspect — and `apply` alone inspects twice (CAS pre-check + readback). Cache
# it so the repair path stops paying ~3.8s of a1 round-trips per inspect; that
# per-inspect cost is what pushes field repair past its timeout when several
# headless sessions run at once.
#
# TTL is deliberately short. This is the candidate set that gates a write, so
# the window in which a removed option could still be offered is kept to
# minutes — long enough to cover the concurrency burst that causes the pile-up,
# short enough that `apply`'s own re-validation and readback stay the real
# guard rather than a stale file.
#
# The workitem itself is deliberately NOT cached: `apply` compares the live
# revision and candidate digest as its concurrency guard, so a stale workitem
# read would silently defeat that check. Set JARVIS_FIELD_META_TTL=0 to bypass.
field_meta_ttl() { echo "${JARVIS_FIELD_META_TTL:-900}"; }

field_list_json() {
    local project_id="$1" type_id="$2" ttl
    ttl="$(field_meta_ttl)"
    if [ "$ttl" -le 0 ] 2>/dev/null; then
        $A1 project workitem field list --project "$project_id" --type "$type_id" -f json
        return
    fi
    bash "$script_dir/cache.sh" get \
        "aonefields_list_${project_id}_${type_id}" "$ttl" -- \
        $A1 project workitem field list --project "$project_id" --type "$type_id" -f json
}

field_options_json() {
    local field_id="$1" project_id="$2" type_id="$3" ttl
    ttl="$(field_meta_ttl)"
    if [ "$ttl" -le 0 ] 2>/dev/null; then
        $A1 project workitem field options "$field_id" \
            --project "$project_id" --type "$type_id" -f json
        return
    fi
    bash "$script_dir/cache.sh" get \
        "aonefields_opts_${project_id}_${type_id}_${field_id}" "$ttl" -- \
        $A1 project workitem field options "$field_id" \
            --project "$project_id" --type "$type_id" -f json
}

# Correlate one already-fetched workitem against its required field definitions.
# Callers that already hold the workitem JSON (inspect) pass it straight in; the
# `missing` subcommand fetches it first. Re-execing the script to reach this had
# every inspect pay for a second `workitem get` of the same object — ~2.5s of
# pure duplication, and two reads that can disagree under replica lag.
compute_missing() {
    local wi_json="$1" project_id="$2" type_id="$3"
    local defs_json wi_tmp defs_tmp missing_json result field_id options_json normalized rc

    if ! defs_json="$(field_list_json "$project_id" "$type_id")"; then
        echo "aone-fields.sh: failed to list fields for project $project_id type $type_id" >&2
        return 1
    fi

    wi_tmp="$(mktemp)"
    defs_tmp="$(mktemp)"
    printf '%s' "$wi_json" > "$wi_tmp"
    printf '%s' "$defs_json" > "$defs_tmp"

    missing_json="$(jq -c -n --slurpfile wi "$wi_tmp" --slurpfile defs "$defs_tmp" '
        def ident: (.fieldIdentifier // .identifier // "" | tostring);
        def current_value:
            if (.value != null and .value != "" and .value != [] and .value != {})
            then .value else (.displayValue // "") end;
        def empty_value: . == null or . == "" or . == [] or . == {};
        def definitions:
            $defs[0] | if type == "array" then .
            elif (.fields | type) == "array" then .fields
            elif (.items | type) == "array" then .items
            elif (.data | type) == "array" then .data
            else [] end;
        (($wi[0].fields // []) | map({key: ident, value: current_value}) | from_entries) as $current |
        [definitions[]
            | select(.isRequired == true or .isRequired == "true")
            | select((.sourceType // "") as $s | ["system", "basic", "service.sprint"] | index($s) == null)
            | ident as $id
            | select($id != "")
            | ($current[$id] // "") as $value
            | select($value | empty_value)
            | {
                id: $id,
                name: (.displayName // .name // ""),
                description: (.description // ""),
                format: (.format // ""),
                sourceType: (.sourceType // ""),
                current: $value,
                options: (if (.options | type) == "array" then .options else [] end),
                optionsSource: "field_list"
              }
        ]
    ')"
    rc=$?
    rm -f "$wi_tmp" "$defs_tmp"
    if [ "$rc" -ne 0 ]; then
        echo "aone-fields.sh: failed to correlate workitem values with field definitions" >&2
        return 1
    fi

    result="$missing_json"
    while IFS= read -r field_id; do
        [ -n "$field_id" ] || continue
        if options_json="$(field_options_json "$field_id" "$project_id" "$type_id" 2>/dev/null)" \
                && normalized="$(printf '%s' "$options_json" | normalize_options 2>/dev/null)"; then
            result="$(printf '%s' "$result" | jq -c --arg id "$field_id" --argjson opts "$normalized" '
                map(if .id == $id then .options = $opts | .optionsSource = "field_options" else . end)
            ')"
        else
            echo "aone-fields.sh: warning: failed to load legal options for field $field_id" >&2
            result="$(printf '%s' "$result" | jq -c --arg id "$field_id" '
                map(if .id == $id then .optionsLookupError = true else . end)
            ')"
        fi
    done < <(printf '%s' "$missing_json" | jq -r '.[] | select((.options | length) == 0) | .id')

    printf '%s\n' "$result"
}

preflight_result() {
    local status="$1" error_type="$2" workitem_id="$3" project_id="$4"
    local type_id="$5" assignments="$6" unresolved="$7" readback="$8"
    local filled="$9" failure_reason="${10:-}"
    jq -c -n --arg status "$status" --arg error_type "$error_type" \
        --arg workitem_id "$workitem_id" --arg project "$project_id" \
        --arg workitem_type "$type_id" --argjson assignments "$assignments" \
        --argjson unresolved "$unresolved" --argjson readback "$readback" \
        --argjson filled "$filled" --arg failure_reason "$failure_reason" '
        {
          status:$status,
          errorType:(if $error_type == "" then null else $error_type end),
          workitemId:$workitem_id,
          project:$project,
          workitemType:$workitem_type,
          assignments:$assignments,
          unresolved:$unresolved,
          readback:$readback,
          filled:$filled
        }
        + (if $failure_reason == "" then {} else {failureReason:$failure_reason} end)
    '
}

cmd="${1:-}"
case "$cmd" in
    inspect)
        [ "$#" -eq 2 ] || [ "$#" -eq 3 ] || { usage; exit 2; }
        workitem_id="$2"
        expected_project="${3:-}"
        valid_id "$workitem_id" || {
            echo "aone-fields.sh: invalid workitem id '$workitem_id'" >&2
            exit 2
        }
        [ -z "$expected_project" ] || valid_id "$expected_project" || {
            echo "aone-fields.sh: invalid expected project id '$expected_project'" >&2
            exit 2
        }

        if ! wi_json="$($A1 project workitem get "$workitem_id" -f json 2>/dev/null)"; then
            jq -c -n --arg id "$workitem_id" \
                '{status:"failed",errorType:"field_inspection_failed",
                  failureReason:"workitem_lookup_error",workitemId:$id}'
            exit 3
        fi
        project_id="$(workitem_value "$wi_json" space)"
        type_id="$(workitem_value "$wi_json" workitemType)"
        title="$(printf '%s' "$wi_json" | jq -r '
            .title // ((.fields // []) |
            map(select((.fieldIdentifier // .identifier // "") == "title")) |
            .[0] | (.displayValue // .value // "")) // ""
        ' 2>/dev/null)"
        description="$(printf '%s' "$wi_json" | jq -r '
            .description // .body // ((.fields // []) |
            map(select((.fieldIdentifier // .identifier // "") == "description")) |
            .[0] | (.displayValue // .value // "")) // ""
        ' 2>/dev/null)"
        revision="$(printf '%s' "$wi_json" | jq -r '
            .modifiedAt // .modified // .gmtModified // .updatedAt //
            .updateTime // .gmtModifiedAt // ""
        ' 2>/dev/null)"
        if [ -z "$project_id" ] || [ -z "$type_id" ]; then
            jq -c -n --arg id "$workitem_id" --arg project "$project_id" \
                --arg type "$type_id" \
                '{status:"failed",errorType:"field_inspection_failed",
                  failureReason:"context_lookup_error",workitemId:$id,
                  project:$project,workitemType:$type}'
            exit 3
        fi
        if [ -n "$expected_project" ] && [ "$project_id" != "$expected_project" ]; then
            jq -c -n --arg id "$workitem_id" --arg project "$project_id" \
                --arg type "$type_id" --arg expected "$expected_project" \
                '{status:"failed",errorType:"field_inspection_failed",
                  failureReason:"project_mismatch",workitemId:$id,
                  project:$project,workitemType:$type,expectedProject:$expected}'
            exit 3
        fi
        if ! missing_json="$(compute_missing "$wi_json" "$project_id" "$type_id")"; then
            jq -c -n --arg id "$workitem_id" --arg project "$project_id" \
                --arg type "$type_id" \
                '{status:"failed",errorType:"field_inspection_failed",
                  failureReason:"required_fields_lookup_error",workitemId:$id,
                  project:$project,workitemType:$type}'
            exit 3
        fi
        if [ "$(printf '%s' "$missing_json" | jq 'length')" -eq 0 ]; then
            jq -c -n --arg id "$workitem_id" --arg project "$project_id" \
                --arg type "$type_id" --arg revision "$revision" \
                --arg title "$title" --arg description "$description" '
                {
                  status:"ready",errorType:null,workitemId:$id,project:$project,
                  workitemType:$type,revision:$revision,title:$title,
                  description:$description,missing:[],assignments:[],unresolved:[]
                }'
            exit 0
        fi
        policy="$(pool_policy "$project_id")"
        if ! resolution="$(resolve_missing "$missing_json" "$policy" "$title")"; then
            jq -c -n --arg id "$workitem_id" --arg project "$project_id" \
                --arg type "$type_id" \
                '{status:"failed",errorType:"field_inspection_failed",
                  failureReason:"resolution_error",workitemId:$id,
                  project:$project,workitemType:$type}'
            exit 3
        fi
        assignments="$(printf '%s' "$resolution" | jq -c '.assignments')"
        unresolved="$(printf '%s' "$resolution" | jq -c '.unresolved')"
        jq -c -n --arg id "$workitem_id" --arg project "$project_id" \
            --arg type "$type_id" --arg revision "$revision" \
            --arg title "$title" --arg description "$description" \
            --argjson missing "$missing_json" --argjson assignments "$assignments" \
            --argjson unresolved "$unresolved" '
            {
              status:"repair_required",errorType:null,workitemId:$id,
              project:$project,workitemType:$type,revision:$revision,
              title:$title,description:$description,missing:$missing,
              assignments:$assignments,unresolved:$unresolved
            }'
        ;;

    apply)
        [ "$#" -ge 7 ] || { usage; exit 2; }
        workitem_id="$2"
        expected_project="$3"
        expected_type="$4"
        expected_revision="$5"
        expected_digest="$6"
        valid_id "$workitem_id" || {
            echo "aone-fields.sh: invalid workitem id '$workitem_id'" >&2
            exit 2
        }
        valid_id "$expected_project" || {
            echo "aone-fields.sh: invalid expected project id '$expected_project'" >&2
            exit 2
        }
        valid_id "$expected_type" || {
            echo "aone-fields.sh: invalid expected workitem type '$expected_type'" >&2
            exit 2
        }
        case "$expected_digest" in
            [0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f][0-9a-f]) ;;
            *) echo "aone-fields.sh: invalid candidate digest" >&2; exit 2 ;;
        esac
        shift 6
        requested='[]'
        cfs_args=()
        for spec in "$@"; do
            case "$spec" in
                *=*) field_id="${spec%%=*}"; value="${spec#*=}" ;;
                *) echo "aone-fields.sh: invalid field assignment '$spec'" >&2; exit 2 ;;
            esac
            if ! valid_id "$field_id" || [ -z "$value" ]; then
                echo "aone-fields.sh: invalid field assignment '$spec'" >&2
                exit 2
            fi
            requested="$(printf '%s' "$requested" | jq -c \
                --arg id "$field_id" --arg value "$value" \
                '. + [{id:$id,value:$value}]')"
            cfs_args+=(--cfs "$field_id=$value")
        done
        if [ "$(printf '%s' "$requested" | jq '[.[].id] | unique | length')" \
                -ne "$(printf '%s' "$requested" | jq 'length')" ]; then
            echo "aone-fields.sh: duplicate field assignment" >&2
            exit 4
        fi
        if ! before="$(bash "$0" inspect "$workitem_id" "$expected_project")"; then
            echo "aone-fields.sh: failed to inspect immediately before apply" >&2
            exit 3
        fi
        if [ "$(printf '%s' "$before" | jq -r '.status')" != "repair_required" ]; then
            echo "aone-fields.sh: workitem no longer requires the proposed repair" >&2
            exit 3
        fi
        actual_type="$(printf '%s' "$before" | jq -r '.workitemType // ""')"
        actual_revision="$(printf '%s' "$before" | jq -r '.revision // ""')"
        actual_digest="$(printf '%s' "$before" | inspection_digest)"
        if [ "$actual_type" != "$expected_type" ] \
                || [ "$actual_revision" != "$expected_revision" ] \
                || [ "$actual_digest" != "$expected_digest" ]; then
            printf 'aone-fields.sh: repair snapshot drift: type=%s revision=%s digest=%s\n' \
                "$actual_type" "$actual_revision" "$actual_digest" >&2
            exit 3
        fi
        validation="$(jq -c -n --argjson before "$before" \
            --argjson requested "$requested" '
            def first_nonempty($values):
              [$values[] | select(. != null) | tostring | select(length > 0)][0] // "";
            def option_value:
              first_nonempty([
                .value, .Value, .identifier, .Identifier, .id, .Id,
                .displayValue, .DisplayValue, .name, .Name, .path, .Path
              ]);
            ($before.missing | map({key:(.id|tostring),value:.}) | from_entries) as $fields |
            ($before.missing | map(.id|tostring) | sort) as $expected |
            ($requested | map(.id|tostring) | sort) as $actual |
            {
              complete:($expected == $actual),
              illegal:[
                $requested[] |
                . as $row |
                ($fields[$row.id].options // []) as $options |
                select(($options | map(option_value) | index($row.value)) == null) |
                {id:$row.id,value:$row.value}
              ]
            }
        ')"
        if [ "$(printf '%s' "$validation" | jq -r '.complete')" != "true" ] \
                || [ "$(printf '%s' "$validation" | jq '.illegal | length')" -ne 0 ]; then
            printf 'aone-fields.sh: assignments rejected by current candidate set: %s\n' \
                "$validation" >&2
            exit 4
        fi
        if ! update_output="$($A1 project workitem update "$workitem_id" \
                "${cfs_args[@]}" 2>&1)"; then
            printf 'aone-fields.sh: repair update failed: %s\n' "$update_output" >&2
            exit 3
        fi
        if ! after="$(bash "$0" inspect "$workitem_id" "$expected_project")"; then
            echo "aone-fields.sh: repair readback failed" >&2
            exit 3
        fi
        if [ "$(printf '%s' "$after" | jq -r '.status')" != "ready" ]; then
            printf 'aone-fields.sh: repair readback still blocked: %s\n' "$after" >&2
            exit 3
        fi
        if ! actual_wi="$($A1 project workitem get "$workitem_id" -f json 2>/dev/null)"; then
            echo "aone-fields.sh: canonical value readback failed" >&2
            exit 3
        fi
        readback="$(jq -c -n --argjson wi "$actual_wi" \
            --argjson requested "$requested" '
            def ident: (.fieldIdentifier // .identifier // "" | tostring);
            def actual_value:
              if (.value != null and .value != "" and .value != [] and .value != {})
              then (.value | tostring)
              else (.displayValue // "" | tostring)
              end;
            (($wi.fields // []) |
             map({key:ident,value:actual_value}) | from_entries) as $actual |
            [$requested[] | {id:.id,value:($actual[.id] // "")}]
        ')"
        mismatches="$(jq -c -n --argjson requested "$requested" \
            --argjson readback "$readback" '
            ($readback | map({key:.id,value:.value}) | from_entries) as $actual |
            [$requested[] | select(($actual[.id] // "") != .value) |
             {id:.id,expected:.value,actual:($actual[.id] // "")}]
        ')"
        if [ "$(printf '%s' "$mismatches" | jq 'length')" -ne 0 ]; then
            jq -c -n --arg id "$workitem_id" --arg project "$expected_project" \
                --arg type "$expected_type" --arg revision "$expected_revision" \
                --argjson assignments "$requested" --argjson readback "$readback" \
                --argjson mismatches "$mismatches" '
                {
                  status:"failed",errorType:"field_apply_readback_mismatch",
                  failureReason:"assignment_conflict_after_readback",
                  workitemId:$id,project:$project,workitemType:$type,
                  revision:$revision,assignments:$assignments,
                  readback:$readback,mismatches:$mismatches,filled:true,
                  missing:[],unresolved:[]
                }'
            exit 3
        fi
        printf '%s' "$after" | jq -c --argjson assignments "$requested" \
            --argjson readback "$readback" \
            '. + {assignments:$assignments,readback:$readback,filled:true}'
        ;;

    missing)
        [ "$#" -eq 2 ] || { usage; exit 2; }
        workitem_id="$2"
        valid_id "$workitem_id" || { echo "aone-fields.sh: invalid workitem id '$workitem_id'" >&2; exit 2; }

        if ! wi_json="$($A1 project workitem get "$workitem_id" -f json)"; then
            echo "aone-fields.sh: failed to get workitem $workitem_id" >&2
            exit 1
        fi
        type_id="$(workitem_value "$wi_json" workitemType)"
        project_id="$(workitem_value "$wi_json" space)"
        if [ -z "$type_id" ] || [ -z "$project_id" ]; then
            echo "aone-fields.sh: cannot resolve workitem type/project for $workitem_id" >&2
            exit 1
        fi
        compute_missing "$wi_json" "$project_id" "$type_id" || exit 1
        ;;

    fill)
        [ "$#" -ge 3 ] || { usage; exit 2; }
        workitem_id="$2"
        valid_id "$workitem_id" || { echo "aone-fields.sh: invalid workitem id '$workitem_id'" >&2; exit 2; }
        shift 2
        cfs_args=()
        for spec in "$@"; do
            case "$spec" in
                *=*) field_id="${spec%%=*}"; value="${spec#*=}" ;;
                *) echo "aone-fields.sh: invalid field assignment '$spec' (want <field>=<value>)" >&2; exit 2 ;;
            esac
            if ! valid_id "$field_id" || [ -z "$value" ]; then
                echo "aone-fields.sh: invalid field assignment '$spec'" >&2
                exit 2
            fi
            cfs_args+=(--cfs "$field_id=$value")
        done
        if output="$($A1 project workitem update "$workitem_id" "${cfs_args[@]}" 2>&1)"; then
            printf 'aone-fields.sh: filled %s:' "$workitem_id"
            printf ' %s' "$@"
            printf '\n'
        else
            rc=$?
            printf '%s\n' "$output" >&2
            exit "$rc"
        fi
        ;;

    auto-fill)
        [ "$#" -eq 2 ] || { usage; exit 2; }
        workitem_id="$2"
        valid_id "$workitem_id" || { echo "aone-fields.sh: invalid workitem id '$workitem_id'" >&2; exit 2; }

        if ! wi_json="$($A1 project workitem get "$workitem_id" -f json)"; then
            echo "aone-fields.sh: failed to get workitem $workitem_id for inference" >&2
            exit 3
        fi
        project_id="$(workitem_value "$wi_json" space)"
        type_id="$(workitem_value "$wi_json" workitemType)"
        if [ -z "$project_id" ] || [ -z "$type_id" ]; then
            echo "aone-fields.sh: cannot resolve workitem type/project for $workitem_id" >&2
            exit 3
        fi
        if ! missing_json="$(compute_missing "$wi_json" "$project_id" "$type_id")"; then
            echo "aone-fields.sh: could not inspect missing fields for $workitem_id" >&2
            exit 3
        fi
        title="$(printf '%s' "$wi_json" | jq -r '
            .title // ((.fields // []) | map(select((.fieldIdentifier // .identifier // "") == "title")) |
            .[0] | (.displayValue // .value // "")) // ""
        ' 2>/dev/null)"
        policy="$(pool_policy "$project_id")"
        if ! resolution="$(resolve_missing "$missing_json" "$policy" "$title")"; then
            echo "aone-fields.sh: failed to resolve deterministic field values" >&2
            exit 3
        fi

        unresolved_count="$(printf '%s' "$resolution" | jq '.unresolved | length')"
        assignment_count="$(printf '%s' "$resolution" | jq '.assignments | length')"
        if [ "$unresolved_count" -ne 0 ] || [ "$assignment_count" -eq 0 ]; then
            printf 'aone-fields.sh: automatic fill unresolved: %s\n' "$resolution" >&2
            printf '%s\n' "$resolution"
            exit 3
        fi

        cfs_args=()
        while IFS=$'\t' read -r field_id value; do
            [ -n "$field_id" ] && [ -n "$value" ] || continue
            cfs_args+=(--cfs "$field_id=$value")
        done < <(printf '%s' "$resolution" | jq -r '.assignments[] | [.id,.value] | @tsv')
        if output="$($A1 project workitem update "$workitem_id" "${cfs_args[@]}" 2>&1)"; then
            printf '%s\n' "$(printf '%s' "$resolution" | jq -c '. + {filled:true}')"
        else
            rc=$?
            printf 'aone-fields.sh: automatic fill update failed: %s\n' "$output" >&2
            printf '%s\n' "$(printf '%s' "$resolution" | jq -c --arg error "$output" \
                '. + {filled:false, updateError:$error}')"
            exit "$rc"
        fi
        ;;

    preflight)
        [ "$#" -eq 2 ] || [ "$#" -eq 3 ] || { usage; exit 2; }
        workitem_id="$2"
        expected_project="${3:-}"
        valid_id "$workitem_id" || {
            echo "aone-fields.sh: invalid workitem id '$workitem_id'" >&2
            exit 2
        }
        [ -z "$expected_project" ] || valid_id "$expected_project" || {
            echo "aone-fields.sh: invalid expected project id '$expected_project'" >&2
            exit 2
        }

        if ! wi_json="$($A1 project workitem get "$workitem_id" -f json 2>/dev/null)"; then
            preflight_result failed preflight_validation_failed "$workitem_id" "" \
                "" '[]' '[]' '[]' false workitem_lookup_error
            exit 3
        fi
        project_id="$(workitem_value "$wi_json" space)"
        type_id="$(workitem_value "$wi_json" workitemType)"
        title="$(printf '%s' "$wi_json" | jq -r '
            .title // ((.fields // []) |
            map(select((.fieldIdentifier // .identifier // "") == "title")) |
            .[0] | (.displayValue // .value // "")) // ""
        ' 2>/dev/null)"
        if [ -z "$project_id" ] || [ -z "$type_id" ]; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" '[]' '[]' '[]' false context_lookup_error
            exit 3
        fi
        if [ -n "$expected_project" ] && [ "$project_id" != "$expected_project" ]; then
            unresolved="$(jq -c -n --arg expected "$expected_project" \
                --arg actual "$project_id" \
                '[{reason:"project_mismatch",expectedProject:$expected,actualProject:$actual}]')"
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" '[]' "$unresolved" '[]' false project_mismatch
            exit 3
        fi
        if ! missing_json="$(compute_missing "$wi_json" "$project_id" "$type_id")"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" '[]' '[]' '[]' false required_fields_lookup_error
            exit 3
        fi
        if [ "$(printf '%s' "$missing_json" | jq 'length')" -eq 0 ]; then
            if ! noop_readback_json="$($A1 project workitem get "$workitem_id" \
                    -f json 2>/dev/null)"; then
                preflight_result failed preflight_validation_failed "$workitem_id" \
                    "$project_id" "$type_id" '[]' '[]' '[]' false \
                    context_read_noop_error
                exit 3
            fi
            if ! context_matches "$noop_readback_json" "$project_id" "$type_id"; then
                preflight_result failed preflight_validation_failed "$workitem_id" \
                    "$project_id" "$type_id" '[]' '[]' '[]' false \
                    context_drift_noop_readback
                exit 3
            fi
            preflight_result ok "" "$workitem_id" "$project_id" "$type_id" \
                '[]' '[]' '[]' false
            exit 0
        fi

        policy="$(pool_policy "$project_id")"
        if ! resolution="$(resolve_missing "$missing_json" "$policy" "$title")"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" '[]' '[]' "$missing_json" false resolution_error
            exit 3
        fi
        assignments="$(printf '%s' "$resolution" | jq -c '.assignments')"
        unresolved="$(printf '%s' "$resolution" | jq -c '.unresolved')"
        if [ "$(printf '%s' "$unresolved" | jq 'length')" -ne 0 ] \
                || [ "$(printf '%s' "$assignments" | jq 'length')" -eq 0 ]; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" "$unresolved" \
                "$missing_json" false unresolved_required_fields
            exit 3
        fi

        cfs_args=()
        while IFS=$'\t' read -r field_id value; do
            [ -n "$field_id" ] && [ -n "$value" ] || continue
            cfs_args+=(--cfs "$field_id=$value")
        done < <(printf '%s' "$assignments" | jq -r '.[] | [.id,.value] | @tsv')
        if ! before_update_json="$($A1 project workitem get "$workitem_id" -f json \
                2>/dev/null)"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" '[]' "$missing_json" \
                false context_read_before_update_error
            exit 3
        fi
        if ! context_matches "$before_update_json" "$project_id" "$type_id"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" '[]' "$missing_json" \
                false context_drift_before_update
            exit 3
        fi
        if ! update_output="$($A1 project workitem update "$workitem_id" \
                "${cfs_args[@]}" 2>&1)"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" '[]' "$missing_json" \
                false update_error
            exit 3
        fi
        # Post-update readback must re-read the workitem itself — the pre-update
        # wi_json cannot show the values we just wrote. Only the field metadata
        # is allowed to come from cache here.
        if ! readback_wi_json="$($A1 project workitem get "$workitem_id" -f json)" \
                || ! readback="$(compute_missing "$readback_wi_json" "$project_id" "$type_id")"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" '[]' '[]' true readback_error
            exit 3
        fi
        if [ "$(printf '%s' "$readback" | jq 'length')" -ne 0 ]; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" "$readback" "$readback" \
                true readback_not_empty
            exit 3
        fi
        if ! final_wi_json="$($A1 project workitem get "$workitem_id" -f json \
                2>/dev/null)"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" '[]' "$readback" \
                true context_read_after_readback_error
            exit 3
        fi
        if ! context_matches "$final_wi_json" "$project_id" "$type_id"; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" '[]' "$readback" \
                true context_drift_after_readback
            exit 3
        fi
        mismatches="$(assignment_mismatches "$final_wi_json" "$assignments")"
        if [ "$(printf '%s' "$mismatches" | jq 'length')" -ne 0 ]; then
            preflight_result failed preflight_validation_failed "$workitem_id" \
                "$project_id" "$type_id" "$assignments" "$mismatches" \
                "$readback" true assignment_conflict_after_readback
            exit 3
        fi
        preflight_result ok "" "$workitem_id" "$project_id" "$type_id" \
            "$assignments" '[]' "$readback" true
        ;;

    help|-h|--help)
        usage
        ;;

    *)
        usage
        exit 2
        ;;
esac

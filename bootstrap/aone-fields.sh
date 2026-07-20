#!/usr/bin/env bash
# Query and fill required custom fields that block updates to legacy Aone workitems.
# `missing` returns legal candidates, `fill` writes explicit assignments, and `auto-fill`
# makes only deterministic choices: a configured pool default, a single legal option, or
# one uniquely matching product option inferred from the workitem title.

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
EOF
}

valid_id() {
    case "$1" in
        ""|*[![:alnum:]_.-]*) return 1 ;;
        *) return 0 ;;
    esac
}

normalize_options() {
    jq -c '
        if type == "array" then .
        elif type == "object" and (.options | type) == "array" then .options
        elif type == "object" and (.items | type) == "array" then .items
        elif type == "object" and (.data | type) == "array" then .data
        elif type == "object" and (has("identifier") or has("value") or has("displayValue")) then [.]
        else [] end
    '
}

cmd="${1:-}"
case "$cmd" in
    missing)
        [ "$#" -eq 2 ] || { usage; exit 2; }
        workitem_id="$2"
        valid_id "$workitem_id" || { echo "aone-fields.sh: invalid workitem id '$workitem_id'" >&2; exit 2; }

        if ! wi_json="$($A1 project workitem get "$workitem_id" -f json)"; then
            echo "aone-fields.sh: failed to get workitem $workitem_id" >&2
            exit 1
        fi
        type_id="$(printf '%s' "$wi_json" | jq -r '
            (.fields // []) | map(select((.fieldIdentifier // .identifier // "") == "workitemType")) |
            .[0] | (.value // .displayValue // empty)
        ' 2>/dev/null)"
        project_id="$(printf '%s' "$wi_json" | jq -r '
            (.fields // []) | map(select((.fieldIdentifier // .identifier // "") == "space")) |
            .[0] | (.value // .displayValue // empty)
        ' 2>/dev/null)"
        if [ -z "$type_id" ] || [ -z "$project_id" ]; then
            echo "aone-fields.sh: cannot resolve workitem type/project for $workitem_id" >&2
            exit 1
        fi

        if ! defs_json="$($A1 project workitem field list --project "$project_id" --type "$type_id" -f json)"; then
            echo "aone-fields.sh: failed to list fields for project $project_id type $type_id" >&2
            exit 1
        fi

        wi_tmp="$(mktemp)"
        defs_tmp="$(mktemp)"
        trap 'rm -f "$wi_tmp" "$defs_tmp"' EXIT
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
        ')" || {
            echo "aone-fields.sh: failed to correlate workitem values with field definitions" >&2
            exit 1
        }

        result="$missing_json"
        while IFS= read -r field_id; do
            [ -n "$field_id" ] || continue
            if options_json="$($A1 project workitem field options "$field_id" \
                    --project "$project_id" --type "$type_id" -f json 2>/dev/null)" \
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

        if ! missing_json="$(bash "$0" missing "$workitem_id")"; then
            echo "aone-fields.sh: could not inspect missing fields for $workitem_id" >&2
            exit 3
        fi
        if ! wi_json="$($A1 project workitem get "$workitem_id" -f json)"; then
            echo "aone-fields.sh: failed to get workitem $workitem_id for inference" >&2
            exit 3
        fi
        project_id="$(printf '%s' "$wi_json" | jq -r '
            (.fields // []) | map(select((.fieldIdentifier // .identifier // "") == "space")) |
            .[0] | (.value // .displayValue // empty)
        ' 2>/dev/null)"
        title="$(printf '%s' "$wi_json" | jq -r '
            .title // ((.fields // []) | map(select((.fieldIdentifier // .identifier // "") == "title")) |
            .[0] | (.displayValue // .value // "")) // ""
        ' 2>/dev/null)"
        defaults="$(jq -c --arg p "$project_id" '
            [.pools[]? | select((.project | tostring) == $p) |
             (.claim_required_field_defaults // {})] | .[0] // {}
        ' "$pools_cfg" 2>/dev/null || printf '{}')"
        if ! resolution="$(jq -c -n --argjson missing "$missing_json" \
                --argjson defaults "$defaults" --arg title "$title" '
            def option_value: (.value // .identifier // .id // .displayValue // .name // "" | tostring);
            def option_label: (.displayValue // .name // .label // .value // .identifier // "" | tostring);
            def title_tokens:
                ascii_downcase | [scan("[a-z0-9]+") as $token |
                $token | select(length >= 3) |
                select(["alicloud","terraform","resource","data","service","instance",
                        "company","information","group","whitelist","named","sharding"] |
                       index($token) == null) |
                select(test("[a-z]"))] | unique;
            ($title | title_tokens) as $tokens |
            [$missing[] |
                . as $field |
                (($defaults[$field.id] // $defaults[$field.name] // "") | tostring) as $default |
                (($field.options // []) | map(select(option_value != ""))) as $options |
                ($options | map(. as $option | . + {
                    _score: ([$tokens[] as $token |
                        select((((($option | option_label) + " " + ($option | option_value)) |
                                 ascii_downcase) | contains($token)))] | length)
                })) as $scored |
                (($scored | map(._score) | max) // 0) as $best |
                ($scored | map(select(._score == $best and $best > 0))) as $matches |
                if $default != "" then
                    {kind:"assignment", id:$field.id, name:$field.name,
                     value:$default, source:"configured_default"}
                elif ($options | length) == 1 then
                    {kind:"assignment", id:$field.id, name:$field.name,
                     value:($options[0] | option_value), source:"single_option"}
                elif ($matches | length) == 1 then
                    {kind:"assignment", id:$field.id, name:$field.name,
                     value:($matches[0] | option_value), source:"title_match"}
                else
                    {kind:"unresolved", id:$field.id, name:$field.name,
                     reason:(if ($options | length) == 0 then "no_legal_options"
                             elif $best == 0 then "no_title_match"
                             else "ambiguous_title_match" end),
                     options:($options | map({value:option_value, displayValue:option_label}))}
                end
            ] as $rows |
            {assignments:[$rows[] | select(.kind=="assignment") | del(.kind)],
             unresolved:[$rows[] | select(.kind=="unresolved") | del(.kind)]}
        ')"; then
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

    help|-h|--help)
        usage
        ;;

    *)
        usage
        exit 2
        ;;
esac

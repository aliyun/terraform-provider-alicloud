#!/usr/bin/env bash
# bootstrap/probe-corpus.sh — F3 Corpus-gen: 场景语料库批量生成器
#
# 从 provider 仓 website/docs/r/<name>.html.markdown 的 Example Usage 抽 HCL,机械改造成一个
# probe 场景三件套(scenario.yaml/main.tf/checks.md),直落 playground git 数据仓 <product>/<id>/。
# 目标:从 website docs 批量生成场景语料,喂给 bootstrap/probe.sh run 做 tier-1 生命周期探测。
#
# 子命令:
#   gen <alicloud_resource> [--force]    单资源生成(--force 覆盖已存在)
#   gen --batch <N>                      批量:全量 website 资源 diff 掉已有,免费族优先,生成 N 个
#   validate [--all | <id>...]           质量门:terraform init+validate+fmt -check;失败移 _quarantine/<id>/
#
# 机械改造(gen):
#   1. pin provider 版本(config .provider.version)、注入 variable "run_id"
#   2. 剥离示例里的 terraform/provider/backend 块(heredoc 安全)
#   3. 引用到的额外 provider(random 等)补进 required_providers(hashicorp/ 命名空间)
#   4. 可命名 name/*_name 字面量/var.name → "probe-${var.run_id}"(已带插值的名称不动)
#   5. 示例既有 tags 块注入 managed_by = "<tags_marker>"
#
# 成本安全门(值敏感,config .tier1_risk_denylist):命中订阅语义(charge_value_fields 字面量值 ∈ subscription_values,
#   或独立 period 取订阅时长)→ 有对应 data source 则 apply:false + ds- 只读变体(引存量);无 ds 则 apply:true +
#   allow_prepaid:true 放行真跑;按量付费大件一律直接 apply(无按量资源名清单)。
#
# 产物落点:playground 解析与 probe.sh 同(env JARVIS_TF_PLAYGROUND > config paths.playground_dir > workspace.sh dir tf_playground > 默认约定)。
# provider 仓解析:env JARVIS_PROBE_PROVIDER_DIR > bootstrap/workspace.sh dir terraform_provider。
# 环境:PROBE_TERRAFORM_BIN — terraform 二进制(默认 terraform;gen 规整/validate 用,便于单测桩)。
#
# 复用 probe.sh 的路径/解析器(source,不触发其 main)。bash 3.2 兼容(无 mapfile;空数组守卫)。
set -uo pipefail

_corpus_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/probe.sh
source "$_corpus_dir/probe.sh"

_corpus_tf() { echo "${PROBE_TERRAFORM_BIN:-terraform}"; }

# ── 文档解析 ─────────────────────────────────────────────────────────
_corpus_subcategory() { # <docfile> — frontmatter subcategory 原值
    sed -n 's/^subcategory:[[:space:]]*"\(.*\)"/\1/p' "$1" 2>/dev/null | head -1
}

# _corpus_product <alicloud_name> <docfile> — 产品目录名:
#   源码 RpcPost 三元组 product 小写(取到即用)> docs subcategory 净化(小写/去括号空格后)> misc
_corpus_product() {
    local name="$1" doc="$2" short src prod sub
    short="${name#alicloud_}"
    src="$(probe_provider_dir)/alicloud/resource_alicloud_${short}.go"
    if [ -f "$src" ]; then
        prod="$(_source_pv "$src" 2>/dev/null | head -1 | cut -f1 | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9')"
        [ -n "$prod" ] && { printf '%s' "$prod"; return; }
    fi
    sub="$(_corpus_subcategory "$doc" | tr '[:upper:]' '[:lower:]' | sed -E 's/[[:space:](].*$//' | tr -cd 'a-z0-9')"
    [ -n "$sub" ] && { printf '%s' "$sub"; return; }
    printf 'misc'
}

# _corpus_extract_example <docfile> — 抽 ## Example Usage 区内第一个含 resource 的代码块(裸 HCL)
_corpus_extract_example() {
    awk '
    /^## Example Usage/ { inex=1; next }
    inex && !incode && /^## / { inex=0 }
    inex && /^```/ {
        if (incode) { incode=0; if (buf ~ /resource[ \t]+"/) { printf "%s", buf; exit }; buf="" }
        else { incode=1; buf="" }
        next
    }
    inex && incode { buf = buf $0 "\n" }
    ' "$1" 2>/dev/null
}

# _corpus_extract_data_block <docfile> — 抽 ## Example Usage 区内第一个 data "alicloud_..." 块(brace 平衡,
#   只取 data 块本身,天然只读、apply 安全)。供订阅类资源的 ds- 只读变体。
_corpus_extract_data_block() {
    awk '
    /^## Example Usage/ { inex=1; next }
    inex && !inc && /^## / { inex=0 }
    inex && /^```/ { inc = !inc; next }
    inex && inc {
        if (!cap && $0 ~ /^[[:space:]]*data[[:space:]]+"alicloud_[a-z0-9_]+"[[:space:]]+"[^"]*"[[:space:]]*\{/) { cap=1; depth=0 }
        if (cap) {
            print
            o=gsub(/\{/,"{"); c=gsub(/\}/,"}"); depth+=o-c
            if (depth<=0) exit
        }
    }
    ' "$1" 2>/dev/null
}

# ── 机械改造管道(stdin→stdout) ─────────────────────────────────────
# 剥离顶层 terraform/provider/backend 块;heredoc 内容原样保留、其 brace 不参与深度计数
_corpus_strip_blocks() {
    awk '
    BEGIN { depth=0; drop=0; dropdepth=0; inhd=0; term="" }
    {
        line=$0
        if (inhd) {
            print line
            t=line; sub(/^[[:space:]]+/,"",t); sub(/[[:space:]]+$/,"",t)
            if (t == term) { inhd=0; term="" }
            next
        }
        if (drop==0 && match(line, /<<-?"?[A-Za-z_][A-Za-z0-9_]*/)) {
            hs=substr(line, RSTART, RLENGTH); tm=hs; sub(/^<<-?"?/,"",tm)
            print line; inhd=1; term=tm; next
        }
        if (depth==0 && drop==0 && line ~ /^[[:space:]]*(terraform|provider|backend)([[:space:]]+"[^"]*")?[[:space:]]*\{/) {
            drop=1; dropdepth=0
        }
        o=gsub(/\{/,"{",line); c=gsub(/\}/,"}",line)
        if (drop==1) {
            dropdepth += o - c
            if (dropdepth<=0) drop=0
            next
        }
        depth += o - c
        if (depth<0) depth=0
        print line
    }
    '
}

# 引用到的非 alicloud provider 短名(random/null/time/tls...)去重
_corpus_extra_providers() {
    grep -oE '(resource|data)[[:space:]]+"[a-z0-9]+_' 2>/dev/null \
        | sed -E 's/.*"([a-z0-9]+)_/\1/' | sort -u | grep -vx 'alicloud' || true
}

# name/*_name 字面量(无插值)或 var.xxx → "probe-${var.run_id}"
_corpus_inject_names() {
    awk '
    {
        line=$0
        if (match(line, /^[[:space:]]*[A-Za-z0-9_]*name[[:space:]]*=[[:space:]]*/)) {
            head=substr(line,1,RLENGTH); rhs=substr(line,RLENGTH+1); sub(/[[:space:]]*$/,"",rhs)
            if (rhs ~ /^"[^"$]*"$/ || rhs ~ /^var\.[A-Za-z0-9_]+$/) {
                k=head; sub(/[[:space:]]*=[[:space:]]*$/,"",k)
                print k " = \"probe-${var.run_id}\""
                next
            }
        }
        print line
    }
    '
}

# 示例既有多行 tags 块 → 注入 managed_by = "<marker>"
_corpus_inject_tags() {
    awk -v marker="$1" '
    {
        print $0
        if ($0 ~ /^[[:space:]]*tags[[:space:]]*=?[[:space:]]*\{[[:space:]]*$/) {
            match($0,/^[[:space:]]*/); ind=substr($0,1,RLENGTH)
            print ind "  managed_by = \"" marker "\""
        }
    }
    '
}

# main.tf 里 alicloud 资源类型并集(逗号连接)
_corpus_resources_list() {
    grep -oE 'resource[[:space:]]+"alicloud_[a-z0-9_]+"' 2>/dev/null \
        | sed -E 's/.*"(alicloud_[a-z0-9_]+)"/\1/' | sort -u | tr '\n' ',' | sed 's/,$//'
}

# 注入头部(pin 版本 + variable run_id + 额外 provider),参数=额外 provider 短名
_corpus_header() {
    local ver req p
    ver="$(cfg '.provider.version')"
    req="$(cfg '.terraform.required_version')"
    { [ -z "$req" ] || [ "$req" = "null" ]; } && req=">= 1.5.0"
    printf 'terraform {\n'
    printf '  required_version = "%s"\n' "$req"
    printf '  required_providers {\n'
    printf '    alicloud = {\n      source  = "aliyun/alicloud"\n      version = "%s"\n    }\n' "$ver"
    for p in "$@"; do
        printf '    %s = {\n      source = "hashicorp/%s"\n    }\n' "$p" "$p"
    done
    printf '  }\n}\n\nvariable "run_id" {\n  type = string\n}\n'
}

# 成本安全门(值敏感):命中订阅语义 → stdout 打 reason + 返回 0(apply:false);否则返回 1(直接 apply)。
#   规则:charge_value_fields 的字面量值 ∈ subscription_values(大小写不敏感),或独立 period 字段取
#   纯整数订阅时长(1..period_subscription_max)。按量值(PostPaid/PayAsYouGo)、秒级 metric period、
#   retention_period 等一律放行。$1=资源名(未用,保留签名),$2=main.tf。
_corpus_gate() {
    local tf="$2" fields subs periods pmax out
    fields="$(cfg '.tier1_risk_denylist.charge_value_fields | join(",")')"
    subs="$(cfg '.tier1_risk_denylist.subscription_values | join(",")')"
    periods="$(cfg '.tier1_risk_denylist.period_fields | join(",")')"
    pmax="$(cfg '.tier1_risk_denylist.period_subscription_max')"
    case "$pmax" in ''|*[!0-9]*) pmax=12 ;; esac
    out="$(awk -v fields="$fields" -v subs="$subs" -v periods="$periods" -v pmax="$pmax" '
    BEGIN{
        n=split(fields,fa,","); for(i=1;i<=n;i++) if(fa[i]!="") F[fa[i]]=1
        m=split(subs,sa,",");   for(i=1;i<=m;i++) if(sa[i]!="") S[tolower(sa[i])]=1
        p=split(periods,pa,","); for(i=1;i<=p;i++) if(pa[i]!="") P[pa[i]]=1
    }
    match($0,/^[[:space:]]*[A-Za-z0-9_]+[[:space:]]*=/){
        key=$0; sub(/^[[:space:]]*/,"",key); sub(/[[:space:]]*=.*/,"",key)
        val=$0; sub(/^[^=]*=[[:space:]]*/,"",val); sub(/[[:space:]]*$/,"",val); gsub(/"/,"",val)
        if((key in P) && val ~ /^[0-9]+$/ && (val+0)>=1 && (val+0)<=pmax){ print "subscription_period:" key "=" val; exit }
        if((key in F) && (tolower(val) in S)){ print "subscription_value:" key "=" val; exit }
    }
    ' "$tf" 2>/dev/null)"
    if [ -n "$out" ]; then printf '%s' "$out"; return 0; fi
    return 1
}

# _corpus_find_ds_doc <short> — 找对应 data source 文档(试 <short>/<short>s/<short>es);打印路径,无则退 1
_corpus_find_ds_doc() {
    local short="$1" prov cand
    prov="$(probe_provider_dir)"
    for cand in "$short" "${short}s" "${short}es"; do
        if [ -f "$prov/website/docs/d/$cand.html.markdown" ]; then echo "$prov/website/docs/d/$cand.html.markdown"; return 0; fi
    done
    return 1
}

# ── ds- 只读变体(订阅类资源规范) ───────────────────────────────────
# _corpus_gen_ds <name> <short> <id> <product> <playdir> <products> <force>
#   有 data source 文档时抽首个 data 块 + 加 output → 落 <playdir>/<product>/ds-<id>/。data 天然只读、apply 安全。
_corpus_gen_ds() {
    local name="$1" short="$2" id="$3" product="$4" playdir="$5" products="$6" force="$7"
    local dsdoc
    dsdoc="$(_corpus_find_ds_doc "$short")" || { echo "  ds: $name 无对应 data source 文档(website/docs/d/),仅留规范注记"; return 0; }

    local block
    block="$(_corpus_extract_data_block "$dsdoc")"
    [ -n "$block" ] || { echo "  ds: $name data source 文档无可抽 data 块,仅留规范注记"; return 0; }

    local body extras_raw header
    body="$(printf '%s' "$block" | _corpus_strip_blocks)"
    extras_raw="$(printf '%s' "$body" | _corpus_extra_providers)"
    if [ -n "$extras_raw" ]; then
        local ea=(); local ln
        while IFS= read -r ln; do [ -n "$ln" ] && ea+=("$ln"); done <<< "$extras_raw"
        header="$(_corpus_header "${ea[@]}")"
    else
        header="$(_corpus_header)"
    fi

    local head1 dtype dname
    head1="$(printf '%s' "$body" | grep -m1 -E '^[[:space:]]*data[[:space:]]+"alicloud_')"
    dtype="$(printf '%s' "$head1" | sed -E 's/.*data[[:space:]]+"([^"]+)".*/\1/')"
    dname="$(printf '%s' "$head1" | sed -E 's/.*data[[:space:]]+"[^"]+"[[:space:]]+"([^"]+)".*/\1/')"

    local dsid="ds-$id" dest
    dest="$playdir/$product/$dsid"
    if [ -n "$(_find_scenario_dirs "$dsid" | head -1)" ] && [ "$force" -ne 1 ]; then
        echo "  ds: skip $dsid (已存在)"; return 0
    fi
    mkdir -p "$dest"
    {
        printf '%s\n\n%s\n' "$header" "$body"
        if [ -n "$dtype" ] && [ -n "$dname" ]; then
            printf '\noutput "ds_result" {\n  value = data.%s.%s\n}\n' "$dtype" "$dname"
        fi
    } > "$dest/main.tf"
    if [ -n "${PROBE_TERRAFORM_BIN:-}" ] || command -v terraform >/dev/null 2>&1; then
        "$(_corpus_tf)" fmt "$dest" >/dev/null 2>&1 || true
    fi
    {
        printf 'id: %s\n' "$dsid"
        printf 'title: %s data source 只读探测(订阅类资源规范:引用存量, 待人工校订)\n' "$name"
        printf 'persona: importer\n'
        printf 'products: %s\n' "$products"
        printf 'resources: %s\n' "${dtype:-data.$short}"
        printf 'cost: free\n'
        printf 'detect: validate_fail,plan_fail\n'
        printf 'update_step: false\n'
        printf 'import_check: false\n'
        printf 'source_docs: https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/%s\n' "$(basename "$dsdoc" .html.markdown)"
        printf 'generated_by: jarvis-probe-corpus\n'
        printf 'origin: generated\n'
    } > "$dest/scenario.yaml"
    {
        printf '# checks: %s(data source 只读探测, jarvis-probe-corpus 生成骨架)\n\n' "$dsid"
        printf '订阅类资源 `%s` 规范:勿真实创建,用 data source 引用存量。本场景纯 data + output,天然 apply 安全。\n' "$name"
        printf '> 待人工核对:data source 参数/过滤条件是否需补(origin: generated)。\n\n'
        printf '## 探测点\n'
        printf '1. validate 通过。\n'
        printf '2. data source 读路径可用(apply 只读、无创建;plan 无副作用)。\n'
    } > "$dest/checks.md"
    echo "  ds: $name → $product/$dsid (data source 只读变体, apply 安全)"
    return 0
}

# ── gen:单资源 ──────────────────────────────────────────────────────
_corpus_gen() {
    local name="" force=0
    while [ $# -gt 0 ]; do
        case "$1" in
            --force) force=1; shift ;;
            -*) echo "gen: 未知参数 '$1'" >&2; return 2 ;;
            *) if [ -z "$name" ]; then name="$1"; shift; else echo "gen: 多余参数 '$1'" >&2; return 2; fi ;;
        esac
    done
    [ -n "$name" ] || { echo "gen: 需要 <alicloud_resource>" >&2; return 2; }
    case "$name" in alicloud_*) ;; *) name="alicloud_$name" ;; esac

    local short id doc
    short="${name#alicloud_}"
    id="$(printf '%s' "$short" | tr '_' '-')"
    doc="$(probe_provider_dir)/website/docs/r/${short}.html.markdown"
    [ -f "$doc" ] || { echo "gen: 找不到文档 $doc" >&2; return 2; }

    local raw
    raw="$(_corpus_extract_example "$doc")"
    [ -n "$raw" ] || { echo "gen: $name Example Usage 无含 resource 的 HCL 块,跳过" >&2; return 2; }

    # 机械改造:strip → inject names → inject tags
    local marker body
    marker="$(cfg '.corpus.tags_marker')"; { [ -z "$marker" ] || [ "$marker" = "null" ]; } && marker="jarvis-probe"
    body="$(printf '%s' "$raw" | _corpus_strip_blocks | _corpus_inject_names | _corpus_inject_tags "$marker")"

    # 额外 provider(空数组守卫,bash 3.2)
    local extras_raw header
    extras_raw="$(printf '%s' "$body" | _corpus_extra_providers)"
    if [ -n "$extras_raw" ]; then
        local ea=(); local ln
        while IFS= read -r ln; do [ -n "$ln" ] && ea+=("$ln"); done <<< "$extras_raw"
        header="$(_corpus_header "${ea[@]}")"
    else
        header="$(_corpus_header)"
    fi

    # diff:同 id 已存在(任意 product)→ 非 --force 跳过
    local existing
    existing="$(_find_scenario_dirs "$id" | head -1)"
    if [ -n "$existing" ] && [ "$force" -ne 1 ]; then
        echo "gen: skip $name (id=$id 已存在 @ $existing)"; return 0
    fi

    local product playdir dest
    product="$(_corpus_product "$name" "$doc")"
    playdir="$(probe_playground_dir)"
    dest="$playdir/$product/$id"
    mkdir -p "$dest"
    printf '%s\n\n%s\n' "$header" "$body" > "$dest/main.tf"
    if [ -n "${PROBE_TERRAFORM_BIN:-}" ] || command -v terraform >/dev/null 2>&1; then
        "$(_corpus_tf)" fmt "$dest" >/dev/null 2>&1 || true
    fi

    # 成本安全门(值敏感)。命中订阅语义后按「有无 data source」再分岔:
    #   有 ds → apply:false + ds- 只读变体(引存量,不真实创建);
    #   无 ds → apply:true + allow_prepaid:true(放行 runner 运行时 prepaid_guard 兜底),真实创建探测全生命周期。
    local gate_reason="" cost="free" apply_yaml="" allow_prepaid_yaml="" has_ds=0
    if gate_reason="$(_corpus_gate "$name" "$dest/main.tf")"; then
        cost="paid"
        if _corpus_find_ds_doc "$short" >/dev/null 2>&1; then
            has_ds=1; apply_yaml="apply: false"
        else
            allow_prepaid_yaml="allow_prepaid: true"
        fi
    fi

    # scenario.yaml
    local products resources
    products="$(_corpus_subcategory "$doc")"; [ -n "$products" ] || products="$(printf '%s' "$product" | tr '[:lower:]' '[:upper:]')"
    resources="$(printf '%s' "$body" | _corpus_resources_list)"; [ -n "$resources" ] || resources="$name"
    {
        printf 'id: %s\n' "$id"
        printf 'title: %s 生成场景(website docs Example Usage 骨架, 待人工校订)\n' "$name"
        printf 'persona: beginner\n'
        printf 'products: %s\n' "$products"
        printf 'resources: %s\n' "$resources"
        printf 'cost: %s\n' "$cost"
        printf 'detect: validate_fail,plan_fail,apply_fail,perpetual_diff,unexpected_replace,destroy_fail\n'
        printf 'update_step: false\n'
        printf 'import_check: false\n'
        [ -n "$apply_yaml" ] && printf '%s\n' "$apply_yaml"
        [ -n "$allow_prepaid_yaml" ] && printf '%s\n' "$allow_prepaid_yaml"
        printf 'source_docs: https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/%s\n' "$short"
        printf 'generated_by: jarvis-probe-corpus\n'
        printf 'origin: generated\n'
    } > "$dest/scenario.yaml"

    # checks.md
    {
        printf '# checks: %s(jarvis-probe-corpus 生成骨架, 待人工校订)\n\n' "$id"
        printf '来源:website/docs/r/%s.html.markdown 的 Example Usage 机械抽取 + 改造' "$short"
        printf '(pin 版本 / 注入 run_id / 标签 / 剥离 provider 块)。\n'
        printf '> 本骨架字段与期望**待人工核对官方文档后收敛**,勿直接当成品(origin: generated)。\n\n'
        printf '## 探测点\n'
        printf '1. validate 通过(官方示例组合应可解析)。\n'
        printf '2. apply(或 plan)后立即 plan 空 diff(幂等性)。\n'
        printf '3. destroy 干净、state 清空。\n'
        if [ "$has_ds" -eq 1 ]; then
            printf '\n> 注:命中成本安全门(订阅语义:%s),有对应 data source → `apply: false`,runner 止步 plan。\n' "$gate_reason"
            printf '> **订阅类资源规范**:勿真实创建,改用 **data source 引用存量**探测读路径(ds- 只读变体见同 product 目录)。\n'
        elif [ -n "$allow_prepaid_yaml" ]; then
            printf '\n> 注:命中成本安全门(订阅语义:%s),但**无对应 data source 文档** → 翻 `apply: true` + `allow_prepaid: true`\n' "$gate_reason"
            printf '(放行 runner 运行时 prepaid_guard 兜底),真实创建探测全生命周期。\n'
            printf '> **告警**:订阅/包年包月资源 **destroy 可能失败**,届时 runner cleanup 触发 `destroy_fail`(S1)流程兜底上报,需人工清理残留。\n'
        fi
    } > "$dest/checks.md"

    if [ -n "$apply_yaml" ]; then
        echo "gen: $name → $product/$id (apply=false 订阅有ds:$gate_reason)"
        _corpus_gen_ds "$name" "$short" "$id" "$product" "$playdir" "$products" "$force"
    elif [ -n "$allow_prepaid_yaml" ]; then
        echo "gen: $name → $product/$id (apply=true+allow_prepaid 订阅无ds:$gate_reason)"
    else
        echo "gen: $name → $product/$id (apply=true)"
    fi
    return 0
}

# ── gen --batch ─────────────────────────────────────────────────────
_corpus_batch() {
    local n="${1:-}"
    case "$n" in ''|*[!0-9]*) echo "gen --batch 需要正整数 N" >&2; return 2 ;; esac
    local prov; prov="$(probe_provider_dir)"
    [ -d "$prov/website/docs/r" ] || { echo "batch: provider docs 不可用($prov/website/docs/r)" >&2; return 2; }

    # free prefixes(空数组守卫)
    local fps_raw; fps_raw="$(cfg '.corpus.free_prefixes[]' 2>/dev/null)"

    local prio="" rest="" f short id skip_ct=0 matched fp
    shopt -s nullglob
    for f in "$prov"/website/docs/r/*.html.markdown; do
        short="$(basename "$f" .html.markdown)"
        id="$(printf '%s' "$short" | tr '_' '-')"
        if [ -n "$(_find_scenario_dirs "$id" | head -1)" ]; then skip_ct=$((skip_ct+1)); continue; fi
        matched=0
        if [ -n "$fps_raw" ]; then
            while IFS= read -r fp; do
                [ -z "$fp" ] && continue
                case "$short" in "$fp"*) matched=1; break ;; esac
            done <<< "$fps_raw"
        fi
        if [ "$matched" -eq 1 ]; then prio="$prio$short"$'\n'; else rest="$rest$short"$'\n'; fi
    done

    # 排序:prio 组字典序在前,rest 组字典序在后
    local ordered gen_ct=0 r
    ordered="$( { [ -n "$prio" ] && printf '%s' "$prio" | grep . | sort; [ -n "$rest" ] && printf '%s' "$rest" | grep . | sort; } )"
    local cand_ct; cand_ct="$(printf '%s' "$ordered" | grep -c . 2>/dev/null)"; [ -z "$cand_ct" ] && cand_ct=0

    while IFS= read -r r; do
        [ -z "$r" ] && continue
        [ "$gen_ct" -ge "$n" ] && break
        if _corpus_gen "alicloud_$r"; then gen_ct=$((gen_ct+1)); fi
    done <<< "$ordered"

    echo "batch: generated=$gen_ct skipped(已存在)=$skip_ct requested=$n candidates=$cand_ct"
    return 0
}

# ── validate ────────────────────────────────────────────────────────
# 质量门:每场景 terraform init(共享插件缓存)+ validate + fmt -check;任一失败 → 移 _quarantine/<id>/ + reason
_corpus_validate() {
    local all=0; local ids=""
    while [ $# -gt 0 ]; do
        case "$1" in
            --all) all=1; shift ;;
            -*) echo "validate: 未知参数 '$1'" >&2; return 2 ;;
            *) ids="$ids $1"; shift ;;
        esac
    done
    local tf; tf="$(_corpus_tf)"
    command -v "$tf" >/dev/null 2>&1 || { echo "validate: 无 terraform 二进制($tf)" >&2; return 2; }
    local playdir; playdir="$(probe_playground_dir)"

    # 目标场景目录清单
    local targets="" d
    if [ "$all" -eq 1 ]; then
        shopt -s nullglob
        for d in "$playdir"/*/*/; do
            case "$d" in *"/_quarantine/"*) continue ;; esac
            [ -f "$d/scenario.yaml" ] && targets="$targets${d%/}"$'\n'
        done
    else
        local one dd
        for one in $ids; do
            dd="$(_find_scenario_dirs "$one" | grep -v '/_quarantine/' | head -1)"
            [ -n "$dd" ] && targets="$targets$dd"$'\n' || echo "validate: 场景 '$one' 未找到,跳过" >&2
        done
    fi
    [ -n "$(printf '%s' "$targets" | grep .)" ] || { echo "validate: 无可校验场景" >&2; return 2; }

    local cache; cache="$(probe_workdir_base)/.plugin-cache"; mkdir -p "$cache"
    export TF_IN_AUTOMATION=1 TF_INPUT=0 TF_PLUGIN_CACHE_DIR="$cache"

    # 共享工作目录:按目标场景的 provider 并集(alicloud + 引用到的 hashicorp/* 如 random)init 一次,
    # 之后逐场景把 main.tf 换进来只 validate(离线复用已装 provider),避免每场景 init 的 registry 往返。
    local sdir wd ver extras
    ver="$(cfg '.provider.version')"
    extras="$(while IFS= read -r sdir; do [ -f "$sdir/main.tf" ] && grep -hoE 'hashicorp/[a-z0-9]+' "$sdir/main.tf"; done <<< "$targets" | sed 's#hashicorp/##' | sort -u)"
    wd="$(mktemp -d)"
    {
        printf 'terraform {\n  required_providers {\n'
        printf '    alicloud = { source = "aliyun/alicloud", version = "%s" }\n' "$ver"
        while IFS= read -r p; do [ -n "$p" ] && printf '    %s = { source = "hashicorp/%s" }\n' "$p" "$p"; done <<< "$extras"
        printf '  }\n}\n'
    } > "$wd/_providers.tf"
    if ! ( cd "$wd" && "$tf" init -input=false >"$wd/.init.log" 2>&1 ); then
        echo "validate: 共享 init 失败(见下),无法离线校验" >&2
        tail -n 20 "$wd/.init.log" >&2 2>/dev/null
        rm -rf "$wd"; return 2
    fi
    rm -f "$wd/_providers.tf"   # 留 .terraform/.terraform.lock.hcl,清模板 tf(避免与场景 terraform{} 块冲突)

    local pass_ct=0 quar_ct=0 id step reason
    while IFS= read -r sdir; do
        [ -z "$sdir" ] && continue
        case "$sdir" in *"/_quarantine/"*) continue ;; esac   # 已隔离体不再校验(防自毁)
        id="$(basename "$sdir")"
        step=""; reason=""
        rm -f "$wd"/*.tf
        cp "$sdir"/*.tf "$wd"/ 2>/dev/null
        # validate(共享 init 已装 provider,离线)
        if ! ( cd "$wd" && "$tf" validate >"$wd/.validate.log" 2>&1 ); then
            step="validate"; reason="$(tail -n 20 "$wd/.validate.log" 2>/dev/null)"
        fi
        # fmt -check(校验存档格式,直接对场景目录)
        if [ -z "$step" ] && ! ( "$tf" fmt -check "$sdir" >"$wd/.fmt.log" 2>&1 ); then
            step="fmt"; reason="terraform fmt -check 不通过(未规整)"$'\n'"$(cat "$wd/.fmt.log" 2>/dev/null)"
        fi

        if [ -n "$step" ]; then
            # 隔离到 _quarantine/<product>/<id>/(深一层,避开 list/_find_scenario_dirs 的两级 glob,不污染 run)
            local product qdir
            product="$(basename "$(dirname "$sdir")")"
            qdir="$playdir/_quarantine/$product/$id"
            mkdir -p "$playdir/_quarantine/$product"
            rm -rf "$qdir"            # 清旧隔离体,mv 才能整目录搬入
            mv "$sdir" "$qdir"
            {
                printf 'quarantine: %s/%s\n' "$product" "$id"
                printf 'failed_step: %s\n' "$step"
                printf 'when: %s\n' "$(date -u +%FT%TZ)"
                printf 'source_dir: %s\n\n' "$sdir"
                printf '%s\n' "$reason"
            } > "$qdir/QUARANTINE_REASON.txt"
            echo "validate: $id QUARANTINE(step=$step)"
            quar_ct=$((quar_ct+1))
        else
            echo "validate: $id PASS"
            pass_ct=$((pass_ct+1))
        fi
    done <<< "$targets"
    rm -rf "$wd"

    echo "validate: pass=$pass_ct quarantine=$quar_ct"
    [ "$quar_ct" -eq 0 ]
}

# ── CLI ─────────────────────────────────────────────────────────────
_corpus_usage() { sed -n '2,20p' "$0"; }

corpus_main() {
    local cmd="${1:-}"; shift 2>/dev/null || true
    case "$cmd" in
        gen)
            if [ "${1:-}" = "--batch" ]; then shift; _corpus_batch "$@"; else _corpus_gen "$@"; fi ;;
        validate) _corpus_validate "$@" ;;
        -h|--help|"") _corpus_usage ;;
        *) echo "probe-corpus: 未知命令 '$cmd'" >&2; _corpus_usage; exit 2 ;;
    esac
}

[[ "${BASH_SOURCE[0]}" == "$0" ]] && corpus_main "$@"

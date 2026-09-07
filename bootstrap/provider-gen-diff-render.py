#!/usr/bin/env python3
"""
provider-gen-diff-render.py — render provider-gen-diff.sh --no-color output
into a polished, localized (Chinese) HTML page with side-by-side AND unified
diff views (side-by-side default, global toggle).

USAGE
  python3 provider-gen-diff-render.py <input.txt> <output.html>
  provider-gen-diff.sh ... --html   # the tool invokes this

The input is the full --no-color stdout of provider-gen-diff.sh, which has a
stable structure: a metadata header, `── <label> ──` diff-section headers,
`diff --git ...` unified hunks, and a `== checklist ==` block. The renderer is
agnostic to the resource/PR — any provider-gen-diff run produces the same shape.
"""
import html
import os
import re
import sys


# ── bug catalogue (Chinese descriptions) ─────────────────────────────────────
BUG_DESC = {
    '1': 'Create-only 属性缺少 ForceNew(手改 = 补 ForceNew: true)',
    '2': 'Read 未反向映射(手改 = 从 objectRaw/nested 做 d.Set 回填)',
    '3': 'Datasource 嵌套 schema 缺字段(手改 = 补 Elem/schema 键)',
    '4': '文档多余 NOTE(手改 = 删除 NOTE)',
    '5': '文档示例是占位文字(手改 = 手写真示例 HCL)',
    '6': 'SetId 后 addDebug 顺序(观察项,无稳定 grep——看上方 diff 里的 Create)',
    '7': 'Update 触发器只对部分字段(手改 = 加 HasChange 分支)',
}

# Strip ANSI color escapes so the renderer tolerates colored input too.
ANSI_RE = re.compile(r'\x1b\[[0-9;]*m')


def esc(s):
    return html.escape(s, quote=False)


# ── parse the tool's text output ──────────────────────────────────────────────
def parse(text):
    lines = [ANSI_RE.sub('', l) for l in text.split('\n')]
    meta_lines, sections = [], []
    i, n = 0, len(lines)

    def is_header(s):
        return s.startswith('── ') or s.startswith('== checklist')

    while i < n and not is_header(lines[i]):
        if lines[i].strip():
            meta_lines.append(lines[i])
        i += 1
    while i < n:
        line = lines[i]
        if line.startswith('── '):
            label = re.sub(r'^──\s*', '', line)
            label = re.sub(r'\s*──\s*$', '', label)
            i += 1
            body = []
            while i < n and not is_header(lines[i]) and not lines[i].startswith('done.'):
                body.append(lines[i])
                i += 1
            sections.append({'kind': 'diff', 'label': label, 'body': body})
        elif line.startswith('== checklist'):
            i += 1
            while i < n and not lines[i].startswith('done.'):
                i += 1
        else:
            i += 1
    return meta_lines, sections


def parse_meta(meta_lines):
    meta = {}
    for ml in meta_lines:
        m = re.match(r'\s+(\S[^:]*?)\s*:\s*(.*)', ml)
        if m:
            meta[m.group(1).strip()] = m.group(2).strip()
    cm = re.search(r'/\S*jarvis-gen-diff/\S+', ' '.join(meta_lines))
    meta['cache'] = cm.group(0) if cm else ''
    pm = re.search(r'popcode=(\S+)', meta.get('resource', ''))
    meta['popcode'] = pm.group(1) if pm else ''
    meta['resource_code'] = re.sub(r'\s*\(popcode=.*\)', '', meta.get('resource', '')).strip()
    return meta


# ── diff extraction ───────────────────────────────────────────────────────────
META_RE = [re.compile(p) for p in (r'^diff --git', r'^index ', r'^--- ', r'^\+\+\+ ', r'^\\ No newline')]


def section_paths(body):
    """Return (baseline_path, actual_path) parsed from the --- a/ and +++ b/ lines."""
    base, actual = '', ''
    for l in body:
        if l.startswith('--- a/'):
            base = l[6:]
        elif l.startswith('+++ b/'):
            actual = l[6:]
    return base, actual


def section_tag(label):
    """Map a section label like 'doc (r): ...' to a Chinese tag."""
    key = label.split(':', 1)[0].strip()
    return {
        'resource': '资源代码',
        'datasource': '数据源代码',
        'doc (r)': '资源文档',
        'doc (d)': '数据源文档',
    }.get(key, key)


def section_file(label):
    """Short filename from the section label (after the colon)."""
    parts = label.split(':', 1)
    return parts[1].strip() if len(parts) > 1 else ''


def parse_hunks(body):
    """Split a diff body into (header_lines, hunks) where each hunk is
    (header_text, left_start, right_start, [content_lines])."""
    pre = []          # file-level meta lines before first @@
    hunks = []
    i = 0
    while i < len(body) and not body[i].startswith('@@'):
        pre.append(body[i])
        i += 1
    while i < len(body):
        m = re.match(r'^@@ -(\d+)(?:,\d+)? \+(\d+)(?:,\d+)? @@', body[i])
        if not m:
            i += 1
            continue
        hdr = body[i]
        l0, r0 = int(m.group(1)), int(m.group(2))
        i += 1
        content = []
        while i < len(body) and not body[i].startswith('@@'):
            content.append(body[i])
            i += 1
        hunks.append((hdr, l0, r0, content))
    return pre, hunks


# ── unified view rows ─────────────────────────────────────────────────────────
def unified_rows(hunks):
    rows = []
    for hdr, l0, r0, content in hunks:
        rows.append(('hunk', hdr))
        for line in content:
            if line.startswith('@@'):
                continue
            if line == '' or line.startswith('\\'):
                rows.append(('meta', '', esc(line)))
                continue
            pfx = line[0]
            if pfx == '+':
                rows.append(('add', '+', esc(line[1:])))
            elif pfx == '-':
                rows.append(('del', '-', esc(line[1:])))
            elif pfx == ' ':
                rows.append(('ctx', ' ', esc(line[1:])))
            else:
                rows.append(('meta', '', esc(line)))
    return rows


# ── side-by-side pairing ──────────────────────────────────────────────────────
def split_rows(hunks):
    """Pair unified hunk lines into (left_lineno, left_content, left_type,
    right_lineno, right_content, right_type) rows. Types: ctx/del/add/blank."""
    rows = []
    for hdr, l0, r0, content in hunks:
        rows.append(('hunk', hdr))
        ln_l, ln_r = l0, r0
        i = 0
        while i < len(content):
            line = content[i]
            if line == '' or line.startswith('\\') or line.startswith('@@'):
                i += 1
                continue
            pfx = line[0]
            if pfx == ' ':
                rows.append((ln_l, esc(line[1:]), 'ctx', ln_r, esc(line[1:]), 'ctx'))
                ln_l += 1
                ln_r += 1
                i += 1
            else:
                dels, adds = [], []
                while i < len(content) and content[i] and content[i][0] == '-':
                    dels.append((ln_l, esc(content[i][1:])))
                    ln_l += 1
                    i += 1
                while i < len(content) and content[i] and content[i][0] == '+':
                    adds.append((ln_r, esc(content[i][1:])))
                    ln_r += 1
                    i += 1
                for k in range(max(len(dels), len(adds))):
                    l = dels[k] if k < len(dels) else (None, '')
                    r = adds[k] if k < len(adds) else (None, '')
                    lt = 'del' if k < len(dels) else 'blank'
                    rt = 'add' if k < len(adds) else 'blank'
                    rows.append((l[0], l[1], lt, r[0], r[1], rt))
    return rows


def render_unified(rows):
    out = []
    for r in rows:
        if r[0] == 'hunk':
            out.append(f'<span class="row hunk">{esc(r[1])}</span>')
        elif r[0] == 'meta':
            out.append(f'<span class="row meta">{r[2]}</span>')
        else:
            cls, sig, txt = r
            out.append(f'<span class="row {cls}"><span class="sig">{sig}</span><span class="txt">{txt}</span></span>')
    return ''.join(out)


def render_split(rows):
    # Max code length across both sides → both code columns get the same width,
    # so the divider stays at a fixed x and rows align, without wrapping.
    maxlen = 0
    for r in rows:
        if r[0] == 'hunk':
            continue
        ll, lc, lt, rl, rc, rt = r
        if lc:
            maxlen = max(maxlen, len(lc))
        if rc:
            maxlen = max(maxlen, len(rc))
    col_w = max(maxlen + 4, 36)
    out = ['<table class="split-table"><colgroup>'
           f'<col class="ln"><col class="code" style="width:{col_w}ch">'
           f'<col class="ln"><col class="code" style="width:{col_w}ch">'
           '</colgroup>']
    for r in rows:
        if r[0] == 'hunk':
            out.append(f'<tr class="hunk-tr"><td colspan="4">{esc(r[1])}</td></tr>')
            continue
        ll, lc, lt, rl, rc, rt = r
        ll_s = str(ll) if ll is not None else ''
        rl_s = str(rl) if rl is not None else ''
        out.append(
            '<tr>'
            f'<td class="ln l {lt}">{ll_s}</td><td class="code l {lt}">{lc}</td>'
            f'<td class="ln r {rt}">{rl_s}</td><td class="code r {rt}">{rc}</td>'
            '</tr>'
        )
    out.append('</table>')
    return ''.join(out)


# ── checklist recompute (Chinese) ────────────────────────────────────────────
def read_lines(path):
    try:
        with open(path, encoding='utf-8', errors='replace') as f:
            return f.read().split('\n')
    except Exception:
        return None


def grep_count(lines, pat):
    if lines is None:
        return 0
    rx = re.compile(pat)
    return sum(1 for l in lines if rx.search(l))


def grep_samples(lines, pat, limit=10):
    if lines is None:
        return []
    rx = re.compile(pat)
    out = []
    for idx, l in enumerate(lines, 1):
        if rx.search(l):
            out.append(f'{idx}:\t{l.strip()}')
            if len(out) >= limit:
                break
    return out


def recompute_checklist(paths):
    """paths = {'resource':(base,actual), 'datasource':(base,actual),
    'doc_r':(base,actual), 'doc_d':(base,actual)}."""
    r_base, r_act = paths.get('resource', ('', ''))
    d_base, d_act = paths.get('datasource', ('', ''))
    doc_r_base, doc_r_act = paths.get('doc_r', ('', ''))
    doc_d_base, doc_d_act = paths.get('doc_d', ('', ''))
    rlines = read_lines(r_act)
    rbase = read_lines(r_base)
    dlines = read_lines(d_act)
    doc_r = read_lines(doc_r_act)
    doc_r_b = read_lines(doc_r_base)
    doc_d = read_lines(doc_d_act)
    doc_d_b = read_lines(doc_d_base)
    bugs = []
    miss = '(文件缺失)'

    # Bug 1
    b = {'num': '1', 'desc': BUG_DESC['1'], 'verdicts': [], 'samples': []}
    b['verdicts'].append(f'实际  ForceNew: true 计数 : {grep_count(rlines, r"ForceNew:\s*true")}')
    if rbase is not None:
        b['verdicts'].append(f'基线  ForceNew: true 计数: {grep_count(rbase, r"ForceNew:\s*true")}')
    b['samples'] = grep_samples(rlines, r'ForceNew:\s*true', 8)
    bugs.append(b)

    # Bug 2
    b = {'num': '2', 'desc': BUG_DESC['2'], 'verdicts': ['Read 中引用 objectRaw/nested 的 d.Set:'], 'samples': []}
    if rlines is not None:
        rx_set = re.compile(r'd\.Set\(')
        kw = re.compile(r'objectRaw|make\(\[\]|for .*range', re.IGNORECASE)
        samples = []
        for idx, l in enumerate(rlines, 1):
            if rx_set.search(l) and kw.search(l):
                samples.append(f'{idx}:\t{l.strip()}')
                if len(samples) >= 8:
                    break
        b['samples'] = samples
    else:
        b['verdicts'].append(miss)
    bugs.append(b)

    # Bug 3
    b = {'num': '3', 'desc': BUG_DESC['3'],
         'verdicts': ['数据源 Read mapping[] 赋值(每个 key 须存在于嵌套 Elem schema):'], 'samples': []}
    b['samples'] = grep_samples(dlines, r'mapping\["', 10)
    if not b['samples'] and dlines is None:
        b['verdicts'].append(miss)
    bugs.append(b)

    # Bug 4
    NOTE_PAT = r'NOTE:.*only evaluated|NOTE:.*immutable|Modifying it in isolation|Changing it after creation has no effect'
    b = {'num': '4', 'desc': BUG_DESC['4'], 'verdicts': [], 'samples': []}
    for name, act, base in (('资源文档', doc_r, doc_r_b), ('数据源文档', doc_d, doc_d_b)):
        base_n = doc_r_base if name == '资源文档' else doc_d_base
        bn = grep_count(base, NOTE_PAT) if base is not None else 0
        line = f'{name}:  NOTE 行  实际={grep_count(act, NOTE_PAT) if act is not None else 0}'
        if base is not None:
            line += f'  基线={bn}'
        b['verdicts'].append(line)
    bugs.append(b)

    # Bug 5
    PH_PAT = r'没有资源测试用例|No resource test case|please pass the resource test case'
    b = {'num': '5', 'desc': BUG_DESC['5'], 'verdicts': [], 'samples': []}
    for name, act, basep in (('资源文档', doc_r, doc_r_base), ('数据源文档', doc_d, doc_d_base)):
        c = grep_count(act, PH_PAT) if act is not None else 0
        tag = '已手改:真示例' if c == 0 else '仍为占位'
        b['verdicts'].append(f'{name}:  占位行 {c}  {tag}')
    bugs.append(b)

    # Bug 6 (observation, no grep)
    bugs.append({'num': '6', 'desc': BUG_DESC['6'], 'verdicts': [], 'samples': []})

    # Bug 7
    b = {'num': '7', 'desc': BUG_DESC['7'],
         'verdicts': ['Update 中 d.HasChange / update=true 出现:'], 'samples': []}
    b['samples'] = grep_samples(rlines, r'd\.HasChange\(|update\s*=\s*true', 8)
    if not b['samples'] and rlines is None:
        b['verdicts'].append(miss)
    bugs.append(b)
    return bugs


def render_checklist(paths):
    bugs = recompute_checklist(paths)
    out = ['<section class="checklist"><h2 class="sec-h">Generator-v4 产物缺陷清单</h2>',
           '<p class="sec-sub">七类已知的 generator-v4 生成后缺陷。每条给出对实际(手改后)文件的 grep 判定,'
           '并附样本代码行;逐条与上方 diff 交叉对照。</p>']
    for b in bugs:
        verdicts = ''.join(f'<div class="vd">{esc(v)}</div>' for v in b['verdicts'])
        samples = ''.join(esc(s) + '\n' for s in b['samples'])
        samp_html = f'<pre class="bug-samples">{samples}</pre>' if samples else ''
        out.append(
            f'<div class="bug"><div class="bug-head">'
            f'<span class="bug-num">Bug {esc(b["num"])}</span>{esc(b["desc"])}</div>'
            f'<div class="bug-verdicts">{verdicts}</div>{samp_html}</div>'
        )
    out.append('</section>')
    return ''.join(out)


# ── diff panel ────────────────────────────────────────────────────────────────
def render_panel(section):
    label = section['label']
    tag = section_tag(label)
    fname = section_file(label)
    base, actual = section_paths(section['body'])
    pre, hunks = parse_hunks(section['body'])
    no_diff = not hunks and not any(l.startswith('@@') for l in section['body'])
    # detect "(no actual file)" / "(no baseline)" bodies
    body_joined = '\n'.join(section['body'])
    if 'no actual file' in body_joined or 'no file matching' in body_joined:
        body_html = '<div class="empty">工作区中无对应的实际文件,跳过该区段 diff。</div>'
    elif 'no baseline' in body_joined:
        body_html = '<div class="empty">无基线(生成器未能产出),跳过该区段 diff。</div>'
    elif no_diff:
        body_html = '<div class="empty">该文件与生成基线一致,无手改差异。</div>'
    else:
        meta_html = ''
        u = render_unified(unified_rows(hunks))
        s = render_split(split_rows(hunks))
        body_html = (
            f'<div class="view-split">{s}</div>'
            f'<div class="view-unified" hidden><pre class="diff-pre unified">{u}</pre></div>'
        )
    return (
        f'<section class="diff-panel"><div class="diff-head">'
        f'<span class="tag">{esc(tag)}</span><span class="path">{esc(fname)}</span></div>'
        f'<div class="diff-body">{body_html}</div></section>'
    )


# ── CSS ───────────────────────────────────────────────────────────────────────
CSS = r"""
:root{
  --ground:#F6F4F0; --panel:#FFFFFF; --ink:#232529; --muted:#6E7681;
  --add-tx:#2E6B3B; --add-bg:#E9F2EB; --del-tx:#9C3A2B; --del-bg:#F6ECE8;
  --blank-bg:#EDEBE6; --accent:#2C6E6B; --accent-tx:#FFFFFF; --border:#E4E1DA;
  --hunk-bg:#EDF1EF; --hunk-tx:#4A5560; --sig:#B6B2A8; --tile-num:#232529;
}
@media (prefers-color-scheme: dark){
  :root:not([data-theme="light"]){
    --ground:#1B1B1E; --panel:#232327; --ink:#DCDBD6; --muted:#8C939D;
    --add-tx:#7DBE8E; --add-bg:rgba(125,190,142,.13); --del-tx:#E08575; --del-bg:rgba(224,133,117,.13);
    --blank-bg:rgba(255,255,255,.04); --accent:#4FA8A3; --accent-tx:#0E2A28; --border:#34343A;
    --hunk-bg:rgba(79,168,163,.13); --hunk-tx:#A0B0AD; --sig:#565660; --tile-num:#DCDBD6;
  }
}
:root[data-theme="dark"]{
  --ground:#1B1B1E; --panel:#232327; --ink:#DCDBD6; --muted:#8C939D;
  --add-tx:#7DBE8E; --add-bg:rgba(125,190,142,.13); --del-tx:#E08575; --del-bg:rgba(224,133,117,.13);
  --blank-bg:rgba(255,255,255,.04); --accent:#4FA8A3; --accent-tx:#0E2A28; --border:#34343A;
  --hunk-bg:rgba(79,168,163,.13); --hunk-tx:#A0B0AD; --sig:#565660; --tile-num:#DCDBD6;
}
body{background:var(--ground); color:var(--ink); font-family:'IBM Plex Sans','Noto Sans SC',system-ui,sans-serif; margin:0; line-height:1.55; -webkit-font-smoothing:antialiased}
.wrap{max-width:1160px; margin:0 auto; padding:36px 24px 56px}
header.page{border-bottom:1px solid var(--border); padding-bottom:22px; margin-bottom:6px}
.eyebrow{font-family:'IBM Plex Mono',monospace; font-size:11px; letter-spacing:.12em; text-transform:uppercase; color:var(--accent); margin:0 0 10px}
h1{font-family:'IBM Plex Sans',sans-serif; font-size:32px; font-weight:600; line-height:1.15; margin:0; text-wrap:balance; letter-spacing:-.01em}
h1 .pop{font-family:'IBM Plex Mono',monospace; font-size:17px; font-weight:500; color:var(--muted); margin-left:10px; vertical-align:middle}
.sub{font-size:14px; color:var(--muted); margin:10px 0 0; max-width:64ch; line-height:1.65}
.sub code{font-family:'IBM Plex Mono',monospace; color:var(--ink); font-size:13px}
.meta-grid{display:grid; grid-template-columns:132px 1fr; gap:5px 18px; margin-top:18px; font-family:'IBM Plex Mono',monospace; font-size:12px}
.mk{color:var(--muted); letter-spacing:.04em; font-size:11px; padding-top:1px}
.mv{color:var(--ink); word-break:break-all}
.bar{display:flex; align-items:center; gap:16px; flex-wrap:wrap; margin-top:18px}
.legend{display:flex; gap:18px; flex-wrap:wrap; font-family:'IBM Plex Mono',monospace; font-size:11px; color:var(--muted)}
.legend .sw{display:inline-block; width:11px; height:11px; border-radius:2px; margin-right:7px; vertical-align:-1px}
.legend .sw.add{background:var(--add-tx)} .legend .sw.del{background:var(--del-tx)} .legend .sw.hunk{background:var(--accent)}
.toggle{display:inline-flex; border:1px solid var(--border); border-radius:7px; overflow:hidden; margin-left:auto}
.toggle button{font-family:'IBM Plex Sans',sans-serif; font-size:12px; padding:5px 14px; background:var(--panel); color:var(--muted); border:none; cursor:pointer; border-right:1px solid var(--border)}
.toggle button:last-child{border-right:none}
.toggle button.active{background:var(--accent); color:var(--accent-tx); font-weight:500}
.toggle button:focus-visible{outline:2px solid var(--accent); outline-offset:-2px}
.stats{display:flex; gap:12px; flex-wrap:wrap; margin:24px 0 28px}
.tile{flex:1 1 128px; background:var(--panel); border:1px solid var(--border); border-radius:9px; padding:14px 16px}
.tile .num{font-family:'IBM Plex Mono',monospace; font-size:23px; font-weight:600; color:var(--tile-num); line-height:1}
.tile .num.add{color:var(--add-tx)} .tile .num.del{color:var(--del-tx)}
.tile .lab{font-size:11px; letter-spacing:.06em; color:var(--muted); margin-top:6px}
.sec-h{font-family:'IBM Plex Sans',sans-serif; font-size:18px; font-weight:600; margin:0 0 6px}
.sec-sub{font-size:13px; color:var(--muted); margin:0 0 18px; max-width:72ch; line-height:1.65}
.diff-panel{background:var(--panel); border:1px solid var(--border); border-radius:9px; overflow:hidden; margin-bottom:20px}
.diff-head{display:flex; align-items:center; gap:11px; padding:10px 16px; border-bottom:1px solid var(--border)}
.tag{font-family:'IBM Plex Mono',monospace; font-size:10.5px; letter-spacing:.06em; padding:3px 9px; border-radius:4px; background:var(--accent); color:var(--accent-tx); white-space:nowrap}
.path{font-family:'IBM Plex Mono',monospace; font-size:12px; color:var(--muted); overflow:hidden; text-overflow:ellipsis; white-space:nowrap; min-width:0}
.diff-body{overflow:auto; max-height:540px; background:var(--panel)}
.empty{padding:18px 20px; font-family:'IBM Plex Mono',monospace; font-size:12.5px; color:var(--muted)}
.diff-pre{margin:0; padding:10px 0; font-family:'IBM Plex Mono',monospace; font-size:12.5px; line-height:1.65; background:var(--panel)}
body.view-split .view-unified{display:none}
body.view-unified .view-split{display:none}
body.view-unified .view-unified{display:block}
/* unified */
.diff-pre.unified .row{display:block; white-space:pre; padding:0 18px}
.diff-pre.unified .sig{display:inline-block; width:1.6ch; color:var(--sig); user-select:none; vertical-align:top}
.diff-pre.unified .add{background:var(--add-bg)} .diff-pre.unified .add .txt{color:var(--add-tx)} .diff-pre.unified .add .sig{color:var(--add-tx)}
.diff-pre.unified .del{background:var(--del-bg)} .diff-pre.unified .del .txt{color:var(--del-tx)} .diff-pre.unified .del .sig{color:var(--del-tx)}
.diff-pre.unified .hunk{background:var(--hunk-bg); color:var(--hunk-tx); border-top:1px solid var(--border); border-bottom:1px solid var(--border)}
.diff-pre.unified .meta{color:var(--muted)} .diff-pre.unified .ctx .txt{color:var(--ink)}
/* split — fixed-layout table; both code columns get the same explicit width
   (the global max line length, set inline per-table) so the divider sits at a
   constant x and rows stay aligned. white-space:pre keeps each line on a single
   row (no wrapping); long content makes the table wider than the panel, which
   .diff-body scrolls horizontally as one unit — no clip, no wrap, equal sides. */
.split-table{border-collapse:collapse; table-layout:fixed; margin:0; font-family:'IBM Plex Mono',monospace; font-size:12.5px; line-height:1.65}
.split-table col.ln{width:6ch}
.split-table col.code{width:auto}
.split-table td{padding:0; vertical-align:top; border:none}
.split-table td.ln{white-space:nowrap; overflow:hidden; text-align:right; padding:0 8px 0 14px; color:var(--sig); user-select:none; font-variant-numeric:tabular-nums}
.split-table td.code{padding:1px 12px 1px 8px; white-space:pre}
.split-table td.code.l{border-right:1px solid var(--border); padding-right:14px}
.split-table td.ln.del, .split-table td.code.del{background:var(--del-bg); color:var(--del-tx)}
.split-table td.ln.add, .split-table td.code.add{background:var(--add-bg); color:var(--add-tx)}
.split-table td.ln.blank, .split-table td.code.blank{background:var(--blank-bg); color:var(--sig)}
.split-table td.ln.ctx{color:var(--sig)} .split-table td.code.ctx{color:var(--ink)}
.split-table tr.hunk-tr td{background:var(--hunk-bg); color:var(--hunk-tx); padding:3px 16px; font-size:11.5px; white-space:pre; border-top:1px solid var(--border); border-bottom:1px solid var(--border)}
.checklist{background:var(--panel); border:1px solid var(--border); border-radius:9px; padding:22px 26px; margin-bottom:20px}
.bug{padding:15px 0; border-bottom:1px solid var(--border)}
.bug:last-child{border-bottom:none; padding-bottom:0} .bug:first-child{padding-top:0}
.bug-head{font-size:14px; font-weight:500; color:var(--ink); line-height:1.55}
.bug-num{font-family:'IBM Plex Mono',monospace; color:var(--accent); margin-right:9px; font-weight:600}
.bug-verdicts{font-family:'IBM Plex Mono',monospace; font-size:12px; color:var(--muted); margin:7px 0 0; line-height:1.85}
.bug-samples{font-family:'IBM Plex Mono',monospace; font-size:12px; background:var(--ground); border:1px solid var(--border); border-radius:6px; padding:10px 14px; margin:9px 0 0; overflow:auto; white-space:pre; color:var(--ink); line-height:1.6}
.footer{margin-top:32px; padding-top:16px; border-top:1px solid var(--border); font-size:12.5px; color:var(--muted); line-height:1.7}
.footer code{font-family:'IBM Plex Mono',monospace; background:var(--panel); border:1px solid var(--border); border-radius:3px; padding:1px 5px; color:var(--ink); font-size:11.5px}
@media (prefers-reduced-motion: reduce){*{transition:none!important}}
"""


def main():
    if len(sys.argv) < 3:
        sys.stderr.write('usage: provider-gen-diff-render.py <input.txt> <output.html>\n')
        sys.exit(2)
    src, out = sys.argv[1], sys.argv[2]
    with open(src, encoding='utf-8', errors='replace') as f:
        text = f.read()
    meta_lines, sections = parse(text)
    meta = parse_meta(meta_lines)
    diff_sections = [s for s in sections if s['kind'] == 'diff']

    # paths for checklist
    paths = {}
    key_map = {'resource': 'resource', 'datasource': 'datasource', 'doc (r)': 'doc_r', 'doc (d)': 'doc_d'}
    for s in diff_sections:
        k = s['label'].split(':', 1)[0].strip()
        pk = key_map.get(k)
        if pk:
            paths[pk] = section_paths(s['body'])

    # counts
    added = removed = 0
    for s in diff_sections:
        for l in s['body']:
            if l.startswith('+') and not l.startswith('+++'):
                added += 1
            elif l.startswith('-') and not l.startswith('---'):
                removed += 1

    panels = ''.join(render_panel(s) for s in diff_sections)
    checklist_html = render_checklist(paths)

    rows_meta = [
        ('资源类型', meta.get('resource_code', '')),
        ('PopCode', meta.get('popcode', '')),
        ('工作区', meta.get('worktree', '')),
        ('snake_case', meta.get('snake', '')),
        ('生成器', meta.get('generator', '')),
        ('基线缓存', meta.get('cache', '')),
    ]
    meta_html = ''.join(f'<div class="mk">{esc(k)}</div><div class="mv">{esc(v)}</div>' for k, v in rows_meta)

    # Prefer the snake_case resource name (e.g. realtime_compute_variable) for
    # the tab title so multiple gen-diff pages are distinguishable; fall back to
    # the resource code when snake is unavailable (degraded mode).
    title_name = (meta.get('snake', '') or '').replace(' / ', '_') or (meta.get('resource_code', 'Provider') or 'Provider')

    PAGE = f"""<title>{esc(title_name)} 生成对比</title>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600&family=IBM+Plex+Sans:wght@400;500;600;700&family=Noto+Sans+SC:wght@400;500;600&display=swap">
<style>{CSS}</style>
<div class="wrap">
  <header class="page">
    <p class="eyebrow">provider-gen-diff · 生成器基线 vs 手写改动</p>
    <h1>{esc(meta.get('resource_code', 'RouteTargetGroup'))}<span class="pop">{esc(meta.get('popcode', ''))}</span></h1>
    <p class="sub"><code>{esc(meta.get('snake', '').replace(' / ', '_'))}</code> 资源在 generator-v4 纯生成基线之上的人工修改。左侧为生成器产物(基线),右侧为实际(手改后)文件;每条<span style="color:var(--add-tx)">新增</span>是人补的修复,每条<span style="color:var(--del-tx)">删除</span>是被替换/删除的生成产物。</p>
    <div class="meta-grid">{meta_html}</div>
    <div class="bar">
      <div class="legend">
        <span><span class="sw add"></span>手写新增</span>
        <span><span class="sw del"></span>生成基线(删除/替换)</span>
        <span><span class="sw hunk"></span>hunk 头</span>
      </div>
      <div class="toggle" role="group" aria-label="diff 视图切换">
        <button id="t-split" class="active" type="button">并排</button>
        <button id="t-unified" type="button">统一</button>
      </div>
    </div>
  </header>
  <div class="stats">
    <div class="tile"><div class="num">{len(diff_sections)}</div><div class="lab">文件区段</div></div>
    <div class="tile"><div class="num add">+{added}</div><div class="lab">新增行</div></div>
    <div class="tile"><div class="num del">−{removed}</div><div class="lab">删除行</div></div>
    <div class="tile"><div class="num">7</div><div class="lab">缺陷检查</div></div>
  </div>
  {panels}
  {checklist_html}
  <div class="footer">基线反映<em>当前</em> CloudSpec pre 环境,由 <code>provider-gen-diff.sh</code> 重新生成。读起来像漂移而非有意手改的 hunk 可能是 pre 漂移——逐条对照上方清单。</div>
</div>
<script>
(function(){{
  var b=document.body, s=document.getElementById('t-split'), u=document.getElementById('t-unified');
  function set(view){{
    b.classList.remove('view-split','view-unified'); b.classList.add('view-'+view);
    s.classList.toggle('active',view==='split'); u.classList.toggle('active',view==='unified');
    try{{localStorage.setItem('gen-diff-view',view);}}catch(e){{}}
  }}
  s.addEventListener('click',function(){{set('split');}});
  u.addEventListener('click',function(){{set('unified');}});
  var saved;try{{saved=localStorage.getItem('gen-diff-view');}}catch(e){{}}
  if(saved==='unified'){{set('unified');}}
}})();
</script>"""
    with open(out, 'w', encoding='utf-8') as f:
        f.write(PAGE)
    print(f'rendered {out} (+{added}/-{removed}, {len(diff_sections)} sections)')


if __name__ == '__main__':
    main()

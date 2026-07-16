# cap: wrap.sh 评论裸 URL 自动包成 markdown 链接（把"知识"变"强制"）

- **类型**：工具人体工学 / 动作点闸门（knowledge → enforcement）
- **关联**：工单 84297352 评论 #124901690（裸 URL 不可点，需补 #124922654 markdown 链接）；aone-triage SKILL §4「Aone 评论渲染 quirk」

## 背景

aone-triage SKILL §4 早写明：Aone 评论区**不 autolink 裸 URL**，唯一可点格式是 markdown `[text](url)`，且点名 `wrap.sh done` 正文同样要写 markdown 链接。规矩白纸黑字、开局加载，但当轮收尾 wrap.sh done 时仍贴了裸 URL（`**PR**: https://github.com/...`、`**内部研发单**：84322082`），评论区渲染成死文本，客户/同事点不动。

## 根因

规矩是**知识级**（skill 散文 + 反模式清单），不是**强制级**——wrap.sh 发评论那一刻没有任何东西拦。SKILL 自己注明 `aone-comment-format.sh`「只管列表项排版空行，不会帮你把裸 URL 转链接」，即工具明确不兜底，全靠 agent 临场记得。规矩与动作点分离，埋头贴内容时顺手就贴了裸 URL。

## 本次改动（动作点闸门：formatter 自动包链接）

`bootstrap/aone-comment-format.sh`（wrap.sh → `format_comment()` 的唯一评论正文格式化入口，覆盖 sync/done/done-no-status 全路径）：

- 加 `_wrap_bare_urls_line(line)`：把裸 `https?://...` URL 包成 `[url](url)` markdown 链接。
- **保护三处不动**：① 已有 markdown 链接 `[text](url)`（占位防重包）；② 行内代码 span `` `url` ``；③ 代码围栏 ``` ``` ``` 内（终末遍历按 ``` 行 toggle 跳过）。
- **URL 边界**：正则排除空白、`)` `]` `<>` 及 CJK 区段（U+3000-303F 符号标点 / U+4E00-9FFF 汉字 / U+FF00-FFEF 全角），让 URL 在中文标点/汉字处停下，避免把「`https://end.com。另一条`」整段当成一个 URL；再 rstrip 尾部 ASCII+CJK 标点。
- **kill switch**：`JARVIS_COMMENT_URL_AUTOLINK=0` 关闭（默认开）。

补 `test/aone_comment_format_test.sh` Test 4：裸 URL 包装 / 已有链接不重包 / 代码内不动 / CJK 处停下 / kill switch 共 12 断言。既有 3 组测试全过（29/29、wrap_check 24/24、a1id 70/70）。

## 残余覆盖盲区（已知，未在本次消除）

- **手工 `a1 comment create` 不走 wrap.sh** → 不经此闸门，仍靠 agent 自觉写 markdown 链接。但 SKILL 已反模式点名此路径，且主流 bookend 全走 wrap.sh，手工路径是边缘。
- **Aone 工单详情 description（非评论）**：经 `a1 update --body-file` 的 description 同样按 markdown 渲染，但 description 不走 aone-comment-format.sh。若 description 里贴裸 URL 也不可点——属 agent 自觉范畴（description 改动频率低，且多数是我方关联单，影响面小）。

## 置信度

high_conf —— 改动点是 wrap.sh→format_comment 的唯一入口（line 118/160），12 条断言覆盖正反场景，且 CJK 边界用例直接复现当轮 84297352 的裸 URL 现场。

## 落地

worktree `worktree-comment-url-gate` → 改 `bootstrap/aone-comment-format.sh` + 补 `test/aone_comment_format_test.sh` → MR 待仓库主人合并。部署 = 合并后 wrap.sh 自动生效（无需重启 bridge，formatter 每次评论即时调用）。

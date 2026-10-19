"use strict";

// Custom markdownlint rule: NOTE attachment in provider docs.
//
// Attribute-level `-> **NOTE:**` callouts are valid CommonMark in two broken
// shapes, so built-in rules cannot see them:
//   1. blank line before + flush-left: renders as a standalone paragraph,
//      detaching the NOTE from its attribute and splitting the argument list;
//   2. no blank line + indented: lazy continuation, the NOTE renders as
//      literal text merged into the attribute line.
//
// For a NOTE whose preceding block contains an attribute bullet, BOTH
// conditions are required (AND, not either):
//   1. a blank line immediately before the NOTE line;
//   2. alignment with the attribute the NOTE follows — indentation at the
//      content column of the NEAREST preceding attribute bullet (its indent
//      + 2), at every nesting level. Plain sub-bullets (value lists such as
//      `  - true`) elaborate the owning attribute and never own the NOTE,
//      so the walk skips over them to the attribute above.
// When no bullet continues the list after the NOTE (end of block), the
// block column (shallowest bullet indent + 2) is equally legal: the NOTE
// may close the block's content or attach to the nearest attribute.
// A NOTE following section-level prose, headings, or code fences is a
// resource-level NOTE and passes.
//
// Bound to the markdownlint-cli 0.47.0 custom-rule API; re-verify when
// upgrading the CLI (markdownlint-cli2 uses a different rule API).

module.exports = {
  names: ["alicloud-note-attachment", "ALICLOUD100"],
  description:
    "NOTE after an attribute bullet needs a blank line before AND alignment with the attribute it follows (nearest attribute bullet content column)",
  tags: ["alicloud"],
  function: function noteAttachment(params, onError) {
    const lines = params.lines;
    const noteRe = /^ *(->|~>) \*\*NOTE/;
    const listRe = /^(\s*)(?:[*+] |\d+\. |- )/;
    const fenceRe = /^\s{0,3}(```+|~~~+)/;
    const headingRe = /^#{1,6} /;
    // House attribute bullet: `* \`attr_name\` - ...`
    const attrBulletRe = /^\s*[*+] `[^`]+`/;
    const indentOf = (line) => line.match(/^ */)[0].length;
    const isBlank = (line) => line.trim() === "";

    // Fence parity: content inside code fences is opaque — a NOTE-looking
    // line in a fence is not an attachment site, and fenced bullets never
    // form an attribute block.
    const fenced = [];
    let inFence = false;
    for (let f = 0; f < lines.length; f++) {
      if (fenceRe.test(lines[f])) {
        inFence = !inFence;
      } else {
        fenced[f] = inFence;
      }
    }

    for (let i = 0; i < lines.length; i++) {
      if (!noteRe.test(lines[i]) || fenced[i]) {
        continue;
      }
      const indent = indentOf(lines[i]);
      const blankBefore = i === 0 || isBlank(lines[i - 1]);

      // Walk back from the previous non-blank line over the attribute block.
      // Fences, headings, and flush-left non-list content are block
      // boundaries: they stop the walk, but bullets already collected still
      // define the list the NOTE visually attaches to.
      let j = i - 1;
      while (j >= 0 && isBlank(lines[j])) {
        j--;
      }
      let shallowest = null;
      let nearestAttr = null;
      let k = j;
      while (k >= 0) {
        const line = lines[k];
        if (fenceRe.test(line) || headingRe.test(line) || fenced[k]) {
          break;
        }
        const m = line.match(listRe);
        if (m) {
          const bulletIndent = m[1].length;
          if (shallowest === null || bulletIndent < shallowest) {
            shallowest = bulletIndent;
          }
          if (nearestAttr === null && attrBulletRe.test(line)) {
            nearestAttr = bulletIndent;
          }
          k--;
          while (k >= 0 && isBlank(lines[k])) {
            k--;
          }
          continue;
        }
        if (indentOf(line) > 0) {
          // Indented non-list content (lazy continuations, sub-paragraphs,
          // preceding NOTEs) belongs to the block above it.
          k--;
          while (k >= 0 && isBlank(lines[k])) {
            k--;
          }
          continue;
        }
        break;
      }

      if (nearestAttr === null) {
        // Resource-level NOTE (no attribute bullet in the preceding block:
        // prose/heading/fence, or a plain list): flush-left is the legal
        // form; not this rule's concern.
        continue;
      }

      // Forward: does a bullet continue the list after the NOTE? A NOTE
      // between bullets interrupts the list and must align with the
      // attribute it follows; a NOTE closing the block may also sit at the
      // block column (shallowest bullet indent + 2).
      let c = i + 1;
      while (c < lines.length && isBlank(lines[c])) {
        c++;
      }
      let nextBullet = null;
      if (c < lines.length) {
        const m2 = lines[c].match(listRe);
        if (m2) {
          nextBullet = m2[1].length;
        }
      }
      const interrupting = nextBullet !== null && nextBullet >= nearestAttr;
      const legal = interrupting
        ? [nearestAttr + 2]
        : Array.from(new Set([shallowest + 2, nearestAttr + 2]));
      const legalText =
        legal.length === 1 ? String(legal[0]) : legal.join(" or ");
      if (!blankBefore) {
        onError({
          lineNumber: i + 1,
          detail:
            "no blank line before NOTE: CommonMark absorbs it as a lazy continuation of the attribute text; " +
            "expected a blank line before AND indent " +
            legalText +
            " (aligned with the attribute it follows)",
          context: lines[i].trim().slice(0, 40),
        });
      } else if (!legal.includes(indent)) {
        onError({
          lineNumber: i + 1,
          detail:
            "NOTE not aligned with the attribute it follows: indent " +
            indent +
            ", expected " +
            legalText +
            " (parent content column)",
          context: lines[i].trim().slice(0, 40),
        });
      }
    }
  },
};

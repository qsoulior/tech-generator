import {
  Decoration,
  MatchDecorator,
  ViewPlugin,
  type DecorationSet,
  type EditorView,
  type ViewUpdate,
} from "@codemirror/view"

// Marks balanced {{ ... }} placeholders so the user can visually distinguish
// template directives from regular markdown. We intentionally do not try to
// validate the contents of the placeholder — backend template parsing is the
// source of truth for syntactic correctness.
const placeholderDeco = Decoration.mark({ class: "cm-tpl-placeholder" })

// Catches obviously broken placeholders that look like attempts at field
// access but contain whitespace right after the leading dot, e.g. `{{ . foo }}`.
// This is the most common typo and is worth flagging in the editor.
const danglingDotRe = /\{\{[\s-]*\.\s+[^}]*\}\}/
const danglingDotDeco = Decoration.mark({ class: "cm-tpl-placeholder cm-tpl-placeholder--invalid" })

const matcher = new MatchDecorator({
  regexp: /\{\{[^{}]*\}\}/g,
  decoration: (match) => (danglingDotRe.test(match[0]) ? danglingDotDeco : placeholderDeco),
})

export const templatePlaceholderExtension = ViewPlugin.fromClass(
  class {
    decorations: DecorationSet
    constructor(view: EditorView) {
      this.decorations = matcher.createDeco(view)
    }
    update(update: ViewUpdate) {
      this.decorations = matcher.updateDeco(update, this.decorations)
    }
  },
  { decorations: (v) => v.decorations },
)

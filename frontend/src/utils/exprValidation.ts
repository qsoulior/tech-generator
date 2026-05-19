// Identifiers that exist in the expression evaluation environment but are not
// variables defined by the user. Covers:
//   - math constants injected by the backend (see variable_process_service/builtins.go)
//   - expr-lang literals
//   - expr-lang keywords and word-like operators (https://expr-lang.org/docs/language-definition)
// Function names are not listed here because the regex below already skips any
// identifier immediately followed by "(".
const KNOWN_NON_VARIABLE_IDENTIFIERS = new Set([
  "pi",
  "e",
  "g",
  "true",
  "false",
  "nil",
  "null",
  "and",
  "or",
  "not",
  "in",
  "if",
  "else",
  "let",
  "contains",
  "startsWith",
  "endsWith",
  "matches",
])

// Matches identifier tokens that are NOT a function call (no "(" after them)
// and NOT prefixed by "$" (so the built-in $env variable is ignored).
const VARIABLE_IDENT_RE = /(?<!\$)\b[A-Za-z_][A-Za-z0-9_]*\b(?!\s*\()/g

// Matches "let foo = ..." so locally-introduced names don't get flagged.
const LET_BINDING_RE = /\blet\s+([A-Za-z_][A-Za-z0-9_]*)\s*=/g

// Replaces every string/char/template literal with an equivalent-length empty
// pair so identifier-like words inside string values are not extracted.
function stripStringLiterals(expression: string): string {
  return expression.replace(/"(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'|`(?:[^`\\]|\\.)*`/g, '""')
}

// Returns the set of identifiers in `expression` that don't match any known
// variable slug, math constant, literal, or local "let" binding. Empty array
// means the expression references nothing unknown.
export function findUnknownIdentifiers(expression: string, knownSlugs: string[]): string[] {
  const source = stripStringLiterals(expression)

  const known = new Set(knownSlugs)
  for (const match of source.matchAll(LET_BINDING_RE)) {
    if (match[1] != null) known.add(match[1])
  }

  const unknown = new Set<string>()
  for (const match of source.matchAll(VARIABLE_IDENT_RE)) {
    const id = match[0]
    if (KNOWN_NON_VARIABLE_IDENTIFIERS.has(id)) continue
    if (known.has(id)) continue
    unknown.add(id)
  }
  return Array.from(unknown)
}

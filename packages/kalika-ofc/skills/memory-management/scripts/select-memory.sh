#!/usr/bin/env bash
set -euo pipefail

FILE="~/.copilot/memory/MEMORY.md"
QUERY=""
TAGS=""
LIMIT=5
MAX_CHARS=6000
MIN_SCORE=1

usage() {
  cat <<'EOF_USAGE'
Usage:
  select-memory.sh --file <path> --query "<keywords>" [options]

Options:
  --file <path>        Memory markdown file. Default: .ai/memory/memories.md
  --query "<text>"     Search query keywords.
  --tags "<tags>"      Optional comma-separated tags.
  --limit <n>          Max memory blocks to return. Default: 5
  --max-chars <n>      Max total output chars. Default: 6000
  --min-score <n>      Minimum relevance score. Default: 1
  --help              Show this help.

Example:
  bash ./scripts/select-memory.sh \
    --file ~/.copilot/memory/MEMORY.md \
    --query "golang mysql database/sql review" \
    --tags "go,mysql,backend,review" \
    --limit 5 \
    --max-chars 6000
EOF_USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
  --file)
    FILE="${2:-}"
    shift 2
    ;;
  --query)
    QUERY="${2:-}"
    shift 2
    ;;
  --tags)
    TAGS="${2:-}"
    shift 2
    ;;
  --limit)
    LIMIT="${2:-5}"
    shift 2
    ;;
  --max-chars)
    MAX_CHARS="${2:-6000}"
    shift 2
    ;;
  --min-score)
    MIN_SCORE="${2:-1}"
    shift 2
    ;;
  --help | -h)
    usage
    exit 0
    ;;
  *)
    echo "Unknown argument: $1" >&2
    usage >&2
    exit 1
    ;;
  esac
done

if [[ ! -f "$FILE" ]]; then
  echo "NO_MEMORY_FILE_FOUND"
  exit 0
fi

build_terms() {
  {
    printf '%s\n' "$QUERY"
    printf '%s\n' "$TAGS"
  } |
    tr '[:upper:]' '[:lower:]' |
    tr ',;|/()[]{}:"'"'"'`' ' ' |
    tr -cs '[:alnum:]_.+-' '\n' |
    awk 'length($0) > 1 && !seen[$0]++ { print }'
}

TERMS="$(build_terms | paste -sd '|' - || true)"

if [[ -z "$TERMS" ]]; then
  echo "NO_QUERY_TERMS_PROVIDED"
  exit 0
fi

awk \
  -v terms="$TERMS" \
  -v limit="$LIMIT" \
  -v max_chars="$MAX_CHARS" \
  -v min_score="$MIN_SCORE" '
function trim(s) {
  gsub(/^[ \t\r\n]+/, "", s)
  gsub(/[ \t\r\n]+$/, "", s)
  return s
}

function field_line(label, rec, lines, n, i, line) {
  n = split(rec, lines, "\n")
  for (i = 1; i <= n; i++) {
    line = lines[i]
    if (index(line, label) == 1) {
      return trim(substr(line, length(label) + 1))
    }
  }
  return ""
}

function process_record(rec, lower, context, tags, score, i, term) {
  if (trim(rec) == "") {
    return
  }

  if (rec !~ /^##[[:space:]]+/) {
    return
  }

  lower = tolower(rec)
  context = tolower(field_line("Contexto:", rec))
  tags = tolower(field_line("Tags:", rec))

  score = 0

  for (i = 1; i <= term_count; i++) {
    term = term_list[i]

    if (term == "") {
      continue
    }

    if (index(tags, term) > 0) {
      score += 10
    }

    if (index(context, term) > 0) {
      score += 5
    }

    if (index(lower, term) > 0) {
      score += 1
    }
  }

  if (score >= min_score) {
    count++
    scores[count] = score
    records[count] = rec
  }
}

function print_top_results(printed, used, k, i, best, best_score, out, remaining) {
  printed = 0
  used = 0

  for (k = 1; k <= limit; k++) {
    best = 0
    best_score = -1

    for (i = 1; i <= count; i++) {
      if (!selected[i] && scores[i] > best_score) {
        best = i
        best_score = scores[i]
      }
    }

    if (best == 0) {
      break
    }

    out = records[best]

    if (used + length(out) > max_chars) {
      remaining = max_chars - used

      if (remaining > 500) {
        out = substr(out, 1, remaining) "\n...[truncated]"
      } else {
        break
      }
    }

    printf "\n<!-- score:%d -->\n%s\n", scores[best], out

    used += length(out)
    selected[best] = 1
    printed++
  }

  if (printed == 0) {
    print "NO_RELEVANT_MEMORY_FOUND"
  }
}

BEGIN {
  term_count = split(terms, term_list, "|")
  record = ""
}

/^##[[:space:]]+/ {
  process_record(record)
  record = $0 "\n"
  next
}

{
  record = record $0 "\n"
}

END {
  process_record(record)
  print_top_results()
}
' "$FILE"

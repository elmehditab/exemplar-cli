# Exemplar CLI

Exemplar CLI is a domain-aware code review CLI.

The project is based on two public references:

- Roblox engineering article:
	https://about.roblox.com/newsroom/2026/01/doubled-ai-code-acceptance-teaching-models-think-like-roblox-engineers
- qmd repository:
	https://github.com/tobi/qmd

The idea is simple: do not review code from the diff alone.
Retrieve relevant context first (code, history, conventions, exemplars), then generate findings with evidence.

## Open and Run the Project

### Prerequisites

- Go 1.26+
- Git

### 1) Open the project

```bash
git clone https://github.com/mehditabet/exemplar-cli.git
cd exemplar-cli
```

Or open the folder directly in VS Code.

### 2) Run the CLI

```bash
go run . review --repo .
```

### 3) Run tests

```bash
go test ./...
```

## Retrieval Architecture (Short)

Exemplar uses a hybrid retrieval strategy inspired by qmd:

1. Lexical search (BM25/FTS) for exact matches (paths, symbols, identifiers).
2. Vector search for semantic similarity.
3. Fusion (RRF) to combine ranked results from multiple sources.
4. Reranking to improve the final top evidence.
5. AST-aware chunking so code is retrieved by semantic units, not random text slices.

This retrieval-first approach is aligned with Roblox's idea:
teach the system with repository history and expert patterns, not only with generic prompts.

## Status

The project is in MVP phase.

Current baseline:

- Go CLI with Cobra
- Git-based change collection
- Initial review pipeline

Next milestones:

- historical PR/comment ingestion
- exemplar engine
- multi-pass reviewers with confidence + citations

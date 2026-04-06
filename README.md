# Exemplar CLI

Exemplar CLI is a local-first tool for context-aware code review.

The goal is to help developers review changes with more context than a raw diff by combining:

- the current git change
- the local repository structure
- review history and team conventions later
- structured review output in the terminal

## Current Status

Exemplar is still in an early MVP stage.

Today, the project includes:

- a Go CLI built with Cobra
- Git-based change collection
- a first review pipeline
- a reusable `ReviewContext` foundation

The next steps are focused on:

- building richer local repository context
- adding historical PR and review retrieval
- introducing selective multi-agent review

## Inspiration

Exemplar is inspired by public engineering work on domain-aware AI code review, especially Roblox's January 12, 2026 article on exemplar alignment, review-history retrieval, and feedback-driven code intelligence.

Exemplar is an independent project and is not affiliated with Roblox.

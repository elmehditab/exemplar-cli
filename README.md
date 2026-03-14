# Exemplar CLI

Exemplar CLI is a Go-based command-line tool designed to bring AI-powered code review directly into the developer workflow.

Its goal is to help engineers review code changes earlier, faster, and with more context than traditional static checks alone. Rather than acting only as a linter or formatter, Exemplar is intended to become an intelligent review assistant capable of understanding code changes, repository structure, and surrounding context before producing useful feedback.

## Project Objective

The objective of Exemplar is to build a review engine that can:

- analyze repository changes from git diffs
- understand code context beyond changed lines
- progressively incorporate repository structure and dependencies
- use caching and contextual processing to scale efficiently
- provide actionable code review feedback from the command line

The long-term vision is to make `exemplar` a developer-facing quality gate that fits naturally into local development, automation pipelines, and AI-assisted coding workflows.

## Vision

Exemplar is being designed as a foundation for future capabilities such as:

- diff-based review workflows
- repository indexing
- contextual code understanding
- code graph exploration
- AI-assisted review reasoning
- cache-backed performance optimizations

The project starts with a clean and production-grade CLI foundation so that these capabilities can be added incrementally without reworking the application architecture.

## Design Principles

Exemplar is intended to follow a few core principles:

- clear separation between CLI, application, and core review logic
- idiomatic Go project structure and maintainability
- composable architecture for future review pipeline stages
- pragmatic integration with local developer tooling
- extensibility for advanced analysis and AI features

## Status

Exemplar CLI is currently in its foundational stage.

The current focus is on establishing the command-line architecture and internal application structure that will support the future review engine.

## Repository

This repository contains the source code for the `exemplar-cli` project.

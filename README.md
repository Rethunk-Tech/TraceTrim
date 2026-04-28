# TraceTrim

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/rethunk-tech/tracetrim)](https://goreportcard.com/report/github.com/rethunk-tech/tracetrim)
[![Release](https://img.shields.io/github/v/release/rethunk-tech/tracetrim.svg)](https://github.com/rethunk-tech/tracetrim/releases)
[![Downloads](https://img.shields.io/github/downloads/rethunk-tech/tracetrim/total.svg)](https://github.com/rethunk-tech/tracetrim/releases)

A cross-platform CLI tool that monitors your clipboard for JavaScript/React stack traces and automatically removes repetitive frames, making error logs clean and readable.

## Features

- **Automatic Detection** — Continuously monitors clipboard for stack traces
- **Smart Cleaning** — Removes only repetitive blocks, preserves all formatting
- **Real-time** — Updates clipboard instantly when traces are detected
- **Script Mode** — Works in shell pipelines and automation
- **Cross-platform** — Windows, macOS, Linux via `golang.design/x/clipboard`
- **Zero-config** — Just run it

## Documentation

| Role | Guide |
|------|-------|
| **Users** | [HUMANS.md](./HUMANS.md) — Installation, usage, configuration, troubleshooting |
| **Developers** | [AGENTS.md](./AGENTS.md) — Architecture, building, testing, extending |
| **Security** | [SECURITY.md](./SECURITY.md) — Vulnerability reporting, threat model, practices |
| **Contributors** | [CONTRIBUTING.md](./CONTRIBUTING.md) — Development workflow, PR guidelines |
| **Releases** | [CHANGELOG.md](./CHANGELOG.md) — Version history and release notes |
| **License** | [LICENSE](./LICENSE) — MIT license |

---

Copyright (c) 2025 Rethunk.Tech, LLC.

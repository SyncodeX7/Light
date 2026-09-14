# Light ⚡

A lightning-fast, minimalist AI companion living right inside your terminal. Built in Go, Light is designed to be your peer-to-peer coding partner, shell assistant, and low-profile digital friend—complete with token streaming, persistent chat memory, and zero corporate bloat.

## Features

- **Blazing Fast Go Binary:** Single-file, compiled binary with zero heavy runtime dependencies.
- **True Token Streaming:** Responses stream into your terminal letter-by-letter for an instantaneous feel.
- **Persistent Local Memory:** Automatically maintains conversation history locally in `~/.config/light/history.json` so context carries across terminal sessions.
- **Hugging Face Powered:** Backed by state-of-the-art open models (defaulting to Llama-3.3-70B-Instruct) via the Hugging Face Serverless API.
- **Dual Modes:** Fire off quick single-shot questions or drop into a continuous interactive REPL chat loop.

## Installation

### From Source
Make sure you have Go installed, then clone and build:

```bash
git clone [https://github.com/yourusername/light.git](https://github.com/yourusername/light.git)
cd light
go build -ldflags="-s -w" -o light main.go
sudo mv light /usr/local/bin/

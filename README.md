# Light ⚡

<p align="center">
  <b>A lightning-fast, minimalist AI companion living right inside your terminal.</b>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#usage">Usage</a> •
  <a href="#license">License</a>
</p>

---

Light is a peer-to-peer coding partner, shell assistant, and low-profile digital friend built completely from scratch in Go. Designed with zero corporate bloat, it features real-time token streaming, persistent local history, and a razor-sharp developer vibe that stays out of your way until you need it.

## Features

- **Blazing Fast Go Binary:** Compiled into a single, lightweight binary with zero heavy runtime dependencies or node module bloat.
- **True Token Streaming:** Responses stream smoothly into your terminal letter-by-letter for an instantaneous, immersive feel.
- **Persistent Local Memory:** Automatically manages conversation context locally in `~/.config/light/history.json` so your sessions carry seamlessly across terminal restarts.
- **Hugging Face Powered:** Driven by state-of-the-art open models (defaulting to `meta-llama/Llama-3.3-70B-Instruct`) through the high-performance Hugging Face Serverless API.
- **Flexible Modes:** Fire off quick single-shot terminal queries or drop straight into a continuous interactive REPL chat session.

## Installation

### From Source
Ensure you have Go installed on your system, then clone, build, and move the binary to your path:

```bash
git clone [https://github.com/SyncodeX7/light.git](https://github.com/SyncodeX7/light.git)
cd light
go build -ldflags="-s -w" -o light main.go
sudo mv light /usr/local/bin/

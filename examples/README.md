# Pulumi fal Provider Examples

This directory contains examples showing how to use the Pulumi fal provider to deploy applications on fal.ai.

## Available Examples

### Sana App Demo

Demonstrates how to deploy the Sana image generation app from the fal-ai-community/fal-demos repository.

Currently available:
- **[Python](./python/)** - Complete Python example with uv environment setup

Coming soon:
- **TypeScript** - TypeScript example with npm/yarn
- **Go** - Go example with modules

Each example creates a fal application that runs the Sana image generation model and exports the app's properties for further use.

## Getting Started

1. Choose your preferred language from the examples above
2. Follow the README instructions in the specific language directory
3. Make sure you have your fal API credentials configured
4. Deploy using `pulumi up`

## Common Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/) installed
- A valid [fal.ai](https://fal.ai) API key
- Language-specific requirements (see individual example READMEs)
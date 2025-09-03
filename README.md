# Pulumi Provider for fal

A Pulumi provider for managing [fal.ai](https://fal.ai) serverless AI applications.

This provider is 🚧 under development.

> **Note**: This is a third-party provider developed independently and is not officially supported by fal.ai

## Installation

Build from source:

```bash
make build
```

### Provider Configuration
```sh
  pulumi config set falKey "your-fal-api-key" --secret
```

## Resources

### `fal:App`

Manages a fal.ai application deployment.

#### Example Usage

```typescript
const app = new fal.App("my-app", {
    name: "test-app",
    entrypoint: "app.py:app",
    authMode: "private",
    strategy: "recreate",
    git: {
        url: "https://github.com/fal-ai-community/fal-demos.git",
    }
});

export const revisionId = app.revisionId;
export const createdAt = app.createdAt;
```

#### Arguments

- `name` (string, required): The name of the application
- `entrypoint` (string, required): The entrypoint for the application (e.g., "app.py:app")
- `strategy` (string, optional): Deployment strategy ("recreate" or "rolling"). Default: "recreate"
- `authMode` (string, optional): Authentication mode ("public", "private", or "shared"). Default: "private"
- `git` (object, optional): Git repository configuration
  - `url` (string, required): Git repository URL
  - `username` (string, optional): Username for HTTP authentication
  - `password` (string, optional): Password for HTTP authentication
  - `privateKey` (string, optional): Private key for SSH authentication
  - `insecureHttpAllowed` (boolean, optional): Allow insecure HTTP connections

#### Outputs

- `revisionId` (string): The revision ID of the deployed application
- `createdAt` (string): When the application was created (ISO 8601 format)
- `updatedAt` (string): When the application was last updated (ISO 8601 format)

## Requirements

- Go 1.24+
- Pulumi CLI
- fal CLI (must be available in PATH)
- Git (for repository operations)

## Architecture

The provider is built using the [Pulumi Go Provider framework](https://github.com/pulumi/pulumi-go-provider) and follows these patterns:

- **Client Layer**: Wraps fal CLI operations in a Go client (`pkg/fal/client.go`)
- **Resource Layer**: Implements Pulumi CRUD operations (`pkg/provider/app.go`)
- **Provider Layer**: Configures authentication and provider metadata (`pkg/provider/config.go`)

## Development

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Build: `go build .`

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

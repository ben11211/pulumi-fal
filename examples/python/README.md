# fal.ai Sana Application Example

This example demonstrates how to deploy a fal.ai application using Pulumi Python, equivalent to the following Terraform configuration:

```hcl
terraform {
  required_providers {
    fal = {
      source = "fal-ai/fal"
      version = "~> 1"
    }
  }
}

provider "fal" {
  fal_key = "<fal key here>"
}

resource "fal_app" "sana_app" {
  entrypoint = "fal_demos/image/sana.py"
  git = {
    url = "https://github.com/fal-ai-community/fal-demos.git"
  }
}
```

## Prerequisites

- Pulumi CLI installed
- Python 3.8+ installed
- uv installed ([installation guide](https://docs.astral.sh/uv/getting-started/installation/))
- fal CLI installed, available in PATH, and authenticated.

## Setup

1. **Clone this repository and navigate to this example:**
   ```bash
   cd examples/python
   ```

2. **Initialize a uv virtual environment:**
   ```bash
   uv venv
   source .venv/bin/activate
   uv pip install pip  # Install pip in the venv for Pulumi compatibility
   ```

3. **Add the fal.ai provider package and generate Python SDK:**
   ```bash
   pulumi package add ../..
   ```

4. **Install the local provider plugin:**
   ```bash
   # First, make sure the provider binary is built
   cd ../../ && make build && cd examples/python
   
   # Clear any pre-existing installations
   pulumi plugin rm resource fal
   
   # Install the plugin locally (match the version the SDK expects)
   pulumi plugin install resource fal 0.0.0 --file ../../pulumi-provider-fal
   ```

5. **Install the generated SDK and dependencies:**
   ```bash
   uv pip install -r requirements.txt
   ```

6. **Set up your fal API key** (choose one method):
   **Option B: Pulumi Configuration**
   ```bash
   pulumi config set fal:falKey "your-fal-api-key-here" --secret
   ```

## Expected Output

After deployment, you should see outputs similar to:

```
Outputs:
    app_name     : "sana-app"
    created_at   : "2024-01-15T10:30:00Z"
    revision_id  : "abc123def456"
    updated_at   : "2024-01-15T10:30:00Z"
```

## What This Example Does

This example:

1. **Deploys a fal.ai application** called "sana-app"
2. **Uses the Sana image generation model** from the fal-demos repository
3. **Clones the repository** `https://github.com/fal-ai-community/fal-demos.git`
4. **Sets the entrypoint** to `fal_demos/image/sana.py`
5. **Exports key information** about the deployed app

The deployed application will be available on fal.ai and can be invoked using the fal CLI or API.

## Customization

You can modify the `__main__.py` file to:

- Change the application name
- Use a different git repository
- Set different entrypoints
- Add authentication modes (`auth_mode`)
- Specify deployment strategies (`strategy`)
- Add git authentication for private repositories

Example with additional options:

```python
sana_app = fal.App("sana_app",
    name="my-custom-sana-app",
    entrypoint="fal_demos/image/sana.py",
    auth_mode="private",  # or "public", "shared"
    strategy="rolling",   # or "recreate"
    git=fal.AppGitArgs(
        url="https://github.com/fal-ai-community/fal-demos.git",
        username="your-username",      # for private repos with HTTP auth
        password="your-password",      # for private repos with HTTP auth
        # private_key="your-ssh-key",  # for private repos with SSH auth
    )
)
```
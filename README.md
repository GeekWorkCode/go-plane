# go-plane

[简体中文](./README.zh-cn.md)

[Plane](https://plane.so/) integration with [GitHub](https://docs.github.com/en/actions) or [Gitea Action](https://docs.gitea.com/usage/actions/overview) for project management.

## Motivation

This project aims to integrate Plane's project management features with GitHub Actions or Gitea Actions, allowing developers to automatically update issues, add comments, change states, and assign users during the CI/CD process based on commit messages.

## Installation

### From Source Code

```bash
# Clone the repository
git clone https://github.com/GeekWorkCode/go-plane
cd go-plane

# Build the binary
make build

# Install to $GOPATH/bin
make install
```

### Using Pre-compiled Binaries

Download the latest version for your platform from [GitHub Releases](https://github.com/GeekWorkCode/go-plane/releases).

```bash
# Linux (amd64)
curl -L https://github.com/GeekWorkCode/go-plane/releases/download/vX.Y.Z/go-plane-vX.Y.Z-linux-amd64 -o go-plane
chmod +x go-plane
sudo mv go-plane /usr/local/bin/

# macOS (amd64)
curl -L https://github.com/GeekWorkCode/go-plane/releases/download/vX.Y.Z/go-plane-vX.Y.Z-darwin-amd64 -o go-plane
chmod +x go-plane
sudo mv go-plane /usr/local/bin/
```

### Docker

```bash
# Using GitHub Container Registry
docker pull ghcr.io/GeekWorkCode/go-plane:latest

# Run
docker run --rm \
  -e PLANE_BASE_URL="https://plane.example.com/api/v1" \
  -e PLANE_TOKEN="your-token" \
  -e PLANE_WORKSPACE_SLUG="your-workspace" \
  -e PLANE_REF="PROJ-123 Fix issue" \
  ghcr.io/GeekWorkCode/go-plane
```

## Usage

### Configuration

You can configure `go-plane` using environment variables:

| Environment Variable | Description                                              | Example                          |
| -------------------- | -------------------------------------------------------- | -------------------------------- |
| PLANE_BASE_URL       | The base URL of your Plane instance API                  | https://plane.example.com/api/v1 |
| PLANE_TOKEN          | Your Plane API token                                     | plane_api_123456789              |
| PLANE_WORKSPACE_SLUG | Your Plane workspace slug                                | my-workspace                     |
| PLANE_REF            | The git commit message or reference to analyze           | "Fix issue PROJ-123: Update API" |
| PLANE_TO_STATE       | The state to transition the issues to (optional)         | "Done"                           |
| PLANE_COMMENT        | Comment to add to the issues (optional)                  | "Fixed in commit: $GITHUB_SHA"   |
| PLANE_ASSIGNEE       | Username to assign the issues to (optional)              | "john.doe"                       |
| PLANE_MARKDOWN       | Set to "true" to format comments as Markdown (optional)  | true                             |
| PLANE_INSECURE       | Set to "true" to skip SSL verification (not recommended) | false                            |
| PLANE_DEBUG          | Set to "true" to enable debug output                     | false                            |

### GitHub Actions Example

```yaml
name: Plane Integration

on:
  push:
    branches: [ main ]

jobs:
  plane-integration:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Update Plane Issues
        uses: docker://ghcr.io/GeekWorkCode/go-plane:latest
        env:
          PLANE_BASE_URL: ${{ secrets.PLANE_BASE_URL }}
          PLANE_TOKEN: ${{ secrets.PLANE_TOKEN }}
          PLANE_WORKSPACE_SLUG: ${{ secrets.PLANE_WORKSPACE_SLUG }}
          PLANE_REF: ${{ github.event.head_commit.message }}
          PLANE_TO_STATE: "Done"
          PLANE_COMMENT: "Resolved in commit: ${{ github.sha }}\n\nAuthor: ${{ github.actor }}"
          PLANE_MARKDOWN: "true"
```

### Gitea Actions Example

```yaml
name: Plane Integration

on:
  push:
    branches: [ main ]

jobs:
  plane-integration:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Update Plane Issues
        uses: docker://ghcr.io/your-username/go-plane:latest
        env:
          PLANE_BASE_URL: ${{ secrets.PLANE_BASE_URL }}
          PLANE_TOKEN: ${{ secrets.PLANE_TOKEN }}
          PLANE_WORKSPACE_SLUG: ${{ secrets.PLANE_WORKSPACE_SLUG }}
          PLANE_REF: ${{ github.event.head_commit.message }}
          PLANE_TO_STATE: "Done"
          PLANE_COMMENT: "Resolved in commit: ${{ github.sha }}\n\nAuthor: ${{ github.actor }}"
          PLANE_MARKDOWN: "true"
```

## FAQ

### Issue format

The tool extracts issue keys in the format "PROJ-123" where "PROJ" is the project identifier and "123" is the sequence ID. The issue key can be anywhere in your commit message, for example:

```
Fix bug in login form PROJ-123
```

or

```
PROJ-123 Implement new feature
```

### API connection problems

If you're having SSL certificate issues, you can use `PLANE_INSECURE=true` to skip verification (not recommended for production).

Make sure your Plane API URL is correctly formatted, usually `https://your-plane-instance.com/api/v1`.

### Can't find issues to update

- Make sure your workspace slug (`PLANE_WORKSPACE_SLUG`) is correct
- Check if your commit message contains valid issue references
- Verify that your API token has permission to access these issues

## License

MIT

## Acknowledgements

This project was inspired by and draws from:
- [go-jira](https://github.com/appleboy/go-jira) for structure and approach 
- [plane-api-go](https://github.com/GeekWorkCode/plane-api-go) for the Plane API client 
# Release Process

This project uses `make` and Shell Scripts for automated release builds.

## Versioning Strategy

Follows [Semantic Versioning 2.0.0](https://semver.org/).

- **Major**: Incompatible API changes.
- **Minor**: Backward-compatible feature additions.
- **Patch**: Backward-compatible bug fixes.

## Build Process

### 1. Local Testing

Ensure full testing is performed before release:

```bash
make lint
make test
```

### 2. Execute Release Build

Use the `scripts/release_build.sh` script for cross-platform compilation:

```bash
./scripts/release_build.sh v1.0.0
```

This script will:

1. Check if the git state is clean.
2. Run tests.
3. Compile versions for macOS (arm64/amd64) and Linux (amd64/arm64).
4. Output binary files to `bin/release/`.
5. Generate SHA256 checksums.

### 3. Publish Artifacts

Generated file structure:

```text
bin/release/
├── vlmcp-v1.0.0-darwin-arm64
├── vlmcp-v1.0.0-linux-amd64
└── checksums.txt
```

Upload these files to the GitHub Releases page.

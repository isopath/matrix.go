## Installation

### Via brew
```bash
# Add the tap
brew tap isopath/tap

# Install matrix
brew install matrix

# Or install directly
brew install isopath/tap/matrix
```

### Manual Installation

Download the latest binary from the [releases page](https://github.com/isopath/matrix.go/releases) and place it in your `$PATH`.

## Usage

```bash
# Default: Shows colorful matrix with "How to Do Great Work" by Paul Graham
matrix

# Show available options and select from list
matrix --options

# Use a specific file
matrix --file hackers-manifesto.txt
```

### Controls

- `q` or `Ctrl+C`: Quit the program
- In `--options` mode: Use arrow keys to navigate, Enter to select

## Building from Source

```bash
# Clone the repository
git clone https://github.com/isopath/matrix.go.git
cd matrix.go

# Build the binary
go build -o matrix .

# Run
./matrix
```

## Development

### Release Process

Releases are automated via GitHub Actions. To create a new release:

1. Go to GitHub → Releases → "Draft a new release"
2. Create a new tag (e.g., `v1.0.0`)
3. Click "Publish release"

This will:
1. Build binaries for all platforms (Linux, macOS, Windows)
2. Create a GitHub release with all binaries attached
3. Create a PR to update the Homebrew formula in `isopath/homebrew-tap`

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

Download the latest binary from the [releases page](https://github.com/sauravmaheshkar/matrix.go/releases) and place it in your `$PATH`.

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

### Release Process

Releases are automated via GitHub Actions. To create a new release:

```bash
# Tag a new version
git tag v0.4.2
git push origin v0.4.2
```

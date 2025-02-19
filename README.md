# imposter-go [![CI](https://github.com/imposter-project/imposter-go/actions/workflows/ci.yml/badge.svg)](https://github.com/imposter-project/imposter-go/actions/workflows/ci.yml)

A Go implementation of the [Imposter Mock Engine](https://www.imposter.sh).

## Features

- Local: Lightweight local HTTP mock server
- AWS Lambda: low latency, high throughput, ~15ms cold start
- [REST/HTTP API mock](https://docs.imposter.sh/rest_plugin/) support
- [SOAP/WSDL mock](https://docs.imposter.sh/soap_plugin/) support

## ⚠️ Limitations

- No support for [scripting](https://docs.imposter.sh/scripting/)
- No support (yet) for [OpenAPI](https://docs.imposter.sh/openapi_plugin/) mocks

## Requirements

- Go 1.23 or later (earlier versions may work but are not tested)
- Make (for building)

## Installation

### Using Pre-built Binaries

Download the latest release for your platform from GitHub:

#### macOS

```bash
# For Intel Macs (x86_64)
curl -L https://github.com/imposter-project/imposter-go/releases/latest/download/imposter-go_Darwin_x86_64.tar.gz | tar xz
sudo mv imposter-go /usr/local/bin/

# For Apple Silicon Macs (arm64)
curl -L https://github.com/imposter-project/imposter-go/releases/latest/download/imposter-go_Darwin_arm64.tar.gz | tar xz
sudo mv imposter-go /usr/local/bin/
```

#### Linux

```bash
# For x86_64 systems
curl -L https://github.com/imposter-project/imposter-go/releases/latest/download/imposter-go_Linux_x86_64.tar.gz | tar xz
sudo mv imposter-go /usr/local/bin/

# For arm64 systems
curl -L https://github.com/imposter-project/imposter-go/releases/latest/download/imposter-go_Linux_arm64.tar.gz | tar xz
sudo mv imposter-go /usr/local/bin/
```

#### Windows

1. Download the latest release from [GitHub Releases](https://github.com/imposter-project/imposter-go/releases/latest)
2. Extract the `imposter-go_Windows_x86_64.zip` file
3. Add the extracted `imposter-go.exe` to your PATH or move it to a directory in your PATH

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/imposter-project/imposter-go.git
   cd imposter-go
   ```

2. Run the server with an example configuration:
   ```bash
   make run ./examples/rest/simple
   ```

3. Visit `http://localhost:8080/hello` in your browser or use `curl`:
   ```bash
   curl http://localhost:8080/hello
   ```

4. Check the server version:
   ```bash
   curl http://localhost:8080/system/status
   ```

## Configuration

Set the `IMPOSTER_PORT` environment variable to change the default port:
```bash
export IMPOSTER_PORT=9090
```

Set the `IMPOSTER_LOG_LEVEL` environment variable to control logging verbosity:
```bash
export IMPOSTER_LOG_LEVEL=DEBUG  # Available levels: TRACE, DEBUG, INFO, WARN, ERROR
```

The default log level is DEBUG. Log levels are processed in order of severity:
- TRACE - Most verbose, logs all messages
- DEBUG - Detailed information for debugging
- INFO - General operational messages
- WARN - Warning messages for potentially harmful situations
- ERROR - Error messages for serious problems

## Development

The project uses Make for building and development. The following targets are available:

- `make run <path>` - Run the server directly (useful during development)
- `make build` - Build the project with version information
- `make install` - Install the binary to your Go bin directory
- `make test` - Run tests with verbose output

For development, use `make run` which will compile and run the server in one step:
```bash
# Run with a specific configuration
make run ./examples/rest/simple

# Run with multiple arguments
make run --debug ./examples/soap/simple
```

For production or installation, use `make build` or `make install`.

The version information is automatically derived from git tags. When building from source:
- Released versions will show the git tag (e.g. "v1.0.0")
- Development builds will show the git commit hash
- If no git information is available, it will show "dev"

## Examples

The repository includes several examples demonstrating different features:

- [Simple REST API](examples/rest/simple) - Basic REST API mocking
- [SOAP Web Service](examples/soap/simple) - SOAP/WSDL-based service mocking
- [AWS Lambda](examples/lambda) - Running Imposter in AWS Lambda

## Testing

Run the test suite using:
```bash
make test
```

This will run all tests with verbose output, showing the progress of each test case.

## Releasing

To create a new release:

1. Tag your commit with a semver tag prefixed with 'v':
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The GitHub Actions workflow will automatically:
   - Build binaries for all supported platforms (Linux, macOS, Windows)
   - Generate a changelog from commit messages
   - Create a GitHub release
   - Upload the built artifacts

The release notes will include:
- Version and release date
- Automatically generated changelog (excluding docs, test, ci, and chore commits)
- Installation instructions
- Link to the full changelog comparing with the previous tag

# stargazer

[![Build Status](https://github.com/samlader/stargazer/actions/workflows/ci.yml/badge.svg)](https://github.com/samlader/stargazer/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Subscribe to an RSS feed to follow what other people are starring, without stalking their GitHub profiles. ⭐

You can use a hosted version of this (for free) [here](https://stargazer.lader.io/feeds/samlader+healeycodes).

## Usage

```
GET /feed/{username} # Get RSS feed for a single user
GET /feeds/{user1+user2+...} # Get combined RSS feed for multiple users
```

### Quick Start:

```bash
make deps
make run
```

### Configuration

Set your GitHub token:

```bash
export GITHUB_TOKEN=your_token_here
```

### Development

```
make test        # Run tests
make lint        # Run linter
make build       # Build binary
make fmt         # Format code
```

## Contributions

Contributions and bug reports are welcome! Feel free to open issues, submit pull requests or contact me if you need any support.

## License

This project is licensed under the [MIT License](LICENSE).

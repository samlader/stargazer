# stargazer

[![Build Status](https://github.com/samlader/stargazer/actions/workflows/ci.yml/badge.svg)](https://github.com/samlader/stargazer/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Subscribe to an RSS feed to follow what other people are starring, without stalking their GitHub profiles. ⭐

You can use a hosted version of this (for free) [here](https://stargazer.lader.io/feeds/samlader+healeycodes).

## Usage

```
GET /feed/{username} # Single user
GET /feeds/{user1}+{user2}+...} # Multiple users
```

### Quick Start

Set your GitHub token:

```bash
export GITHUB_TOKEN=your_token_here
```

Install dependencies and run:
```bash
make deps
make run
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

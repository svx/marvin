# marvin
Documentation QA Overview

## Checks

- [Vale](vale.sh)
- [markdownlint](https://github.com/DavidAnson/markdownlint)

## Project Structure

### CLI

### Web

### Documentation

The documentation is build with [VitePress](https://vitepress.dev/) using [Bun](https://bun.sh/).

## Development and Contributing

### Requirements

- [Devbox](https://www.jetify.com/docs/devbox)

Open a terminal with:

```shell
devbox shell
```

### Working on Documentation

Make sure that you have [Devbox](https://www.jetify.com/docs/devbox) installed and do the following:

- Open a terminal in the root of the repository and run `devbox shell`
- Change into the documentation directory and run `bun install`
- Run `bun run docs:dev` to start a local development build of the documentation

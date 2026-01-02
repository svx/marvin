---
title: CLI
description:
  Comprehensive guide for contributing to the CLI, including setup,
  development and testing
outline: deep
---

# CLI

## Directory Structure

### `/cmd`

Main applications for this project.

### `/internal`

You can optionally add a bit of extra structure to your internal packages to separate your shared and non-shared internal code.
Your actual application code can go in the `/internal/app` directory (e.g., `/internal/app/marvin`) and the code shared by those apps in the `/internal/pkg` directory (e.g., `/internal/pkg/myprivlib`).

### `/test`

Additional external test apps and test data.  Note that Go will also ignore directories or files that begin with "." or "_", so you have more flexibility in terms of how you name your test data directory.

### `/tools`

Supporting tools for this project. Note that these tools can import code from the `/internal` directory.

## Tests

## Writing tests

Consider whether new tests are required. These tests
should ensure that the functionality you are adding will continue to work in the
future. Existing tests may also need updating.

You may also consider adding unit tests for any new functions you have added.
The unit tests should follow the Go convention of being location in a file named
`*_test.go` in the same package as the code being tested.

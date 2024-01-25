# How to contribute

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

Check out the [Diamcircle Contribution Guide](https://github.com/diamcircle/.github/blob/master/CONTRIBUTING.md) that apply to all Diamcircle projects.

## Style guides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")

### Issues

* Issues and PR titles start with:
  * The package name most affected, ex. `ingest: fix...`.
  * Or, for services and tools, the service most affected, ex. `services/aurora: fix...`, `services/ticker: add...`
  * Or, multiple package names separated by a comma when the fix addresses multiple packages worth noting, ex. `services/aurora, services/friendbot: fix...`.
  * Or, `all:` when changes or an issue are broad, ex. `all: update...`.
  * Or, `doc:` when changes or an issue are isolated to non-code documentation not limited to a single package.
* Label issues with `bug` if they're clearly a bug.
* Label issues with `feature request` if they're a feature request.

### Pull Requests

* PR titles follow the same rules as described in the [Issues](#Issues) section above.
* PRs must update the [CHANGELOG](CHANGELOG.md) with a small description of the change
* PRs are merged into master or release branch using squash merge
* Carefully think about where your PR fits according to [semver](https://semver.org). Target it at master if it’s only a patch change, otherwise if it contains breaking change or significant feature additions, set the base branch to the next major or minor release.
* Keep PR scope narrow. Expectation: 20 minutes to review max
* Explicitly differentiate refactoring PRs and feature PRs. Refactoring PRs don’t change functionality. They usually touch a lot more code, and are reviewed in less detail. Avoid refactoring in feature PRs.

### Go Style Guide

* Use `gofmt` or preferably `goimports` to format code
* Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Go Coding conventions

- Always document exported package elements: vars, consts, funcs, types, etc.
- Tests are better than no tests.

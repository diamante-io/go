<!-- If you're making a doc PR or something tiny where the below is irrelevant, delete this
template and use a short description, but in your description aim to include both what the
change is, and why it is being made, with enough context for anyone to understand. -->

<details>
  <summary>PR Checklist</summary>
  
### PR Structure

* [ ] This PR has reasonably narrow scope (if not, break it down into smaller PRs).
* [ ] This PR avoids mixing refactoring changes with feature changes (split into two PRs
  otherwise).
* [ ] This PR's title starts with name of package that is most changed in the PR, ex.
  `services/friendbot`, or `all` or `doc` if the changes are broad or impact many
  packages.

### Thoroughness

* [ ] This PR adds tests for the most critical parts of the new functionality or fixes.
* [ ] I've updated any docs ([developer docs](https://developers.diamcircle.org/api/), `.md`
  files, etc... affected by this change). Take a look in the `docs` folder for a given service,
  like [this one](https://go/tree/master/services/aurora/internal/docs).

### Release planning

* [ ] I've updated the relevant CHANGELOG ([here](services/aurora/CHANGELOG.md) for aurora) if
  needed with deprecations, added features, breaking changes, and DB schema changes.
* [ ] I've decided if this PR requires a new major/minor version according to
  [semver](https://semver.org/), or if it's mainly a patch change. The PR is targeted at the next
  release branch if it's not a patch change.
</details>

### What

[TODO: Short statement about what is changing.]

### Why

[TODO: Why this change is being made. Include any context required to understand the why.]

### Known limitations

[TODO or N/A]

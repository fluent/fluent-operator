# Releases

This page describes the release process and the currently planned schedule for upcoming releases as well as the respective release shepherd.

## Release schedule

| release series | date  (year-month-day) | release shepherd                            |
|----------------|--------------------------------------------|---------------------------------------------|
| v0.1.0           | 2020-02-17                                 | Guangzhe Huang (GitHub: @huanggze) |
| v0.2.0           | 2020-08-27                                 | Guangzhe Huang (GitHub: @huanggze)         |
| v0.3.0           | 2020-11-10                                 | Guangzhe Huang (GitHub: @huanggze)     |
| v0.4.0           | 2021-04-01                                 | Wanjun Lei (GitHub: @wanjunlei) |
| v0.5.0           | 2021-04-14                                 | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.6.0           | 2021-06-03                                 | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.6.1           | 2021-06-11                                 | Benjamin Huo (GitHub: @benjaminhuo)         |
| v0.6.2           | 2021-06-11                                 | Benjamin Huo (GitHub: @benjaminhuo)     |
| v0.7.0           | 2021-06-29                                 | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.7.1           | 2021-07-09                                 | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.8.0           | 2021-07-23                                 | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.9.0           | 2021-08-13                                 | Wanjun Lei (GitHub: @wanjunlei)         |

## How to cut an individual release

These instructions are currently valid for the [fluentbit operator repository](https://github.com/kubesphere/fluentbit-operator.git).
### Branch management and versioning strategy

We use [Semantic Versioning](https://semver.org/).

We maintain a separate branch for each minor release, named `release-<major>.<minor>`, e.g. `release-1.1`, `release-2.0`.

Note that branch protection kicks in automatically for any branches whose name starts with `release-`. Never use names starting with `release-` for branches that are not release branches.

The usual flow is to merge new features and changes into the main branch and to merge bug fixes into the latest release branch. Bug fixes are then merged into main from the latest release branch. The main branch should always contain all commits from the latest release branch. As long as main hasn't deviated from the release branch, new commits can also go to main, followed by merging main back into the release branch.

If a bug fix got accidentally merged into main after non-bug-fix changes in main, the bug-fix commits have to be cherry-picked into the release branch, which then have to be merged back into main. Try to avoid that situation.

Maintaining the release branches for older minor releases happens on a best effort basis.

### 1. Prepare your release

At the start of a new major or minor release cycle create the corresponding release branch based on the main branch. 
For example if we're releasing `0.6.0` and the previous stable release is `0.5.0` we need to create a `release-0.6` branch. 
Note that all releases are handled in protected release branches, see the above `Branch management and versioning` section. 
Release candidates and patch releases for any given major or minor release happen in the same `release-<major>.<minor>` branch. 
Do not create `release-<version>` for patch or release candidate releases.

Note that it need to update the image version of `fluentbit-operator` to the `release-<version>`.

Changes for a patch release or release candidate should be merged into the previously mentioned release branch via pull request.

Bump the version in the `VERSION` file. Do this in a proper PR pointing to the release branch as this gives others the opportunity to chime in on the release in general and on the addition to the changelog in particular. 
For a release candidate, append something like `-rc.0` to the version (with the corresponding changes to the tag name, the release name etc.).

For release candidates still update `CHANGELOG.md`, but when you cut the final release later, merge all the changes from the pre-releases into the one final update.

Entries in the `CHANGELOG.md` are meant to be in this order:

* `[CHANGE]`
* `[FEATURE]`
* `[ENHANCEMENT]`
* `[BUGFIX]`

### 2. Draft the new release

Tag the new release via the following commands:

```bash
$ tag="v$(< VERSION)"
$ git tag -s "${tag}" -m "${tag}"
$ git push origin "${tag}"
```

Optionally, you can use this handy `.gitconfig` alias.

```ini
[alias]
  tag-release = "!f() { tag=v${1:-$(cat VERSION)} ; git tag -s ${tag} -m ${tag} && git push origin ${tag}; }; f"
```

Then release with `git tag-release`.

### 3. Build and push image

```bash
make build
```
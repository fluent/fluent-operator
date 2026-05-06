# Releases

This page describes the release process and the currently planned schedule for upcoming releases as well as the respective release shepherd.

## Release schedule

| release series | date (year-month-day) | release shepherd                        |
| -------------- | --------------------- | --------------------------------------- |
| v0.1.0         | 2020-02-17            | Guangzhe Huang (GitHub: @huanggze)      |
| v0.2.0         | 2020-08-27            | Guangzhe Huang (GitHub: @huanggze)      |
| v0.3.0         | 2020-11-10            | Guangzhe Huang (GitHub: @huanggze)      |
| v0.4.0         | 2021-04-01            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.5.0         | 2021-04-14            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.6.0         | 2021-06-03            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.6.1         | 2021-06-11            | Benjamin Huo (GitHub: @benjaminhuo)     |
| v0.6.2         | 2021-06-11            | Benjamin Huo (GitHub: @benjaminhuo)     |
| v0.7.0         | 2021-06-29            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.7.1         | 2021-07-09            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.8.0         | 2021-07-23            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.9.0         | 2021-08-13            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.10.0        | 2021-08-20            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.11.0        | 2021-09-01            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.12.0        | 2021-09-13            | Wanjun Lei (GitHub: @wanjunlei)         |
| v0.13.0        | 2022-03-14            | Elon Cheng (GitHub: @wenchajun)         |
| v1.0.0         | 2022-03-25            | Elon Cheng (GitHub: @wenchajun)         |
| v1.0.1         | 2022-05-12            | Elon Cheng (GitHub: @wenchajun)         |
| v1.0.2         | 2022-05-17            | Elon Cheng (GitHub: @wenchajun)         |
| v1.1.0         | 2022-06-15            | Elon Cheng (GitHub: @wenchajun)         |
| v1.5.0         | 2022-09-24            | Elon Cheng (GitHub: @wenchajun)         |
| v1.5.1         | 2022-09-30            | Elon Cheng (GitHub: @wenchajun)         |
| v1.6.0         | 2022-10-25            | Elon Cheng (GitHub: @wenchajun)         |
| v1.6.1         | 2022-10-31            | Elon Cheng (GitHub: @wenchajun)         |
| v1.7.0         | 2022-11-23            | Elon Cheng (GitHub: @wenchajun)         |
| v2.0.0         | 2023-02-03            | Elon Cheng (GitHub: @wenchajun)         |
| v2.0.1         | 2023-02-08            | Elon Cheng (GitHub: @wenchajun)         |
| v2.1.0         | 2023-03-13            | Elon Cheng (GitHub: @wenchajun)         |
| v2.2.0         | 2023-04-07            | Elon Cheng (GitHub: @wenchajun)         |
| v2.3.0         | 2023-06-05            | Elon Cheng (GitHub: @wenchajun)         |
| v2.4.0         | 2023-07-19            | Elon Cheng (GitHub: @wenchajun)         |
| v2.5.0         | 2023-09-13            | Elon Cheng (GitHub: @wenchajun)         |
| v2.6.0         | 2023-11-22            | Elon Cheng (GitHub: @wenchajun)         |
| v2.7.0         | 2023-12-19            | Anthony Treuillier (GitHub: @antrema)   |
| v2.8.0         | 2024-04-22            | Zhang Peng (GitHub: @Gentleelephant)    |
| v2.9.0         | 2024-06-13            | Elon Cheng (GitHub: @wenchajun)         |
| v3.0.0         | 2024-07-09            | Elon Cheng (GitHub: @wenchajun)         |
| v3.1.0         | 2024-08-14            | Zhang Peng (GitHub: @Gentleelephant)    |
| v3.2.0         | 2024-09-21            | Chengwei Guo (GitHub: @cw-Guo)          |
| v3.3.0         | 2025-02-27            | Chengwei Guo (GitHub: @cw-Guo)          |
| v3.4.0         | 2025-05-08            | Marco Franssen (GitHub: @marcofranssen) |
| v3.5.0         | 2025-10-24            | Josh Baird (GitHub: @joshuabaird)       |
| v3.6.0         | 2025-10-24            | Josh Baird (GitHub: @joshuabaird)       |
| v3.7.0         | 2026-02-27            | Josh Baird (GitHub: @joshuabaird)       |

## Versioning strategy

We use [Semantic Versioning](http://semver.org/).

Version bumps are driven automatically by [Conventional Commits](https://www.conventionalcommits.org/) in PR titles:

| Commit prefix | Effect |
|---|---|
| `feat:` | Minor version bump |
| `fix:` | Patch version bump |
| `feat!:` or `BREAKING CHANGE:` | Major version bump |
| `chore:`, `docs:`, `ci:`, etc. | No version bump (excluded from changelog) |

## How to cut a new release

Releases are fully automated via [release-please](https://github.com/googleapis/release-please-action). The only human action required is **merging the Release PR**.

### Minor and major releases

1. As PRs merge to `master`, release-please continuously updates a **Release PR** that accumulates the changelog and bumps `version.txt` and all `Chart.yaml` `appVersion` fields.
2. When you are ready to ship, review the Release PR — edit the changelog body if needed — and merge it.
3. release-please automatically pushes the git tag and creates the GitHub Release.
4. The [`build-op-image`](./github/workflows/build-op-image.yaml) workflow fires on the tag and builds and pushes multi-arch operator images to GHCR and Docker Hub.
5. The [`upload-release-assets`](.github/workflows/upload-release-assets.yaml) workflow fires on the release publication and attaches `setup.yaml` (with the correct image tag stamped in) as a release asset.

To force a specific version instead of the auto-calculated one, dispatch the [Release Please](https://github.com/fluent/fluent-operator/actions/workflows/release-please.yaml) workflow manually with the `release-as` input set to the desired version (e.g. `3.8.0`).

### Patch releases

For a patch release against an older minor line (e.g. `v3.7.1`):

1. Create a `release-3.7` branch from the `v3.7.0` tag if it does not already exist.
2. Cherry-pick the relevant fix commits onto that branch.
3. Dispatch the [Release Please](https://github.com/fluent/fluent-operator/actions/workflows/release-please.yaml) workflow targeting the `release-3.7` branch with `release-as: 3.7.1`.
4. Review and merge the resulting Release PR on the `release-3.7` branch.
5. Cherry-pick the fix commits back into `master` if they are not already there.

### Publish updated Helm chart

The `appVersion` fields in the three chart `Chart.yaml` files in this repo are bumped automatically by the Release PR. After merging the Release PR, the chart still needs to be manually synced to [fluent/helm-charts](https://github.com/fluent/helm-charts/tree/main/charts/fluent-operator/) where users install from:

- Copy the updated `charts/fluent-operator/` directory to a branch in `fluent/helm-charts`
- Bump the chart `version` field (separate from `appVersion`) per the helm-charts versioning convention
- Open a PR in `fluent/helm-charts`

Do the same for the `fluent-operator-fluentd-crds` and `fluent-operator-fluent-bit-crds` charts.

## Fluentd/Fluent-bit Images

Fluent Operator uses [custom builds](./cmd/fluent-watcher/README.md) of both Fluentd and Fluent Bit. These images can (and often should) be published out-of-band of Fluent Operator releases.

### Fluent Bit

1. Execute the [bump-fluent-bit-version](https://github.com/fluent/fluent-operator/actions/workflows/bump-fluent-bit-version.yaml) workflow dispatch to generate a PR updating all Fluent Bit version references in this repo.
2. Merge the PR. This automatically triggers the [build-fluentbit-image](https://github.com/fluent/fluent-operator/actions/workflows/build-fluentbit-image.yaml) workflow which builds and pushes the new image.

### Fluentd

1. Open a PR updating all references to the image tag in this repo, including `cmd/fluent-watcher/fluentd/VERSION`.
2. Merge the PR.
3. Execute the [publish-fluentd-image](https://github.com/fluent/fluent-operator/actions/workflows/publish-fluentd-image.yaml) workflow dispatch to build and push the new image.

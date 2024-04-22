# Releases

This page describes the release process and the currently planned schedule for upcoming releases as well as the respective release shepherd.

## Release schedule

| release series | date  (year-month-day) | release shepherd                      |
|----------------|------------------------|---------------------------------------|
| v0.1.0         | 2020-02-17             | Guangzhe Huang (GitHub: @huanggze)    |
| v0.2.0         | 2020-08-27             | Guangzhe Huang (GitHub: @huanggze)    |
| v0.3.0         | 2020-11-10             | Guangzhe Huang (GitHub: @huanggze)    |
| v0.4.0         | 2021-04-01             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.5.0         | 2021-04-14             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.6.0         | 2021-06-03             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.6.1         | 2021-06-11             | Benjamin Huo (GitHub: @benjaminhuo)   |
| v0.6.2         | 2021-06-11             | Benjamin Huo (GitHub: @benjaminhuo)   |
| v0.7.0         | 2021-06-29             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.7.1         | 2021-07-09             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.8.0         | 2021-07-23             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.9.0         | 2021-08-13             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.10.0        | 2021-08-20             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.11.0        | 2021-09-01             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.12.0        | 2021-09-13             | Wanjun Lei (GitHub: @wanjunlei)       |
| v0.13.0        | 2022-03-14             | Elon Cheng (GitHub: @wenchajun)       |
| v1.0.0         | 2022-03-25             | Elon Cheng (GitHub: @wenchajun)       |
| v1.0.1         | 2022-05-12             | Elon Cheng (GitHub: @wenchajun)       |
| v1.0.2         | 2022-05-17             | Elon Cheng (GitHub: @wenchajun)       |
| v1.1.0         | 2022-06-15             | Elon Cheng (GitHub: @wenchajun)       |
| v1.5.0         | 2022-09-24             | Elon Cheng (GitHub: @wenchajun)       |
| v1.5.1         | 2022-09-30             | Elon Cheng (GitHub: @wenchajun)       |
| v1.6.0         | 2022-10-25             | Elon Cheng (GitHub: @wenchajun)       |
| v1.6.1         | 2022-10-31             | Elon Cheng (GitHub: @wenchajun)       |
| v1.7.0         | 2022-11-23             | Elon Cheng (GitHub: @wenchajun)       |
| v2.0.0         | 2023-02-03             | Elon Cheng (GitHub: @wenchajun)       |
| v2.0.1         | 2023-02-08             | Elon Cheng (GitHub: @wenchajun)       |
| v2.1.0         | 2023-03-13             | Elon Cheng (GitHub: @wenchajun)       |
| v2.2.0         | 2023-04-07             | Elon Cheng (GitHub: @wenchajun)       |
| v2.3.0         | 2023-06-05             | Elon Cheng (GitHub: @wenchajun)       |
| v2.4.0         | 2023-07-19             | Elon Cheng (GitHub: @wenchajun)       |
| v2.5.0         | 2023-09-13             | Elon Cheng (GitHub: @wenchajun)       |
| v2.6.0         | 2023-11-22             | Elon Cheng (GitHub: @wenchajun)       |
| v2.7.0         | 2023-12-19             | Anthony Treuillier (GitHub: @antrema) |
| v2.8.0         | 2024-04-22             | Zhang Peng (GitHub: @Gentleelephant)  |

# How to cut a new release

> This guide is strongly based on the [Prometheus release instructions](https://github.com/prometheus/prometheus/blob/master/RELEASE.md).

## Branch management and versioning strategy

We use [Semantic Versioning](http://semver.org/).

We maintain a separate branch for each minor release, named `release-<major>.<minor>`, e.g. `release-1.1`, `release-2.0`.

The usual flow is to merge new features and changes into the master branch and to merge bug fixes into the latest release branch. Bug fixes are then merged into master from the latest release branch. The master branch should always contain all commits from the latest release branch.

If a bug fix got accidentally merged into master, cherry-pick commits have to be created in the latest release branch, which then have to be merged back into master. Try to avoid that situation.

Maintaining the release branches for older minor releases happens on a best effort basis.

## Prepare your release

For a new major or minor release, work from the `main` branch. For a patch release, work in the branch of the minor release you want to patch (e.g. `release-0.1` if you're releasing `v0.1.1`).

Add an entry for the new version to the `CHANGELOG.md` file. Entries in the `CHANGELOG.md` should be in this order:

* `[CHANGE]`
* `[FEATURE]`
* `[ENHANCEMENT]`
* `[BUGFIX]`

Create a PR for the changes to be reviewed.

## Publish the new release

For new minor and major releases, create the `release-<major>.<minor>` branch starting at the PR merge commit.
From now on, all work happens on the `release-<major>.<minor>` branch.

Bump the version in the `VERSION` file in the root of the repository.

Regenerate setup.yaml based on latest code and then commit the changed bundle.yaml to the `release-<major>.<minor>` branch:
```bash
make manifests
git add ./
git commit -s -m "regenerate setup.yaml"
git push
```

Images will be automatically built and pushed whenever code changes or a tag is created. If users want to build images manually, use the following command:

```bash
make build-op-amd64 -e FO_IMG=<image of fluent operator>
docker push <image of fluent operator>
make build-fb-amd64 -e FB_IMG=<image of fluent bit>
docker push <image of fluent bit>
make build-fd-amd64 -e FD_IMG=<image of fluentd>
docker push <image of fluentd>
```

Tag the new release with a tag named `v<major>.<minor>.<patch>`, e.g. `v2.1.3`. Note the `v` prefix. You can do the tagging on the commandline:

```bash
tag="$(< VERSION)"
git tag -a "${tag}" -m "${tag}"
git push origin "${tag}"
```
Commit all the changes.

Finally, create a new release:

- Go to https://github.com/fluent/fluent-operator/releases/new.
- Associate the new release with the previously pushed tag.
- Add release notes based on `CHANGELOG.md`.
- Add file `setup.yaml` from `manifests/setup/setup.yaml` and then click "Publish release".

For patch releases, cherry-pick the commits from the release branch into the master branch.


How do we test a pre-release safely?

* We need a way to build main (eg, latest)
* What is building "latest"?
* "latest" gets pushed on merge to main -- test this
* We shold push SHA-based tags on pull requests 

What should happen when a PR is opened?

* Image built and pushed with tag=SHA 

What should happen when a PR is merged to main?

* Image build and pushed with tag=SHA 

What should happen when a tag is pushed?

* Image build and pushed with tag=v<major>.<minor>.<patch>

Releasing v3.5.0 of fluent-operator:

* Cut "release-3.5" brancg
* Update VERSION 
* Update README.md 
* Update RELEASE.md
* Update CHANGELOG.md (generated changelog from draft release) 
* Update helm chart appVersion/chartVersions 
* Run "make manifests" (may have to update image tag in Deployment manually)
* Merge PR
* Generate and push tag to upstream

```
export tag="v3.5.0" 
git tag -a $tag -m $tag 
git push upstream $tag 
```
* Draft new release (copy CHANGELOG) 
* Attach setup.yaml 
* Convert release from draft to published
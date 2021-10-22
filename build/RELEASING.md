# Relasing a new `cas` version

This document provides all steps required by a maintainer to release a new `cas` version.

We assume all commands are entered from the root of the `cas` working directory.
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html), in this document we will use vX.Y.V as placeholder for the version we are going to relase.

Before relasing, ensure that all modifications have been tested and pushed. When the release introduce a user-facing change, please ensure the documentation is updated accordingly. Last but not least, remind to run `make docs/cmd` to regenerate the Markdown CLI documentation when the CLI interface or the help messages have been modified.

### About the branching model

Although `cas` aims to have a "[OneFlow](https://www.endoflineblog.com/oneflow-a-git-branching-model-and-workflow)" git branching model, [release branches](https://www.endoflineblog.com/oneflow-a-git-branching-model-and-workflow#release-branches) have never been used.
Thus, the instructions on the current document assume that just the `master` branch is used for the release process (with all modifications previously merged in). However the whole process can be easily adapter to a release branch if needed.

## 1. Bump version (vX.Y.Z)

Edit `Makefile` and modify the `VERSION` variable:

```
VERSION=X.Y.Z
```
> N.B. Omit the `v` prefix.

Then run:

```
make CHANGELOG.md.next-tag
```

## 2. Commit and tag (locally)

Add the files modified above to the git index:

```
git add Makefile
git add CHANGELOG.md
```

Then:
```
git commit -m "release: vX.Y.Z"
git tag vX.Y.Z
```
> Do not push now.

Finally, sign the commit using `cas` itself:
```
cas n -p git://.
```

## 3. Make dist files

In order to sign the Windows binaries with a digital certificate, you will need an `.spc` and `.pvk` files (and the password to unlook the `.pvk`).
Make sure the path of those files is accessible.

Build all dist files:

```
SIGNCODE_PVK_PASSWORD=<pvk password> SIGNCODE_PVK=<path to vchian.pvk> SIGNCODE_SPC=<path to vchain.spc> make dist/all
```
> Distribution files will be created into the `dist` directory.


Check that everthing worked as expected and finally sign all files using `cas` itself:
```
make dist/sign
```

## 4. Push and edit the release on github

Push your commits and tag:
```
git push
git push --tags
```
> From now on, your relase will be publicy visible, and [dockerhub](https://hub.docker.com/repository/docker/codenotary/cas/builds) should start building docker images for `cas`.

Now you can edit the vX.Y.Z newly created [release on GitHub](https://github.com/codenotary/cas/releases), using the follwing template

```
# Changelog

<!-- copy and past here the latest part from CHANGELOG.md -->

# Downloads

**Docker image**
https://hub.docker.com/r/codenotary/cas

**Binaries**

File | SHA256
------------- | -------------
<!-- use `make dist/binary.md` to generate the downloads links and checksums -->
```

Finally, uploads all files from `dist`.

## 5. Sign docker images

Once [dockerhub](https://hub.docker.com/repository/docker/codenotary/cas/builds) has finished to build the images, then pull, check and finally sign them.

Pull:
```
docker pull codenotary/cas:vX.Y.Z
docker pull codenotary/cas:vX.Y.Z-docker
```

Check:
```
docker run --rm -it codenotary/cas:vX.Y.Z info
docker run --rm -it codenotary/cas:vX.Y.Z-docker info
```

Sign:
```
cas n -p docker://codenotary/cas:X.Y.Z
cas n -p docker://codenotary/cas:X.Y.Z-docker
```

## 6. Miscellaneous

### Homebrew

Once a new `cas` version has been released, the [cas Homebrew formula](https://github.com/vchain-us/homebrew-brew/blob/master/Formula/cas.rb) must be updated accordingly.

### CodeNotary platform

Once a new `cas` version has been released, the latest version should be set on the CodeNotary platform backend.
Then you can check the latest version has been set properly [here](https://api.codenotary.io/foundation/v1/version/cas/latest).

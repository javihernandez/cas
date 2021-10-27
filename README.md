# Community Attestation Service (CAS)  <img align="right" src="docs/img/cn-color.eeadbabe.svg" width="160px"/>

[![Build and run testsuite](https://github.com/codenotary/cas/actions/workflows/pull.yml/badge.svg)](https://github.com/codenotary/cas/actions/workflows/pull.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/codenotary/cas)](https://goreportcard.com/report/github.com/codenotary/cas)
[![Docker pulls](https://img.shields.io/docker/pulls/codenotary/cas?style=flat-square)](https://hub.docker.com/r/codenotary/cas)
[![Changelog](https://img.shields.io/badge/CHANGELOG-.md-blue?style=flat-square)](CHANGELOG.md)
[![Release](https://img.shields.io/github/release/codenotary/cas.svg?style=flat-square)](https://github.com/codenotary/cas/releases/latest)

Give any digital asset a meaningful, globally-unique, immutable identity that is authentic, verifiable, traceable from anywhere.

<img align="right" src="docs/img/codenotary_mascot.png" width="256px"/>
When using Codenotary CAS in source code, release, deployment or at runtime, you allow a continuous trust verification that can be used to detect unusual or unwanted activity in your workload and act on it.
<br/>
Powered by Codenotary's digital identity infrastructure, CAS lets you Attest all your digital assets that add a trust level of your choice, custom attributes and meaningful status without touching or appending anything (unlike digital certificates). That allows change and revocation post-release without breaking any customer environment.
<br/>
Everything is done in a global, collaborative way to break the common silo solution architecture. Leveraging an immutable always-on platform allows you to avoid complex setup of Certificate authorities or digital certificates (that are unfit for DevOps anyway).

----
> :warning: **From version v0.10 a major refactory has replaced the old VCN CLI. While the old VCN versions are available to download in the release section, we don't provide support and maintenance anymore.** 
----

## Table of contents

- [Quick start](#quick-start)
- [DevSecOps in mind](#devsecops-in-mind)
- [What kind of behaviors can Codenotary cas detect](#what-kind-of-behaviors-can-codenotary-cas-detect)
- [Installation](#installation)
- [Usage](#usage)
- [Integrations](#integrations)
- [Documentation](#documentation)
- [Advanced Usage](#advanced-usage)
- [License](#license)

## Quick start

1. [**Create your identity (free)**](https://cas.codenotary.com) - You will get an `API_KEY` from our free cloud CAS Cloud.  


2. **Download Codenotary CAS**


   ```
   bash <(curl http://getcas.codenotary.io -L)
   ```
   
> For Windows users, donwload your binay [here](https://github.com/codenotary/cas/releases/latest).


3. **Login**

   ```bash
   export CAS_API=<your API KEY>; cas login
   ```
   

4. **Create a Software Bill of Materials (SBOM)**

   ```bash
   cas bom docker://wordpress
   ```

4. **Attest your assets** Attestation is the combination of Notarization (creating digital proof of an asset) and Authentication (getting the authenticity of an asset).

    Notarize an asset: 

   ```bash
   cas notarize docker://wordpress
   ```
   
   Authenticate an asset: 

   ```bash
   cas authenticate docker://wordpress
   ```
----
&nbsp;


## Table of contents

- [DevSecOps in mind](#devsecops-in-mind)
- [What kind of behaviors can Codenotary cas detect](#what-kind-of-behaviors-can-codenotary-cas-detect)
- [Installation](#installation)
- [Usage](#usage)
- [Integrations](#integrations)
- [Documentation](#documentation)
- [More Examples](#examples)
- [License](#license)

&nbsp;


## DevSecOps in mind
Codenotary cas is a solution written by devops-obsessed engineers for Devops engineers to bring better trust and security to the the CloudNative source to deployment process

## What kind of behaviors can Codenotary cas detect
cas (and its extensions for Docker, Kubernetes, documents or CI/CD) can detect, authenticate and alert on any behavior that involves using unauthentic digital assets. cas verification can be embedded anywhere and can be used to trigger alerts, updates or workflows.

cas detects or acts on the following (but not limited to):
* Immutable tagging of source code, builds, and container images with version number, owner, timestamp, organization, trust level, and much more
* Simple and tamper-proof extraction of notarized tags like version number, owner, timestamp, organization, and trust level from any source code, build and container (based on the related image)
* Quickly discover and identify untrusted, revoked or obsolete libraries, builds, and containers in your application
* Detect the launch of an authorized or unknown container immediately
* Prevent untrusted or revoked containers from starting in production
* Verify the integrity and the publisher of all the data received over any channel

and more
* Enable application version checks and actions
* Buggy or rogue libraries can be traced by simple revoke or unsupport
* Revoke or unsupport your build or build version post-deployment (no complex certificate revocation that includes delivery of newly signed builds)
* Stop unwanted containers from being launched
* Make revocation part of the remediation process
* Use revocation without impairing customer environments
* Trace source code to build to deployment by integration into CI/CD or manual workflow
* Tag your applications for specific use cases (alpha, beta - non-commercial aso).

not just containers, also virtual machines -  [check out vCenter Connector, in case you're running VMware vSphere](https://github.com/openfaas-incubator/vcenter-connector)
* Newly created or existing virtual machines automatically get a unique identity that can be trusted or untrusted
* Prevent launch of untrusted VMs
* Stop or suspend running outdated or untrusted VMs
* Detect the cloning or export of VMs and alert

&nbsp;


## Installation

### Download binary

It's easiest to download the latest version for your platform from the [release page](
https://github.com/codenotary/cas/releases).

Once downloaded, you can rename the binary to `cas`, then run it from anywhere.
> For Linux and macOS you need to mark the file as executable: `chmod +x cas`

### Homebrew / Linuxbrew

If you are on macOS and using [Homebrew](https://brew.sh/) (or on Linux and using [Linuxbrew](https://linuxbrew.sh/)), you can install `cas` with the following:

```
brew tap codenotary/brew
brew install cas
```

### Build from Source

After having installed [golang](https://golang.org/doc/install) 1.12 or newer clone this
repository into your working directory.

Now, you can build `cas` in the working directory by using `make cas` and then run `./cas`.

Alternatively, you can install `cas` in your system simply by running `make install`. This will put the `cas` executable into `GOBIN` which is
accessible throughout the system.

### yum and deb (TBD)

&nbsp;


## Usage

Basically, `cas` can notarize or authenticate any of the following kind of assets:

- a **file**
- a **git commit** (by prefixing the local git working directory path with `git://`)
- a **container image** (by using `docker://` or `podman://` followed by the name of an image present in the local registry of docker or podman, respectively)

> It's possible to provide a hash value directly by using the `--hash` flag.

For detailed **command line usage** see [docs/cmd/cas.md](docs/cmd/cas.md) or just run `cas help`.

### Wildcard support and recursive notarization

 It's also possible to notarize assets using wildcard.
 With `--recursive` flag is possible to iterate over inner directories.
```shell script
./cas n "*.md" --recursive
```


### Notarization

Register an account with [codenotary.com](https://cas.codenotary.com) first.

Then start with the `login` command. `cas` will walk you through login and importing up your secret upon initial use.
```
cas login
```

Once your secret is set you can notarize assets like in the following examples:

```
cas notarize <file>
cas notarize docker://<imageId>
cas notarize podman://<imageId>
cas notarize git://<path_to_git_repo>
cas notarize --hash <hash>
```

By default all assets are notarized private, so not much information is disclosed about the asset. If you want to make that public and therefore, more trusted, please use the `--public` flag.

```
cas notarize --public <asset>
```

Change the asset's status:

```
cas unsupport <asset>
cas untrust <asset>
```

Finally, to fetch all assets you've notarized:

```
cas list
```

### Authentication

```
cas authenticate <file>
cas authenticate docker://<imageId>
cas authenticate podman://<imageId>
cas authenticate git://<path_to_git_repo>
cas authenticate --hash <hash>
```

:bulb: Public authentication is also possible without having an CAS_API_KEY - more info here [Public Authentication](#public-authentication)

To output results in `json` or `yaml` formats:
```
cas authenticate --output=json <asset>
cas authenticate --output=yaml <asset>
```
> Check out the [user guide](docs/user-guide/formatted-output.md) for further details.

&nbsp;


## Integrations

* [Github Action](https://github.com/marketplace/actions/verify-commit) - An action to verify the authenticity of your commits within your Github workflow
* [docker](docs/user-guide/schemes/docker.md) - Out of the box support for notarizing and authenticating Docker images.
* [hub.docker.com/r/codenotary/cas](https://hub.docker.com/r/codenotary/cas) - The `cas`'s DockerHub repository.

&nbsp;

## Documentation

* [Community Attestation Service jumpstart](docs/cas-cnlc-jumpstart.md)
* [Command line usage](docs/cmd/cas.md)
* [Configuration](docs/user-guide/configuration.md)
* [Environments](docs/user-guide/environments.md)
* [Formatted output (json/yaml)](docs/user-guide/formatted-output.md)
* Notarization explained (TBD)

&nbsp;

## Advanced Usage

#### Authenticate a Docker image automatically prior to running it

First, you’ll need to pull the image by using:

```
docker pull hello-world
```

Then use the below command to put in place an automatic safety check. It allows only verified images to run.

```
cas authenticate docker://hello-world && docker run hello-world
```
If an image was not verified, it will not run and nothing will execute.


#### Authenticate multiple assets
You can authenticate multiple assets by piping other command outputs into `cas`:
```
ls | xargs cas authenticate
```
> The exit code will be `0` only if all the assets in you other command outputs are verified.

#### Authenticate by a specific signer
By adding `--signerID`, you can authenticate that your asset has been signed by a specific SignerID.
> A SignerID is the signer public address (represented as a 40 hex characters long string prefixed with `0x`).

```
cas authenticate --signerID 0x8f2d1422aed72df1dba90cf9a924f2f3eb3ccd87 docker://hello-world
```

#### Authenticate by a list of signers

If an asset you or your organization wants to trust needs to be verified against a list of signers as a prerequisite, then use the `cas authenticate` command and the following syntax:

- Add a `--signerID` flag in front of each SignerID you want to add
(eg. `--signerID 0x0...1 --signerID 0x0...2`)
- Or set the env var `cas_SIGNERID` correctly by using a space to separate each SignerID (eg. `cas_SIGNERID=0x0...1 0x0...2`)
> Be aware that using the `--signerID` flag will take precedence over `cas_SIGNERID`.

The asset authentication will succeed only if the asset has been signed by at least one of the signers.

#### Authenticate using the asset's hash

If you want to authenticate an asset using only its hash, you can do so by using the command as shown below:

```
cas authenticate --hash fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e
```

#### Unsupport/untrust an asset you do not have anymore

In case you want to unsupport/untrust an asset of yours that you no longer have, you can do so using the asset hash(es) with the following steps below.

First, you’ll need to get the hash of the asset using the `cas list` command. Then, in the CLI, use:

```
cas untrust --hash <asset's hash>
# or
cas unsupport --hash <asset's hash>
```


#### TLS

By default, cas will try to establish a secure connection (TLS) with Community Attestation Service.

The user can also provide a custom TLS certificate for the server, in case cas is not able to download it automatically:

```shell script
cas login --port 443 --host cas.codenotary.com --cert mycert.pem
```

For testing purposes or in case the provided certificate should be always trusted by the client, the user can
configure cas to skip TLS certificate verification with the `--skip-tls-verify` option:

```shell script
cas login --port 443 --host cas.codenotary.com --cert mycert.pem --skip-tls-verify
```

Finally in case the Community Attestation Service is not exposed through a TLS endpoint, the user can request a cleartext
connection using the `--no-tls` option:

```shell script
cas login --port 80 --host cas.codenotary.com  --no-tls
```

#### Verify CAS server identity
Every message returned by CAS is cryptographically signed.
In order to verify the identity of the server you can calculate locally the fingerprint and compare it with the following:

`SHA256:Re5IAHGkYk32xfnG8txbwJuJPVFe8Mf5AOv3bLg6XsY`

To generate local fingerprint use the following commands:
```shell
ssh-keygen -i -m PKCS8 -f ~/.cas-trusted-signing-pub-key > mykey.pem.pub
ssh-keygen -l -v -f mykey.pem.pub
rm mykey.pem.pub
```

### Add custom metadata when signing assets
The user can upload custom metadata when doing an asset notarization using the `--attr` option, e.g.:

```shell script
cas n README.md --attr Testme=yes --attr project=5 --attr pipeline=test
```

This command would add the custom asset metadata Testme: yes, project: 5, pipeline: test.

The user can read the metadata back on asset authentication, i.e. using the `jq` utility:

```shell script
cas a README.md -o json | jq .metadata
```

### Inspect
Inspect has been extended with the addition of new filter: `--last`, `--first`, `--start` and `--end`.
With `--last` and `--first` are returned the N first or last respectively.

```shell script
cas inspect document.pdf --last 10
```

With `--start` and `--end` it's possible to use a time range filter:

```shell script
cas inspect document.pdf --start 2020/10/28-08:00:00 --end 2020/10/28-17:00:00
```

If no filters are provided only maximum 10 items are returned.

### Signer Identifier
It's possible to filter results by signer identifier:

```shell script
cas inspect document.pdf --signerID CygBE_zb8XnprkkO6ncIrbbwYoUq5T1zfyEF6DhqcAI=
```

### Public Authentication

The authentication is performed by a user possessing an `CAS_API_KEY` issued by the Community Attestation Service. But there are situations in which an anonymous authentication is needed: for example the authentication is performed by a GitHub action in an Open Source repository. For such scenarios, a public authentication is possible, where the authentication process does not need an `CAS_API_KEY` - nevetheless the `SIGNER_ID` has to be defined. Example:


```
cas authenticate --signerID 0xxxxxxxxxxxxxxxxxxxxxxxxxxx docker://hello-world
```
 
## License

This software is released under [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0).

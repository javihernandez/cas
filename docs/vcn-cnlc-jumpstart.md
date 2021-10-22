# cas - Community Attestation Service jumpstart

## Table of contents

- [Community Attestation Service](#community-attestation-service)
- [Quick start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [Documentation](#documentation)



## Community Attestation Service

cas has been extended in order to be compatible with [Community Attestation Service](https://cas.codenotary.com/).
Notarized assets informations are stored in a tamperproof ledger with cryptographic verification backed by [immudb](https://codenotary.com/technologies/immudb/), the immutable database.
Thanks to this `cas` is faster and provides more powerful functionalities like local data inclusion, consistency verification and enhanced CLI filters.

### Obtain an API Key
To provide access to Immutable Ledger a valid API Key is required. The key can be obtained from [Community Attestation Service](https://cas.codenotary.com/).

## Quick start

1. **Download CodeNotary cas.** There are releases for different platforms:

- [Download the latest release](https://github.com/codenotary/cas/releases/latest) and then read the [Usage](#usage) section below.
- We recommend storing `cas` in your `PATH` - Linux example:
   ```bash
   cp cas-v<version>-linux-amd64 /usr/local/bin/cas
   ```

2. **Authenticate digital objects** You can use the command as a starting point.

   ```bash
   cas login --host cnlc-host.com --port 443
   cas authenticate <file|dir://directory|docker://dockerimage|git://gitdirectory|javacom://javacomponent|gocom://gocomponent|pythoncom://pythoncomponent|dotnetcom://dotnetcomponent>
   ```


3. **Notarize existing digital objects** Once you have an account you can start notarizing digital assets to give them an identity.

   ```bash
   # cas login can be skipped, if already performed
   cas login --host cnlc-host.com --port 443
   cas notarize <file|dir://directory|docker://dockerimage|git://gitdirectory>
   ```

### Login

To login in Immutable Ledger provides `--port` and `--host` flags, also the user submit API Key when requested.
Once host, port and API Key are provided, it's possible to omit them in following commands. Otherwise, the user can provide them in other commands like `notarize`, `verify` or `inspect`.

```shell script
cas login --port 443 --host cnlc-host.com
```

> One time password (otp) is not mandatory

Alternatively, for using cas in non-interactive mode, the user can supply the API Key via the `CAS_API_KEY` environment variable, e.g.:

```shell script
export CAS_API_KEY=apikeyhere

# No cas login command needed

# Other cas commands...
cas notarize asset.txt --host cnlc-host.com --port 443
```

#### TLS

By default, cas will try to establish a secure connection (TLS) with a Immutable Ledger server.

The user can also provide a custom TLS certificate for the server, in case cas is not able to download it automatically:

```shell script
cas login --port 443 --host cnlc-host.com --cert mycert.pem
```

For testing purposes or in case the provided certificate should be always trusted by the client, the user can
configure cas to skip TLS certificate verification with the `--skip-tls-verify` option:

```shell script
cas login --port 443 --host cnlc-host.com --cert mycert.pem --skip-tls-verify
```

Finally in case the Immutable Ledger Server is not exposed through a TLS endpoint, the user can request a cleartext
connection using the `--no-tls` option:

```shell script
cas login --port 80 --host cnlc-host.com --no-tls
```

## Installation

### Download binary

It's easiest to download the latest version for your platform from the [release page](
https://github.com/codenotary/cas/releases).

Once downloaded, you can rename the binary to `cas`, then run it from anywhere.
> For Linux and macOS you need to mark the file as executable: `chmod +x cas`

### Homebrew / Linuxbrew

If you are on macOS and using [Homebrew](https://brew.sh/) (or on Linux and using [Linuxbrew](https://linuxbrew.sh/)), you can install `cas` with the following:

```
brew tap vchain-us/brew
brew install cas
```

### Build from Source

After having installed [golang](https://golang.org/doc/install) 1.15 or newer clone this
repository into your working directory.

Now, you can build `cas` in the working directory by using `make cas` and then run `./cas`.

Alternatively, you can install `cas` in your system simply by running `make install`. This will put the `cas` executable into `GOBIN` which is
accessible throughout the system.

## Usage

Basically, `cas` can notarize or authenticate any of the following kind of assets:

- a **file**
- an entire **directory** (by prefixing the directory path with `dir://`)
- a **git commit** (by prefixing the local git working directory path with `git://`)
- a **container image** (by using `docker://` or `podman://` followed by the name of an image present in the local registry of docker or podman, respectively)

It's possible to provide a hash value directly by using the `--hash` flag.

For detailed **command line usage** see [docs/cmd/cas.md](https://github.com/codenotary/cas/blob/master/docs/cmd/cas.md) or just run `cas help`.

### Wildcard support and recursive notarization

It's also possible to notarize assets using a wildcard pattern.

With `--recursive` flag the utility can recursively notarize inner directories.
```shell script
./cas n "*.md" --recursive
```
### Notarization

Start with the `login` command. `cas` will walk you through login and importing up your secret upon initial use.

```
cas login --host cnlc-host.com --port 443
```

Once your secret is set you can notarize assets like in the following examples:

```
cas notarize <file>
cas notarize dir://<directory>
cas notarize docker://<imageId>
cas notarize podman://<imageId>
cas notarize git://<path_to_git_repo>
cas notarize --hash <hash>
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
cas authenticate dir://<directory>
cas authenticate docker://<imageId>
cas authenticate podman://<imageId>
cas authenticate git://<path_to_git_repo>
cas authenticate --hash <hash>
```

To output results in `json` or `yaml` formats:
```
cas authenticate --output=json <asset>
cas authenticate --output=yaml <asset>
```
> Check out the [user guide](https://github.com/codenotary/cas/blob/master/docs/user-guide/formatted-output.md) for further details.

## Documentation

* [Command line usage](https://github.com/codenotary/cas/blob/master/docs/cmd/cas.md)
* [Configuration](https://github.com/codenotary/cas/blob/master/docs/user-guide/configuration.md)
* [Environments](https://github.com/codenotary/cas/blob/master/docs/user-guide/environments.md)
* [Formatted output (json/yaml)](https://github.com/codenotary/cas/blob/master/docs/user-guide/formatted-output.md)
* [Notarization explained](https://github.com/codenotary/cas/blob/master/docs/user-guide/notarization.md)

## Examples

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
ls | xargs -n 1 cas authenticate
```
> The exit code will be `0` only if all the assets in you other command outputs are verified.

#### Authenticate by a specific signer
By adding `--signerID`, you can authenticate that your asset has been signed by a specific SignerID.
> A SignerID is the signer public address (represented as a 40 hex characters long string prefixed with `0x`).

```
cas authenticate --signerID 0x8f2d1422aed72df1dba90cf9a924f2f3eb3ccd87 docker://hello-world
```

#### Authenticate using the asset's hash

If you want to authenticate an asset using only its hash, you can do so by using the command as shown below:

```
cas authenticate --hash fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e
```

#### Unsupport/untrust an asset you do not have anymore

In case you want to unsupport/untrust an asset of yours that you no longer have, you can do so using the asset hash(es) with the following steps below.

First, you’ll need to get the hash of the asset from your Community Attestation Service dashboard or alternatively you can use the `cas list` command. Then, in the CLI, use:

```
cas untrust --hash <asset's hash>
# or
cas unsupport --hash <asset's hash>
```

#### Notarization within automated environments

Simply, set up your environment accordingly using the following commands:

```bash
export CAS_API_KEY=apikeyhere
```

Once done, you can use `cas` in your non-interactive environment using:

```
cas login --host cnlc-host.com --port 443
cas notarize <asset>
```

> Other commands like `untrust` and `unsupport` will also work.


#### Add custom metadata when signing assets
The user can upload custom metadata when doing an asset notarization using the `--attr` option, e.g.:

```shell script
cas n README.md --attr Testme=yes --attr project=5 --attr pipeline=test
```

This command would add the custom asset metadata Testme: yes, project: 5, pipeline: test.

The user can read the metadata back on asset authentication, i.e. using the `jq` utility:

```shell script
cas a README.md -o json | jq .metadata
```

#### Inspect
Inspect has been extended with the addition of new filter: `--last`, `--first`, `--start` and `--end`.
With `--last` and `--first` are returned the N first or last respectively.

```shell script
cas inspect document.pdf --last 10
```

With `--start` and `--end` it's possible to use a time range filter:

```shell script
cas inspect document.pdf --start 2020/10/28-08:00:00 --end 2020/10/28-17:00:00
```

If no filters are provided only maximum 100 items are returned.

#### Signer Identifier
It's possible to filter results by signer identifier:

```shell script
cas inspect document.pdf --signerID CygBE_zb8XnprkkO6ncIrbbwYoUq5T1zfyEF6DhqcAI=
```

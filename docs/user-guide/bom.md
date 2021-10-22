# Bill of Materials (BoM)

`cas` can identify, authenticate and notarize dependencies of the software assets.

## Supported languages/environments


| Language(s)/environment | Component scheme | Source | Package manager |
|-|-|-|--|
| Docker image (RPM-based distro) | `docker` | RPM database (`/var/lib/rpm/Packages`) | rpm |
| Docker image (Debian-based distro) | `docker` | DPKG database (`/var/lib/dpkg`) | dpkg |
| Docker image (Alpine) | `docker` | APK database (`/lib/apk/db/installed`) | apk |V


The following enviroments are supported only for Codenotary Enterprise Edition, more info at [Codenotary.com](https://codenotary.com)

| Language(s)/environment | Component scheme | Source | Package manager |
|-|-|-|--|
| Go | `gocom` |compiled binary | |
| | | directory with `go.sum` file | |
| | | directory with `*.go` file(s) | |
| Python | `pythoncom` | `Pipfile.lock` file or directory containing this file | pipenv |
| | | `poetry.lock` file or directory containing this file | poetry |
| | | `requirements.txt` file or directory containing this file | pip |
| JVM (Java, Scala, Kotlin) | `javacom` | JAR file containing `pom.xml` | maven |
| .Net (C#, F#, Visual Basic) | `dotnet` | `*.sln` file or directory containing this file | NuGet |
| &nbsp; - C#, F# only| | `*.csproj` file or directory containing this file | |
| &nbsp; - Visual Basic only| | `*.vbproj` file or directory containing this file | |
| JavaScript | `nodecom` | `package-lock.json` file or directory containing this file | npm |


## Working with builds

### Resolving dependencies

`cas bom <asset> [bom output options]`

This command resolves the dependencies for the asset and prints out the list of dependencies.


See [output options](#output-options) for details about outputting BoM in standard formats.

Examples:
```
cas bom docker://alpine
cas bom docker://ubuntu:20.04 --bom-spdx ubuntu.spdx
```

### Authentication

`cas a --bom <asset> [bom options] [bom output options]`

This command resolves the dependencies for the asset, authenticates the dependencies and the asset, and prints out the list of dependencies with their trust levels.

Following `bom options` modify the behavior of this command:

| Option | Default | Description |
|-|-|-|
| `--signerID` | current user | Signer ID to use for dependency and asset authentication. This isn't a BoM-specific options, but it has a special meaning for BoM |
| `--bom-trust-level` | `trusted` | Minimal accepted trust level for the dependencies (or its abbreviation), one of: |
||| `untrusted` (`unt`) |
||| `unsupported` (`uns`) |
||| `unknown` (`unk`) |
||| `trusted` (`t`) |
| `--bom-max-unsupported` | `0` | Max number of unsupported/unknown dependencies to accept, in percent. If number of unsupported/unknown dependencies doesn't exceed this threshold, authentication is considered successful |
| `--bom-batch-size` |`10` | Send requests to server in batches of specified size |

Any of this options (except) implies `--bom` mode.

See [output options](#output-options) for details about outputting BoM in standard formats.

This command returns one of the following exit codes:

- `0` - success
- `1` - any dependency or BoM source is untrusted
- `2` - any dependency or BoM source is unknown and there are no untrusted or unsupported dependencies
- `3` - any dependency or BoM source is unknown and there are no untrusted dependencies

Examples:
```
cas a --bom docker://alpine --signerID auditor
cas a docker://ubuntu:20.04 --bom-trust-level unknown --bom-spdx ubuntu.spdx
```

### Notarization

`cas n --bom <asset> [bom options] [bom output options]`

This command resolves the dependencies for the asset, authenticates and notarizes the dependencies (only the unknown one, is `--bom-force` is not specified) and the asset, and prints out the list of dependencies with their trust levels.

Following options modify the behavior of this command:

| Option | Default | Description |
|-|-|-|
| `--bom-signerID` | current user | Signer ID to use for dependency authentication |
| `--bom-batch-size` |`10` | Send requests to server in batches of specified size |

Any of this options () implies `--bom` mode.

Examples:
```
cas n docker://alpine --bom-signerID auditor
cas n docker://ubuntu:20.04 --bom-spdx ubuntu.spdx
```

### Output options

User can specify one or several options to output BoM in different supported standard formats.

| Option | Description |
|-|-|
| `--bom-spdx` | Name of output SPDX tag-value file |
| `--bom-cyclonedx-json` | Name of output CycloneDX JSON file |
| `--bom-cyclonedx-xml` | Name of output CycloneDX XML file |

Any of this options implies `--bom` mode.

## Working with individual dependencies

`cas a|n|ut|us <scheme>://<name>@<version> | --hash <hash>`

Individual components are authenticated/notarized/unsupported/untrusted as any other asset, but you need to specify either component hash with `--hash` option.

Examples:
```
cas n --hash 691631371bfa886425c956999a4e998181036be260d7c0f179b3d2adde9b8353
cas ut --hash 6dbb9cc54074106d46d4ccb330f2a40a682d49dda5f4844962b7dce9fe44aaec
```

## Support for Docker

`cas <command> docker://<image>[:<tag>] [command options]`

When asset has `docker` scheme, `cas` starts the container for the specified `<image>:<tag>` and finds the dependencies, therefore docker daemon must be running and required image is already pulled. `cas` supports Linux distributions that use `apk` (Alpine), `dpkg` (Debian, Ubuntu) or `rpm` (RedHat, Fedora, CentOS, AlmaLinux, openSUSE etc.) package managers.

As always with Docker, missing image `tag` implies `latest`.

Examples:
```
cas bom docker://alpine --bom-spdx docker.spdx
cas a --bom docker://debian
cas n --bom docker://nginx:stable-alpine
```

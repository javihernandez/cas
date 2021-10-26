# CHANGELOG
All notable changes to this project will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
<a name="unreleased"></a>
## [Unreleased]

<a name="v0.10.0-RC-2"></a>
## [v0.10.0-RC-2] - 2021-10-26
### Bug Fixes
- fixed integration tests
- updated documentation
- correctly limit the number of inspected items

### Changes
- added support for plain e-mail as signer ID

<a name="v0.10.0-RC-1"></a>
## [v0.10.0-RC-1] - 2021-10-22
### Bug Fixes
- added ability to do not write on read only dir (fixes [#45](https://github.com/codenotary/cas/issues/45))
- gitignore
- dashboard cmd has no args
- temporary dirs creation at startup
- remove log.Fatal in favour of returing errors
- Make test file content unique to avoid parallel notarization
- add missing header comments
- improve err msgs when auth is required
- token deletion when already logged out
- ensure default keystore before looking for keys
- do not quit when executed by Win context menu
- fix return error codes on cn.io
- do not exit early when login is required
- unsuccessful status code
- find home dir in the right way
- vcn docker notarization readme.md section
- restore info when no user is logged
- prefix for docker assets
- return proper error by checking quota before
- remove test-resources from dockerignore
- apikey is removed from userlist when logout
- fix autodiscover dir artifact
- if logged in always show user own signature
- clear context when logging out
- verify fallback
- do not print verification error message if --output=json
- load artifact wait for immudb set to be indexed
- fix state service and verified methods
- correct err msg when token has expired
- restore verification with the highest level available when not logged
- fix integration tests and add circleci job for go 1.15
- fix dist release go1.15.6
- do not read in env var for config
- print correct config file name
- clean other context when user logs in
- fix ledger-compliance-go dependency
- ca-certificates for Dockerfile
- `BlockchainVerification` json unmarshalling
- correct `Signer` field semantic
- correct error msg when entering empty email (fixes [#43](https://github.com/codenotary/cas/issues/43))
- make fields name consistent across output formats (json, yaml)
- contentType was not shown when `ArtifactResponse` was not present
- simplied error message when wrong notarization password is provided
- Pull latest xgo (go cross compilation tools) when creating the full stack
- trust level was not shown when equals to DISABLED
- correct cross compile deps
- update go.sum for xgo
- attachments internal map is iterated with a side sorted string slice. Add some cleanup
- ask for confirmation before quitting
- **api:** correct alert config JSON name
- **api:** correct secret download content format
- **api:** nil pointer dereference
- **api:** correct tracing log label for `VerifyMatchingSignerIDWithFallback`
- **api:** correct alert error message
- **build:** correct NSIS installation directory
- **build/makensis:** correct debian image version for NSIS
- **bundle:** use uint64 for size
- **cmd:** show config file path if not default one
- **cmd:** enforce lowercase hash string (fixes [#67](https://github.com/codenotary/cas/issues/67))
- **cmd:** typo and language improvement
- **cmd:** correct flags collision
- **cmd/inspect:** enforce lowercase hash string
- **cmd/serve:** add missing event tracking
- **cmd/serve:** CORS preflight (fixes [#60](https://github.com/codenotary/cas/issues/60))
- **cmd/serve:** enforce lowercase hash string
- **cmd/sign:** try again when notarization password is wrong
- **cmd/sign:** fix missing alert sign nil pointer error
- **cmd/sign:** command must accept exactly one arg
- **cmd/sign:** cleanup unused arg to make --hash work again
- **cmd/sign:** do not track events when errored
- **cmd/verify:** show error when alert fails
- **dir:** OS agnostic paths
- **docker:** accept images but not other objects
- **extractor:** correct empty files handling
- **extractor:** nil pointer dereference
- **extractor/dir:** runtime error when passing a regular file instead of a directory (fixes [#56](https://github.com/codenotary/cas/issues/56))
- **list:** show all assets with pagination (fixes [#28](https://github.com/codenotary/cas/issues/28))
- **login:** silence usage message when errored
- **migrate:** do not create a new key if it already exists
- **pkg:** add version metadata header and fix signer ID extraction from API key to support email addresses as signer IDs
- **pkg:** fix vcn blockchain login
- **pkg/api:** swap signerID position with the hash in the vcn lc key format
- **pkg/api:** fix artifact status json unmarshall
- **pkg/api:** return vcn lc apikey error missing only in lc context
- **pkg/api:** handle JSON marshal error when creating artifact
- **pkg/api/verify:** fix pkg api verify
- **pkg/cicontext:** remove sensitive data
- **pkg/cmd:** version check is enabled only for codenotary.io context
- **pkg/cmd:** don't track untrusted verification
- **pkg/cmd/inspect:** fix inspect ux, examples and errors
- **pkg/cmd/internal/cli:** return an error if empty api key is provided
- **pkg/cmd/login:** fix env var on login
- **pkg/cmd/serve:** improving error if no content is submitted while signing
- **pkg/cmd/sign:** improve attachment container structure
- **pkg/cmd/sign:** sign error handler fix
- **pkg/cmd/sign:** fix cnlc hash signing
- **pkg/cmd/verify:** attachments can be used only on CNIL context. On blockchain error is returned
- **pkg/cmd/verify:** fix attachments printing when artifact is not notarized
- **pkg/cmd/verify:** fix attachments download with multiple version
- **pkg/cmd/verify:** load artifact error handler fix
- **pkg/cmd/verify:** fix empty signerid handling
- **pkg/meta:** add apikey env var
- **pkg/store:** fix api keys clean on load config
- **sign:** cross platform loading spinner
- **sign:** remove ownership disclaimer
- **sign:** flag accessed but not defined: hash
- **store:** correct key funcs and tests
- **tracking:** send events at right place in time
- **verify:** cross platform coloured printing
- **verify:** only print size when available
- **verify:** switch to single key for current user

### Changes
- update dependencies
- attach-full-download flag is removed and full download is default behaviour
- add dev stage
- generate markdown docs for commands
- use two-stage build
- use make install for Dockerfile
- rename CNLC to CNIL
- fix badge url
- upgrade cnil sdk
- upgrade cnil sdk in order to fix file lock issue
- remove circleCI
- move file hashing func to internal
- fix missing apikey error
- remove unused import
- add artifact print date method
- add proper error handling on notFound or notVerified
- add lcLedger and signerID parameters on PublicCNLCVerify
- expose lc verify api and explicit environment vars
- introducing Cobra
- update ledger-compliance-go dependency to fix missing metadata bug
- remove 'blockchain' from login message
- vcn context ux improvements. closes [#134](https://github.com/codenotary/cas/issues/134)
- increase pipeline timeout
- fix attachments download and add friendly print of attachments when downloading
- update immudb dependency to latest 0.9.2 version and ledger-compliance-go to latest master
- bump sdk version, it fixes notarization on empty cnlc ledger
- improve logout message
- no action on logout if there is no user logged
- authenticate return 2 exit code if hash is not found
- if --lc-no-tls is provided default port is 80
- new no api key configured error message
- renaming flag config in vcnpath and set default vcnpath in os.tmp folder
- bump lc sdk
- bump lc sdk
- verify if config is present in context check
- clean configuration file by removing api keys
- user saving removed on lc login
- disable version check in lc context
- adjust colors settings with fgreen
- adjust colors settings
- handle blue screen colors
- minor changes on releasing guide, makefile and minor dependencies updates
- update lc-sdk dependency
- list of signers is not supported in CNLC
- wildcard instructions in their own section
- add a separated CNLC jumpstart
- code clean up
- partial api logging refactor
- switch main to cobra
- bumb lc sdk version
- restore lc inspect capabilities
- enable InsecureSkipVerify when tls-skip-verify is active
- disable verified set in lc artifact creation
- bump lc sdk version
- bump lc sdk version
- bump sdk version, remove useless package
- bump lc sdk version
- update gitignore
- correct printing func
- improve wallet not synced message
- fix inspect full history message
- upgrade go version to 1.15 and sdk
- move logs to internal
- introduce store package
- add swagger.json
- bumb version of go-ethereum and smart contracts
- remove all whitespaces from otp
- add store context and key utils
- code cleanup and minor improvements
- update messages
- distinct event for alert verification
- improve naming
- move cli package to internal
- no race test for make install
- bump version
- printing stuff
- code cleanup and fix comments
- reduce public api surface
- disable config file message
- improve "your assets will not be uploaded" notice
- update copyright year
- move tracking call outside api pkg
- reorg publisher fields
- minor printing and marshaling improvements
- make createArtifact private
- removed profile migration from v0.3.x
- change metadata naming
- update Smart Contract addresses to latest version
- improve API and user messages to reflect the new terminology
- deprecate `--key` in favor of `--signerID`
- use go v1.13 for the building system
- add "your assets will not be uploaded" message
- fix login help message
- deprecate `KEYSTORE_PASSWORD` env var in favour of `VCN_NOTARIZATION_PASSWORD`
- add Platform and Architecture fields
- go mod tidy
- correct typos
- promote artifact fields to root
- remove onboarding message (closes [#52](https://github.com/codenotary/cas/issues/52))
- executable file sniffer
- new AssetsRelay smart contract
- platform managed secret is not stored locally anymore
- minor refactoring of login/logout
- **api:** add Artifact.Copy()
- **api:** omit empy field in `ArtifactResponse`
- **api:** add `BlockchainVerification.UnmarshalYAML`
- **api:** removed `BlockchainVerification.LevelName()` method
- **api:** size must be valid, name can be empty
- **api:** allow empty size
- **api:** fine-tune blockchain verify funcs
- **api:** more fields for Error struct
- **api:** add name to alert struct
- **api:** refactor to single address and code cleanup
- **api:** deprecate publisher fields
- **api:** improve alert error messages
- **api:** move type Error to its own source file
- **api:** deprecate public publisher APIs
- **api:** load user's artifact by hash only
- **build:** make makefile more resilient
- **build:** minor improvements
- **build:** improve makefile
- **build:** enforce xgo image digest
- **bundle:** rework diff
- **bundle:** return error if distinct sizes are found for the same digest
- **bundle:** descriptor test
- **ci:** add verify-commit action
- **ci:** github action for docs
- **ci:** initial circleci setup
- **cli:** spinner changed and code cleaned up
- **cmd:** update cobra to 0.5.0
- **cmd:** internal types
- **cmd:** added ability to use empty notarization password
- **cmd:** gracefully print error message if artifact is nil (fixes [#57](https://github.com/codenotary/cas/issues/57))
- **cmd:** no echo when asking for confirmation before quitting
- **cmd:** use extractors
- **cmd:** always show local extracted metadata
- **cmd/info:** improved message for expired tokens
- **cmd/internal/cli:** reader and writer for YAML files
- **cmd/serve:** notarization API will allow only registered schemes for the kind field (closes [#51](https://github.com/codenotary/cas/issues/51))
- **cmd/sign:** return full output when notarizing
- **cmd/sign:** allow notarization when no dir with write access by default
- **cmd/verify:** if api-key is revoked status will be shown as revoked in authenticate
- **cmd/verify:** trigger alert on asset error, if any
- **cmd/verify:** do not show "diff unavailable" when no manifest is found
- **cmd/verify:** manifest lookup from storage, fallback to target dir
- **dir/extractor:** ignore irregular files
- **dist:** omit symbol table and debug info
- **env:** remove VERIFY_ prefix
- **extractor:** do not store version in metadata when empty
- **extractor/dir:** option to skip ignore file creation
- **extractor/dir:** do not create `.vcnignore` on directory authentication (refs [#45](https://github.com/codenotary/cas/issues/45))
- **help:** add info about assets and env vars
- **internal:** move cli package
- **internal:** correct errors funcs and tests
- **log:** improve metahash tracing
- **meta:** use Stringer interface for stage environment
- **meta:** add StaticBool()
- **meta:** increase tx verification rounds
- **meta:** removed event tracking for keystore creation
- **meta:** use Stringer interface for status, level, and visibility
- **meta:** increse signing timeout
- **output:** correct WriteTo interface
- **output:** silence printing when formatted output
- **pkg/api:** handle parallel requests with same hash on blockchain
- **pkg/api:** fix uid friendly label extraction from key
- **pkg/api:** artifact load uses getAt
- **pkg/cicontext:** add context probes and minor improvements
- **pkg/cmd:** show status unknown if tamper is detected
- **pkg/cmd/login:** fix typo in example
- **pkg/cmd/serve:** handle missing api key error
- **pkg/cmd/sign:** add '-' char to enforce pipe mode
- **pkg/extractor:** fix typo
- **pkg/extractor/wildcard:** improve max notarized files number with wildcard
- **pkg/extractor/wildcard:** remove max notarization limit
- **pkg/meta:** fix sign attachment flag description
- **pkg/meta:** fix friendly label sign command desc
- **sign:** prompt login if needed when started by explorer
- **sign:** dropping support for multiple keys
- **store:** switch from multi-key to single secret config
- **store:** read alerts
- **terminology:** switch to `notarize` and `authenticate`
- **verify:** minor cmd usage improvements
- **verify:** add explanatory output messages
- **verify:** improve error labels

### Code Refactoring
- asset and metadata fields naming
- change api in order to make atomic indexing of notarized items in lc
- docker compose
- sign, dashboard cmds and APIs
- functional options for signing method
- improve secret handling API
- untrust and unsupport cmds
- list cmd
- use new smart contract functions
- gofmt and golint
- vcn lc mode uses as default env var for apikey
- reduce artifact API surface
- login cmd
- package oriented design
- **api:** always check for user existence when login
- **cmd:** simplify cobra commands naming
- **cmd/sign:** move alert handling to another function
- **meta:** rationalize config functions
- **pkg/extractor:** handle array of artifacts instead of a single one in extractor
- **sign:** printing

### Features
- enable git repo notarization and authentication
- add --force flag to silently overwrite files
- add attach-full-download flag to download all attachments for a label
- add notarize and verify friendly labels
- provide Dockerfile
- upgrade to latest assetsrelay sc
- add ci-attr flag. It permits to save ci context env as metadata
- add --lc-ledger flag for authenticate and inspect commands when used with cross-ledgers API keys
- when apikey is revoked authenticate exits with revokedApikey status - 4
- vcn handle new apikey format{friendlySignerID}.{secret} CNCL
- revokation infos on signed artifact are returned
- logout command
- vcn handles multi attachment to the notarized item. closes [#133](https://github.com/codenotary/cas/issues/133)
- enable file upload
- add multi files support in pipe
- add support for ledger compliance env vars. Closes [#628](https://github.com/codenotary/cas/issues/628)
- add verbose flag to add extra data to the UI
- config file and multi-keys support
- add vcn exit codes
- add manifest management for dir artifacts
- profile migration from v0.3.x
- key flag for sign and verify
- add metadatas in lc sign
- add lc sdk file lock
- add no-tls mode
- add support for tls insecure connections
- align to new immudb implementation
- vcn accepts a certificate to setup a tls connection
- if no limit or filter are submitted inspect return last 100 items
- add multiple notarizations capability with wildcard
- add wildcard extractor
- add serve lc inspect
- subscription limit enforcement
- add verify filter by LC signerID. fix [#104](https://github.com/codenotary/cas/issues/104)
- add vcn lc login. fix [#103](https://github.com/codenotary/cas/issues/103)
- URI parser for assets
- add inspect for ledger compliance
- modularize metadata extractors
- add lc operations for vcn server
- add verify asset on ledger compliance
- add ledger compliance authentication and notarization
- new kind and mimeType attributes for assets
- add otp_empty environment var
- partial parallel notarizations support
- Add otp support. close [#87](https://github.com/codenotary/cas/issues/87) close [#88](https://github.com/codenotary/cas/issues/88)
- user defined attributes for assets
- improve mime and metadata extractor for executables
- added --alert-email option
- store TX hash
- add --hash for verify
- get metadata from docker image
- infer asset version
- json output format (--output=json)
- static build and docker cli support in dockerfile
- improve printing and add yaml output format
- add test env
- profile dir per stage
- json and yaml output for sign
- promote --output to global flag and return formatted errors accordingly
- vcn user agent
- add ARM build
- autogenerated changelog
- Add builds for linux/s390x and linux/ppc64le architectures It is required to refresh the techknowlogick/xgo docker images before this can work as a specific patch has been applied for s390x and ppc64le support
- directory signing and verify with manifest
- add automatic checking for newer versions
- automatically get secret from platform at login
- add trial expiration message
- improve inspect result with timestamp for each element
- **api:** platform alerting system APIs
- **api:** added ability to upload the secret to the platform
- **api:** added get alert API
- **api:** add blockchain organisation support
- **api:** public AlertConfig struct
- **api:** add `SignerID()` method
- **api:** add functional options for signing
- **api:** list assets grouped by hash with pagination
- **build:** build system with Makefile
- **build:** improved build system
- **bundle:** diff
- **bundle:** enforce manifest specs
- **bundle:** arrange multiple items in a Merkle Directed Acyclic Graph
- **cli:** improved column printing
- **cmd:** added --silent flag
- **cmd:** augmenting alerts metadata (full path and diff)
- **cmd:** use stderr for errors and show warning instead when notarization is committed
- **cmd:** added `vcn info` command
- **cmd:** add `vcn serve` command
- **cmd:** added command for setting a custom notarization password (closes [#53](https://github.com/codenotary/cas/issues/53))
- **cmd:** alerts list command
- **cmd:** improved help messages
- **cmd/inspect:** add ability to print only locally extracted info
- **cmd/serve:** allow to pass credentials via HTTP
- **cmd/serve:** TLS support (closes [#48](https://github.com/codenotary/cas/issues/48))
- **cmd/serve:** print environment info at startup
- **cmd/sign:** added --read-only flag for dir notarization
- **cmd/sign:** add `--create-alert` for notarization
- **cmd/sign:** add hostname to alert's metadata by default
- **cmd/verify:** added --alerts to authenticate configured alerts
- **cmd/verify:** add `--alert-config` for authentication
- **extractor:** add scheme for git repository
- **extractor:** add functional options
- **extractor:** scheme for directories
- **extractor:** add support for podman
- **extractor/dir:** .vcnignore file support
- **extractor/dir:** default ignore file
- **inspect:** add new `vcn inspect` command
- **list:** support for --output json/yaml
- **pkg:** move exit-code flag to authenticate command
- **pkg/api:** lc artifact contains uid
- **pkg/cmd:** add lc info and list method
- **pkg/cmd:** authenticate command can download all linked attachments
- **pkg/cmd:** improve lc login method in order to check if connection and api key are valid
- **pkg/cmd:** add inspect api
- **pkg/cmd/inspect:** improve inspect in lc
- **pkg/cmd/sign:** notarization accepts pipe. closes [#579](https://github.com/codenotary/cas/issues/579)
- **pkg/cmd/verify:** verify accepts uid to specify an artifact in time
- **pkg/cmd/verify:** manage verify by hash in lc mode
- **sign:** signing by --hash
- **sign:** `VCN_KEY` env variable for signing
- **sign:** untrust/unsupport by --hash flag
- **sign:** alert creation
- **store:** alert configuration storage
- **store:** manifests centralized storage
- **verify:** with multiple key by VCN_VERIFY_KEYS env var
- **verify:** allow multiple --key flags
- **verify:** show asset visibility
- **verify:** add --org for verify
- **verify/dir:** automatically check manifest integrity


<a name="v0.9.9"></a>
## [v0.9.9] - 2021-09-24

<a name="v0.9.10"></a>
## [v0.9.10] - 2021-09-24

<a name="v0.9.8"></a>
## [v0.9.8] - 2021-09-24

<a name="v0.9.7"></a>
## [v0.9.7] - 2021-06-30
### Bug Fixes
- attachments internal map is iterated with a side sorted string slice. Add some cleanup
- **pkg/cmd/verify:** fix attachments download with multiple version
- **pkg/cmd/verify:** attachments can be used only on CNIL context. On blockchain error is returned
- **pkg/cmd/verify:** fix attachments printing when artifact is not notarized

### Changes
- fix attachments download and add friendly print of attachments when downloading
- attach-full-download flag is removed and full download is default behaviour
- **pkg/api:** fix uid friendly label extraction from key
- **pkg/meta:** fix sign attachment flag description
- **pkg/meta:** fix friendly label sign command desc

### Features
- add verbose flag to add extra data to the UI
- add --force flag to silently overwrite files
- add attach-full-download flag to download all attachments for a label
- add notarize and verify friendly labels


<a name="v0.9.6"></a>
## [v0.9.6] - 2021-06-01
### Bug Fixes
- Make test file content unique to avoid parallel notarization
- **cmd/sign:** fix missing alert sign nil pointer error

### Changes
- rename CNLC to CNIL
- fix badge url
- upgrade cnil sdk
- upgrade cnil sdk in order to fix file lock issue
- remove circleCI
- **pkg/api:** handle parallel requests with same hash on blockchain


<a name="v0.9.5"></a>
## [v0.9.5] - 2021-05-07
### Bug Fixes
- unsuccessful status code
- fix return error codes on cn.io
- **pkg:** add version metadata header and fix signer ID extraction from API key to support email addresses as signer IDs
- **pkg/api:** handle JSON marshal error when creating artifact
- **pkg/api/verify:** fix pkg api verify
- **pkg/cicontext:** remove sensitive data
- **pkg/cmd/login:** fix env var on login
- **pkg/cmd/sign:** improve attachment container structure

### Changes
- vcn context ux improvements. closes [#134](https://github.com/codenotary/cas/issues/134)
- fix missing apikey error
- remove unused import
- add artifact print date method
- add proper error handling on notFound or notVerified
- add lcLedger and signerID parameters on PublicCNLCVerify
- expose lc verify api and explicit environment vars
- increase pipeline timeout
- update ledger-compliance-go dependency to fix missing metadata bug
- remove 'blockchain' from login message
- **cmd/verify:** if api-key is revoked status will be shown as revoked in authenticate
- **pkg/cicontext:** add context probes and minor improvements
- **pkg/cmd/sign:** add '-' char to enforce pipe mode

### Features
- add support for ledger compliance env vars. Closes [#628](https://github.com/codenotary/cas/issues/628)
- add ci-attr flag. It permits to save ci context env as metadata
- add --lc-ledger flag for authenticate and inspect commands when used with cross-ledgers API keys
- when apikey is revoked authenticate exits with revokedApikey status - 4
- vcn handle new apikey format{friendlySignerID}.{secret} CNCL
- revokation infos on signed artifact are returned
- vcn handles multi attachment to the notarized item. closes [#133](https://github.com/codenotary/cas/issues/133)
- enable file upload
- add multi files support in pipe
- **pkg/api:** lc artifact contains uid
- **pkg/cmd:** authenticate command can download all linked attachments
- **pkg/cmd/sign:** notarization accepts pipe. closes [#579](https://github.com/codenotary/cas/issues/579)
- **pkg/cmd/verify:** verify accepts uid to specify an artifact in time


<a name="v0.9.4"></a>
## [v0.9.4] - 2021-02-17
### Bug Fixes
- apikey is removed from userlist when logout
- vcn docker notarization readme.md section
- restore info when no user is logged
- **pkg/api:** return vcn lc apikey error missing only in lc context
- **pkg/cmd:** version check is enabled only for codenotary.io context
- **pkg/cmd/sign:** fix cnlc hash signing
- **pkg/store:** fix api keys clean on load config

### Changes
- update immudb dependency to latest 0.9.2 version and ledger-compliance-go to latest master
- bump sdk version, it fixes notarization on empty cnlc ledger
- improve logout message
- no action on logout if there is no user logged
- authenticate return 2 exit code if hash is not found
- if --lc-no-tls is provided default port is 80
- new no api key configured error message
- renaming flag config in vcnpath and set default vcnpath in os.tmp folder
- bump lc sdk
- bump lc sdk
- verify if config is present in context check
- clean configuration file by removing api keys
- user saving removed on lc login
- disable version check in lc context

### Code Refactoring
- vcn lc mode uses as default env var for apikey

### Features
- add vcn exit codes
- **pkg:** move exit-code flag to authenticate command


<a name="v0.9.3"></a>
## [v0.9.3] - 2021-02-01
### Changes
- adjust colors settings with fgreen
- adjust colors settings
- handle blue screen colors
- minor changes on releasing guide, makefile and minor dependencies updates

### Features
- add manifest management for dir artifacts


<a name="v0.9.2"></a>
## [v0.9.2] - 2021-01-28
### Bug Fixes
- fix autodiscover dir artifact
- **pkg/cmd/internal/cli:** return an error if empty api key is provided

### Changes
- update lc-sdk dependency
- list of signers is not supported in CNLC
- wildcard instructions in their own section
- add a separated CNLC jumpstart
- fix login help message

### Features
- **pkg/cmd:** add lc info and list method
- **pkg/cmd:** improve lc login method in order to check if connection and api key are valid


<a name="v0.9.1"></a>
## [v0.9.1] - 2021-01-18
### Bug Fixes
- load artifact wait for immudb set to be indexed
- fix state service and verified methods
- **pkg:** fix vcn blockchain login
- **pkg/cmd/serve:** improving error if no content is submitted while signing
- **pkg/cmd/sign:** sign error handler fix
- **pkg/cmd/verify:** load artifact error handler fix

### Changes
- bump lc sdk version
- bump lc sdk version
- bumb lc sdk version
- restore lc inspect capabilities
- enable InsecureSkipVerify when tls-skip-verify is active
- disable verified set in lc artifact creation
- **pkg/api:** artifact load uses getAt
- **pkg/cmd/login:** fix typo in example

### Features
- add metadatas in lc sign
- add lc sdk file lock
- add no-tls mode
- add support for tls insecure connections
- align to new immudb implementation


<a name="v0.9.0"></a>
## [v0.9.0] - 2020-12-14
### Bug Fixes
- restore verification with the highest level available when not logged
- fix integration tests and add circleci job for go 1.15
- fix dist release go1.15.6
- fix ledger-compliance-go dependency
- **pkg/api:** swap signerID position with the hash in the vcn lc key format
- **pkg/api:** fix artifact status json unmarshall
- **pkg/cmd:** don't track untrusted verification
- **pkg/cmd/inspect:** fix inspect ux, examples and errors
- **pkg/cmd/verify:** fix empty signerid handling
- **pkg/meta:** add apikey env var

### Changes
- bump sdk version, remove useless package
- bump lc sdk version
- add swagger.json
- bumb version of go-ethereum and smart contracts
- fix inspect full history message
- upgrade go version to 1.15 and sdk
- remove all whitespaces from otp
- **pkg/cmd:** show status unknown if tamper is detected
- **pkg/cmd/serve:** handle missing api key error
- **pkg/extractor:** fix typo
- **pkg/extractor/wildcard:** remove max notarization limit
- **pkg/extractor/wildcard:** improve max notarized files number with wildcard

### Code Refactoring
- change api in order to make atomic indexing of notarized items in lc
- **pkg/extractor:** handle array of artifacts instead of a single one in extractor

### Features
- vcn accepts a certificate to setup a tls connection
- if no limit or filter are submitted inspect return last 100 items
- add multiple notarizations capability with wildcard
- add wildcard extractor
- add serve lc inspect
- improve inspect result with timestamp for each element
- add verify filter by LC signerID. fix [#104](https://github.com/codenotary/cas/issues/104)
- add vcn lc login. fix [#103](https://github.com/codenotary/cas/issues/103)
- partial parallel notarizations support
- add inspect for ledger compliance
- add otp_empty environment var
- add lc operations for vcn server
- add verify asset on ledger compliance
- add ledger compliance authentication and notarization
- Add otp support. close [#87](https://github.com/codenotary/cas/issues/87) close [#88](https://github.com/codenotary/cas/issues/88)
- **pkg/cmd:** add inspect api
- **pkg/cmd/inspect:** improve inspect in lc
- **pkg/cmd/verify:** manage verify by hash in lc mode


<a name="v0.8.3"></a>
## [v0.8.3] - 2020-02-19
### Changes
- distinct event for alert verification
- **cmd/sign:** allow notarization when no dir with write access by default
- **cmd/verify:** do not show "diff unavailable" when no manifest is found
- **extractor/dir:** option to skip ignore file creation


<a name="v0.8.2"></a>
## [v0.8.2] - 2020-02-11
### Changes
- **cmd/verify:** manifest lookup from storage, fallback to target dir
- **meta:** add StaticBool()

### Features
- **cmd/sign:** added --read-only flag for dir notarization
- **store:** manifests centralized storage


<a name="v0.8.1"></a>
## [v0.8.1] - 2020-02-03
### Bug Fixes
- **cmd/sign:** cleanup unused arg to make --hash work again


<a name="v0.8.0"></a>
## [v0.8.0] - 2020-01-27
### Bug Fixes
- Pull latest xgo (go cross compilation tools) when creating the full stack
- **api:** correct alert error message
- **api:** correct alert config JSON name
- **cmd/verify:** show error when alert fails

### Changes
- improve "your assets will not be uploaded" notice
- update copyright year
- **api:** add name to alert struct
- **api:** improve alert error messages
- **api:** more fields for Error struct
- **cmd/internal/cli:** reader and writer for YAML files
- **cmd/verify:** trigger alert on asset error, if any
- **store:** read alerts

### Code Refactoring
- **cmd/sign:** move alert handling to another function

### Features
- added --alert-email option
- add ARM build
- **api:** added get alert API
- **api:** public AlertConfig struct
- **api:** platform alerting system APIs
- **cmd:** augmenting alerts metadata (full path and diff)
- **cmd:** alerts list command
- **cmd:** use stderr for errors and show warning instead when notarization is committed
- **cmd/sign:** add `--create-alert` for notarization
- **cmd/sign:** add hostname to alert's metadata by default
- **cmd/verify:** add `--alert-config` for authentication
- **cmd/verify:** added --alerts to authenticate configured alerts
- **sign:** alert creation
- **store:** alert configuration storage


<a name="v0.7.4"></a>
## [v0.7.4] - 2019-12-19
### Bug Fixes
- **build/makensis:** correct debian image version for NSIS
- **cmd:** correct flags collision
- **cmd/sign:** do not track events when errored

### Changes
- **build:** enforce xgo image digest

### Features
- Add builds for linux/s390x and linux/ppc64le architectures It is required to refresh the techknowlogick/xgo docker images before this can work as a specific patch has been applied for s390x and ppc64le support
- **cmd:** added --silent flag


<a name="v0.7.3"></a>
## [v0.7.3] - 2019-10-30
### Bug Fixes
- **cmd:** enforce lowercase hash string (fixes [#67](https://github.com/codenotary/cas/issues/67))
- **cmd/inspect:** enforce lowercase hash string
- **cmd/serve:** enforce lowercase hash string
- **cmd/sign:** try again when notarization password is wrong

### Changes
- **cmd:** no echo when asking for confirmation before quitting


<a name="v0.7.2"></a>
## [v0.7.2] - 2019-10-24
### Bug Fixes
- **cmd/serve:** add missing event tracking

### Changes
- **ci:** add verify-commit action
- **ci:** github action for docs
- **cmd/sign:** return full output when notarizing

### Features
- store TX hash
- **cmd/inspect:** add ability to print only locally extracted info


<a name="v0.7.1"></a>
## [v0.7.1] - 2019-10-08
### Bug Fixes
- added ability to do not write on read only dir (fixes [#45](https://github.com/codenotary/cas/issues/45))
- **cmd:** typo and language improvement
- **cmd/serve:** CORS preflight (fixes [#60](https://github.com/codenotary/cas/issues/60))


<a name="v0.7.0"></a>
## [v0.7.0] - 2019-10-04
### Bug Fixes
- simplied error message when wrong notarization password is provided
- correct cross compile deps
- trust level was not shown when equals to DISABLED
- contentType was not shown when `ArtifactResponse` was not present
- make fields name consistent across output formats (json, yaml)
- **api:** correct tracing log label for `VerifyMatchingSignerIDWithFallback`
- **cmd/sign:** command must accept exactly one arg

### Changes
- update messages
- update dependencies
- platform managed secret is not stored locally anymore
- remove onboarding message (closes [#52](https://github.com/codenotary/cas/issues/52))
- use go v1.13 for the building system
- promote artifact fields to root
- move tracking call outside api pkg
- **api:** removed `BlockchainVerification.LevelName()` method
- **api:** add `BlockchainVerification.UnmarshalYAML`
- **api:** omit empy field in `ArtifactResponse`
- **cmd:** gracefully print error message if artifact is nil (fixes [#57](https://github.com/codenotary/cas/issues/57))
- **cmd:** added ability to use empty notarization password
- **cmd/info:** improved message for expired tokens
- **cmd/serve:** notarization API will allow only registered schemes for the kind field (closes [#51](https://github.com/codenotary/cas/issues/51))
- **extractor:** do not store version in metadata when empty
- **extractor/dir:** do not create `.vcnignore` on directory authentication (refs [#45](https://github.com/codenotary/cas/issues/45))
- **meta:** use Stringer interface for stage environment
- **meta:** use Stringer interface for status, level, and visibility

### Code Refactoring
- functional options for signing method
- improve secret handling API
- **api:** always check for user existence when login
- **cmd:** simplify cobra commands naming

### Features
- **api:** add `SignerID()` method
- **api:** added ability to upload the secret to the platform
- **api:** add functional options for signing
- **cmd:** improved help messages
- **cmd:** added `vcn info` command
- **cmd:** add `vcn serve` command
- **cmd:** added command for setting a custom notarization password (closes [#53](https://github.com/codenotary/cas/issues/53))
- **cmd/serve:** TLS support (closes [#48](https://github.com/codenotary/cas/issues/48))
- **cmd/serve:** print environment info at startup
- **cmd/serve:** allow to pass credentials via HTTP
- **extractor:** add functional options

### BREAKING CHANGE

`api.Sign` method signature has been changed.

`meta.StageName` has been removed; please use the `.String()` on value instead.

All command's factory methods have been renamed.

`BlockchainVerification.LevelName()` has been removed, please use `BlockchainVerification.Level.String()` instead.

`meta.StatusName`, `meta.LevelName`, `meta.VisibilityName` have been removed; please use the `.String()` on values instead.

`artifact` field is not present anymore in json/yaml results.

Extractors method signature have changed to accomodate functional options.

Secret storage is not used anymore for platform managed secrets.
`store.PublicAddress()` has been renamed to `store.SignerIDFromSecret()`, and will return a value only when a local secret is available


<a name="v0.6.3"></a>
## [v0.6.3] - 2019-09-25
### Bug Fixes
- **extractor/dir:** runtime error when passing a regular file instead of a directory (fixes [#56](https://github.com/codenotary/cas/issues/56))


<a name="v0.6.2"></a>
## [v0.6.2] - 2019-09-17
### Changes
- new AssetsRelay smart contract

### Code Refactoring
- use new smart contract functions


<a name="v0.6.1"></a>
## [v0.6.1] - 2019-09-10
### Bug Fixes
- correct error msg when entering empty email (fixes [#43](https://github.com/codenotary/cas/issues/43))

### Features
- enable git repo notarization and authentication
- **extractor:** add scheme for git repository


<a name="v0.6.0"></a>
## [v0.6.0] - 2019-09-04
### Bug Fixes
- correct `Signer` field semantic
- **api:** correct secret download content format
- **verify:** switch to single key for current user

### Changes
- update Smart Contract addresses to latest version
- removed profile migration from v0.3.x
- deprecate `KEYSTORE_PASSWORD` env var in favour of `VCN_NOTARIZATION_PASSWORD`
- add "your assets will not be uploaded" message
- deprecate `--key` in favor of `--signerID`
- improve API and user messages to reflect the new terminology
- correct typos
- **api:** refactor to single address and code cleanup
- **help:** add info about assets and env vars
- **meta:** removed event tracking for keystore creation
- **sign:** dropping support for multiple keys
- **store:** switch from multi-key to single secret config
- **terminology:** switch to `notarize` and `authenticate`

### Code Refactoring
- gofmt and golint
- reduce artifact API surface
- **meta:** rationalize config functions

### Features
- add trial expiration message
- automatically get secret from platform at login
- add automatic checking for newer versions

### BREAKING CHANGE

Env variable `KEYSTORE_PASSWORD` has been renamed to `VCN_NOTARIZATION_PASSWORD`.

All `api.BlockChainVerify*()` funcs have been renamed to `Verify*()`

`api.BlockchainVerification.Key()` has been renamed to `SignerID()`
`api.BlockChainOrganisation.MembersKeys()` has been renamed to `MembersIDs()`
`api.BlockChainGetOrganisation()` has been renamed to `GetBlockChainOrganisation()`

`meta`'s endpoint funcs has been removed in favour of single `meta.APIEndpoint()`. `meta.MainNetEndpoint()` has been renamed to `meta.MainNet()`.

Artifact field `Signer` will not contain the user pub key anymore.

`api.LoadArtifactForHash` has been renamed to `api.LoadArtifact`.
`api.ArtifactRequest` has been privatized. Artifact shall be created only thru the notarization process.

`--key` has been removed from `vcn sign`

Dropping `LoadAllArtifacts` and `LoadArtifacts` and wallets/keys/keystore related functions from `api.User`.
`api.Sign` won't accept anymore a pub key as parameter.

support for multiple keystores in config file and related APIs within the `store` have been removed.


<a name="v0.5.4"></a>
## [v0.5.4] - 2019-07-30
### Features
- **extractor:** add support for podman


<a name="v0.5.3"></a>
## [v0.5.3] - 2019-07-19
### Bug Fixes
- **build:** correct NSIS installation directory


<a name="v0.5.2"></a>
## [v0.5.2] - 2019-07-17
### Bug Fixes
- update go.sum for xgo
- **bundle:** use uint64 for size
- **dir:** OS agnostic paths
- **list:** show all assets with pagination (fixes [#28](https://github.com/codenotary/cas/issues/28))

### Changes
- **api:** add Artifact.Copy()
- **bundle:** return error if distinct sizes are found for the same digest
- **bundle:** descriptor test
- **bundle:** rework diff
- **cmd:** always show local extracted metadata
- **dir/extractor:** ignore irregular files

### Features
- directory signing and verify with manifest
- **api:** list assets grouped by hash with pagination
- **bundle:** enforce manifest specs
- **bundle:** arrange multiple items in a Merkle Directed Acyclic Graph
- **bundle:** diff
- **extractor:** scheme for directories
- **extractor/dir:** default ignore file
- **extractor/dir:** .vcnignore file support
- **verify/dir:** automatically check manifest integrity


<a name="v0.5.1"></a>
## [v0.5.1] - 2019-07-02
### Bug Fixes
- `BlockchainVerification` json unmarshalling

### Changes
- **cmd:** update cobra to 0.5.0
- **verify:** add explanatory output messages
- **verify:** minor cmd usage improvements

### Features
- **inspect:** add new `vcn inspect` command
- **sign:** `VCN_KEY` env variable for signing


<a name="0.5.0"></a>
## [0.5.0] - 2019-06-25

<a name="v0.5.0"></a>
## [v0.5.0] - 2019-06-25
### Bug Fixes
- ca-certificates for Dockerfile
- clean other context when user logs in
- **cmd:** show config file path if not default one
- **docker:** accept images but not other objects
- **extractor:** correct empty files handling

### Changes
- reduce public api surface
- code cleanup and fix comments
- bump version
- **api:** move type Error to its own source file
- **api:** deprecate public publisher APIs
- **api:** deprecate publisher fields
- **api:** allow empty size
- **cmd:** internal types
- **env:** remove VERIFY_ prefix
- **internal:** move cli package
- **output:** silence printing when formatted output
- **verify:** improve error labels

### Features
- autogenerated changelog
- vcn user agent
- promote --output to global flag and return formatted errors accordingly
- json and yaml output for sign
- **api:** add blockchain organisation support
- **list:** support for --output json/yaml
- **sign:** signing by --hash
- **verify:** add --org for verify

### BREAKING CHANGE

some methods are now private, tracking functions are now deprecated.

Some publisher and auth related methods and structs have been privatized and will removed in future

ArtifactResponse fields (related to the publisher) as been renamed as following:
```
Publisher -> Signer
PublisherCompany -> Company
PublisherWebsiteUrl -> Website
PublisherCount -> SignerCount
```
ArtifactRequest's `Url` has been fixed to `URL` too.

meta.VcnClientName() has been removed in favor of meta.UserAgent()

`VCN_VERIFY_KEYS` has been removed in favour of `VCN_KEY` with the same functionality.


<a name="0.4.3"></a>
## [0.4.3] - 2019-05-21
### Bug Fixes
- **api:** nil pointer dereference
- **extractor:** nil pointer dereference
- **sign:** flag accessed but not defined: hash
- **sign:** cross platform loading spinner

### Changes
- minor refactoring of login/logout
- generate markdown docs for commands
- **api:** load user's artifact by hash only
- **api:** size must be valid, name can be empty
- **sign:** prompt login if needed when started by explorer

### Features
- profile dir per stage
- add test env
- **sign:** untrust/unsupport by --hash flag


<a name="0.4.2"></a>
## [0.4.2] - 2019-05-14
### Bug Fixes
- print correct config file name
- do not read in env var for config
- **sign:** remove ownership disclaimer
- **verify:** cross platform coloured printing
- **verify:** only print size when available

### Changes
- disable config file message
- reorg publisher fields
- minor printing and marshaling improvements
- **api:** fine-tune blockchain verify funcs
- **build:** make makefile more resilient
- **log:** improve metahash tracing
- **output:** correct WriteTo interface

### Code Refactoring
- **sign:** printing

### Features
- improve printing and add yaml output format
- static build and docker cli support in dockerfile
- **verify:** with multiple key by VCN_VERIFY_KEYS env var
- **verify:** allow multiple --key flags
- **verify:** show asset visibility


<a name="0.4.1"></a>
## [0.4.1] - 2019-05-08
### Bug Fixes
- correct err msg when token has expired
- do not print verification error message if --output=json
- verify fallback
- clear context when logging out
- if logged in always show user own signature
- remove test-resources from dockerignore
- return proper error by checking quota before

### Changes
- executable file sniffer
- go mod tidy
- change metadata naming
- add Platform and Architecture fields
- **build:** minor improvements
- **ci:** initial circleci setup
- **cmd:** use extractors
- **meta:** increse signing timeout
- **meta:** increase tx verification rounds

### Code Refactoring
- asset and metadata fields naming

### Features
- user defined attributes for assets
- json output format (--output=json)
- infer asset version
- get metadata from docker image
- add --hash for verify
- improve mime and metadata extractor for executables
- new kind and mimeType attributes for assets
- modularize metadata extractors
- URI parser for assets
- subscription limit enforcement
- **build:** improved build system


<a name="0.4.0"></a>
## [0.4.0] - 2019-04-30
### Bug Fixes
- dashboard cmd has no args
- prefix for docker assets
- find home dir in the right way
- do not exit early when login is required
- do not quit when executed by Win context menu
- ensure default keystore before looking for keys
- remove log.Fatal in favour of returing errors
- gitignore
- ask for confirmation before quitting
- token deletion when already logged out
- improve err msgs when auth is required
- add missing header comments
- temporary dirs creation at startup
- **login:** silence usage message when errored
- **migrate:** do not create a new key if it already exists
- **store:** correct key funcs and tests
- **tracking:** send events at right place in time

### Changes
- use two-stage build
- make createArtifact private
- printing stuff
- no race test for make install
- add dev stage
- move cli package to internal
- improve naming
- code cleanup and minor improvements
- add store context and key utils
- introduce store package
- move logs to internal
- improve wallet not synced message
- correct printing func
- update gitignore
- switch main to cobra
- partial api logging refactor
- code clean up
- introducing Cobra
- move file hashing func to internal
- use make install for Dockerfile
- **build:** improve makefile
- **cli:** spinner changed and code cleaned up
- **dist:** omit symbol table and debug info
- **internal:** correct errors funcs and tests

### Code Refactoring
- login cmd
- list cmd
- untrust and unsupport cmds
- sign, dashboard cmds and APIs
- docker compose
- package oriented design

### Features
- key flag for sign and verify
- profile migration from v0.3.x
- config file and multi-keys support
- logout command
- upgrade to latest assetsrelay sc
- provide Dockerfile
- **build:** build system with Makefile
- **cli:** improved column printing

### BREAKING CHANGE

this commit changes the usage of verify and sign methods

this commit introduce the config file with multi-keys support, and a huge refactoring of vcn code.


<a name="0.3.6"></a>
## [0.3.6] - 2019-04-08

<a name="0.3.5"></a>
## [0.3.5] - 2019-03-28

<a name="0.3.4"></a>
## [0.3.4] - 2019-03-20

<a name="0.3.3"></a>
## [0.3.3] - 2019-03-13

<a name="0.3.2"></a>
## [0.3.2] - 2019-03-11

<a name="0.3.1"></a>
## [0.3.1] - 2019-03-11

<a name="0.3.0"></a>
## [0.3.0] - 2019-03-08

<a name="0.2.2"></a>
## [0.2.2] - 2019-03-07

<a name="v.0.2-beta.0"></a>
## [v.0.2-beta.0] - 2019-02-25

<a name="v.0.1-beta.2"></a>
## v.0.1-beta.2 - 2019-02-19

[Unreleased]: https://github.com/codenotary/cas/compare/v0.10.0-RC-1...HEAD
[v0.10.0-RC-1]: https://github.com/codenotary/cas/compare/v0.9.9...v0.10.0-RC-1
[v0.9.9]: https://github.com/codenotary/cas/compare/v0.9.10...v0.9.9
[v0.9.10]: https://github.com/codenotary/cas/compare/v0.9.8...v0.9.10
[v0.9.8]: https://github.com/codenotary/cas/compare/v0.9.7...v0.9.8
[v0.9.7]: https://github.com/codenotary/cas/compare/v0.9.6...v0.9.7
[v0.9.6]: https://github.com/codenotary/cas/compare/v0.9.5...v0.9.6
[v0.9.5]: https://github.com/codenotary/cas/compare/v0.9.4...v0.9.5
[v0.9.4]: https://github.com/codenotary/cas/compare/v0.9.3...v0.9.4
[v0.9.3]: https://github.com/codenotary/cas/compare/v0.9.2...v0.9.3
[v0.9.2]: https://github.com/codenotary/cas/compare/v0.9.1...v0.9.2
[v0.9.1]: https://github.com/codenotary/cas/compare/v0.9.0...v0.9.1
[v0.9.0]: https://github.com/codenotary/cas/compare/v0.8.3...v0.9.0
[v0.8.3]: https://github.com/codenotary/cas/compare/v0.8.2...v0.8.3
[v0.8.2]: https://github.com/codenotary/cas/compare/v0.8.1...v0.8.2
[v0.8.1]: https://github.com/codenotary/cas/compare/v0.8.0...v0.8.1
[v0.8.0]: https://github.com/codenotary/cas/compare/v0.7.4...v0.8.0
[v0.7.4]: https://github.com/codenotary/cas/compare/v0.7.3...v0.7.4
[v0.7.3]: https://github.com/codenotary/cas/compare/v0.7.2...v0.7.3
[v0.7.2]: https://github.com/codenotary/cas/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/codenotary/cas/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/codenotary/cas/compare/v0.6.3...v0.7.0
[v0.6.3]: https://github.com/codenotary/cas/compare/v0.6.2...v0.6.3
[v0.6.2]: https://github.com/codenotary/cas/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/codenotary/cas/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/codenotary/cas/compare/v0.5.4...v0.6.0
[v0.5.4]: https://github.com/codenotary/cas/compare/v0.5.3...v0.5.4
[v0.5.3]: https://github.com/codenotary/cas/compare/v0.5.2...v0.5.3
[v0.5.2]: https://github.com/codenotary/cas/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/codenotary/cas/compare/0.5.0...v0.5.1
[0.5.0]: https://github.com/codenotary/cas/compare/v0.5.0...0.5.0
[v0.5.0]: https://github.com/codenotary/cas/compare/0.4.3...v0.5.0
[0.4.3]: https://github.com/codenotary/cas/compare/0.4.2...0.4.3
[0.4.2]: https://github.com/codenotary/cas/compare/0.4.1...0.4.2
[0.4.1]: https://github.com/codenotary/cas/compare/0.4.0...0.4.1
[0.4.0]: https://github.com/codenotary/cas/compare/0.3.6...0.4.0
[0.3.6]: https://github.com/codenotary/cas/compare/0.3.5...0.3.6
[0.3.5]: https://github.com/codenotary/cas/compare/0.3.4...0.3.5
[0.3.4]: https://github.com/codenotary/cas/compare/0.3.3...0.3.4
[0.3.3]: https://github.com/codenotary/cas/compare/0.3.2...0.3.3
[0.3.2]: https://github.com/codenotary/cas/compare/0.3.1...0.3.2
[0.3.1]: https://github.com/codenotary/cas/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/codenotary/cas/compare/0.2.2...0.3.0
[0.2.2]: https://github.com/codenotary/cas/compare/v.0.2-beta.0...0.2.2
[v.0.2-beta.0]: https://github.com/codenotary/cas/compare/v.0.1-beta.2...v.0.2-beta.0

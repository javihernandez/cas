# Environments

By default `cas` will put the config file and secrets within the a directory called `.cas` within your [home directory](https://en.wikipedia.org/wiki/Home_directory) (e.g. `$HOME/.cas` or `%USERPROFILE%\.cas` on Windows).

However, `cas` can work with distinct environments (eg. for testing purpose).

The following environments are supported by setting the `STAGE` environment var:

Stage | Directory | Note
------------ | ------------- | -------------
`STAGE=PRODUCTION` | `.cas` | *default* 
`STAGE=STAGING` | `.cas.staging` |
`STAGE=TEST` | `.cas.test` | *`CAS_TEST_DASHBOARD`, `CAS_TEST_NET`, `CAS_TEST_CONTRACT`, `CAS_TEST_API` must be set accordingly to your test environment*


## Other environment variables

Name | Description | Example 
------------ | ------------- | -------------
`CAS_SIGNERID` | For `cas authenticate` acts as a list of SignerID(s) (separated by space) to authenticate against | `CAS_SIGNERID="0x0...0 0x0...1" cas authenticate <asset>` or `CAS_SIGNERID="0x0...0 <asset>` 
`CAS_ORG` | Organization's ID to authenticate against | `CAS_ORG="vchain.us" cas authenticate <asset>`
`LOG_LEVEL` | Logging verbosity. Accepted values: `TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC`  | `LOG_LEVEL=TRACE cas login` 
`HTTP_PROXY` | HTTP Proxy configuration | `HTTP_PROXY=http://localhost:3128 cas authenticate <asset>`
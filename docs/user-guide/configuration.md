# Configuration

By default, the `cas` command line stores its config file (`config.json`) in temporary directory, for Linux the path is `/tmp/.cas/config.json`, for Windows - `c:\temp\config.json` or `c:\windows\temp\config.json`.
> If the `STAGE` environment variable has been set, the default configuration directory can be different. See [environments](environments.md).

However, you can specify a different location for the config file via the `--caspath` command line option. For example:

```
cas --caspath /path/to/your/config.json
```

<!-- The config file contains paths to keystore directories, and stores credentials of the current authenticated user.

`cas` manages these files and directories and you should not modify them.
However, *you can modify* the config file to control where keys are stored. -->

## Config file

### Example of `config.json`

```
{
  "currentcontext": {
    "LcHost": "cas.codenotary.com",
    "LcPort": "443"
  },
  "schemaversion": 3,
  "users": null
}
```

### Breakdown of `config.json`'s currenctcontext section

The property `currentcontext` holds the connection details for active session. Supported fields:


 - `LcHost` - server's hostname or IP address
 - `LcPort` - server's port
 - `LcCert` - absolute path to a certificate file needed to set up TLS connection
 - `LcSkipTlsVerify` - boolean flag instructing to skip server TLS verification (`false` by default)
 - `LcNoTls` - boolean flag instructing to not to use TLS

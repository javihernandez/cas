# Formatted output (json/yaml)

`cas` can output results in `json` or `yaml` formats by using the [--output global flag](../cmd/cas.md#options).
> Although all commands support `--output`, some could return an empty results (ie. `cas login`, `cas logout`, and `cas dashboard`).

## Examples

```
cas authenticate docker://nginx --output yaml
```

```
cas list --output json
```

```
cas notarize file.txt
```
> You need to set `CAS_NOTARIZATION_PASSWORD` [environment variable](environments.md#other-environment-variables) to make `cas` work in non-interactive mode

## Dealing with errors

When an error is encountered, `cas` will print the usual error message to the *Standard error* but also will return an error object (formatted accordingly to `--output`) to the *Standard output*.

### Example of mixed *Standard error* and *Standard output*
```
$ cas authenticate non-existing.txt --output json
Error: open non-existing.txt: no such file or directory
{
  "error": "open non-existing.txt: no such file or directory"
}
```

### Example of redirecting the *Standard output* to get the formatted result
```
$ cas authenticate non-existing.txt --output json > output.json
Error: open non-existing.txt: no such file or directory

$ cat output.json
{
  "error": "open non-existing.txt: no such file or directory"
}
```
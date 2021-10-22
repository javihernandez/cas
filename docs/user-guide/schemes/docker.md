# Docker Integration

`cas` supports local docker installations out of the box using `docker://` as a location. You just need to point to the correct container image name or the container image id.

If you prefer [podman](https://podman.io/), just use `podman://` instead.


## Notarize a local docker image

`cas` uses docker default schemes, so the latest tag is automatically used, if no tag is given

```
cas notarize docker://hello-world
```

or using a tag

```
cas notarize docker://hello-world:v1
```

To be able to notarize, you need to register at [CodeNotary](https://dashboard.codenotary.io) and get an account.

## Authenticate a local docker image

```
cas authenticate docker://hello-world
```

or using a tag

```
cas authenticate docker://hello-world:v1
```

## Docker Sidecar Integration

`cas` also offers a sidecar project, you can use to automatically authenticate used container images during runtime.

Check out (https://github.com/codenotary/cas-watchdog) on your server. The tool continuously verifies the integrity of your containers:

```
git clone https://github.com/codenotary/cas-watchdog.git
```

Edit the verify file and set the alerting/monitoring tool you are using (see the following instructions), if you want to change the alerting

Make sure `/var/run/docker.sock` is accessible and run the following command on your server within the [cas-watchdog](https://github.com/codenotary/cas-watchdog.git) directory.
```
docker-compose build && docker-compose up
```

To modify the verify file, hook up your alerting tool into the err() function.

Example using Slack, do the following:

* Create a Slack Bot (Slack documentation here)
* Use the following code:

```
function err() {
    echo "Container ${1} (${2}) verification failed" >&2
    curl -q -X POST \
        -H 'Content-type: application/json' \
        --data "{\"text\":\"Container ${1} (${2}) verification failed\"}" \
        "https://hooks.slack.com/services/$TOKEN/$KEY" > /dev/null 2>&1}
```

If all works well, you should receive slack messages in your slack channel

![Slack alert based on cas verify](https://www.vchain.us/wp-content/uploads/2019/04/002_Alerting-on-Slack-example-768x129.png "Slack alert based on cas verify")

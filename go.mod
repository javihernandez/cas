module github.com/codenotary/cas

go 1.13

require (
	github.com/CycloneDX/cyclonedx-go v0.4.0
	github.com/anchore/go-rpmdb v0.0.0-20210602151223-1f0f707a2894
	github.com/blang/semver v3.5.1+incompatible
	github.com/caarlos0/spin v1.1.0
	github.com/codenotary/immudb v1.2.1
	github.com/containerd/containerd v1.5.7 // indirect
	github.com/dghubble/sling v1.3.0
	github.com/docker/docker v20.10.8+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.13.0
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/h2non/filetype v1.0.10
	github.com/mattn/go-colorable v0.1.12
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0
	github.com/package-url/packageurl-go v0.1.0
	github.com/schollz/progressbar/v3 v3.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	github.com/vchain-us/ledger-compliance-go v0.9.3-0.20220118134549-9591b15eb645
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	google.golang.org/grpc v1.43.0
	gopkg.in/src-d/go-git.v4 v4.13.1
)

replace github.com/spf13/afero => github.com/spf13/afero v1.5.1

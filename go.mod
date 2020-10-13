module github.com/iqlusioninc/relayer

go 1.14

require (
	github.com/99designs/keyring v1.1.5 // indirect
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/avast/retry-go v2.6.0+incompatible
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/confio/ics23-iavl v0.6.0 // indirect
	github.com/confio/ics23-tendermint v0.6.1 // indirect
	github.com/confio/ics23/go v0.0.0-20200604202538-6e2c36a74465 // indirect
	github.com/containerd/continuity v0.0.0-20200228182428-0f16d7a0959c // indirect
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200622203133-4716260a6e2d
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/datachainlab/cross v0.0.0-20200507061335-6fded9aa2c2e
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/ory/dockertest/v3 v3.5.5
	github.com/otiai10/copy v1.2.0 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/prometheus/client_golang v1.7.0 // indirect
	github.com/regen-network/cosmos-proto v0.3.0 // indirect
	github.com/sirkon/goproxy v1.4.8
	github.com/sirupsen/logrus v1.5.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.7
	github.com/tendermint/tm-db v0.5.1
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => /Users/zero/learn/cosmos-sdk
	github.com/datachainlab/cross => /Users/zero/learn/cross
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
)

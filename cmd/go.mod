module github.com/situ-vault/situ-vault/cmd

go 1.17

require (
	github.com/situ-vault/situ-vault/pkg v0.0.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/situ-vault/situ-vault/pkg v0.0.0 => ./../pkg

module github.com/situ-vault/situ-vault/cmd

go 1.17

require (
	github.com/situ-vault/situ-vault/pkg v0.0.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/nicksnyder/basen v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20211023085530-d6a326fbbf70 // indirect
)

replace github.com/situ-vault/situ-vault/pkg v0.0.0 => ./../pkg

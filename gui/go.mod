module github.com/polarctos/situ-vault/gui

go 1.16

require (
	fyne.io/fyne v1.4.3
	github.com/polarctos/situ-vault/pkg v0.0.0
	github.com/stretchr/testify v1.6.1
)

replace github.com/polarctos/situ-vault/pkg v0.0.0 => ./../pkg

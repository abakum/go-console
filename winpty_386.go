//go:build windows && 386
// +build windows,386

package console

import (
	"embed"

	"github.com/abakum/embed-encrypt/encryptedfs"
)

// go:embed winpty/386/*
//var winpty_deps embed.FS

//encrypted:embed winpty/386/*
var winpty_deps encryptedfs.FS

//go:embed winpty/386/*.enc
var winpty_depsEnc embed.FS

func init() {
	winpty_deps = encryptedfs.InitFS(winpty_depsEnc, key)
}

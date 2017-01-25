package multiModulePage

import (
	"net/http"
	"path/filepath"

	"github.com/cryptix/go/goutils"
)

var Assets = http.Dir(
	filepath.Join(
		goutils.LocatePackage("github.com/cryptix/exp/multiModulePage"),
		"tmpl"),
)

package static

import (
	"bytes"
	"strings"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata@v3.18.0+incompatible -nometadata -ignore=\.go -pkg=static -o=bindata.go ./...

// Schema reads the .gql schema files from the generated _bindata.go file, concatenating the
// files together into one string.
func Schema() string {
	buf := bytes.Buffer{}

	for _, name := range AssetNames() {
		if strings.Contains(name, ".gql") {
			b := MustAsset(name)
			buf.Write(b)

			// Add a newline if the file does not end in a newline.
			if len(b) > 0 && b[len(b)-1] != '\n' {
				buf.WriteByte('\n')
			}
		}
	}

	return buf.String()
}

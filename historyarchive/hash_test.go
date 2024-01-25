// Copyright 2016 Diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashDecode(t *testing.T) {
	d, err := DecodeHash("f8acbd8c9a901e17ed5488d7ebf44781dc82924457f488b75ba577911c4e2656")
	assert.Nil(t, err)
	assert.Equal(t, []byte{0xf8, 0xac, 0xbd}, d[:3])
}

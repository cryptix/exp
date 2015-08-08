package git

import (
	"os"
	"testing"

	"github.com/cheekybits/is"
)

func TestUnpackObject_blob(t *testing.T) {
	is := is.New(t)
	tcases := []struct {
		Fname  string
		Object *Object
	}{
		// blobs
		{
			Fname:  "tests/blob/1f7a7a472abf3dd9643fd615f6da379c4acb3e3a",
			Object: &Object{"blob", []byte("version 2\n")},
		},
		{
			Fname:  "tests/blob/d670460b4b4aece5915caf5c68d12f560a9fe3e4",
			Object: &Object{"blob", []byte("test content\n")},
		},

		// trees
		{
			Fname:  "tests/tree/3c4e9cd789d88d8d89c1073707c3585e41b0e614",
			Object: &Object{"tree", []byte("test123")},
		},
		{
			Fname:  "tests/tree/0155eb4229851634a0f03eb265b69f5a2d56f341",
			Object: &Object{"tree", []byte("test123")},
		},

		// commit
		{
			Fname:  "tests/commit/3c4e9cd789d88d8d89c1073707c3585e41b0e614",
			Object: &Object{"commit", []byte("test123")},
		},
		// commit
		{
			Fname:  "tests/commit/ad8fdc888c6f6caed63af0fb08484901e4e7e41e",
			Object: &Object{"commit", []byte("test123")},
		},
	}

	for _, tc := range tcases {
		f, err := os.Open(tc.Fname)
		is.Nil(err)
		obj, err := DecodeObject(f)
		is.Nil(err)
		is.Equal(obj, tc.Object)
		is.Nil(f.Close())
	}

	t.Fail()
}

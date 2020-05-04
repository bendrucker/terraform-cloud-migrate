package configwrite

import (
	"testing"

	"github.com/lithammer/dedent"
	"github.com/spf13/afero"
)

func newTestWriter(t *testing.T, path string, setup func(afero.Fs)) *Writer {
	fs := afero.NewMemMapFs()
	setup(fs)
	writer, diags := newWriter(path, fs)
	if len(diags) != 0 {
		for _, diag := range diags {
			t.Error(diag.Error())
		}
		t.FailNow()
	}
	return writer
}

func newTestModule(t *testing.T, files map[string]string) *Writer {
	return newTestWriter(t, "", func(fs afero.Fs) {
		for name, content := range files {
			if err := afero.WriteFile(fs, name, []byte(dedent.Dedent(content)), 0644); err != nil {
				t.Error(err)
			}
		}
	})
}

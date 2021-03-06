package gen

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/logical"
	"gitlab.geniesolutions.com.au/itops/vault-plugin-secrets-oauth2/version"
)

func testBackend(tb testing.TB) (*backend, logical.Storage) {
	tb.Helper()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Factory(context.Background(), config)
	if err != nil {
		tb.Fatal(err)
	}
	return b.(*backend), config.StorageView
}

func TestBackend(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		t.Parallel()

		b, storage := testBackend(t)
		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.ReadOperation,
			Path:      "info",
		})
		if err != nil {
			t.Fatal(err)
		}

		if v, exp := resp.Data["version"].(string), version.Version; v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}

		if v, exp := resp.Data["commit"].(string), version.GitCommit; v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}
	})
}

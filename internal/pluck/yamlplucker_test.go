package pluck_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/pluck"
)

const (
	enclaveSEVYAMLEnclaveName    = "enclave"
	enclaveSEVYAMLEnclaveSnippet = `enclave:
  addr: "http://127.0.0.1:8083"
  addr_tls: "https://127.0.0.1:8444"
  args:
    domain: "bearclave.tee"`
	nonclaveSEVYAMLMeasurementName    = "nonclave.measurement"
	nonclaveSEVYAMLMeasurementSnippet = `measurement: |
  {
    "version": 5,
    "guest_svn": 0,
    "policy": 196608,
    "family_id": "AAAAAAAAAAAAAAAAAAAAAA==",
    "image_id": "AAAAAAAAAAAAAAAAAAAAAA==",
    "vmpl": 0,
    "current_tcb": 16004667175767900164,
    "platform_info": 37,
    "signer_info": 0,
    "measurement": "FBlf3jaFK2nyrBcWbr8yIzUbHDSBTxEDOUEmUoGQc+Bh7XxH5uANAqgXrrSZLOIN",
    "host_data": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
    "id_key_digest": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "author_key_digest": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "report_id": "",
    "report_id_ma": "//////////////////////////////////////////8=",
    "reported_tcb": 16004667175767900164,
    "chip_id": "",
    "committed_tcb": 16004667175767900164,
    "current_build": 0,
    "current_minor": 58,
    "current_major": 1,
    "committed_build": 0,
    "committed_minor": 58,
    "committed_major": 1,
    "launch_tcb": 16004667175767900164,
    "cpuid_1eax_fms": 10489617
  }`
)

//go:embed testdata/enclave-sev.yaml
var enclaveSEVYAML string

//go:embed testdata/nonclave-sev.yaml
var nonclaveSEVYAML string

func TestYAMLPlucker_Pluck(t *testing.T) {
	t.Run("happy path - nonclave sev (scalar node)", func(t *testing.T) {
		// given
		ctx := context.Background()
		code := nonclaveSEVYAML
		name := nonclaveSEVYAMLMeasurementName
		kind := pluck.Node
		want := nonclaveSEVYAMLMeasurementSnippet + "\n"
		plucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		// when
		got, err := plucker.Pluck(ctx, code, name, kind)

		// then
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("happy path - enclave sev (mapping node)", func(t *testing.T) {
		// given
		ctx := context.Background()
		code := enclaveSEVYAML
		name := enclaveSEVYAMLEnclaveName
		kind := pluck.Node
		want := enclaveSEVYAMLEnclaveSnippet + "\n"
		plucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		// when
		got, err := plucker.Pluck(ctx, code, name, kind)

		// then
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("happy path - enclave sev (file)", func(t *testing.T) {
		// given
		ctx := context.Background()
		code := enclaveSEVYAML
		name := enclaveSEVYAMLEnclaveName
		kind := pluck.File
		want := enclaveSEVYAML
		plucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		// when
		got, err := plucker.Pluck(ctx, code, name, kind)

		// then
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}

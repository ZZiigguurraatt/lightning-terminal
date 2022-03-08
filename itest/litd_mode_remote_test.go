package itest

import (
	"context"
	"os"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/stretchr/testify/require"
)

// testModeRemote makes sure that in remote mode all daemons work correctly.
func testModeRemote(net *NetworkHarness, t *harnessTest) {
	ctx := context.Background()

	// Some very basic functionality tests to make sure lnd is working fine
	// in remote mode.
	net.SendCoins(t.t, btcutil.SatoshiPerBitcoin, net.Bob)

	// We expect a non-empty alias (truncated node ID) to be returned.
	resp, err := net.Bob.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	require.NoError(t.t, err)
	require.NotEmpty(t.t, resp.Alias)
	require.Contains(t.t, resp.Alias, "0")

	t.t.Run("certificate check", func(tt *testing.T) {
		cfg := net.Bob.Cfg

		// In remote mode we expect the LiT HTTPS port (8443 by default)
		// and to have its own certificate
		litCerts, err := getServerCertificates(cfg.LitAddr())
		require.NoError(tt, err)
		require.Len(tt, litCerts, 1)
		require.Equal(
			tt, "litd autogenerated cert",
			litCerts[0].Issuer.Organization[0],
		)
	})
	t.t.Run("gRPC macaroon auth check", func(tt *testing.T) {
		cfg := net.Bob.Cfg

		for _, endpoint := range endpoints {
			endpoint := endpoint
			tt.Run(endpoint.name+" lit port", func(ttt *testing.T) {
				if !endpoint.supportsMacAuthOnLitPort {
					return
				}

				runGRPCAuthTest(
					ttt, cfg.LitAddr(), cfg.LitTLSCertPath,
					endpoint.macaroonFn(cfg),
					endpoint.requestFn,
					endpoint.successPattern,
				)
			})
		}
	})

	t.t.Run("UI password auth check", func(tt *testing.T) {
		cfg := net.Bob.Cfg

		for _, endpoint := range endpoints {
			endpoint := endpoint
			tt.Run(endpoint.name+" lit port", func(ttt *testing.T) {
				runUIPasswordCheck(
					ttt, cfg.LitAddr(), cfg.LitTLSCertPath,
					cfg.UIPassword,
					endpoint.requestFn, false,
					!endpoint.supportsUIPasswordOnLitPort,
					endpoint.successPattern,
				)
			})
		}
	})

	t.t.Run("UI index page fallback", func(tt *testing.T) {
		runIndexPageCheck(tt, net.Bob.Cfg.LitAddr())
	})

	t.t.Run("grpc-web auth", func(tt *testing.T) {
		cfg := net.Bob.Cfg

		for _, endpoint := range endpoints {
			endpoint := endpoint
			tt.Run(endpoint.name+" lit port", func(ttt *testing.T) {
				runGRPCWebAuthTest(
					ttt, cfg.LitAddr(), cfg.UIPassword,
					endpoint.grpcWebURI,
				)
			})
		}
	})

	t.t.Run("gRPC super macaroon auth check", func(tt *testing.T) {
		cfg := net.Bob.Cfg

		superMacFile, err := bakeSuperMacaroon(cfg, true)
		require.NoError(tt, err)

		defer func() {
			_ = os.Remove(superMacFile)
		}()

		for _, endpoint := range endpoints {
			endpoint := endpoint
			tt.Run(endpoint.name+" lit port", func(ttt *testing.T) {
				if !endpoint.supportsMacAuthOnLitPort {
					return
				}

				runGRPCAuthTest(
					ttt, cfg.LitAddr(), cfg.LitTLSCertPath,
					superMacFile,
					endpoint.requestFn,
					endpoint.successPattern,
				)
			})
		}
	})

	t.t.Run("REST auth", func(tt *testing.T) {
		cfg := net.Bob.Cfg

		for _, endpoint := range endpoints {
			endpoint := endpoint

			if endpoint.restWebURI == "" {
				continue
			}

			tt.Run(endpoint.name+" lit port", func(ttt *testing.T) {
				runRESTAuthTest(
					ttt, cfg.LitAddr(), cfg.UIPassword,
					endpoint.macaroonFn(cfg),
					endpoint.restWebURI,
					endpoint.successPattern,
				)
			})
		}
	})

	t.t.Run("lnc auth", func(tt *testing.T) {
		cfg := net.Bob.Cfg

		for _, endpoint := range endpoints {
			endpoint := endpoint
			tt.Run(endpoint.name+" lit port", func(ttt *testing.T) {
				runLNCAuthTest(
					ttt, cfg.LitAddr(), cfg.UIPassword,
					cfg.LitTLSCertPath,
					endpoint.requestFn,
					endpoint.successPattern,
					endpoint.allowedThroughLNC,
				)
			})
		}
	})
}
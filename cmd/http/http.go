package http

import (
	"context"
	"github.com/ProtocolONE/go-core/entrypoint"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/grpc"
	"github.com/paysuper/paysuper-management-api/cmd"
	"github.com/paysuper/paysuper-management-api/internal/daemon"
	"github.com/paysuper/paysuper-management-api/pkg/http"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:           "http",
		Short:         "HTTP API daemon",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			var (
				s *http.HTTP
				c func()
				e error
			)
			defer func() {
				if c != nil {
					c()
				}
			}()
			cmd.Slave.Executor(func(ctx context.Context) error {
				initial, _ := entrypoint.CtxExtractInitial(ctx)
				s, c, e = daemon.BuildHTTP(ctx, initial, cmd.Observer)
				if e != nil {
					return e
				}
				return nil
			}, func(ctx context.Context) error {
				if e := s.ListenAndServe(); e != nil {
					return e
				}
				return nil
			})
		},
	}
)

func init() {
	// pflags
	Cmd.PersistentFlags().StringP(http.UnmarshalKeyBind, "b", ":8081", "bind address")
}

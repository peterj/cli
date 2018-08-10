package commands

import (
	"os"

	"github.com/fnproject/cli/client"
	"github.com/fnproject/cli/common"
	"github.com/fnproject/cli/run"
	fnclient "github.com/fnproject/fn_go/client"
	"github.com/fnproject/fn_go/provider"
	"github.com/urfave/cli"
)

type invokeCmd struct {
	provider provider.Provider
	client   *fnclient.Fn
}

// InvokeFnFlags used to invoke and fn
var InvokeFnFlags = append(run.RunFlags,
	cli.BoolFlag{
		Name:  "display-call-id",
		Usage: "whether display call ID or not",
	},
)

// InvokeCommand returns call cli.command
func InvokeCommand() cli.Command {
	cl := invokeCmd{}
	return cli.Command{
		Name:    "invoke",
		Usage:   "\tInvoke a remote function",
		Aliases: []string{"iv"},
		Before: func(c *cli.Context) error {
			var err error
			cl.provider, err = client.CurrentProvider()
			if err != nil {
				return err
			}
			cl.client = cl.provider.APIClient()
			return nil
		},
		ArgsUsage:   "<fn-id>",
		Flags:       InvokeFnFlags,
		Category:    "DEVELOPMENT COMMANDS",
		Description: "This command explicitly invokes a function.",
		Action:      cl.Invoke,
	}
}

func (cl *invokeCmd) Invoke(c *cli.Context) error {
	var contentType string

	fnID := c.Args().Get(0)

	content := run.Stdin()
	wd := common.GetWd()

	if c.String("content-type") != "" {
		contentType = c.String("content-type")
	} else {
		_, ff, err := common.FindAndParseFuncFileV20180707(wd)
		if err == nil && ff.Content_type != "" {
			contentType = ff.Content_type
		}
	}

	return client.Invoke(cl.provider, fnID, content, os.Stdout, c.String("method"), c.StringSlice("e"), contentType, c.Bool("display-call-id"))
}

package account

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/AlexandreMarcq/gozimbra/internal/cmd_utils"
	"github.com/AlexandreMarcq/gozimbra/internal/models"
	client "github.com/AlexandreMarcq/gozimbra/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	attributes []string
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get account information",
		RunE: func(cmd *cobra.Command, args []string) error {
			defaults, err := cmd_utils.GetDefaults(cmd)
			if err != nil {
				return err
			}

			client, err := cmd_utils.AuthC(defaults.Platform)
			if err != nil {
				return err
			}

			log.Printf("Desired attributes are: %v", attributes)

			out, err := cmd_utils.SetupOutput(defaults.Output, defaults.Stdout)
			if err != nil {
				return err
			}
			defer out.Close()

			sort.Strings(attributes)
			_, err = out.WriteString(fmt.Sprintf("account;%s\n", strings.Join(attributes, ";")))
			if err != nil {
				return err
			}

			list, err := getAccountList(client)
			if err != nil {
				return err
			}

			m := models.NewBatchModel(
				"Getting accounts information",
				list,
				out,
				func(s string) tea.Cmd {
					return getAccount(client, s, attributes)
				},
			)

			return cmd_utils.RunModel(m, defaults.NoUI, defaults.Stdout)
		},
	}

	cmd.Flags().StringSliceVar(&attributes, "attributes", nil, "attributs to get for the account(s)")
	cmd.MarkFlagRequired("attributes")

	return cmd
}

func getAccount(client *client.Client, account string, attributes []string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Getting account information for %s", account)
		attrs, err := client.GetAccount(account, attributes)
		if err != nil {
			log.Printf("Error getting information for %s: %v", account, err)
			return models.NewGetMsg(account, nil, err)
		}

		return models.NewGetMsg(account, attrs, nil)
	}
}

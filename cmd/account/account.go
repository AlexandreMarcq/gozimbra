package account

import (
	"log"
	"os"

	"github.com/AlexandreMarcq/gozimbra/internal/cmd_utils"
	client "github.com/AlexandreMarcq/gozimbra/pkg"
	"github.com/spf13/cobra"
)

var (
	accounts  []string
	domains   []string
	inputFile string
)

func NewAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Account related commands",
	}

	cmd.PersistentFlags().StringSliceVar(&accounts, "accounts", nil, "account to operate on")
	cmd.PersistentFlags().StringSliceVar(&domains, "domains", nil, "list of domains to get the accounts from")
	cmd.PersistentFlags().StringVar(&inputFile, "input", "", "input file to get the accounts from")

	cmd.MarkFlagFilename("input")
	cmd.MarkFlagsMutuallyExclusive("accounts", "domains", "input")
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewModifyCmd())

	return cmd
}

func getAccountList(client *client.Client) ([]string, error) {
	log.Println("Getting list of accounts")
	var file *os.File = nil
	var err error

	if inputFile != "" {
		log.Printf("Opening input file '%s'", inputFile)
		file, err = os.Open(inputFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}

	list, err := cmd_utils.ReadInput(client, accounts, domains, file)
	if err != nil {
		return nil, err
	}

	return list, nil
}

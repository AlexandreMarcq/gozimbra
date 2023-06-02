package account

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AlexandreMarcq/gozimbra/internal/cmd_utils"
	"github.com/AlexandreMarcq/gozimbra/internal/models"
	"github.com/AlexandreMarcq/gozimbra/internal/utils"
	client "github.com/AlexandreMarcq/gozimbra/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	data   string
	lock   bool
	unlock bool
)

func NewModifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify",
		Short: "modify account",
		RunE: func(cmd *cobra.Command, args []string) error {
			platform, err := cmd.Flags().GetString("platform")
			if err != nil {
				return err
			}

			outputFile, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			noUI, err := cmd.Flags().GetBool("no-ui")
			if err != nil {
				return err
			}

			stdout, err := cmd.Flags().GetBool("stdout")
			if err != nil {
				return err
			}

			client, err := cmd_utils.AuthC(platform)
			if err != nil {
				return err
			}

			formattedData := make(utils.AttrsMap, 0)

			if lock {
				formattedData["zimbraAccountStatus"] = "locked"
			}

			if unlock {
				formattedData["zimbraAccountStatus"] = "active"
			}

			if data != "" {
				log.Println("Opening data file")
				dataFile, err := os.Open(data)
				if err != nil {
					return err
				}
				defer dataFile.Close()

				log.Println("Reading data file")
				sc := bufio.NewScanner(dataFile)
				for sc.Scan() {
					keyValue := strings.Split(sc.Text(), ";")
					if len(keyValue) < 2 {
						return errors.New("data file is invalid")
					}
					if formattedData[keyValue[0]] != "" {
						continue
					}
					formattedData[keyValue[0]] = keyValue[1]
				}
				if err := sc.Err(); err != nil {
					log.Printf("Error reading data file: %v", err)
					return err
				}
			}
			if len(formattedData.Keys()) == 0 {
				return errors.New("no attribute to modify, use a flag or provide a data file with --data")
			}

			log.Printf("Desired attributes are: %v", formattedData.Keys())

			out, err := cmd_utils.SetupOutput(outputFile, stdout)
			if err != nil {
				return err
			}
			defer out.Close()

			var sb strings.Builder
			for _, k := range formattedData.Keys() {
				_, err := sb.WriteString(fmt.Sprintf(";old_%s;new_%s", k, k))
				if err != nil {
					return err
				}
			}

			_, err = out.WriteString(fmt.Sprintf("account%s\n", sb.String()))
			if err != nil {
				return err
			}

			list, err := getAccountList(client)
			if err != nil {
				return err
			}

			m := models.NewBatchModel(
				"Modifying accounts",
				list,
				out,
				func(account string) tea.Cmd {
					return modifyAccount(client, account, formattedData)
				},
			)

			return cmd_utils.RunModel(m, noUI, stdout)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "data file to modify the accounts (CSV format without headers)")
	cmd.Flags().BoolVar(&lock, "lock", false, "locks the accounts")
	cmd.Flags().BoolVar(&unlock, "unlock", false, "unlocks the accounts")

	cmd.MarkFlagFilename("data")
	cmd.MarkFlagsMutuallyExclusive("lock", "unlock")

	return cmd
}

func modifyAccount(client *client.Client, account string, attributes utils.AttrsMap) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Getting information for %s", account)
		oldAttrs, err := client.GetAccount(account, attributes.Keys())
		if err != nil {
			log.Printf("Error getting information for %s: %v", account, err)
			return models.NewModifyMsg(account, attributes, nil, err)
		}

		log.Printf("Modifying information for %s", account)
		newAttrs, err := client.ModifyAccount(account, attributes)
		if err != nil {
			log.Printf("Error modifying information for %s: %v", account, err)
			return models.NewModifyMsg(account, attributes, nil, err)
		}

		return models.NewModifyMsg(account, oldAttrs, newAttrs, nil)
	}
}

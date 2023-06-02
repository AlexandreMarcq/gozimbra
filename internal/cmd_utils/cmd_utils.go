package cmd_utils

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"time"

	client "github.com/AlexandreMarcq/gozimbra/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AuthC(configSub string) (*client.Client, error) {
	log.Println("Reading configuration file")
	config := viper.Sub(configSub)
	if config == nil {
		log.Println("Error reading configuration file")
		return nil, errors.New("error getting client configuration")
	}

	c := client.NewClient(config.GetString("url"))
	log.Println("Authenticating client")
	if err := c.Auth(config.GetString("username"), config.GetString("password")); err != nil {
		log.Printf("Error authenticating client: %v", err)
		return nil, err
	}

	return c, nil
}

func ReadInput(c *client.Client, accounts, domains []string, inputFile io.Reader) ([]string, error) {
	res := []string{}

	if accounts == nil && domains == nil && inputFile == nil {
		log.Println("Error reading input: missing accounts, domains or input")
		return nil, errors.New("at least one of [accounts domains input] must be set")
	}

	if accounts != nil {
		return accounts, nil
	} else if domains != nil {
		for _, domain := range domains {
			log.Printf("Getting quota usage of domain %v", domain)
			quotas, err := c.GetQuotaUsage(domain)
			if err != nil {
				log.Printf("Error getting quota usage: %v", err)
				return nil, err
			}
			for _, quota := range quotas {
				res = append(res, quota.Account)
			}
		}
	} else if inputFile != nil {
		log.Println("Reading input file")
		sc := bufio.NewScanner(inputFile)
		for sc.Scan() {
			res = append(res, sc.Text())
		}

		if err := sc.Err(); err != nil {
			log.Printf("Error reading input file: %v", err)
			return nil, err
		}
	}

	return res, nil
}

func GetFormattedTime() string {
	return time.Now().Format("20060102-150405")
}

func SetupOutput(outputFile string, stdout bool) (*os.File, error) {
	var out *os.File
	var err error
	if stdout {
		out = os.Stdout
	} else {
		log.Printf("Creating output file '%s'", outputFile)
		out, err = os.Create(outputFile)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func RunModel(model tea.Model, noUI, stdout bool) error {
	var opts []tea.ProgramOption
	if noUI || stdout {
		opts = []tea.ProgramOption{tea.WithoutRenderer()}
	}

	if _, err := tea.NewProgram(model, opts...).Run(); err != nil {
		return err
	}

	log.Println("All done ! ✔")

	return nil
}

type Defaults struct {
	NoLog    bool
	NoUI     bool
	Output   string
	Platform string
	Stdout   bool
}

func GetDefaults(cmd *cobra.Command) (*Defaults, error) {
	def := &Defaults{
		NoLog:    viper.GetBool("defaults.no-log"),
		NoUI:     viper.GetBool("defaults.no-ui"),
		Platform: viper.GetString("defaults.platform"),
		Stdout:   viper.GetBool("defaults.stdout"),
	}

	if !def.NoUI && def.Stdout {
		return nil, errors.New("cannot have both UI and stdout")
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, err
	}

	def.Output = output

	return def, nil
}

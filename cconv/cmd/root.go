package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	targetDate    string
	targetAddress string
)

var rootCmd = &cobra.Command{
	Use:   "cconv",
	Short: "A tool for fetching currency rates and performing currency conversions",
	Long: `
cconv has 3 modes of operation.
1) Returning the exchange rate of a given base currency into one or more target currencies.

> cconv rate EUR to USD GBP CAD

2) Returning a conversion between a value of a given currency and a target currency

> cconv value 123.45 GPB to EUR USD CAD

3) Returning the best exchange rate of the last 7 days

> cconv best CAD to EUR
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&targetAddress, "address", "", "The address to email results to")
	rootCmd.PersistentFlags().StringVar(&targetDate, "date", "", "Target date for rates and conversions")
}

package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/kujtimiihoxha/opencode/internal/llm/models"
	"github.com/spf13/cobra"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List available LLM providers and models",
	Long:  `Display all supported LLM providers and their available models.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return displayAllModels()
	},
}

func displayAllModels() error {
	providerModels := make(map[models.ModelProvider][]models.Model)

	for _, model := range models.SupportedModels {
		providerModels[model.Provider] = append(providerModels[model.Provider], model)
	}

	providers := make([]models.ModelProvider, 0, len(providerModels))
	for provider := range providerModels {
		providers = append(providers, provider)
	}
	sort.Slice(providers, func(i, j int) bool {
		return string(providers[i]) < string(providers[j])
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROVIDER\tMODEL ID\tNAME")

	for _, provider := range providers {
		modelList := providerModels[provider]

		sort.Slice(modelList, func(i, j int) bool {
			return modelList[i].Name < modelList[j].Name
		})

		for _, model := range modelList {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				model.Provider,
				model.ID,
				model.Name,
			)
		}
	}

	return w.Flush()
}

func init() {
	rootCmd.AddCommand(providersCmd)
}

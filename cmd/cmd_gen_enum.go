package cmd

import (
	"github.com/spf13/cobra"

	"github.com/artisanhe/tools/codegen"
	"github.com/artisanhe/tools/courier/enumeration/gen"
)

var cmdGenEnum = &cobra.Command{
	Use:   "enum",
	Short: "generate enum stringify",
	Run: func(cmd *cobra.Command, args []string) {
		enumGenerator := gen.EnumGenerator{
			Filters: args,
		}
		codegen.Generate(&enumGenerator)
	},
}

func init() {
	cmdGen.AddCommand(cmdGenEnum)
}

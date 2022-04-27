package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

var rootCmd = &cobra.Command{
	Use:   "json-schema-validator",
	Short: "A small application to validate json documents with json schema",
	Long:  `A small application to validate json documents with json schema`,
	Run: func(cmd *cobra.Command, args []string) {
		var validationPassed bool = true

		if !viper.IsSet("json-schema") {
			fmt.Println("Json schema must be provided")
			os.Exit(10)
		}

		schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", viper.GetString("json-schema")))

		for _, file := range args {
			err := validateJsonDocument(file, schemaLoader)
			if err != nil {
				fmt.Println(err)
				validationPassed = false
			}
		}

		if !validationPassed {
			os.Exit(10)
		}
		fmt.Println("All documents valid")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("json-schema", "j", "", "Json schema for validation")
	viper.BindPFlag("json-schema", rootCmd.Flags().Lookup("json-schema"))
}

func validateJsonDocument(documentPath string, schemaLoader gojsonschema.JSONLoader) error {

	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", documentPath))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		return err
	}

	if !result.Valid() {
		schemaViolations := ""
		fmt.Printf("Invalid document: %s\n", documentPath)
		for _, desc := range result.Errors() {
			schemaViolations += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("%s", schemaViolations)
	}
	return nil
}

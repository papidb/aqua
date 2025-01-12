package cli

import (
	"fmt"
	"log"

	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/resources"
	"github.com/spf13/cobra"
)

var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with cloud resources",
	Long:  fmt.Sprintf(`Generate and seed the database with %d randomly generated cloud resources for testing purposes.`, resources.DefaultMaxResources),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var env config.Env
		if err := config.LoadEnv(&env); err != nil {
			panic(err)
		}
		app, err := config.New(env)
		if err != nil {
			panic(err)
		}

		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}
		defer app.Database.DB.Close()

		// Retrieve the flag value for max resources
		maxResources, err := cmd.Flags().GetInt("max-resources")
		if err != nil {
			log.Fatalf("Error retrieving max-resources flag: %v", err)
		}

		err = resources.SeedResources(app, maxResources)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var RootCmd = &cobra.Command{
	Use:   "Aqua CLI",
	Short: "Aqua CLI",
	Long: `Aqua CLI is a command-line interface for managing cloud resources.
	You can use it to seed the database with sample cloud resources.
	Start by running "aqua seed" to generate and seed the database with 100 cloud resources.
	You can also use it to start the server by running "aqua server".
	`,
}

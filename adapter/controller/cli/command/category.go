package command

import (
	"github.com/urfave/cli/v2"

	"go-api-arch-clean-template/adapter/controller/cli/action"
	"go-api-arch-clean-template/adapter/controller/cli/presenter"
	"go-api-arch-clean-template/pkg/logger"
)

var CategoryName string

func SetCategoryCommand(app *cli.App, categoryAction *action.CategoryAction) {
	cliFlag := []cli.Flag{
		&cli.StringFlag{
			Name:        "category_name",
			Aliases:     []string{"c"},
			Usage:       "Name for the category",
			Destination: &CategoryName,
		},
	}
	app.Flags = append(app.Flags, cliFlag...)

	cliCommand := []*cli.Command{
		{
			Name:  "category",
			Usage: "Select a category",
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "Name for create a category",
					Action: func(c *cli.Context) error {
						category, err := categoryAction.CreateCategory(CategoryName)
						if err != nil {
							logger.Error(err.Error())
							return err
						}
						presenter.PrettyPrintStructToJson(category)
						return nil
					},
				},
			},
		},
	}
	app.Commands = append(app.Commands, cliCommand...)

}

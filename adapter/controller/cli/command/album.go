package command

import (
	"sort"

	"github.com/urfave/cli/v2"

	"go-api-arch-clean-template/adapter/controller/cli/action"
	"go-api-arch-clean-template/adapter/controller/cli/presenter"
	"go-api-arch-clean-template/pkg/logger"
)

var AlbumTitle string

func SetAlbumCommand(app *cli.App, albumAction *action.AlbumAction) {
	cliFlag := []cli.Flag{
		&cli.StringFlag{
			Name:        "album_title",
			Aliases:     []string{"a"},
			Usage:       "Title for the album",
			Destination: &AlbumTitle,
		},
	}
	app.Flags = append(app.Flags, cliFlag...)

	cliCommand := []*cli.Command{
		{
			Name:    "album",
			Aliases: []string{"a"},
			Usage:   "Select a album",
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "Create for album",
					Action: func(c *cli.Context) error {
						album, err := albumAction.CreateAlbum(AlbumTitle, CategoryName)
						if err != nil {
							logger.Error(err.Error())
							return err
						}
						presenter.PrettyPrintStructToJson(album)
						return nil
					},
				},
			},
		},
	}
	app.Commands = append(app.Commands, cliCommand...)

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))
}

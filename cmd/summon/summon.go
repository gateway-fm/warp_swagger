package summon

import (
	"embed"
	"errors"
	"fmt"
	"github.com/gateway-fm/warp_swagger/internal"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
	"github.com/spf13/cobra"
	"os"
)

func Cmd(app *internal.App, header embed.FS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summon",
		Short: "summon",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := InitTemplates(header); err != nil {
				return err
			}
			if err := app.Summon(); err != nil {
				return err
			}
			if err := cleanUp(); err != nil {
				return fmt.Errorf("cleanUp failed to run:%w", err)
			}
			return nil
		},
	}
	return cmd
}
func cleanUp() error {
	if err := os.RemoveAll("./templates"); err != nil {
		return fmt.Errorf("failed to remove all:%w", err)
	}
	return nil
}
func InitTemplates(header embed.FS) error {

	fetched, err := header.ReadDir("templates")
	if err != nil {
		return fmt.Errorf("an error occurred when exec ReadDir %w", err)
	}
	if err = createTmplFiles(header); err != nil {
		return fmt.Errorf("an error occurred when exec createTmplFiles %w", err)
	}

	for i := range fetched {
		b, err := header.ReadFile("templates/" + fetched[i].Name())
		if err != nil {
			return fmt.Errorf("an error occurred while ReadFile %w", err)
		}
		if err = os.WriteFile("templates/"+fetched[i].Name(), b, 0777); err != nil {
			return fmt.Errorf("an error occurred while WriteFile %w", err)
		}
	}

	return nil
}

func createTmplFiles(header embed.FS) error {
	fetched, err := header.ReadDir("templates")
	if err != nil {
		return fmt.Errorf("an error occurred when exec ReadDir %w", err)
	}
	for i := range fetched {
		if err = templater.CopyTemplatesToTemp(fetched[i].Name()); err != nil {
			return fmt.Errorf("an error occurred when exec CopyTemplatesToTemp %w", err)
		}
	}
	check, err := os.ReadDir("./templates")
	if err != nil {
		return fmt.Errorf("an error occurred when exec ReadDir %w", err)
	}
	if len(check) <= 0 {
		return errors.New("/templates directory is empty")
	}
	fmt.Println("dir len is", len(check))
	return nil
}

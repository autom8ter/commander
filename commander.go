package commander

import (
	"context"
	"fmt"
	"github.com/autom8ter/util"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

type Author struct {
	FullName string
	Email    string
}

type Context struct {
	Name       string    `json:"name"`
	Download   string    `json:"download"`
	Version    string    `json:"version"`
	Authors    []*Author `json:"authors"`
	ConfigPath string    `json:"config_path"`
	EnvPrefix  string
	Meta       map[string]string `json:"meta"`
}

var ctx = context.Background()

var rootCmd = &cobra.Command{}
var fs = &afero.Afero{
	afero.NewOsFs(),
}

func Init(c *Context) {
	rootCmd.Use = c.Name
	c.Version = c.Version
	rootCmd.Long = fmt.Sprintf("Context: \n%s\nCurrent Config: \n%s", util.ToPrettyJsonString(c), util.ToPrettyJson(viper.AllSettings()))
	viper.SetFs(fs)
	viper.SetEnvPrefix(c.EnvPrefix)
	viper.SetConfigFile(c.ConfigPath)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err.Error())
	}
}

func Sub(name string, info string, fn func(ctx context.Context) error) {
	newCtx := context.WithValue(ctx, "settings", viper.AllSettings())
	rootCmd.AddCommand(&cobra.Command{
		Use:  name,
		Short: info,
		Run: func(cmd *cobra.Command, args []string) {
			if err := fn(newCtx); err != nil {
				log.Fatalln(errors.WithStack(err))
			}
		},
	})
}

func Execute() error {
	return rootCmd.Execute()
}

func FS() *afero.Afero {
	return fs
}

func Config() *viper.Viper {
	return viper.GetViper()
}
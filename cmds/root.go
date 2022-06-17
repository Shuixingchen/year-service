package cmds

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Shuixingchen/year-service/handlers"
	"github.com/Shuixingchen/year-service/models"
	"github.com/Shuixingchen/year-service/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/config.yaml", "config file (default is $HOME/./configs/config.yaml)")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parser",
	Short: "",
	Long:  ``,
	Run:   server,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		utils.InitConfig(cfgFile)
	}
}

// initLog represents the init logger.
func initLog() {
	utils.InitLog(utils.Config.Log.Level, utils.Config.Log.Path, utils.Config.Log.Filename)
}

func server(cmd *cobra.Command, args []string) {
	log.Infof("producer start, app name: %#v. (Go version %s %s/%s), gomaxprocs: %d", utils.Config.AppName,
		runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0))
	fmt.Printf("%+v\n", utils.Config)

	// db init
	config := make(map[string]utils.MySQLDSN)
	for i := range utils.Config.StatsDatabase {
		databaseConfig := utils.Config.StatsDatabase[i]
		utils.AddDatabaseConfig(&databaseConfig, config)
	}
	models.InitMySQLDB(config)
	handler := handlers.NewWebHandler()
	handler.Handle()
}

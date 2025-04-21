package main

import (
	"fmt"
	"os"

	"github.com/n101661/maney/server/impl/iris"
)

const configPath = "config.toml"

func main() {
	config, err := LoadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = CreateDefaultConfig(configPath); err != nil {
				fmt.Printf("failed to create %s: %v", configPath, err)
				os.Exit(1)
			}
			fmt.Printf("the %s has been created, please setup first", configPath)
			return
		}
		fmt.Printf("failed to load config: %v", err)
		os.Exit(1)
	}

	repos, err := newBoltRepositories(config.Storage.BoltDBDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer repos.Close()

	services, err := newServices(repos, config.Auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := iris.NewServer(config.App.Config, newIrisController(services))
	if err := s.ListenAndServe(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)); err != nil {
		fmt.Printf("failed to listen and serve: %v", err)
		os.Exit(1)
	}
}

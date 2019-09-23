package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type detail struct {
	SelfBroadcastIP   string //The IP to broadcast to network which host can accept traffic
	SelfBroadcastPort int    //The Port to broadcast to network which host can accept traffic
	Protocol          string //Only `tcp` supported at the moment
	ChainDBPath       string
	StateDBPath       string
	BlockDBPath       string
	BadgerDB          string
	LevelDB           string
	NodeKeyDir        string
	S3Bucket          string
	CheckerBtcURL     string
}

// GetConfiguration ...
func GetConfiguration(env string) *detail {
	var configuration *detail
	var once sync.Once

	if env != "staging" {
		env = "dev"
	}
	once.Do(func() {
		dirname := os.Getenv("GOPATH")
		viper.SetConfigName("config")                                                // Config file name without extension
		viper.AddConfigPath(dirname + "/src/github.com/herdius/herdius-core/config") // Path to config file
		err := viper.ReadInConfig()
		if err != nil {
			log.Printf("Config file not found: %v", err)
		} else {
			configuration = &detail{
				SelfBroadcastIP:   viper.GetString(env + ".selfbroadcastip"),
				SelfBroadcastPort: viper.GetInt(env + ".selfbroadcastport"),
				Protocol:          viper.GetString(env + ".protocol"),
				ChainDBPath:       viper.GetString(env + ".chaindbpath"),
				StateDBPath:       viper.GetString(env + ".statedbpath"),
				BlockDBPath:       viper.GetString(env + ".blockdbpath"),
				BadgerDB:          viper.GetString(env + ".badgerdb"),
				LevelDB:           viper.GetString(env + ".leveldb"),
				NodeKeyDir:        viper.GetString(env + ".nodekeydir"),
				S3Bucket:          viper.GetString(env + ".s3backupbucket"),
				CheckerBtcURL:     viper.GetString(env + ".checkerbtcurl"),
			}
		}
	})

	return configuration
}

func (d *detail) ConstructTCPAddress() string {
	return d.Protocol + "://" + d.SelfBroadcastIP + ":" + fmt.Sprint(d.SelfBroadcastPort)
}

package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tagPrefix = "viper"

//ReadConfig reads the config from the file and merges the local file if present
//config files must be in toml format
func ReadConfig(cmd *cobra.Command, configName string, localConfigName string, config interface{}) error {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}
	configDir, err := cmd.Flags().GetString("config-dir")
	if err != nil {
		return err
	}

	isDebug := os.Getenv("icop_debug")

	viper.SetConfigName(configName)
	viper.AddConfigPath("./")

	if configDir == "" {
		//read config-dir env and use that if set
		configDir = os.Getenv("ICOP_CONFIG_DIR")
		if configDir == "" {
			configDir = "./data"
		}
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigType("toml")

	// fmt.Println("Using config dir <", configDir, ">")

	if isDebug != "1" {
		//listening only if not in debugger. vscode will not stop debug, if we watch the config-file
		//add this to the launch json:
		// "env": {"icop_debug": "1"}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			_ = readConfig(localConfigName, config)
		})
	}

	return readConfig(localConfigName, config)
}

func readConfig(localConfigName string, config interface{}) error {
	// fmt.Println("Reading config file")
	var err error
	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	// merge in local.json if exists
	localfile := filepath.Join(filepath.Dir(viper.ConfigFileUsed()), localConfigName) + ".toml"
	if _, err = os.Stat(localfile); err == nil {
		viper.SetConfigFile(localfile)
		if err := viper.MergeInConfig(); err != nil {
			return err
		}
	}
	config, err = populateConfig(config)
	if err != nil {
		return err
	}

	return nil
}

//PopulateConfig reads the configuration via reflect and sets it into the config struct
func populateConfig(config interface{}) (interface{}, error) {
	err := recursivelySet(reflect.ValueOf(config), "")
	if err != nil {
		return nil, err
	}

	return config, nil
}

func recursivelySet(val reflect.Value, prefix string) error {
	if val.Kind() != reflect.Ptr {
		return errors.New("WTF")
	}

	// dereference
	val = reflect.Indirect(val)
	if val.Kind() != reflect.Struct {
		return errors.New("FML")
	}

	// grab the type for this instance
	vType := reflect.TypeOf(val.Interface())

	// go through child fields
	for i := 0; i < val.NumField(); i++ {
		thisField := val.Field(i)
		thisType := vType.Field(i)
		tag := prefix + getTag(thisType)

		switch thisField.Kind() {
		case reflect.Struct:
			if err := recursivelySet(thisField.Addr(), tag+"."); err != nil {
				return err
			}
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			configVal := float64(viper.GetFloat64(tag))
			thisField.SetFloat(configVal)
		case reflect.Int:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			// you can only set with an int64 -> int
			configVal := int64(viper.GetInt(tag))
			thisField.SetInt(configVal)
		case reflect.String:
			// fmt.Println("setting field ", tag, " to ", viper.GetString(tag))
			thisField.SetString(viper.GetString(tag))
		case reflect.Bool:
			thisField.SetBool(viper.GetBool(tag))
		case reflect.Slice:
			// only string slices allowed
			if thisField.Type() == reflect.TypeOf(([]string)(nil)) {
				strs := viper.GetStringSlice(tag)
				thisField.Set(reflect.ValueOf(strs))
			}
		default:
			//fmt.Println("config: unexpected type detected ~  %s", thisField.Kind())
		}
	}

	return nil
}

func getTag(field reflect.StructField) string {
	// check if maybe we have a special magic tag
	tag := field.Tag
	if tag != "" {
		for _, prefix := range []string{tagPrefix, "mapstructure", "json"} {
			if v := tag.Get(prefix); v != "" {
				return v
			}
		}
	}

	return field.Name
}

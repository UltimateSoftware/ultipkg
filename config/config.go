package config

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	_ "github.com/joho/godotenv/autoload"
)

// Logger is the default logger that envconfig will use. Override this with
// your customized instance of Logrus (if you want).
var Logger = logrus.New()

func init() {
	Logger.Level = logrus.WarnLevel
}

func Get(key string, fallback ...string) string {
	val := os.Getenv(keyToEnv(key))
	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}
	return val
}

func MustGet(key string) string {
	val := Get(key)
	if val == "" {
		Logger.WithFields(logrus.Fields{
			"key":     key,
			"env_var": keyToEnv(key),
		}).Panic("Missing configuration value")
	}

	return val
}

func keyToEnv(key string) string {
	return strings.ToUpper(strings.Replace(key, ".", "_", -1))
}

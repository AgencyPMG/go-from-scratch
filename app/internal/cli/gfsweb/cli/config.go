package cli

import (
	"regexp"
	"strings"

	"github.com/gogolfing/config"
	"github.com/gogolfing/config/loaders/env"
)

func newConfig() (*config.Config, error) {
	c := config.New()
	c.KeyParser = config.PeriodSeparatorKeyParser
	c.AddLoaders(envLoader)

	_, err := c.LoadAll()
	return c, err
}

const (
	//EnvPrefix is the prefix an environment variable must have to be loaded.
	//This prefix is not part of the config key.
	EnvPrefix = "GFSAPP__"
)

var (
	//envLoader is the loader to use for environment files.
	envLoader = env.NewPrefixParserLoader(EnvPrefix, config.KeyParserFunc(EnvKeyParser))

	//envSplitter is a regexp that will split keys around single underscores.
	envSplitter = regexp.MustCompile(`[^_](_)[^_]`)

	//envReplace is a regexp that can replace multiple underscores with one less thatn present.
	envReplacer = regexp.MustCompile(`_(_+)`)
)

//EnvKeyParser parses key according the following rules:
//key is coverted to lower case via strings.ToLower.
//Then it is split around SINGLE underscores.
//Within those values, consecutive underscores are reduced by a count of one.
func EnvKeyParser(key string) config.Key {
	key = strings.ToLower(key)
	indices := envSplitter.FindAllStringSubmatchIndex(key, -1)

	result := []string{}
	last := 0

	for _, index := range indices {
		current := index[2]
		result = append(result, key[last:current])
		last = current + 1
	}

	result = append(result, key[last:])
	for i, r := range result {
		result[i] = envReplacer.ReplaceAllString(r, "$1")
	}

	return config.Key(result)
}

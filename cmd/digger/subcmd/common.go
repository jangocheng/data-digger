package subcmd

import (
	"plugin"
	"strings"

	dig "github.com/segmentio/data-digger/pkg/digger"
	log "github.com/sirupsen/logrus"
)

type commonConfig struct {
	Debug        bool   `flag:"--debug"             help:"turn on debug logging" default:"false"`
	Filter       string `flag:"-f,--filter"         help:"filter regexp to apply before generating stats" default:"-"`
	K            int    `flag:"-k,--num-categories" help:"number of top values to show" default:"25"`
	Numeric      bool   `flag:"--numeric"           help:"treat values as numbers instead of strings" default:"false"`
	PathsStr     string `flag:"--paths"             help:"comma-separated list of paths to generate stats for" default:"-"`
	Plugins      string `flag:"--plugins"           help:"comma-separated list of golang plugins to load at start" default:"-"`
	PrintMissing bool   `flag:"--print-missing"     help:"print out messages that missing all paths" default:"false"`
	Raw          bool   `flag:"--raw"               help:"show raw messages that pass filters" default:"false"`
	RawExtended  bool   `flag:"--raw-extended"      help:"show extended info about messages that pass filters" default:"false"`
	SortByName   bool   `flag:"--sort-by-name"      help:"sort top k values by their category/key names" default:"false"`
}

func makeProcessors(config commonConfig, protoTypes []string) ([]dig.Processor, error) {
	liveStats, err := dig.NewLiveStats(
		dig.LiveStatsConfig{
			K:            config.K,
			Filter:       config.Filter,
			Numeric:      config.Numeric,
			PrintMissing: config.PrintMissing,
			ProtoTypes:   protoTypes,
			PathsStr:     config.PathsStr,
			Raw:          config.Raw,
			RawExtended:  config.RawExtended,
			SortByName:   config.SortByName,
		},
	)
	if err != nil {
		return nil, err
	}

	return []dig.Processor{liveStats}, nil
}

func loadPlugins(pathsStr string) error {
	if pathsStr == "" {
		return nil
	}

	paths := strings.Split(pathsStr, ",")
	for _, path := range paths {
		log.Debugf("Loading plugin at path %s", path)
		_, err := plugin.Open(path)
		if err != nil {
			return err
		}
	}

	return nil
}

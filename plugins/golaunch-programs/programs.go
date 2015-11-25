package main

import (
	"encoding/json"
	"fmt"
	sdk "golaunch/sdk/go"
	"golaunch/sdk/go/idletime"
	"golaunch/sdk/go/system"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	defaultExtensions = []string{".exe", ".lnk", ".bat", ".appref-ms"}
)

type Source struct {
	Path       string   `toml:"path"`
	BonusScore int      `toml:"bonus_score"`
	MaxDepth   int      `toml:"max_depth"`
	Extensions []string `toml:"extensions"`
}

func (s *Source) containsExt(ext string) bool {
	for x := 0; x < len(s.Extensions); x++ {
		if s.Extensions[x] == ext {
			return true
		}
	}

	return false
}

type Config struct {
	MaxResults           int           `toml:"max_results"`
	CacheInMemory        bool          `toml:"cache_in_memory"`
	ScanInterval         string        `toml:"scan_interval"`
	scanInterval         time.Duration `toml:"-"`
	ScanOnStartup        bool          `toml:"scan_on_startup"`
	ScanOnStartupIfEmpty bool          `toml:"scan_on_startup_if_empty"`
	ScanWhenIdle         bool          `toml:"scan_when_idle"`
	MaxScansPerRun       int           `toml:"max_scans_per_run"`
	Sources              []Source      `toml:"sources"`
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)

	var cfg Config
	if _, err := toml.DecodeFile("settings.toml", &cfg); err != nil {
		log.Fatal(err)
	}

	var err error
	cfg.scanInterval, err = time.ParseDuration(cfg.ScanInterval)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range cfg.Sources {
		if len(cfg.Sources[i].Extensions) == 0 {
			cfg.Sources[i].Extensions = defaultExtensions
		}
	}

	system := system.NewSystem()
	defer system.Close()

	// empty for now, but will be filled by "init"
	var metadata sdk.Metadata

	catalog := NewCatalog(&metadata, &cfg, system)
	if err := catalog.Init(); err != nil {
		log.Fatal(err)
	}
	defer catalog.Shutdown()

	if cfg.ScanOnStartupIfEmpty && catalog.IsEmpty() {
		go catalog.Index()
	}

	if cfg.ScanOnStartup {
		go catalog.Index()
	}

	if cfg.scanInterval != 0 {
		go func() {
			scanCount := 0
			for {
				time.Sleep(cfg.scanInterval)
				if cfg.ScanWhenIdle {
					for {
						idle, _ := idletime.Get()
						if idle < cfg.scanInterval {
							time.Sleep(cfg.scanInterval - idle)
						} else {
							break
						}
					}
				}

				catalog.Index()
				scanCount += 1

				if scanCount >= cfg.MaxScansPerRun {
					return
				}
			}
		}()
	}

	// var ppid int
	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Runs the backend",
		Long:  `Runs the backend`,
		Run: func(cmd *cobra.Command, args []string) {
			dec := json.NewDecoder(os.Stdin)
			enc := json.NewEncoder(os.Stdout)

			for {
				var v sdk.Request
				if err := dec.Decode(&v); err != nil {
					continue
				}

				switch v.Method {
				case "init":
					json.Unmarshal(v.Params, &metadata)
				case "query":
					var query string
					json.Unmarshal(v.Params, &query)

					results := catalog.Query(query)
					if results != nil && len(results) > 0 {
						msg := sdk.Response{
							Result: results,
						}

						start := time.Now()
						enc.Encode(msg)
						fmt.Fprintf(os.Stderr, "json encoding: %v\n", time.Now().Sub(start))
					}

					// start = time.Now()
					// a.Reset()
					// fresults := []flatbuffers.UOffsetT{}
					// for x := 0; x < len(results); x++ {
					// 	fresults = append(fresults, a.CreateQueryResult(
					// 		results[x].Title,
					// 		results[x].Subtitle,
					// 		results[x].Image,
					// 		results[x].Query,
					// 		results[x].Score))
					// }
					// a.CreateResponse(v.ID, fresults)
					// fmt.Fprintf(os.Stderr, "flatbuffers encoding: %v\n", time.Now().Sub(start))
				case "action":
					var param sdk.Action
					json.Unmarshal(v.Params, &param)

					if param.Type == "contextmenu" {
						switch param.Name {
						case "Copy path":
							clipboard.WriteAll(param.QueryResult.Path)
						case "Open containing folder":
							system.OpenFolder(filepath.Dir(param.QueryResult.Path))
						case "Run as admin":
						}
					} else {
						catalog.used(param.QueryResult)
						if err := system.RunProgram(param.QueryResult.Path, "", "", ""); err != nil {
							log.Println(err)
						}
					}
				}
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdRun)
	rootCmd.Execute()
}

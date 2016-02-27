package main

import (
	sdk "golaunch/sdk/go"
	"golaunch/sdk/go/idletime"
	"golaunch/sdk/go/plugin"
	"golaunch/sdk/go/system"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/atotto/clipboard"
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
	// MaxScansPerRun       int           `toml:"max_scans_per_run"`
	Sources []Source `toml:"sources"`
}

type Plugin struct {
	metadata sdk.Metadata
	client   *plugin.Client
	system   *system.System
	cfg      *Config
	catalog  *Catalog
}

func NewPlugin() *Plugin {
	p := &Plugin{
		client: plugin.NewClient(),
		// FIXME: how to defer system.Close()??
		system: system.NewSystem(),
	}

	return p
}

func (p *Plugin) Init(m sdk.Metadata) {
	p.metadata = m

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

	catalog := NewCatalog(&p.metadata, &cfg, p.system)
	if err := catalog.Init(); err != nil {
		log.Fatal(err)
	}

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

				// if cfg.MaxScansPerRun > 0 && scanCount >= cfg.MaxScansPerRun {
				// 	return
				// }
			}
		}()
	}

	p.cfg = &cfg
	p.catalog = catalog

	// FIXME: how to defer catalog.Shutdown()??
}

func (p *Plugin) Query(q string) {
	results := p.catalog.Query(q)
	if results != nil && len(results) > 0 {
		p.client.Call("queryresults", results)
		return
	}

	p.client.Call("noqueryresults", nil)
}

func (p *Plugin) Action(a sdk.Action) {
	if a.Type == "contextmenu" {
		switch a.Name {
		case "Copy path":
			clipboard.WriteAll(a.QueryResult.Path)
		case "Open containing folder":
			p.system.OpenFolder(filepath.Dir(a.QueryResult.Path))
		case "Run as admin":
		}
	} else {
		p.catalog.used(a.QueryResult)
		if err := p.system.RunProgram(a.QueryResult.Path, "", "", ""); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)

	s := plugin.NewServer()
	s.Register(NewPlugin())
	s.Serve()
}

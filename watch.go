package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bmatcuk/doublestar"

	"gopkg.in/fsnotify.v1"
)

func throttle(limit time.Duration) func(func()) bool {
	var last int64
	lims := limit.Seconds()

	return func(cb func()) bool {
		now := time.Now().Unix()
		l := atomic.LoadInt64(&last)

		if l+int64(lims) < now {
			cb()
			atomic.StoreInt64(&last, now)
			return true
		}
		return false
	}
}

type Target struct {
	Dir      string
	Globs    []string
	Callback func(string)
}

type Watcher struct {
	targets []Target
}

func NewWatcher() *Watcher {
	return &Watcher{}
}

func (w *Watcher) Add(targets ...Target) {
	w.targets = append(w.targets, targets...)
}

func (w *Watcher) Run() error {
	wr, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, target := range w.targets {
		wr.Add(target.Dir)
		files, err := doublestar.Glob(target.Dir + string(filepath.Separator) + "**")
		if err != nil {
			return err
		}
		for _, file := range files {
			if stat, err := os.Stat(file); err == nil && stat.IsDir() {
				wr.Add(file)
			}
		}
	}

	throttled := throttle(50 * time.Millisecond)
	for {
		select {
		case event := <-wr.Events:
			throttled(func() {
				for _, target := range w.targets {
					//fmt.Println(event.Name, target.Dir, event.Name+string(filepath.Separator))
					if strings.HasPrefix(event.Name, target.Dir+string(filepath.Separator)) {
						for _, glob := range target.Globs {
							if match, _ := doublestar.Match(glob, event.Name); match {
								go target.Callback(event.Name)
								break
							}
						}
					}
				}
			})
		case err := <-wr.Errors:
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func main() {
	w := NewWatcher()
	w.Add(Target{
		Dir:   "src",
		Globs: []string{"*.js", "*.html"},
		Callback: func(file string) {
			cmd := exec.Command("gobble", "build", "-f", "build")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Println(err)
			}
		},
	})

	w.Add(Target{
		Dir:   "plugins",
		Globs: []string{"*.go"},
		Callback: func(file string) {
			pkgname := strings.Split(file, string(filepath.Separator))[1]
			pkg := "." + string(filepath.Separator) + filepath.Join("plugins", pkgname)
			fmt.Printf("Compiling %s...\n", pkg)
			output := filepath.Join(pkg, pkgname)
			if runtime.GOOS == "windows" {
				output += ".exe"
			}
			cmd := exec.Command("go", "build", "-o", output, pkg)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Println(err)
			}
			fmt.Println("Done")
		},
	})
	w.Run()
}

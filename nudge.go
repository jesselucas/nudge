package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/naoina/toml"
)

type Nudge struct {
	Patterns []Pattern
}

type Pattern struct {
	Regex    string
	Template string
}

// CreateNudge parses a TOML file and creates a timeline object
func CreateNudge(configPath string) (*Nudge, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var n Nudge
	if err := toml.Unmarshal(buf, &n); err != nil {
		return nil, err
	}

	return &n, nil
}

func (n *Nudge) Parse(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("Must of more than 0 args")
	}

	// Create a file for each arg
	for _, a := range args {
		// First see if the file exists
		if fileExists(a) {
			continue
		}

		fout, err := os.Create(a)
		if err != nil {
			return err
		}
		defer fout.Close()

		// Check each arg to see if there is a pattern match
		for _, p := range n.Patterns {
			if p.Regex == "" {
				continue
			}

			matched, err := regexp.MatchString(p.Regex, a)
			if err != nil {
				return err
			}

			// If there was a match copy the template to the new file
			if matched {
				fin, err := os.Open(p.Template)
				if err != nil {
					return err
				}
				defer fin.Close()

				r := bufio.NewReader(fin)
				w := bufio.NewWriter(fout)

				buf := make([]byte, 1024)
				for {
					n, err := r.Read(buf)
					if err != nil && err != io.EOF {
						return err
					}
					if n == 0 {
						break
					}

					if _, err := w.Write(buf[:n]); err != nil {
						return err
					}
				}

				if err = w.Flush(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

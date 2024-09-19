/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"k8s.io/release/pkg/release"
	"sigs.k8s.io/release-utils/http"
	"sigs.k8s.io/release-utils/util"
)

const (
	numberOfVersions = 4
	baseURL          = release.ProductionBucketURL + "/release"
)

var errMajorVersion = errors.New("assuming that latest stable major version is 1")

type Binary struct {
	Version         string
	OperatingSystem string
	Architecture    string
	Name            string
}
type Binaries []Binary

func (b Binaries) AllArch() []string {
	allVersions := map[string]struct{}{}
	for _, bin := range b {
		allVersions[bin.Architecture] = struct{}{}
	}
	out := []string{}
	for version := range allVersions {
		out = append(out, version)
	}
	return out
}

func (b Binaries) AllOSes() []string {
	allVersions := map[string]struct{}{}
	for _, bin := range b {
		allVersions[bin.OperatingSystem] = struct{}{}
	}
	out := []string{}
	for version := range allVersions {
		out = append(out, version)
	}
	return out
}

func (b Binaries) AllBins() []string {
	allVersions := map[string]struct{}{}
	for _, bin := range b {
		split := strings.Split(bin.Name, ".")
		allVersions[split[0]] = struct{}{}
	}
	out := []string{}
	for version := range allVersions {
		out = append(out, version)
	}
	return out
}

func (b Binaries) Len() int { return len(b) }

func (b Binaries) Less(i, j int) bool {
	iVersion, err := b[i].version()
	if err != nil {
		log.Fatal(err)
	}
	jVersion, err := b[j].version()
	if err != nil {
		log.Fatal(err)
	}

	if iVersion.Major != jVersion.Major {
		return iVersion.Major > jVersion.Major
	}
	if iVersion.Minor != jVersion.Minor {
		return iVersion.Minor > jVersion.Minor
	}
	if iVersion.Patch != jVersion.Patch {
		return iVersion.Patch > jVersion.Patch
	}

	if b[i].OperatingSystem != b[j].OperatingSystem {
		return b[i].OperatingSystem < b[j].OperatingSystem
	}
	if b[i].Architecture != b[j].Architecture {
		return b[i].Architecture < b[j].Architecture
	}
	return b[i].Name < b[j].Name
}

func (b Binaries) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b Binary) String() string {
	return fmt.Sprintf("%s %s %s %s", b.Name, b.Version, b.OperatingSystem, b.Architecture)
}

func (b Binary) Link() string {
	return fmt.Sprintf(
		"%s/%s/bin/%s/%s/%s",
		strings.TrimPrefix(release.ProductionBucketURL, "https://"),
		b.Version,
		b.OperatingSystem,
		b.Architecture,
		b.Name,
	)
}

func (b Binary) version() (semver.Version, error) {
	tag, err := util.TagStringToSemver(b.Version)
	if err != nil {
		return semver.Version{}, fmt.Errorf("parse tag %s: %w", b.Version, err)
	}

	return tag, nil
}

type arguments struct {
	templateFile string
	outputFile   string
	dataFile     string
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Unable to update index: %v", err)
	}
}

func run() error {
	args := &arguments{}
	fs := flag.NewFlagSet("arguments", flag.ExitOnError)
	fs.StringVar(&args.templateFile, "index-template", "./cmd/update-index/data/index.html.template", "path to the index.html template file")
	fs.StringVar(&args.outputFile, "index-output", "./dist/index.html", "the location of the file this program writes")
	fs.StringVar(&args.dataFile, "binary_details", "./dist/release_binaries.json", "the location of the json file this program writes")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to parsing the flags: %w", err)
	}

	agent := http.NewAgent()

	latestStable, err := agent.Get(baseURL + "/stable.txt")
	if err != nil {
		return fmt.Errorf("get latest stable version: %w", err)
	}

	log.Printf("Got latest stable version: %s", latestStable)

	latestStableSemver, err := util.TagStringToSemver(string(latestStable))
	if err != nil {
		return fmt.Errorf("convert latest stable version to semver: %w", err)
	}

	if latestStableSemver.Major != 1 {
		return fmt.Errorf("%w, but it's %d", errMajorVersion, latestStableSemver.Major)
	}

	stableVersions := []string{string(latestStable)}
	for range numberOfVersions - 1 {
		latestStableSemver.Minor--

		url := fmt.Sprintf("%s/stable-1.%d.txt", baseURL, latestStableSemver.Minor)
		log.Printf("Getting previous stable from: %s", url)
		previousStable, err := agent.Get(url)
		if err != nil {
			return fmt.Errorf("unable to get previous stable: %w", err)
		}

		log.Printf("Got version: %s", previousStable)
		stableVersions = append(stableVersions, string(previousStable))
	}

	binaries := []Binary{}
	for _, version := range stableVersions {
		url := fmt.Sprintf("%s/%s/SHA256SUMS", baseURL, version)
		shaSums, err := agent.Get(url)
		if err != nil {
			return fmt.Errorf("get SHA256SUMS from %q: %w", url, err)
		}

		scanner := bufio.NewScanner(bytes.NewReader(shaSums))
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) != 2 {
				log.Printf("Skipping unknown SHA256SUMS line for version %s: %v", version, fields)
				continue
			}

			binPath := fields[1]
			if !strings.HasPrefix(binPath, "bin/") {
				continue
			}

			parts := strings.Split(binPath, "/")
			if len(parts) < 4 {
				log.Printf("Skipping unknown bin path for version %s: %s", version, binPath)
				continue
			}

			if !shouldInclude(parts[len(parts)-1]) {
				log.Printf("Excluding binary for version %s: %s", version, binPath)
				continue
			}

			binaries = append(binaries, Binary{
				Version:         version,
				OperatingSystem: parts[1],
				Architecture:    parts[2],
				Name:            parts[3],
			})
		}
	}
	sort.Sort(Binaries(binaries))

	tmpl, err := template.New("index.html.template").Funcs(map[string]interface{}{
		"clean": interface{}(clean),
	}).ParseFiles(args.templateFile)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer

	// sort those bins
	bins := Binaries(binaries).AllBins()
	sort.Strings(bins)
	oses := Binaries(binaries).AllOSes()
	sort.Strings(oses)
	arch := Binaries(binaries).AllArch()
	sort.Strings(arch)

	if err := tmpl.Execute(&buf, struct {
		Binaries    Binaries
		AllOSes     []string
		AllBins     []string
		AllVersions []string
		AllArch     []string
		Year        int
	}{
		Binaries:    binaries,
		AllOSes:     oses,
		AllBins:     bins,
		AllVersions: stableVersions[:numberOfVersions],
		AllArch:     arch,
		Year:        time.Now().Year(),
	}); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	err = os.WriteFile(args.outputFile, buf.Bytes(), os.FileMode(0o644)) //nolint:gocritic
	if err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	binaryDetails := struct {
		Binaries    Binaries
		AllOSes     []string
		AllVersions []string
		AllArch     []string
	}{
		Binaries:    binaries,
		AllOSes:     oses,
		AllVersions: stableVersions[:numberOfVersions],
		AllArch:     arch,
	}

	// Store the binaryDetails data in a JSON file
	jsonData, err := json.MarshalIndent(binaryDetails, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	err = os.WriteFile(args.dataFile, jsonData, os.FileMode(0o644))
	if err != nil {
		return fmt.Errorf("write data file: %w", err)
	}

	return nil
}

func shouldInclude(path string) bool {
	if strings.HasSuffix(path, ".exe") {
		return true
	}

	if strings.Contains(path, ".") {
		return false
	}

	return true
}

func clean(item string) string {
	if strings.Contains(item, ".") {
		return strings.ReplaceAll(item, ".", "-")
	}

	if item[0] < 'a' {
		return "a-" + item
	}

	return item
}

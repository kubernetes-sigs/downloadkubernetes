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
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"cloud.google.com/go/storage"
	"github.com/blang/semver/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"sigs.k8s.io/release-utils/util"
)

const (
	numberOfVersions = 4
	defaultTimeout   = 10 * time.Minute
)

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
	return fmt.Sprintf("dl.k8s.io/%s/bin/%s/%s/%s", b.Version, b.OperatingSystem, b.Architecture, b.Name)
}

func (b Binary) version() (semver.Version, error) {
	tag, err := util.TagStringToSemver(b.Version)
	if err != nil {
		return semver.Version{}, fmt.Errorf("parse tag %s: %w", b.Version, err)
	}

	return tag, nil
}

type versions []string

func (v versions) Len() int { return len(v) }

func (v versions) Less(i, j int) bool {
	tagI, err := util.TagStringToSemver(v[i])
	if err != nil {
		return false
	}
	tagJ, err := util.TagStringToSemver(v[j])
	if err != nil {
		return false
	}
	return tagI.Minor > tagJ.Minor
}

func (v versions) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

type arguments struct {
	templateFile string
	outputFile   string
	dataFile     string
}

func main() {
	args := &arguments{}
	fs := flag.NewFlagSet("arguments", flag.ExitOnError)
	fs.StringVar(&args.templateFile, "index-template", "./cmd/update-index/data/index.html.template", "path to the index.html template file")
	fs.StringVar(&args.outputFile, "index-output", "./dist/index.html", "the location of the file this program writes")
	fs.StringVar(&args.dataFile, "binary_details", "./dist/release_binaries.json", "the location of the json file this program writes")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to parsing the flags: %v", err)
	}

	re := regexp.MustCompile(`release/(v\d+\.\d+\.\d+)/bin/([a-zA-Z]+)/([a-zA-Z0-9-]+)/([a-zA-Z0-9-\.]+)`)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err) //nolint:gocritic // no need to cancel the context here
	}

	bh := client.Bucket("kubernetes-release")
	oi := bh.Objects(ctx, &storage.Query{Prefix: "release/stable-"})
	stableVersions := []string{}
	for {
		attr, err := oi.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("%+v\n", err)
			break
		}
		if strings.HasSuffix(attr.Name, "-1.txt") {
			continue
		}
		oh := bh.Object(attr.Name)
		reader, err := oh.NewReader(ctx)
		if err != nil {
			fmt.Printf("%+v\n", err)
			break
		}
		out, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("%+v\n", err)
			break
		}
		stableVersions = append(stableVersions, strings.TrimSpace(string(out)))
	}

	sort.Sort(versions(stableVersions))
	binaries := []Binary{}
	for _, version := range stableVersions[:numberOfVersions] {
		query := &storage.Query{Prefix: fmt.Sprintf("release/%s/bin", version)}
		oi := bh.Objects(ctx, query)
		for {
			attr, err := oi.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Printf("%+v\n", err)
				break
			}
			parts := strings.Split(attr.Name, "/")
			if !shouldInclude(parts[len(parts)-1]) {
				continue
			}
			groups := re.FindStringSubmatch(attr.Name)
			binaries = append(binaries, Binary{
				Version:         groups[1],
				OperatingSystem: groups[2],
				Architecture:    groups[3],
				Name:            groups[4],
			})
		}
	}

	sort.Sort(Binaries(binaries))
	tmpl, err := template.New("index.html.template").Funcs(map[string]interface{}{
		"clean": interface{}(clean),
	}).ParseFiles(args.templateFile)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	err = os.WriteFile(args.outputFile, buf.Bytes(), os.FileMode(0o644)) //nolint:gocritic
	if err != nil {
		panic(err)
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
		panic(err)
	}

	err = os.WriteFile(args.dataFile, jsonData, os.FileMode(0o644))
	if err != nil {
		panic(err)
	}
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

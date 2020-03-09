package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	numberOfVersions = 5
)

type version struct {
	major, minor, patch int
}

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
	iVersion := b[i].version()
	jVersion := b[j].version()
	if iVersion.major != jVersion.major {
		return iVersion.major > jVersion.major
	}
	if iVersion.minor != jVersion.minor {
		return iVersion.minor > jVersion.minor
	}
	if iVersion.patch != jVersion.patch {
		return iVersion.patch > jVersion.patch
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
func (b Binary) Row() string {
	tr := fmt.Sprintf(`<tr class="%s %s %s %s">%%s</tr>`, clean(b.Version), b.OperatingSystem, clean(b.Architecture), b.Name)
	rows := fmt.Sprintf(`
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
    <td><span title="download"><a href="https://%s">%s</a></span></td>
    <td><span class="icon"><i class="fa fa-copy"></i></span><span title="copy to clipboard"><a class="copy" href="https://%s">  %s</a></span></td>`, b.Version, b.OperatingSystem, b.Architecture, b.downloadLink(), b.Name, b.downloadLink(), b.downloadLink())
	return fmt.Sprintf(tr, rows)
}
func (b Binary) downloadLink() string {
	return fmt.Sprintf("dl.k8s.io/%s/bin/%s/%s/%s", b.Version, b.OperatingSystem, b.Architecture, b.Name)
}

func (b Binary) version() version {
	items := strings.Split(b.Version[1:], ".")
	major, _ := strconv.Atoi(items[0])
	minor, _ := strconv.Atoi(items[1])
	patch, _ := strconv.Atoi(items[2])
	return version{major, minor, patch}
}

type versions []string

func (v versions) Len() int { return len(v) }
func (v versions) Less(i, j int) bool {
	items := strings.Split(v[i][1:], ".")
	iminor, _ := strconv.Atoi(items[1])
	items = strings.Split(v[j][1:], ".")
	jminor, _ := strconv.Atoi(items[1])
	return iminor > jminor
}
func (v versions) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

type arguments struct {
	templateFile string
	outputFile   string
}

func main() {
	args := &arguments{}
	fs := flag.NewFlagSet("arguments", flag.ExitOnError)
	fs.StringVar(&args.templateFile, "index-template", "./cmd/update-index/data/index.html.template", "path to the index.html template file")
	fs.StringVar(&args.outputFile, "index-output", "./dist/index.html", "the location of the file this program writes")
	fs.Parse(os.Args[1:])

	re := regexp.MustCompile(`release/(v\d+\.\d+\.\d+)/bin/([a-zA-Z]+)/([a-zA-Z0-9-]+)/([a-zA-Z0-9-\.]+)`)

	client, err := storage.NewClient(context.Background(), option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	bh := client.Bucket("kubernetes-release")
	oi := bh.Objects(context.Background(), &storage.Query{Prefix: "release/stable-"})
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
		reader, err := oh.NewReader(context.Background())
		if err != nil {
			fmt.Printf("%+v\n", err)
			break
		}
		out, err := ioutil.ReadAll(reader)
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
		oi := bh.Objects(context.Background(), query)
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
	arch := Binaries(binaries).AllArch()
	sort.Strings(arch)

	if err := tmpl.Execute(&buf, struct {
		Binaries    Binaries
		AllOSes     []string
		AllBins     []string
		AllVersions []string
		AllArch     []string
	}{
		Binaries:    binaries,
		AllOSes:     Binaries(binaries).AllOSes(),
		AllBins:     bins,
		AllVersions: stableVersions[:numberOfVersions],
		AllArch:     arch,
	}); err != nil {
		panic(err)
	}
	ioutil.WriteFile(args.outputFile, buf.Bytes(), os.FileMode(0644))
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
	if strings.Index(item, ".") != -1 {
		return strings.ReplaceAll(item, ".", "-")
	}
	if item[0] < 'a' {
		return "a-" + item
	}
	return item
}

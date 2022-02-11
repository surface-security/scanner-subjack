package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/fopina/subjack/subjack"
)

func main() {
	defaultConfig := "/fingerprints.json"

	o := subjack.Options{}

	// store results in a temporary file and move to final destination upon completion
	// otherwise Surface file sync will import (and delete) the file in chunks...
	var finalOutput string
	var hideNotVulnerable bool

	flag.IntVar(&o.Threads, "t", 10, "Number of concurrent threads (Default: 10).")
	flag.IntVar(&o.Timeout, "timeout", 10, "Seconds to wait before connection timeout (Default: 10).")
	flag.BoolVar(&o.Ssl, "ssl", false, "Force HTTPS connections (May increase accuracy (Default: http://).")
	flag.BoolVar(&o.All, "a", false, "Find those hidden gems by sending requests to every URL. (Default: Requests are only sent to URLs with identified CNAMEs).")
	// reverse original verbose flag as usual parser will need to see "Not Vulnerable" to disable those results
	// flag.BoolVar(&o.Verbose, "v", false, "Display more information per each request.")
	flag.BoolVar(&hideNotVulnerable, "q", false, "Hide non vulnerable targets.")
	flag.StringVar(&finalOutput, "o", "/output/output.txt", "Output results to file (Subjack will write JSON if file ends with '.json').")
	flag.StringVar(&o.Config, "c", defaultConfig, "Path to configuration file.")
	flag.BoolVar(&o.Manual, "m", false, "Flag the presence of a dead record, but valid CNAME entry.")
	flag.BoolVar(&o.IncludeEdge, "e", false, "Include edge takeover cases.")
	flag.BoolVar(&o.Follow, "follow", false, "Follow redirects.")

	flag.Parse()

	if flag.NArg() > 0 {
		o.Wordlist = flag.Arg(0)
	} else {
		o.Wordlist = "/input/input.txt"
	}
	o.NoColor = true

	file, err := ioutil.TempFile("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	o.Output = file.Name()
	o.Verbose = !hideNotVulnerable

	subjack.Process(&o)

	copyFile(file.Name(), finalOutput)
}

func copyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	return nil
}

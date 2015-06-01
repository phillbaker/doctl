package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/digitalocean/doctl/Godeps/_workspace/src/gopkg.in/yaml.v1"
)

var Output *os.File

func init() {
	Output = os.Stdout
}

type Outputable interface {
	// Headers used in table format
	Headers() []string
	// Format string used in table format
	FormatString() string
	// A function that given a single object, returns the appropriate values for a table row
	RowValues(datum interface{}) []interface{}
	// A function that given a row index returns object for the row
	RowObject(int) interface{}
	// Returns the length of the rows
	Len() int
}

func WriteOutputable(data Outputable) {
	switch OutputFormat {
	case "table":
		tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		defer tw.Flush()
		fmt.Fprintln(tw, strings.Join(data.Headers(), "\t"))

		for i := 0; i < data.Len(); i++ {
			datum := data.RowObject(i)
			fields := data.RowValues(datum)
			fmt.Fprintf(tw, data.FormatString(), fields...)
		}
	case "json":
		output, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("JSON Encoding Error: %s", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, string(output))
	case "yaml":
		output, err := yaml.Marshal(data)
		if err != nil {
			fmt.Printf("YAML Encoding Error: %s", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, string(output))
	}
}

//
// Legacy printers
//

type CLIOutput struct {
	w *tabwriter.Writer
}

func NewCLIOutput() *CLIOutput {
	return &CLIOutput{
		w: tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0),
	}
}

func WriteOutput(data interface{}) {
	var output []byte
	var err error

	switch OutputFormat {
	case "json":
		output, err = json.Marshal(data)
		if err != nil {
			fmt.Printf("JSON Encoding Error: %s", err)
			os.Exit(1)
		}

	case "yaml":
		output, err = yaml.Marshal(data)
		if err != nil {
			fmt.Printf("YAML Encoding Error: %s", err)
			os.Exit(1)
		}
	}
	fmt.Printf("%s", string(output))
}

func (c *CLIOutput) Header(a ...string) {
	fmt.Fprintln(c.w, strings.Join(a, "\t"))
}

func (c *CLIOutput) Writeln(format string, a ...interface{}) {
	fmt.Fprintf(c.w, format, a...)
}

func (c *CLIOutput) Flush() {
	c.w.Flush()
}

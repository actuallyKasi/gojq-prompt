package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/actuallykasi/gojq-prompt/completer"

	"github.com/c-bata/go-prompt"
	"github.com/itchyny/gojq"
	flag "github.com/spf13/pflag"
	"github.com/tidwall/pretty"
)

var (
	// Version .
	Version = "unset"
	// Revision .
	Revision = "unset"
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "gojq-promt <jq filter> [file...]")
		flag.PrintDefaults()
	}
	// Flags after the first non flag are parsed as Argsz
	// if flag.NArg() > 0 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	var file = flag.StringP("file", "f", "", "json file name to parse")
	var url = flag.StringP("url", "u", "", "url to parse")
	flag.Parse()
	if *file == "" && *url == "" {
		// flag.Usage()
		// os.Exit(1)
	} else if *file != "" && *url != "" {
		log.Fatalf("Either one flag is supported. -f or -url")
		flag.Usage()
		os.Exit(1)
	}

	var lines []string
	if *file != "" { // file mode
		f, err := os.Open(*file)
		defer f.Close()
		if err != nil {
			log.Fatalf("Error opening file %s : %s", *file, err)
		}
		reader := bufio.NewReader(f)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	} else if *url != "" { // url mode

	} else { // standard input mode
		// Check for os.stdin for input
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalf("Error opening stdin: %s", err)
		}
		if fi.Size() > 0 {
			reader := bufio.NewReader(os.Stdin)
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
		}
	}

	inputString := strings.Join(lines, "\n")

	fmt.Printf("gojq-prompt %s (rev-%s)\n", Version, Revision)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")

	fmt.Println("gh-prompt started")

	c, err := completer.NewCompleter(Version)
	if err == completer.ErrNotFoundRemotes {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to get remote informations on your git directory.\n")
		os.Exit(1)
	} else if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Initialization error: %s\n", err)
		_, _ = fmt.Fprintf(os.Stderr, "You current directory might not be a git repository.")
		os.Exit(1)
	}
	p := prompt.New(
		executorFunc,
		c.Complete,
		prompt.OptionTitle("gh-prompt: interactive GitHub CLI"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(gpc.FilePathCompletionSeparator),
	)
	p.Run()
	// prettified, err := parseJQ(inputString, query)
	// fmt.Println(prettified)
}

func parseJQ(inputString string, query *gojq.Query) (prettified string, err error) {
	var input map[string]interface{}
	err = json.Unmarshal([]byte(inputString), &input)
	if err != nil {
		log.Fatalf("Error parsing json from stdin: %s", err)
	}

	iter := query.Run(input) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			log.Fatalln(err)
		}

		json, err := getJSONFormatted(v)
		if err != nil {
			log.Fatalf("Error formating json: %s", err)
		}

		prettifiedBytes := pretty.Pretty(json)
		// prettifiedBytes = pretty.Color(prettifiedBytes, nil)
		prettified = string(prettifiedBytes)
	}
	return prettified, err
}

func getJSONFormatted(v interface{}) (bytes []byte, err error) {
	bytes, err = json.Marshal(v)
	if err != nil {
		return bytes, err
	}
	return bytes, err
}

func executorFunc(inputString string, queryString string) {
	queryString = strings.TrimSpace(queryString)
	if queryString == "" {
		return
	} else if queryString == "quit" || queryString == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}
	var query *gojq.Query
	var err error
	if query, err = gojq.Parse(os.Args[1]); err != nil {
		fmt.Printf("Error parsing query expression: %s", err)
		return
	}
	var prettified string
	if prettified, err = parseJQ(inputString, query); err != nil {
		fmt.Printf("Error parsing query input json: %s", err)
		return
	}
	fmt.Println(prettified)
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/pflag"
)

type docOutput struct {
	title        string
	funcName     string
	body         []string
	path         string
	goFilename   string
	jsonFilename string
	apidocs      bool
	postTitle    string
	post         []string
}

func newDocOutput(filePath string, goFilename string) docOutput {
	base := strings.Split(goFilename, "_test.go")
	return docOutput{
		path:         filePath,
		goFilename:   goFilename,
		jsonFilename: fmt.Sprintf("%s.apisamples", base[0]),
	}
}

func (d *docOutput) getGoInclude() string {
	return fmt.Sprintf("include::%s[tag=%s,indent=0]", path.Join(d.path, d.goFilename), d.funcName)
}

func (d *docOutput) getJsonInclude() string {
	return fmt.Sprintf("include::%s[tag=%s]", path.Join(d.path, d.jsonFilename), d.funcName)
}

func (d *docOutput) functionName() string {
	return fmt.Sprintf("Test_%s", d.funcName)
}

func (d *docOutput) processPost() {
	var body []string
	found := false
	end := false
	for _, line := range d.body {
		if strings.Contains(line, "startpostdocs") {
			t := strings.TrimPrefix(line, "startpostdocs")
			d.postTitle = strings.Trim(t, " ")
			found = true
			continue
		}
		if strings.Contains(line, "endpostdocs") {
			end = true
			continue
		}
		if found && !end {
			d.post = append(d.post, line)
			continue
		}

		body = append(body, line)
	}
	d.body = body
}

func (d *docOutput) getAsciiDoc() string {
	d.processPost()
	var out strings.Builder
	out.WriteString("// tag::" + d.funcName + "[]\n")
	out.WriteString("<<<\n== ")
	out.WriteString(d.title)
	out.WriteString("\n")
	for _, line := range d.body {
		out.WriteString(line)
		out.WriteString("\n")
	}
	out.WriteString("\n[#" + strings.ToLower(d.title) + "_" + strings.ToLower(d.funcName) + "_go]")
	out.WriteString("\n.Go ")
	out.WriteString(d.title)
	// out.WriteString("\n[%collapsible]\n")
	// out.WriteString("====\n")
	out.WriteString("\n[source,go]\n")
	out.WriteString("----\n")
	out.WriteString(d.getGoInclude())
	out.WriteString("\n----\n")
	// out.WriteString("====\n")

	if d.apidocs {
		out.WriteString("\n[#" + strings.ToLower(d.title) + "_" + strings.ToLower(d.funcName) + "_json]")
		out.WriteString("\n.JSON ")
		out.WriteString(d.title)
		// out.WriteString("\n[%collapsible]\n")
		// out.WriteString("====\n")
		out.WriteString("\n[source,json]\n")
		out.WriteString("----\n")
		out.WriteString(d.getJsonInclude())
		out.WriteString("\n----\n")
		// out.WriteString("====\n")
	}

	if d.postTitle != `` {
		out.WriteString("\n=== ")
		out.WriteString(d.postTitle)
		out.WriteString("\n")
		for _, line := range d.post {
			out.WriteString(line)
			out.WriteString("\n")
		}
		out.WriteString("\n")
	}
	out.WriteString("// end::" + d.funcName + "[]\n")
	return out.String()
}

var (
	flags      *pflag.FlagSet
	sourceFile string
	outFile    string
	overwrite  bool
	runTests   bool
	noheader   bool
)

func init() {
	/*
		Used by documentation for the manpage
		tag::options[]
		*--source string*
			Source file to parse into AsciiDoc, recommended is to set the absolute path.

		*--out*
			File to write to, if left empty writes to stdout

		*--f*
			Overwrite the existing out file

		*--run*
			Run the tests to produce the output file for the JSON samples.
			The JSON samples need to be written to a file called the same as
			the source Go file minus the _test with a .apisamples extension

		*--no-header*
			Do not set a document header and ToC
		end::options[]
	*/
	flags = pflag.NewFlagSet("AsciiDoc Generator for Distrib", pflag.ContinueOnError)
	flags.StringVar(&sourceFile, "source", "", "Source file to parse into AsciiDoc, recommended is to set the absolute path.")
	flags.StringVar(&outFile, "out", "", "File to write to, if left empty writes to stdout")
	flags.BoolVar(&overwrite, "f", false, "Overwrite the existing out file")
	flags.BoolVar(&runTests, "run", false, "Run the tests to produce the output file for the JSON samples. "+
		"The JSON samples need to be written to a file called the same as the source Go file minus the _test with a .apisamples extension")
	flags.BoolVar(&noheader, "no-header", false, "Do not set a document header and ToC")
}

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		if err != pflag.ErrHelp {
			fmt.Fprint(os.Stderr, err.Error()+"\n")
			flags.PrintDefaults()
		}
		os.Exit(100)
	}

	if sourceFile == "" {
		fmt.Fprint(os.Stderr, "Sourcefile must be set\n")
		flags.PrintDefaults()
		os.Exit(100)
	}

	filePath, goFilename := path.Split(sourceFile)
	data, err := ioutil.ReadFile(path.Join(filePath, goFilename))
	if err != nil {
		exitError("Could not read source file", err)
	}
	var f []string
	f = append(f, strings.Split(string(data), "\n")...)

	var funcs []int
	var docbuf []docOutput

	// first find all tests in the file
	for i, line := range f {
		if strings.Contains(line, "Test_") {
			funcs = append(funcs, i)
		}
	}

	// now we know where the funcs are we can test the read file for our docs
	for _, i := range funcs {
		doc := newDocOutput(filePath, goFilename)

		split := strings.Split(f[i], "Test_")
		split = strings.Split(split[1], "(")
		doc.funcName = split[0]

		search := i

		if search+1 == len(f) {
			break
		}

		// check if it's the first after the func declare, if it means this func has no docs
		if !strings.Contains(f[search+1], "startdocs") && !strings.Contains(f[search+1], "startapidocs") {
			continue
		}

		for {
			search++
			if search == len(f) {
				break
			}

			buf := f[search]
			buf = strings.Trim(buf, "\t")

			// title found
			if strings.Contains(buf, "startdocs") || strings.Contains(buf, "startapidocs") {
				if strings.Contains(buf, "startapidocs") {
					doc.apidocs = true
				}
				buf = strings.TrimPrefix(buf, "// startdocs")
				buf = strings.TrimPrefix(buf, "// startapidocs")
				doc.title = strings.Trim(buf, " ")
				continue
			}

			// end
			if strings.Contains(buf, "// enddocs") {
				break
			}

			// contents
			if strings.Contains(buf, "//") {
				if strings.Contains(buf, "// tag::") {
					continue
				}
				if strings.Contains(buf, "// end::") {
					continue
				}
				buf = strings.TrimPrefix(buf, "//")
				doc.body = append(doc.body, strings.Trim(buf, " "))
				continue
			}
		}

		if doc.title != "" {
			docbuf = append(docbuf, doc)
		}
	}

	_, title := path.Split(sourceFile)
	title = strings.Replace(title, "test.go", "", 1)
	title = strings.ReplaceAll(title, "_", " ")

	if outFile != `` {
		if stat, _ := os.Stat(outFile); stat != nil {
			if !overwrite {
				exitError("File already exists", nil)
			}
			if err = os.Remove(outFile); err != nil {
				exitError("Could not delete file", err)
			}
		}
		o, err := os.Create(outFile)
		if err != nil {
			exitError("Could not create file for writing: "+outFile, err)
		}

		if !noheader {
			o.WriteString("= " + title + "\n")
			o.WriteString(":toc: left\n\n")
		}
		o.WriteString("// THIS FILE IS GENERATED. DO NOT EDIT.\n")
		for _, out := range docbuf {
			_, err = o.WriteString(out.getAsciiDoc())
			if err != nil {
				exitError("Could not write string to file: "+outFile, err)
			}
		}
	} else {
		if !noheader {
			fmt.Printf("= " + title + "\n")
			fmt.Printf(":toc: left\n\n")
		}
		fmt.Printf("// THIS FILE IS GENERATED. DO NOT EDIT.\n")
		for _, out := range docbuf {
			fmt.Printf("%s\n", out.getAsciiDoc())
		}
	}

	if runTests {
		var testsToRun []string
		for _, out := range docbuf {
			testsToRun = append(testsToRun, out.functionName())
		}

		tests := strings.Join(testsToRun, "|")
		packagedir, _ := path.Split(sourceFile)
		// command := "/usr/local/go/bin/go"
		args := []string{
			"test", "-timeout", "30s",
			packagedir,
			"-run",
			fmt.Sprintf("^(%s)$", tests),
			"-count", "1",
		}
		// fmt.Printf("Executing: \n%s %s\n", command, args)
		cmd := exec.Command("go", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("\nFailed to run for sourcefile: %s\n", sourceFile)
			fmt.Printf("Command: %s ", "go")
			for _, arg := range args {
				fmt.Printf("%s ", arg)
			}
			fmt.Printf("\n")
			exitError("Could not run tests", err)
		}
	}

	os.Exit(0)
}

func exitError(reason string, err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, reason+": "+err.Error()+"\n")
	} else {
		fmt.Fprint(os.Stderr, reason+"\n")
	}
	os.Exit(1)
}

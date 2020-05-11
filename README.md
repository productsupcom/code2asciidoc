# code2asciidoc

## About

code2asciidoc solves the problem where code samples for API’s are not properly
tested in most documentation.
It leverages the Go testing framework by writing the sample in there, using
the AsciiDoc tagging system and this tool joins it all together.

When using gRPC best use it in conjunction with [proto2asciidoc](https://github.com/productsupcom/proto2asciidoc)

## Usage

**--source string**
Source file to parse into AsciiDoc, recommended is to set the absolute path.

**--out**
File to write to, if left empty writes to stdout

**--f**
Overwrite the existing out file

**--run**
Run the tests to produce the output file for the JSON samples.
The JSON samples need to be written to a file called the same as
the source Go file minus the \_test with a .apisamples extension

**--no-header**
Do not set a document header and ToC

## Output

You can view the output of the example below under
[Examples Documentation](https://github.com/productsupcom/code2asciidoc/blob/master/docs/generated/examples.adoc)

## Example

The following file shows a complete example, it’s provided inside the `examples/`
directory including a Protobuf definition plus the compiled output.

The output can be seen at `docs/generated/examples.adoc`.

**Test Example.**

``` go
package examples

import (
    "testing"

    "github.com/productsupcom/code2asciidoc/cmd/code2asciidoc/sample"
)

func Test_ExampleSample1(t *testing.T) {
    // startapidocs Example Hello World
    // A simple demonstration to show the syntax of code2asciidoc.
    // You can use any form of asciidoc inside the comments, as it will be
    // properly parsed by AsciiDoc later.
    //
    // What is important is that the function of the test method is the same
    // as the Tags set below.
    // Also the name for the CreateJsonSample must match (the name of the current)
    // package!
    //
    // Be sure to place it inside a filename you can later ingest into your
    // API docs. Recommended is to keep it per gRPC Service.
    //
    // [NOTE.small]
    // ====
    // We can even use blocks in here if we want.
    // ====
    // startpostdocs Response
    // And in some cases you also want to have docs behind the samples to make
    // data more clear.
    // endpostdocs
    // enddocs
    // tag::ExampleSample1[]
    ex := Example{
        SomeString: "Foo",
        SomeInt:    1,
    }
    // end::ExampleSample1[]

    f, err := sample.Setup("examples.apisamples")
    if err != nil {
        t.Errorf("%v", err)
    }
    err = sample.CreateJsonSample(&ex, "ExampleSample1", f)
    if err != nil {
        t.Errorf("%v", err)
    }
}
```

You can produce the same output by doing:

``` shell
code2asciidoc --source ${PWD}/examples/examples_test.go --out docs/generated/examples.adoc --run --f
```

<div class="informalexample">

The `${PWD}` (for Makefile `${CURDIR}`) is very important, it has to be an absolute path for the source.

</div>

The `--run` causes the tool to call the Go test-suite which will produce the
output files.

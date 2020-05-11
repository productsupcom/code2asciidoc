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

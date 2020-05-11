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

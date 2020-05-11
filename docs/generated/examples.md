# examples

## Example Hello World

A simple demonstration to show the syntax of code2asciidoc.
You can use any form of asciidoc inside the comments, as it will be
properly parsed by AsciiDoc later.

What is important is that the function of the test method is the same
as the Tags set below.
Also the name for the CreateJsonSample must match (the name of the current)
package\!

Be sure to place it inside a filename you can later ingest into your
API docs. Recommended is to keep it per gRPC Service.

<div class="note">

We can even use blocks in here if we want.

</div>

**Go Example Hello World.**

``` go
ex := Example{
    SomeString: "Foo",
    SomeInt:    1,
}
```

**JSON Example Hello World.**

``` json
{
  "someString": "Foo",
  "someInt": 1
}
{
  "someString": "Foo",
  "someInt": 1
}
```

### Response

And in some cases you also want to have docs behind the samples to make
data more clear.

# LXT

*lxt* is the L. XSL Templates language.
It is designed to present a more “natural” language to programmers,
rather than the highly verbose and difficult to comprehend XSLT language.

So, it’s design is to take the Domain-Specific Language LXT and compile that into XSLT,
which can then be used for all things that XSLT can be used for.

## Grammar

### Keywords

#### Top-level Directives
* output: Defines the format of the output document via the given `( param => "value" )` map.

#### Variables and Parameters
* var: define an `xsl:variable` with the given value.
* param: define an `xsl:param` that defaults to the given value, but can be overridden by arguments.

#### Subfunctions and Templates
* sub: define a named `xsl:template`: `sub name ( param => <default> ) body`.
* call: call a named `xsl:template`: `call name ( argument => <value> )`.
* template: define an anonymous `xsl:template` used for template matching.
* apply-templates: automatically match and apply matching templates.

#### Control flow:
* when/otherwise: these are chained together to construct an `xsl:choose` block. An `otherwise` always terminates the `xsl:choose` block.
* if: constructs a simple if-then `xsl:if` block from the given XPath and expression.
* foreach/for-each: constructs a `xsl:for-each` to loop over a given XPath selector, executing the given body.

#### HTML/XHTML sugar
* tag: constructs an `xsl:element` with the given name and body.
* attribs: constructs a map of `key => value` attributes for the current block using `xsl:attribute`.
* div: constructs the XSL appropriate to output a `<div class="name">body</div>` with the given class name, and body.
* span: constructs the XSL appropriate to output a `<span class="name">body</span>` with the given class name, and body.

### Strings

There are three kinds of strings: `"double quote"`, `"single quote"`, and back-tick quotes.

There is currently no distinction between them, except that each quote format does not need to escape any of the others.
(NOTE: This will likely change, as the language is made more strict.)

Most times, when a quoted string is included, it will automatically put into an `xsl:text` block.

### XPath statements

XPath statements are:
* `$variable-identifier`
* `@attribute-identifier`
* `<./namespace::simple/*/xpath/statements/@identifier>`
* `<{ full < extended/xpath[statements = 0] }>`

Any abitrary XPath statement can be written by surrounding it with angled brackets around braces: `<{xpath}>`.

If an XPath statement only consists of selectors, it can be shorthanded by surrounding it by only angled brackets: `<xpath>`

If an XPath statement only consists of a variable reference, it can be used directly: `$variable`.

If an XPath statement only consists of an attribute reference, it can be used directly: `@attribute`.

The common use of an XPath statement of just `.`, it requires being placed in angled brackets: `<.>`

In most cases, where an XPath appears as an expression, it is automatically turned into a `xsl:value-of` block.

### Blocks

A block is any group of expressions between a block-start, and a corresponding block-end.
The three kinds are: `( … )`, `[ … ]` and `{ … }`.

Some keywords have different semantics according to which specific block is used,
though often any pair of block-start and block-end can be used to define most blocks.
(NOTE: This will likely change, as the language is made more strict.)

Example: The `sub` keyword, will interpret a `( block )` as a map of param names to default values.
The other two block types will be interpreted as expressions to define the body of the subfunction.

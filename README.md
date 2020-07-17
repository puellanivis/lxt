# LXT

*lxt* is the L. XSL Templates language.
It is designed to present a more “natural” language to programmers,
rather than the highly verbose and difficult to comprehend XSLT language.

So, it’s design is to take the Domain-Specific Language LXT and compile that into XSLT,
which can then be used for all things that XSLT can be used for.

## Grammar

### Keywords

* output: Defines the format of the output document.
* var: define an `xsl:variable` with the given value.
* param: define an `xsl:param` that defaults to the given value, but can be overridden by arguments.

#### Subfunctions and Templates
* sub: define a named `xsl:template`
* call: call a named `xsl:template`
* template: define an unnamed matching `xsl:template`
* apply-templates: automatically match and apply matching `xsl:template`s

#### Control flow:
* when/otherwise: construct an `xsl:choose` block.
* if: construct a conditional `xsl:if`
* foreach/for-each: construct a `xsl:for-each` to loop over an XPath selector.

#### HTML/XHTML sugar
* tag: creates a given XHTML tag with the given name, and body.
* attribs: a map of `key => value` attributes for the current block using `xsl:attribute`.
* div: creates an XHTML `<div></div>` with the given class name, and body.
* span: creates an XHTML `<span></span>` with the given class name, and body.

### Strings

There are three kinds of strings: `"double quote"`, `"single quote"`, and back-tick quotes.

There is currently no distinction between them, except that each quote format does not need to escape any of the others.
(NOTE: This will likely change, as the language is made more strict.)

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

### Blocks

A block is any group of tokens between a block-start, and a block end. The three kinds are: `( … )`, `[ … ]` and `{ … }`.

Some keywords expect different semantics to a specific block, though often any set of characters can be used to define a block.
(NOTE: This will likely change, as the language is made more strict.)

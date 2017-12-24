# Compact JSON Test suite format

The compact test suite JSON looks like so:

```
{
     "format": "compact",
     "tests": [
        [input, final, left, right, transformed, rebased]
        ...
     ]
}
```

where `inpug` and `final` is a simple strings.  `left`, `right` are
arrays of encoded operations and `transformed` and `rebased` are the
result of transforming `left` against `right`.  That is, applying
`left` and then `transformed` to `input` will result in `final` and
similarly applying `right` and `rebased` to `input` will also result
in `final`.

## Encoding Splices

A splice can be thought of as a string section that is being removed
with another being inserted in its place.  So, if "rob" is being
replaced with "roy", the following encoding will specify that:

```
  "hello (rob=roy)!"
```

The brackets identify the area of the string that is being
modified and the equal separates the before and after.  Note that
either before or after can be empty to indicate insertion and deletion
respectively.

## Encoding Moves

Moves are encoded using brackets to identify the section of
the string that is being moved.  The location where they end up in is
encoded using the `=` sign.  For example if `Bad Big Wolf` is being
fixed up to be `Big Bad Wolf`, this would be encoded as:

```
  "=Bad (Big) Wolf"
  or
  "(Bad )Big =Wolf"
```

## Encoding Ranges

Ranges does not work on strings.  It only works on arrays where each
element itself can be further operated upon. Visual representation of
arrays uses square brackets to indicate the contents of an element.

```
  "oh[el][lo]" represents an array with elements "o", "h", "el" and "lo"
```

Note that characters without square brackets stand for their own
string representations. Note also that it is legitimate to have an
empty array element.

Now we can use the same scheme of regular brackets to identify the
range and the actual operation being applied can be provided to the
right. For example:

```
    "he([ll][lo]=(l=l ))"
```

That is a range operation which applies to the two elements "ll" and
"lo" and the actual operation is represented by `[l=l ]` which is
effectively replacing the first "l" with "l ".  Note that the actual
operation operations on one string -- the actual input string is not
validated.  So, inserting "x" at offset 1 can always be encoded as
`a[x]` without regard to the size of the elements in the range being
encoded.

## Set operations

Sets are represented via curly brackets and commas:
`{key1:value1,key2:value2}`

Set changes are reqpresented via a plus sign followed by the acutal
changes.

```
   {key1:calue1,key2:value2}+{key3:addition,key1:}
```

In the example above, `key1` is getting deleted while `key3` is
getting added.

## Nested paths.

Actual changes can happen in a deep path.  See the following examples:

```
   {Key1:[Big=Bad] Wolf} ===>  splice happens on path = [Key1]
   {Key1:{Inner1:1}+{Inner2:2}} ==> set happens on path = [Key1]
   he({key:a}+{key:b})lo ==> set happens on path = [2]
```
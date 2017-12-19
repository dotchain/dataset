# Compact JSON Test suite format

The compact test suite JSON looks like so:

```json
{
     "format": "compact",
     "tests": [
        [final, changes1, changes2...],
        ...
     ]
}
```

where final is a simple string.  All changes can be either strings
or arrays of strings which use the following encoding scheme to
encoding various changes.

## Encoding Splices

A splice can be thought of as a string section that is being removed
with another being inserted in its place.  So, if "rob" is being
replaced with "roy", the following encoding will specify that:

```
  "hello [rob:roy]!"
```

The square brackets identify the area of the string that is being
modified and the colon separates the before and after.  Note that
either before or after can be empty to indicate insertion and deletion
respectively.

## Encoding Moves

Moves are encoding using square brackets to identify the section of
the string that is being moved.  The location where they end up in is
encoded using the pipe symbol.  For example if `Bad Big Wolf` is being
fixex up to be `Big Bad Wolf`, this would be encoded as:

```
  "|Bad [Big] Wolf"
```


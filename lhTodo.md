# TODOs
See `LH TODO`s in code as well

## Bugs
- Ctrl+Q binding crashes in transbar

## Now
- Support for more transform types
- Diff-ing of transform outputs

## Long term
- Maintain line numbers between transforms

## Considerations
- Would tab groups be a better way to show transforms?
- Syntax highlighting/ maintaining line numbers: This won't always make sense depending on the nature of the transform
	- EG grep may filter out some lines and syntax/line nums make sense
	- grep may also return a count of a word, syntax and line nums would make no sense
	- If a line num prefix was added to each input buf, if that line num prefix was in the outbuf would that be a good sign to maintain line nums and syntax highlighting? Ie. if the lines are intact enough to have the prefix they are intact enough for line nums/syntax highlighting to still make sense?


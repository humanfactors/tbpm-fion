# TBPM E-prime Extraction
My first released Go project. This extracts an E-prime log file and outputs to wide form CSV.

**Note that E-prime exports to UTF-16 but this program only accepts UTF8**. Therefore you will need to manually convert
these (I do have a program for this as well).


## Contributing [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwyl/esta/issues)
To developers, this might contain some pointers as to how to extract E-prime logs using Golang. Nothing fancy - but I'm
happy to accept pull requests or future collaborations to make a more general library.


## Change Log

**Thu 30/11/2017**:

  * Fixed 0's being included in RT
  * Updated to include non-response trial frequency
  * Updated so status based on first number in digit string
  * Added response to the processing chain


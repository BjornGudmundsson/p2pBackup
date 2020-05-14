package files

import (
	"time"
)

//Here I write all of the constants or default values used

//DefaultMaxSize is the default maximum
//amount of bytes a backed up file can be
const DefaultMaxSize = 500

//DefaultMinSize is the default
//minimum amount of bytes a file can be
const DefaultMinSize = 0

//DefaultExcludedRule is a regexp
//that says what files should not be included
const DefaultExcludedRule = ""

//DefaultMaxElapsedTime is the default maximum
//amount of time that can pass before a file can
//be considered for backing up
const DefaultMaxElapsedTime = time.Second

const metadatasize = 2 * KEYLEN
const compress = true

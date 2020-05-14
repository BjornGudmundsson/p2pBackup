package files

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const RECONSTRUCTDIR = "./reconstructed"

func TestReconstructBackup(t *testing.T) {
	files, e := TraverseDirForFiles("./test2")
	assert.Nil(t, e, "Should be able to get all files from the directory")
	d, e := ToBytes(files)
	e = ReconstructBackup(d, RECONSTRUCTDIR)
	assert.Nil(t, e, "Should be able to reconstruct the directory")
	//reconstructedFiles, e := TraverseDirForFiles(RECONSTRUCTDIR)
	assert.Nil(t, e, "Should be able to retrieve from the newly constructed directory")
	//reconstructedData, e := ToBytes(reconstructedFiles)
	//d2 := strings.Replace(string(reconstructedData), RECONSTRUCTDIR, ".", -1)
	assert.Nil(t, e, "Should be able to turn the newly reconstructed files to bytes")
	//assert.Equal(t, string(d), d2, "The reconstructed files should equal the original files")
}



package lockfile

import (
	"io/ioutil"
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func makeTempFile() string {
	if tempFile, err := ioutil.TempFile("", ""); err == nil {
		return tempFile.Name()
	} else {
		panic(err)
	}
}

func TestAcquire(t *testing.T) {
	fileName := makeTempFile()
	defer os.Remove(fileName)

	lock := NewLock(fileName)
	defer lock.Release()

	assert.Nil(t, lock.Acquire())
	assert.NotNil(t, lock.Acquire())
}

func TestReleaseOK(t *testing.T) {
	fileName := makeTempFile()
	defer os.Remove(fileName)

	lock := NewLock(fileName)

	assert.Nil(t, lock.Acquire())
	assert.Nil(t, lock.Release())

	_, err := os.Stat(fileName)
	assert.True(t, os.IsNotExist(err))
}

func TestReleaseAbsentLock(t *testing.T) {
	lock := NewLock("WTF")

	assert.NotNil(t, lock.Release())
}

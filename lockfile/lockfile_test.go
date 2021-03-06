package lockfile

import (
	"io/ioutil"
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func makeTempFile() string {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

	return tempFile.Name()
}

func TestAcquire(t *testing.T) {
	fileName := makeTempFile()
	defer os.Remove(fileName)

	lock := NewLock(fileName)
	defer lock.Release()

	assert.Nil(t, lock.Acquire())
	assert.NotNil(t, lock.Acquire())
}

func TestAcquireNotCreated(t *testing.T) {
	name := "TestAcquireNotCreated"

	lock := NewLock(name)
	defer os.Remove(name)
	defer lock.finish()

	assert.Nil(t, lock.Acquire())
	_, err := os.Stat(name)
	assert.Nil(t, err)

	assert.NotNil(t, lock.Acquire())
	_, err = os.Stat(name)
	assert.Nil(t, err)

	assert.Nil(t, lock.Release())
	_, err = os.Stat(name)
	assert.True(t, os.IsNotExist(err))
}

func TestReleaseOK(t *testing.T) {
	fileName := makeTempFile()
	defer os.Remove(fileName)

	lock := NewLock(fileName)
	defer lock.finish()

	assert.Nil(t, lock.Acquire())
	assert.Nil(t, lock.Release())

	_, err := os.Stat(fileName)
	assert.Nil(t, err)
}

func TestReleaseAbsentLock(t *testing.T) {
	lock := NewLock("WTF")

	assert.NotNil(t, lock.Release())
}

func TestLockFileIsNotHarmuful(t *testing.T) {
	content := []byte("content")
	fileName := makeTempFile()
	defer os.Remove(fileName)

	ioutil.WriteFile(fileName, content, os.FileMode(0666))

	lock := NewLock(fileName)
	defer lock.finish()

	lock.Acquire()

	readContent, err := ioutil.ReadFile(fileName)
	assert.Nil(t, err)
	assert.Equal(t, content, readContent)

	lock.Release()

	readContent, err = ioutil.ReadFile(fileName)
	assert.Nil(t, err)
	assert.Equal(t, content, readContent)
}

func TestCannotAcquireWithWrongPermissions(t *testing.T) {
	fileName := makeTempFile()
	defer os.Remove(fileName)

	os.Chmod(fileName, os.FileMode(0200))

	lock := NewLock(fileName)
	assert.NotNil(t, lock.Acquire())
}

func TestLockDirectory(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)

	lock := NewLock(tempDir)
	assert.Nil(t, lock.Acquire())
	assert.Nil(t, lock.Release())
}

func TestStringer(t *testing.T) {
	assert.True(t, NewLock("").String() != "")
}

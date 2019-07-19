package helpers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	REMOVE_ATTEMPTS int = 20
	MV_ATTEMPTS     int = 20
)

func IsPathExists(path string) (bool, error) {
	if _, e := os.Stat(path); e != nil {
		if os.IsNotExist(e) {
			return false, nil
		} else {
			return true, e
		}
	}
	return true, nil
}

func MkdirAll(path string, perm os.FileMode) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{"mkdir", path, syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = MkdirAll(path[0:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = os.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}

	err = os.Chmod(path, perm)
	if err != nil {
		return err
	}

	return nil
}

func PathEnsure(path string, perm os.FileMode) error {
	if ok, e := IsPathExists(path); e != nil {
		return e
	} else if !ok {
		if e = MkdirAll(path, perm); e != nil {
			return e
		}
	}

	return nil
}

func PathRemove(path string) (e error) {
	attempt := REMOVE_ATTEMPTS

	for attempt > 0 {
		e = os.RemoveAll(path)
		if e != nil {
			attempt -= 1
			continue
		}

		return nil
	}

	return e
}

func PathMV(oldpath, newpath string) (e error) {
	attempt := MV_ATTEMPTS

	for attempt > 0 {
		e = os.Rename(oldpath, newpath)
		if e != nil {
			attempt -= 1
			continue
		}

		return nil
	}

	return e
}

func PathIsDirEmpty(name string) (bool, error) {
	f, e := os.Open(name)
	if e != nil {
		return false, e
	}
	defer f.Close()

	// read in ONLY one file
	_, e = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if e == io.EOF {
		return true, nil
	}
	return false, e
}

func PathsAbsolutize(paths []*string) error {
	for _, p := range paths {
		if v, e := filepath.Abs(*p); e != nil {
			return e
		} else {
			*p = v
		}
	}

	return nil
}

type pathAndMode struct {
	p    string
	mode os.FileMode
}

type PathAndModes []pathAndMode

func (self *PathAndModes) Append(path string, mode os.FileMode) {
	*self = append(*self, pathAndMode{path, mode})
}

func NewPathAndModes() PathAndModes { return PathAndModes{} }

func PathsEnsure(paths PathAndModes) error {
	for _, r := range paths {
		if _, e := os.Stat(r.p); os.IsNotExist(e) {
			if e = os.MkdirAll(r.p, r.mode); e != nil {
				return e
			}
		} else if e != nil {
			return e
		}
	}

	return nil
}

func PathsAbsent(paths []string) (err error) {
	for _, path := range paths {
		if e := PathRemove(path); e != nil {
			err = e
		}
	}
	return
}

func PathClear(path string) (e error) {
	if ok, e := IsPathExists(path); e != nil {
		return e
	} else if !ok {
		return nil
	}

	fs, e := ioutil.ReadDir(path)
	if e != nil {
		return e
	}

	for _, f := range fs {
		e = PathRemove(filepath.Join(path, f.Name()))
	}
	return
}

func PathToFile(path string, mode os.FileMode, reader io.Reader) error {
	dir := filepath.Dir(path)
	if e := PathEnsure(dir, mode); e != nil {
		return e
	}

	f, e := os.Create(path)
	if e != nil {
		return e
	}
	defer f.Close()

	if _, e := io.Copy(f, reader); e != nil {
		return e
	}

	return nil
}

func PathWriteFile(path string, dirMode, fileMode os.FileMode, src io.Reader) (int64, error) {
	dir := filepath.Dir(path)
	if e := PathEnsure(dir, dirMode); e != nil {
		return 0, e
	}

	f, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if e != nil {
		return 0, e
	}
	defer f.Close()

	sz, e := io.Copy(f, src)
	if e != nil {
		return 0, e
	}

	return sz, nil
}

func PathRename(oldpath, newpath string, mode os.FileMode) error {
	if e := PathEnsure(filepath.Dir(newpath), mode); e != nil {
		return nil
	}

	return os.Rename(oldpath, newpath)
}

func PathCopyFile(source string, dest string, mode os.FileMode) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {

		if mode != 0 { // custom mode
			err = os.Chmod(dest, mode)
		} else { // copy mode from source file
			sourceinfo, err := os.Stat(source)
			if err != nil {
				return err
			} else {
				err = os.Chmod(dest, sourceinfo.Mode())
			}
		}
	}

	return
}

func PathCopyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = PathCopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = PathCopyFile(sourcefilepointer, destinationfilepointer, 0)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func PathCopyFileEnsureDir(src string, dst string, fileMode os.FileMode, dirMode os.FileMode) error {
	dirPath := filepath.Dir(dst)
	if e := PathEnsure(dirPath, dirMode); e != nil {
		return e
	}

	if e := PathCopyFile(src, dst, fileMode); e != nil {
		return e
	}

	return nil
}

func PathRemoveFileOKEvenIfNotExists(path string) (e error) {
	if e := os.Remove(path); e != nil {
		if pathErr, ok := e.(*os.PathError); ok {
			if pathErr.Err.Error() == "no such file or directory" {
				return nil
			} else {
				return e
			}
		}
	}
	return nil
}

type PathMD5 struct {
	Latency time.Duration // calculation latency
	Sum     [md5.Size]byte
	Sz      uint64
	B       *bytes.Buffer
}

var pathMD5Pool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		buf := new(bytes.Buffer)
		return buf
	},
}

func (it *PathMD5) Release() { pathMD5Pool.Put(it.B) }

func (it *PathMD5) Do(path string) error {
	start := time.Now()
	defer func() { it.Latency = time.Since(start) }()

	if fi, e := os.Stat(path); e != nil {
		return e
	} else {
		it.Sz = uint64(fi.Size())
	}

	data, e := ioutil.ReadFile(path)
	if e != nil {
		return nil
	}

	it.Sum = md5.Sum(data)

	b := it.B
	if b == nil {
		b = pathMD5Pool.Get().(*bytes.Buffer)
	}

	b.Reset()
	b.Write(data)

	it.B = b

	return nil
}

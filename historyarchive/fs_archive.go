// Copyright 2016 Diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type FsArchiveBackend struct {
	prefix string
}

func (b *FsArchiveBackend) GetFile(pth string) (io.ReadCloser, error) {
	return os.Open(path.Join(b.prefix, pth))
}

func (b *FsArchiveBackend) Exists(pth string) (bool, error) {
	pth = path.Join(b.prefix, pth)
	log.WithField("path", pth).Trace("fs: check exists")
	_, err := os.Stat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			log.WithField("path", pth).WithField("exists", false).Trace("fs: check exists")
			return false, nil
		} else {
			log.WithField("path", pth).WithError(err).Error("fs: check exists")
			return false, err
		}
	}
	log.WithField("path", pth).WithField("exists", true).Trace("fs: check exists")
	return true, nil
}

func (b *FsArchiveBackend) Size(pth string) (int64, error) {
	pth = path.Join(b.prefix, pth)
	log.WithField("path", pth).Trace("fs: get size")
	fi, err := os.Stat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			log.WithField("path", pth).WithError(err).Warn("fs: get size")
			return 0, nil
		} else {
			log.WithField("path", pth).WithError(err).Error("fs: get size")
			return 0, err
		}
	}
	log.WithField("path", pth).WithField("size", fi.Size()).Trace("fs: got size")
	return fi.Size(), nil
}

func (b *FsArchiveBackend) PutFile(pth string, in io.ReadCloser) error {
	dir := path.Join(b.prefix, path.Dir(pth))
	log.WithField("path", pth).Trace("fs: put file")
	exists, err := b.Exists(dir)
	if err != nil {
		log.WithField("path", pth).WithError(err).Error("fs: put file (check exists)")
		return err
	}

	if !exists {
		if e := os.MkdirAll(dir, 0755); e != nil {
			log.WithField("path", pth).WithError(err).Error("fs: put file (mkdir)")
			return e
		}
	}

	pth = path.Join(b.prefix, pth)
	out, e := os.Create(pth)
	if e != nil {
		log.WithField("path", pth).WithError(err).Error("fs: put file (create)")
		return e
	}
	defer in.Close()
	defer out.Close()
	_, e = io.Copy(out, in)
	if e != nil {
		log.WithField("path", pth).WithError(err).Error("fs: put file (copy)")
	}
	return e
}

func (b *FsArchiveBackend) ListFiles(pth string) (chan string, chan error) {
	ch := make(chan string)
	errs := make(chan error)
	go func() {
		log.WithField("path", pth).Trace("fs: list files")
		exists, err := b.Exists(pth)
		if err != nil {
			errs <- err
			return
		}
		if exists {
			filepath.Walk(path.Join(b.prefix, pth),
				func(p string, info os.FileInfo, err error) error {
					if err != nil {
						log.WithField("path", pth).WithError(err).Error("fs: list files (walk)")
						errs <- err
						return nil
					}
					if info != nil && !info.IsDir() {
						log.WithField("found", p).Trace("fs: list files (walk)")
						ch <- p
					}
					return nil
				})
		}
		close(ch)
		close(errs)
	}()
	return ch, errs
}

func (b *FsArchiveBackend) CanListFiles() bool {
	return true
}

func makeFsBackend(pth string, opts ConnectOptions) ArchiveBackend {
	return &FsArchiveBackend{
		prefix: pth,
	}
}

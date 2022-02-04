// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
package main

import (
	"archive/zip"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

// Encapsulates state and operations for a virtual file inside zip archive on HDFS file system
type ZipFile struct {
	Attrs      Attrs
	zipFile    *zip.File
	FileSystem *FileSystem
}

// Verify that *Dir implements necesary FUSE interfaces
var _ fs.Node = (*ZipFile)(nil)
var _ fs.NodeOpener = (*ZipFile)(nil)

// Responds on FUSE Attr request to retrieve file attributes
func (zipfile *ZipFile) Attr(ctx context.Context, fuseAttr *fuse.Attr) error {
	return zipfile.Attrs.ConvertAttrToFuse(fuseAttr)
}

// Responds on FUSE Open request for a file inside zip archive
func (zipfile *ZipFile) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	contentStream, err := zipfile.zipFile.Open()
	if err != nil {
		logerror("Opened zip file failed", Fields{Operation: OpenArch, Archive: zipfile.Attrs.Name, Error: err})
		return nil, err
	}
	// reporting to FUSE that the stream isn't seekable
	resp.Flags |= fuse.OpenNonSeekable
	return NewZipFileHandle(contentStream), nil
}

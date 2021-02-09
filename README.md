# embeddedfs #
Embed file system in the source code, and use the file system in memory.

**the new “embed” package which allows you to embed a file contents as part of the Go application binary, in Go 1.16 Release.**
please use "embed" to replace embeddedfs package.

## Install ##

	go get -u github.com/whiler/embeddedfs

## Usage ##
### quick start ###
#### demo.go ####

	package main
	
	import (
		"io/ioutil"
		"log"
	
		"github.com/whiler/embeddedfs"
	)
	
	func main() {
		var content = []byte("hello embeddedfs.")
		var fs embeddedfs.FileSystem = embeddedfs.EmbeddedFileSystem(map[string]*embeddedfs.EmbeddedFile{
			"/index": &embeddedfs.EmbeddedFile{
				Info: &embeddedfs.FileInfo{
					RawName: "index",
					RawSize: int64(len(content)),
				},
				Content: content,
			},
		})
	
		file, err := fs.Open("/index")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("content: %s\n", string(bs))
	}

### use embeddedfs.EmbeddedFileSystem as http.FileSystem ###

	var fs embeddedfs.FileSystem = embeddedfs.EmbeddedFileSystem(map[string]*embeddedfs.EmbeddedFile{})
	http.ListenAndServe(":8080", http.FileServer(fs))

### generate embeddedfs.EmbeddedFileSystem from local file system ###
#### install embeddedfs command ####

	go install github.com/whiler/embeddedfs/cmd/embeddedfs

#### generate embeddedfs.EmbeddedFileSystem ####

	${GOPATH}/bin/embeddedfs -name ImageFS -out resource/images.go -package resource -root ./images

## embeddedfs.FileSystem ##
the `embeddedfs.FileSystem` is an interface.

	type FileSystem interface {
		Open(name string) (http.File, error)
		Stat(name string) (os.FileInfo, error)
	}

### Operations ###
| Operation | Original     | embeddedfs.DefaultFileSystem                                   | embeddedfs.EmbeddedFileSystem                                     |
| --------- | ------------ | -------------------------------------------------------------- | ----------------------------------------------------------------- |
| -         |              | var fs embeddedfs.FileSystem = &embeddedfs.DefaultFileSystem{} | var fs embeddedfs.FileSystem = embeddedfs.EmbeddedFileSystem(...) |
| Open      | os.Open(...) | fs.Open(...)                                                   | fs.Open(...)                                                      |
| Stat      | os.Stat(...) | fs.Stat(...)                                                   | fs.Stat(...)                                                      |


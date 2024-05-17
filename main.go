package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}

type appArgs struct {
	local    string
	prefixes []string

	listOnly bool
	write    bool
	diffOnly bool
}

func run() error {
	var args appArgs

	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage: gocomments [flags] [path ...]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.StringVar(&args.local, "local", "", "put imports beginning with this string after 3rd-party package")
	flag.Var((*ArrayStringFlag)(&args.prefixes), "prefix", "relative local prefix to from a new import group (can be given several times)")

	flag.BoolVar(&args.listOnly, "l", false, "list files whose formatting differs from goimport's")
	flag.BoolVar(&args.write, "w", false, "write result to (source) file instead of stdout")
	flag.BoolVar(&args.diffOnly, "d", false, "display diffs instead of rewriting files")

	flag.Parse()

	return process(&args, flag.Args()...)
}

type fileSource uint8

const (
	fileSourceStdin fileSource = iota
	fileSourceFilepath
)

func process(args *appArgs, paths ...string) error {
	cache := NewCommentConfigCache(args.local, args.prefixes)

	if len(paths) == 0 {
		return processFile(cache, "<standard input>", fileSourceStdin, os.Stdin, os.Stdout, args)
	}

	for _, path := range paths {
		switch dir, err := os.Stat(path); {
		case err != nil:
			return err
		case dir.IsDir():
			return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
				return visitFile(cache, args, path, d, err)
			})
		default:
			if err := processFile(cache, path, fileSourceFilepath, nil, os.Stdout, args); err != nil {
				return err
			}
		}
	}

	return nil
}

func processFile(cache *CommentConfigCache, filename string, source fileSource, in io.Reader, out io.Writer, args *appArgs) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		defer func() {
			_ = f.Close()
		}()

		in = f
	}

	src, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := processComments(filename, src, cache)
	if err != nil {
		return err
	}

	if !bytes.Equal(src, res) {
		if args.listOnly {
			_, _ = fmt.Fprintln(out, filename)
		}
		if args.write {
			if source == fileSourceStdin {
				return errors.New("can't use -w on stdin")
			}
			return os.WriteFile(filename, res, 0o644)
		}

		if args.diffOnly {
			if source == fileSourceStdin {
				filename = "stdin.go" // because <standard input>.orig looks silly
			}

			data, err := diff(src, res, filename)
			if err != nil {
				return fmt.Errorf("computing diff: %v", err)
			}

			_, _ = out.Write(data)
		}
	}

	if !args.listOnly && !args.write && !args.diffOnly {
		if _, err := out.Write(res); err != nil {
			return err
		}
	}

	return nil
}

func visitFile(cache *CommentConfigCache, args *appArgs, path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	name := d.Name()

	if d.IsDir() {
		if name == "vendor" {
			return fs.SkipDir
		}

		return nil
	}

	if strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".go") {
		return nil
	}

	return processFile(cache, path, fileSourceFilepath, nil, os.Stdout, args)
}

func diff(b1, b2 []byte, filename string) ([]byte, error) {
	f1, err := writeTempFile("", "gocomments", b1)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = os.Remove(f1)
	}()

	f2, err := writeTempFile("", "gocomments", b2)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = os.Remove(f2)
	}()

	data, err := exec.Command("diff", "-u", f1, f2).CombinedOutput()
	if len(data) >= 0 {
		data, err = replaceTempFilename(data, filename)
	}
	return data, err
}

func writeTempFile(dir string, prefix string, data []byte) (string, error) {
	f, err := os.CreateTemp(dir, prefix)
	if err != nil {
		return "", err
	}

	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}

	if err != nil {
		_ = os.Remove(f.Name())
		return "", err
	}

	return f.Name(), nil
}

// replaceTempFilename replaces temporary filenames in diff with actual one.
//
// --- /tmp/gofmt316145376	2017-02-03 19:13:00.280468375 -0500
// +++ /tmp/gofmt617882815	2017-02-03 19:13:00.280468375 -0500
// ...
// ->
// diff -u path/to/file.go.orig path/to/file.go
// --- path/to/file.go.orig	2017-02-03 19:13:00.280468375 -0500
// +++ path/to/file.go	2017-02-03 19:13:00.280468375 -0500
// ...
func replaceTempFilename(diff []byte, filename string) ([]byte, error) {
	bs := bytes.SplitN(diff, []byte{'\n'}, 3)
	if len(bs) < 3 {
		return nil, fmt.Errorf("got unexpected diff for %s", filename)
	}

	// Preserve timestamps.
	var t0, t1 []byte
	if i := bytes.LastIndexByte(bs[0], '\t'); i != -1 {
		t0 = bs[0][i:]
	}
	if i := bytes.LastIndexByte(bs[1], '\t'); i != -1 {
		t1 = bs[1][i:]
	}

	// Always print filepath with slash separator.
	f := filepath.ToSlash(filename)
	bs[0] = []byte(fmt.Sprintf("--- %s%s", f+".orig", t0))
	bs[1] = []byte(fmt.Sprintf("+++ %s%s", f, t1))

	// Insert diff header.
	header := fmt.Sprintf("diff -u %s %s", f+".orig", f)
	bs = append([][]byte{[]byte(header)}, bs...)

	return bytes.Join(bs, []byte{'\n'}), nil
}

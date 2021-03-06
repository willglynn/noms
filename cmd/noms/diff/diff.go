// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package diff

import (
	"io"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"
	"github.com/dustin/go-humanize"
)

func shouldDescend(v1, v2 types.Value) bool {
	kind := v1.Type().Kind()
	return !types.IsPrimitiveKind(kind) && kind == v2.Type().Kind() && kind != types.RefKind
}

func Diff(w io.Writer, v1, v2 types.Value) error {
	return d.Try(func() {
		diff(w, types.NewPath(), nil, v1, v2)
	})
}

func diff(w io.Writer, p types.Path, key, v1, v2 types.Value) {
	if !v1.Equals(v2) {
		if shouldDescend(v1, v2) {
			switch v1.Type().Kind() {
			case types.ListKind:
				diffLists(w, p, v1.(types.List), v2.(types.List))
			case types.MapKind:
				diffMaps(w, p, v1.(types.Map), v2.(types.Map))
			case types.SetKind:
				diffSets(w, p, v1.(types.Set), v2.(types.Set))
			case types.StructKind:
				diffStructs(w, p, v1.(types.Struct), v2.(types.Struct))
			default:
				panic("Unrecognized type in diff function")
			}
		} else {
			line(w, DEL, key, v1)
			line(w, ADD, key, v2)
		}
	}
}

func diffLists(w io.Writer, p types.Path, v1, v2 types.List) {
	spliceChan := make(chan types.Splice)
	stopChan := make(chan struct{}, 1) // buffer size of 1, so this won't block if diff already finished

	defer stop(stopChan)
	go func() {
		v2.Diff(v1, spliceChan, stopChan)
		close(spliceChan)
	}()

	wroteHeader := false

	for splice := range spliceChan {
		if splice.SpRemoved == splice.SpAdded {
			for i := uint64(0); i < splice.SpRemoved; i++ {
				lastEl := v1.Get(splice.SpAt + i)
				newEl := v2.Get(splice.SpFrom + i)
				if shouldDescend(lastEl, newEl) {
					idx := types.Number(splice.SpAt + i)
					diff(w, p.AddIndex(idx), idx, lastEl, newEl)
				} else {
					wroteHeader = writeHeader(w, wroteHeader, p)
					line(w, DEL, nil, v1.Get(splice.SpAt+i))
					line(w, ADD, nil, v2.Get(splice.SpFrom+i))
				}
			}
		} else {
			for i := uint64(0); i < splice.SpRemoved; i++ {
				wroteHeader = writeHeader(w, wroteHeader, p)
				line(w, DEL, nil, v1.Get(splice.SpAt+i))
			}
			for i := uint64(0); i < splice.SpAdded; i++ {
				wroteHeader = writeHeader(w, wroteHeader, p)
				line(w, ADD, nil, v2.Get(splice.SpFrom+i))
			}
		}
	}

	writeFooter(w, wroteHeader)
}

func diffMaps(w io.Writer, p types.Path, v1, v2 types.Map) {
	changeChan := make(chan types.ValueChanged)
	stopChan := make(chan struct{}, 1) // buffer size of 1, so this won't block if diff already finished

	defer stop(stopChan)
	go func() {
		v2.Diff(v1, changeChan, stopChan)
		close(changeChan)
	}()

	wroteHeader := false

	for change := range changeChan {
		switch change.ChangeType {
		case types.DiffChangeAdded:
			wroteHeader = writeHeader(w, wroteHeader, p)
			line(w, ADD, change.V, v2.Get(change.V))
		case types.DiffChangeRemoved:
			wroteHeader = writeHeader(w, wroteHeader, p)
			line(w, DEL, change.V, v1.Get(change.V))
		case types.DiffChangeModified:
			c1, c2 := v1.Get(change.V), v2.Get(change.V)
			if shouldDescend(c1, c2) {
				wroteHeader = writeFooter(w, wroteHeader)
				diff(w, p.AddIndex(change.V), change.V, c1, c2)
			} else {
				wroteHeader = writeHeader(w, wroteHeader, p)
				line(w, DEL, change.V, c1)
				line(w, ADD, change.V, c2)
			}
		default:
			panic("unknown change type")
		}
	}

	writeFooter(w, wroteHeader)
}

func diffStructs(w io.Writer, p types.Path, v1, v2 types.Struct) {
	changeChan := make(chan types.ValueChanged)
	stopChan := make(chan struct{}, 1) // buffer size of 1, so this won't block if diff already finished

	defer stop(stopChan)
	go func() {
		v2.Diff(v1, changeChan, stopChan)
		close(changeChan)
	}()

	wroteHeader := false

	for change := range changeChan {
		fn := string(change.V.(types.String))
		switch change.ChangeType {
		case types.DiffChangeAdded:
			wroteHeader = writeHeader(w, wroteHeader, p)
			field(w, ADD, change.V, v2.Get(fn))
		case types.DiffChangeRemoved:
			wroteHeader = writeHeader(w, wroteHeader, p)
			field(w, DEL, change.V, v1.Get(fn))
		case types.DiffChangeModified:
			f1 := v1.Get(fn)
			f2 := v2.Get(fn)
			if shouldDescend(f1, f2) {
				diff(w, p.AddField(fn), types.String(fn), f1, f2)
			} else {
				wroteHeader = writeHeader(w, wroteHeader, p)
				field(w, DEL, change.V, f1)
				field(w, ADD, change.V, f2)
			}
		}
	}
	writeFooter(w, wroteHeader)
}

func diffSets(w io.Writer, p types.Path, v1, v2 types.Set) {
	changeChan := make(chan types.ValueChanged)
	stopChan := make(chan struct{}, 1) // buffer size of 1, so this won't block if diff already finished

	defer stop(stopChan)
	go func() {
		v2.Diff(v1, changeChan, stopChan)
		close(changeChan)
	}()

	wroteHeader := false

	for change := range changeChan {
		wroteHeader = writeHeader(w, wroteHeader, p)
		switch change.ChangeType {
		case types.DiffChangeAdded:
			line(w, ADD, nil, change.V)
		case types.DiffChangeRemoved:
			line(w, DEL, nil, change.V)
		default:
			panic("unknown change type")
		}
	}

	writeFooter(w, wroteHeader)
}

func line(w io.Writer, op prefixOp, key, val types.Value) {
	pw := newPrefixWriter(w, op)
	if key != nil {
		writeEncodedValue(pw, key)
		write(w, []byte(": "))
	}
	writeEncodedValue(pw, val)
	write(w, []byte("\n"))
}

func field(w io.Writer, op prefixOp, name, val types.Value) {
	pw := newPrefixWriter(w, op)
	write(pw, []byte(name.(types.String)))
	write(w, []byte(": "))
	writeEncodedValue(pw, val)
	write(w, []byte("\n"))
}

func writeHeader(w io.Writer, wroteHeader bool, p types.Path) bool {
	if !wroteHeader {
		if len(p) == 0 {
			write(w, []byte("(root)"))
		} else {
			write(w, []byte(p.String()))
		}
		write(w, []byte(" {\n"))
	}
	return true
}

func writeFooter(w io.Writer, wroteHeader bool) bool {
	if wroteHeader {
		write(w, []byte("  }\n"))
	}
	return false
}

func write(w io.Writer, b []byte) {
	_, err := w.Write(b)
	d.PanicIfError(err)
}

func writeEncodedValue(w io.Writer, v types.Value) {
	if v.Type().Kind() == types.BlobKind {
		w.Write([]byte("Blob ("))
		w.Write([]byte(humanize.Bytes(v.(types.Blob).Len())))
		w.Write([]byte(")"))
	} else {
		d.PanicIfError(types.WriteEncodedValue(w, v))
	}
}

func writeEncodedValueWithTags(w io.Writer, v types.Value) {
	d.PanicIfError(types.WriteEncodedValueWithTags(w, v))
}

func stop(ch chan<- struct{}) {
	ch <- struct{}{}
}

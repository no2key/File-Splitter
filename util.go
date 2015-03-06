// copy_md5
package main

import (
	"fmt"
	"hash"
	"io"
)

func copyAndHash(w io.Writer, r io.Reader, h hash.Hash) (string, error) {
	buf := make([]byte, 32*1024)
	n := 0
	var er, ew error

	for {
		n, er = r.Read(buf)
		if n > 0 {
			_, ew = w.Write(buf[:n])
			if ew != nil {
				return "", ew
			}

			_, ew = h.Write(buf[:n])
			if ew != nil {
				return "", ew
			}
		}

		if er == io.EOF {
			break
		} else if er != nil {
			return "", er
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

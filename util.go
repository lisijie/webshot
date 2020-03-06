package main

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"image/png"
	"io"
)

// width, height 表示调整后的宽高，当其中一个设为0时，表示等比缩放
func ImageResize(src io.Reader, width, height uint) ([]byte, error) {
	img, err := png.Decode(src)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var dst bytes.Buffer
	m := resize.Resize(width, height, img, resize.Lanczos3)
	if err := png.Encode(&dst, m); err != nil {
		return nil, errors.WithStack(err)
	}
	return dst.Bytes(), nil
}

func SizeFormat(size float64) string {
	units := []string{"Byte", "KB", "MB", "GB", "TB"}
	n := 0
	for size >= 1024 {
		size /= 1024
		n += 1
	}
	return fmt.Sprintf("%.2f %s", size, units[n])
}

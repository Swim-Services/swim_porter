package compressor

import (
	"bytes"
	"encoding/json"
	"image/png"
	"maps"
	"path"

	"github.com/gameparrot/fastpng"
	"github.com/klauspost/compress/flate"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/jsonnewline"
	"github.com/swim-services/swim_porter/utils"
	stripjsoncomments "github.com/trapcodeio/go-strip-json-comments"
)

func Compress(in []byte) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	fs := utils.NewMapFS(zipMap)
	CompressRaw(fs)
	out, err := utils.ZipCompressionLevel(fs.RawMap(), flate.BestCompression)
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func CompressRaw(in *utils.MapFS) {
	internal.ParallelMap(maps.Clone(in.RawMap()), func(file string, data []byte) {
		ext := path.Ext(file)
		switch ext {
		case ".png":
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				return
			}
			e := fastpng.Encoder{CompressionLevel: fastpng.BestCompression}
			buf := bytes.NewBuffer([]byte{})
			if e.Encode(buf, img) != nil {
				return
			}
			in.Write(buf.Bytes(), file)
		case ".json", ".mcmeta":
			jsonStr := string(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")))
			commentStripped := stripjsoncomments.Strip(string(jsonStr))
			newlineFixed := jsonnewline.NewLineToEscape(commentStripped)
			buf := bytes.NewBuffer([]byte{})
			if json.Compact(buf, []byte(newlineFixed)) != nil {
				return
			}
			in.Write(buf.Bytes(), file)
		}
	})
}

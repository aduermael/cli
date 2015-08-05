package chrootarchive

import "github.com/docker/docker/pkg/archive"

// ApplyLayer parses a diff in the standard layer format from `layer`,
// and applies it to the directory `dest`. The stream `layer` can only be
// uncompressed.
// Returns the size in bytes of the contents of the layer.
func ApplyLayer(dest string, layer archive.Reader) (size int64, err error) {
	return applyLayerHandler(dest, layer, true)
}

// ApplyUncompressedLayer parses a diff in the standard layer format from
// `layer`, and applies it to the directory `dest`. The stream `layer`
// can only be uncompressed.
// Returns the size in bytes of the contents of the layer.
func ApplyUncompressedLayer(dest string, layer archive.Reader) (int64, error) {
	return applyLayerHandler(dest, layer, false)
}

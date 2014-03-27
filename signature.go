package magick

// #include <magick/api.h>
import "C"

// Signature computes a message digest from an image pixel stream with an
// implementation of the NIST SHA-256 Message Digest algorithm. This
// signature uniquely identifies the image and is convenient for determining
// if an image has been modified or whether two images are identical.
func (im *Image) Signature() string {
	if !im.HasProperty("signature") {
		C.SignatureImage(im.image)
	}
	return im.Property("signature")
}

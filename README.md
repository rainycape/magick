magick image library for Go (golang)
===================================

ImageMagick bindings for Go (golang)

Requires Go 1.2 (or a 1.2 release candidate) due to C function
pointer support.

Supports both ImageMagick and GraphicsMagick. The former
is used by default, while the latter can be enabled by
building the package with the gm build tag.

Installation requires either ImageMagick or GraphicsMagick and libgif, with their
respective development packages.

Then, install this package by running.

```
    go get gopkgs.com/magick.v1
```

For documentation and available versions, see http://gopkgs.com/magick.
Some functions are not documented. For those, see the MagickCore documentation
at http://www.imagemagick.org/script/magick-core.php.

package utills

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func ResizeImage(f *os.File, width uint) []byte {
    //encoding message is discarded, because OP wanted only jpg, else use encoding in resize function

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

    m := resize.Resize(width, 0, img, resize.Lanczos3)
	
	return imgToBytes(m)
}


func imgToBytes(img image.Image) []byte {
    var opt jpeg.Options
    opt.Quality = 80

    buff := bytes.NewBuffer(nil)
    err := jpeg.Encode(buff, img, &opt)
    if err != nil {
        log.Fatal(err)
    }

    return buff.Bytes()
}

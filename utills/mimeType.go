package utills

import "mime/multipart"

// Check if file format is supported by our API.
func IsFileFormatSupported(fileHandler *multipart.FileHeader) bool {
	switch mimeType := fileHandler.Header.Get("Content-Type") ; mimeType{
	case "image/jpeg":
		return true
	case "image/png":
		return true	
	case "image/tiff":
			return true
	default :
		return false
	}
}
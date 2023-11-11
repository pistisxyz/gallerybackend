package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var ImageFormats = []string{".jpeg", ".jpg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".ico", ".webp", ".heif", ".heic", ".svg", ".psd", ".eps", ".ai", ".pdf", ".jfif", ".jp2", ".pbm", ".pgm", ".ppm", ".hdr", ".exr", ".tga", ".dds", ".wmf", ".emf"}

var VideoExtensions = []string{".mp4", ".mov", ".avi", ".wmv", ".mkv", ".flv", ".webm", ".mpeg", ".mpg", ".3gp", ".m4v", ".ogg", ".ogv", ".qt", ".vob", ".swf", ".mp2", ".mpe", ".asf", ".m2v", ".divx", ".rm", ".rmvb", ".mpg4", ".3g2", ".m2ts", ".mts", ".ts", ".wmx", ".xvid", ".f4v", ".mov", ".m2p", ".dvr-ms", ".mxf", ".mpg2", ".mpeg1", ".mpeg2", ".mpeg4", ".m1v", ".vro", ".dat", ".amv", ".bik", ".csf", ".dav", ".dce", ".dpg", ".dvr", ".dzm", ".eye", ".f4p", ".fbr", ".gbx", ".grasp", ".ivf", ".m21", ".mjp", ".mks", ".mod", ".moov", ".mqv", ".mts", ".nsv", ".nuv", ".pva", ".qt", ".ratdvd", ".rav", ".roq", ".smk", ".spl", ".ssm", ".svi", ".tivo", ".tod", ".tp", ".trp", ".vdr", ".vfw", ".vid", ".vmw", ".vp3", ".weba", ".xlmv"}

func CatchErr(err error, msg ...string) bool {
	if err != nil {
		if len(msg) > 0 {
			fmt.Println(msg[0], err)
		} else {
			fmt.Println(err)
		}
		return true
	}
	return false
}

func RemoveDuplicates(arr []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range arr {
		if encountered[v] {
			continue
		} else {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func GenerateUniqueFilename() string {
	return uuid.New().String()
}

func ContainsString(arr []string, target string) bool {
	for _, s := range arr {
		if strings.EqualFold(s, target) {
			return true
		}
	}
	return false
}

func MakeMediaThumbnail(file_location string, file_type string) {

	width := "250"
	height := "250"

	var cmdStruct *exec.Cmd
	switch file_type {
	case "image":
		cmdStruct = exec.Command("ffmpeg", "-i", file_location, "-vf", "scale='if(gt(a,1),"+width+",-1)':'if(gt(a,1),-1,"+height+")'", "-c:v", "libwebp", "-q:v", "80", "compressed/"+RemoveExtension(path.Base(file_location))+".webp")
	case "video":
		cmdStruct = exec.Command("ffmpeg", "-i", file_location, "-ss", "00:00:01", "-vframes", "1", "-vf", "scale='if(gt(a,1),"+width+",-1)':'if(gt(a,1),-1,"+height+")'", "-c:v", "libwebp", "-q:v", "80", "compressed/"+RemoveExtension(path.Base(file_location))+".webp")
	}

	_, err := cmdStruct.Output()
	if err != nil {
		fmt.Println(err)
	}
}

var EXIFTOOL = os.Getenv("EXIFTOOL")

func GetMetaData(path string) string {
	cmdStruct := exec.Command(EXIFTOOL, path)

	metadata, err := cmdStruct.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(metadata)
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func RemoveExtension(filename string) string {
	return filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
}

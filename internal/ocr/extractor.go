package ocr

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

type DNI struct {
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	IDNumber       string `json:"id_number"`
	Birthdate      string `json:"birthdate,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
	DocumentNumber string `json:"document_number,omitempty"`
	IssueDate      string `json:"issue_date,omitempty"`
	Nationality    string `json:"nationality,omitempty"`
	Gender         string `json:"gender,omitempty"`
}

func ExtracDataFromDNI(processedImg gocv.Mat) (*DNI, error) {
	dni := &DNI{}
	// width, height := processedImg.Cols(), processedImg.Rows()

	// rois := map[string]image.Rectangle{
	// 	"Surname":        percentageRect(width, height, 0.25, 0.14, 0.75, 0.22),
	// 	"Name":           percentageRect(width, height, 0.25, 0.22, 0.75, 0.30),
	// 	"Nationality":    percentageRect(width, height, 0.25, 0.35, 0.6, 0.4),
	// 	"Gender":         percentageRect(width, height, 0.6, 0.35, 0.75, 0.4),
	// 	"IDNumber":       percentageRect(width, height, 0.03, 0.6, 0.3, 0.7),
	// 	"Birthdate":      percentageRect(width, height, 0.4, 0.6, 0.75, 0.65),
	// 	"IssueDate":      percentageRect(width, height, 0.03, 0.75, 0.3, 0.8),
	// 	"ExpirationDate": percentageRect(width, height, 0.5, 0.75, 0.75, 0.8),
	// 	"DocumentNumber": percentageRect(width, height, 0.7, 0.68, 0.98, 0.72),
	// }

	rois := map[string]image.Rectangle{
		"Surname":        image.Rect(1010, 443, 1551, 652),
		"Name":           image.Rect(1004, 716, 2325, 833),
		"Nationality":    image.Rect(989, 999, 1458, 1105),
		"Gender":         image.Rect(2032, 1000, 2299, 1109),
		"Birthdate":      image.Rect(979, 1191, 1615, 1311),
		"DocumentNumber": image.Rect(2032, 1192, 2671, 1303),
		"IssueDate":      image.Rect(986, 1379, 1626, 1493),
		"ExpirationDate": image.Rect(2035, 1375, 2674, 1512),
		"IDNumber":       image.Rect(230, 1691, 881, 1827),
	}
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("spa")

	for name, rect := range rois {
		roiMat := processedImg.Region(rect)
		roi := roiMat.Clone()
		roiMat.Close()

		defer roi.Close()

		img, err := roi.ToImage()
		if err != nil {
			fmt.Printf("error while at convert ROI %s to image: %v\n", name, err)
			continue
		}

		client.SetImageFromBytes(ImgToBytes(img))

		switch name {
		case "Name", "Surname":
			client.SetPageSegMode(gosseract.PSM_SINGLE_BLOCK)
		case "Gender", "Nationality":
			client.SetPageSegMode(gosseract.PSM_SINGLE_WORD)
		default:
			client.SetPageSegMode(gosseract.PSM_SINGLE_LINE)
		}

		switch name {
		case "IDNumber", "DocumentNumber":
			client.SetWhitelist("0123456789kK-.")
		case "Birthdate", "IssueDate", "ExpirationDate":
			client.SetWhitelist("0123456789./-")
		case "Gender":
			client.SetWhitelist("MF")
		default:
			client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÑ ")
		}

		text, err := client.Text()
		if err != nil {
			fmt.Printf("OCR error for %s: %v\n", name, err)
		}
		cleanText := strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
		fmt.Printf("🔍 [%s] OCR: '%s'\n", name, cleanText)

		switch name {
		case "Name":
			dni.Name = cleanText
		case "Surname":
			dni.Surname = cleanText
		case "IDNumber":
			dni.IDNumber = cleanText
		case "Birthdate":
			dni.Birthdate = cleanText
		case "ExpirationDate":
			dni.ExpirationDate = cleanText
		case "IssueDate":
			dni.IssueDate = cleanText
		case "Nationality":
			dni.Nationality = cleanText
		case "Gender":
			dni.Gender = cleanText
		}
	}

	debugImg := processedImg.Clone()
	for name, rect := range rois {
		gocv.Rectangle(&debugImg, rect, color.RGBA{255, 0, 0, 0}, 2)
		gocv.PutText(&debugImg, name, image.Pt(rect.Min.X, rect.Min.Y-5), gocv.FontHersheySimplex, 0.5, color.RGBA{0, 255, 0, 0}, 1)
	}
	gocv.IMWrite("testdata/debug_rois.png", debugImg)

	return dni, nil
}

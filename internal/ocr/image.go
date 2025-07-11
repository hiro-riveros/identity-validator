package ocr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"sort"

	"gocv.io/x/gocv"
)

func ImgToBytes(img image.Image) []byte {
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	return buf.Bytes()
}

func FindDocumentContour(img gocv.Mat) (gocv.PointVector, error) {
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	blur := gocv.NewMat()
	defer blur.Close()
	gocv.GaussianBlur(gray, &blur, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	canny := gocv.NewMat()
	defer canny.Close()
	gocv.Canny(blur, &canny, 76, 200)

	contours := gocv.FindContours(canny, gocv.RetrievalTree, gocv.ChainApproxSimple)

	var largestContour gocv.PointVector
	found := false
	maxArea := 0.0

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		area := gocv.ContourArea(contour)
		if area < 1000 {
			contour.Close()
			continue
		}
		peri := gocv.ArcLength(contour, true)
		approx := gocv.ApproxPolyDP(contour, 0.02*peri, true)

		if len(approx.ToPoints()) == 4 && area > maxArea {
			if found {
				largestContour.Close()
			}
			maxArea = area
			largestContour = approx
			found = true
		} else {
			approx.Close()
		}
		contour.Close()
	}

	if !found {
		return gocv.PointVector{}, fmt.Errorf("no se encontró un contorno válido del documento")
	}

	return largestContour, nil
}

// debug_contour.go
func SaveContoursDebugImage(imgPath string, outPath string) error {
	img := gocv.IMRead(imgPath, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("no se pudo abrir imagen: %s", imgPath)
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	blur := gocv.NewMat()
	defer blur.Close()
	gocv.GaussianBlur(gray, &blur, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	canny := gocv.NewMat()
	defer canny.Close()
	gocv.Canny(blur, &canny, 76, 200)

	contours := gocv.FindContours(canny, gocv.RetrievalTree, gocv.ChainApproxSimple)
	gocv.DrawContours(&img, contours, -1, color.RGBA{0, 255, 0, 0}, 2)

	if ok := gocv.IMWrite(outPath, img); !ok {
		return fmt.Errorf("no se pudo guardar imagen de debug en %s", outPath)
	}

	return nil
}

func PerspectiveTransform(img gocv.Mat, contour gocv.PointVector) gocv.Mat {
	points := contour.ToPoints()
	sort.Slice(points, func(i, j int) bool {
		return points[i].X < points[j].X
	})

	width := 600
	height := 400

	src := gocv.NewPointVectorFromPoints(points)
	dst := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {width, 0}, {0, height}, {width, height}})

	transform := gocv.GetPerspectiveTransform(src, dst)
	defer transform.Close()

	warped := gocv.NewMat()
	gocv.WarpPerspective(img, &warped, transform, image.Pt(width, height))
	return warped
}

func percentageRect(imgWidth, imgHeight int, x1f, y1f, x2f, y2f float64) image.Rectangle {
	return image.Rect(
		int(x1f*float64(imgWidth)),
		int(y1f*float64(imgHeight)),
		int(x2f*float64(imgWidth)),
		int(y2f*float64(imgHeight)),
	)
}

func PreprocessImage(imgPath string) (gocv.Mat, error) {
	original := gocv.IMRead(imgPath, gocv.IMReadColor)
	if original.Empty() {
		return gocv.NewMat(), fmt.Errorf("no se pudo leer la imagen: %s", imgPath)
	}

	// 1. Desenfoque
	blurred := gocv.NewMat()
	gocv.GaussianBlur(original, &blurred, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	// 2. Convertir a escala de grises
	gray := gocv.NewMat()
	gocv.CvtColor(blurred, &gray, gocv.ColorBGRToGray)

	// 3. Deskew: opcional, aquí asumimos que el documento ya está bien orientado
	// Podríamos estimar el ángulo con HoughLines y rotar, si es necesario

	// 4. Umbral adaptativo fuerte
	thresh := gocv.NewMat()
	gocv.AdaptiveThreshold(gray, &thresh, 255, gocv.AdaptiveThresholdMean, gocv.ThresholdBinaryInv, 15, 8)

	// 5. Dilation opcional para resaltar texto
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(2, 2))
	gocv.Dilate(thresh, &thresh, kernel)

	// 6. Guardar para debug
	ok := gocv.IMWrite("testdata/debug_processed.png", thresh)
	if !ok {
		return gocv.NewMat(), fmt.Errorf("no se pudo guardar imagen procesada")
	}

	return thresh, nil
}

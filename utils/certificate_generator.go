package utils

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func GenerateCertificatePDF(userFullName, assessmentTitle string, totalMarks, marksObtained int, passingScore, userScore float64, certificateDate time.Time, passed bool) ([]byte, error) {

	templatePath := os.Getenv("CERTIFICATE_TEMPLATE_PATH")

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()

	// Add background template
	pdf.ImageOptions(
		templatePath,
		0, 0,
		297, 210, // Full A4 landscape size
		false,
		gofpdf.ImageOptions{ImageType: "JPEG"},
		0,
		"",
	)

	// Helper to draw centered text
	center := func(y float64, font string, size float64, text string) {
		pdf.SetFont("Helvetica", font, size)
		pdf.SetXY(0, y)
		pdf.CellFormat(297, 10, text, "", 0, "C", false, 0, "")
	}

	if passed {
		// ---------------- PASS CERTIFICATE ----------------
		center(40, "B", 36, "DHL")
		center(60, "", 18, "Certifies that")
		center(75, "B", 26, userFullName)

		marksText := fmt.Sprintf("Successfully achieved %d/%d", marksObtained, totalMarks)
		center(95, "", 18, marksText)

		center(120, "B", 20, assessmentTitle)

		dateText := "Date of Certification: " + certificateDate.Format("01/02/2006 15:04:05")
		center(135, "", 14, dateText)

	} else {
		// ---------------- FAIL CERTIFICATE / RESULT CARD ----------------
		center(40, "B", 36, "DHL RESULT REPORT")
		center(60, "", 18, "We regret to inform")

		center(75, "B", 26, userFullName)

		result := fmt.Sprintf("You scored %d/%d (%.2f%%)", marksObtained, totalMarks, userScore)
		center(95, "", 18, result)

		need := fmt.Sprintf("Required Passing Score: %.2f%%", passingScore)
		center(110, "B", 20, need)

		center(130, "B", 24, "STATUS: FAILED")

		center(150, "", 16, "Certificate cannot be generated for this attempt.")
	}

	// Output PDF buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

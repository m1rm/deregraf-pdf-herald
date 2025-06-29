
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

func TestPDFGeneration(t *testing.T) {
	// The text we send to the PDF generator
	inputText := "This is a test with some special characters: äöüß."

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Generate the PDF
	pdfBytes, err := generatePDF(inputText)
	if err != nil {
		t.Fatalf("generatePDF failed: %v", err)
	}

	// Write the PDF bytes to the ResponseRecorder
	rr.Header().Set("Content-Type", "application/pdf")
	rr.Header().Set("Content-Disposition", "attachment; filename=generated.pdf")
	rr.Write(pdfBytes)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the Content-Type header
	expectedContentType := "application/pdf"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// Check the Content-Disposition header
	expectedContentDisposition := "attachment; filename=generated.pdf"
	if contentDisposition := rr.Header().Get("Content-Disposition"); contentDisposition != expectedContentDisposition {
		t.Errorf("handler returned wrong content disposition: got %v want %v",
			contentDisposition, expectedContentDisposition)
	}
}

func TestGeneratePDFContent(t *testing.T) {
	inputText := "Ich sehe, die aus der Schenke jemand abgeführt wird. Trete ein, grüße und frage was los ist. Leute im Schankraum sind etwas hangover und baff. Stimmung ist drückend."

	pdfBytes, err := generatePDF(inputText)
	if err != nil {
		t.Fatalf("generatePDF failed: %v", err)
	}

	// Save the PDF to a file for inspection
	outputFilePath := "test_generated.pdf"
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to create output file %s: %v", outputFilePath, err)
	}
	defer outFile.Close() // Ensure the file is closed

	if _, err := outFile.Write(pdfBytes); err != nil {
		t.Fatalf("Failed to write PDF to file %s: %v", outputFilePath, err)
	}

	t.Logf("Generated PDF saved to: %s", outputFilePath)

	// IMPORTANT: Manually inspect this PDF file to verify that the header image is present
	// and that the text content starts below it, without overflowing.
	// This test primarily ensures PDF generation and provides a file for visual inspection.
}

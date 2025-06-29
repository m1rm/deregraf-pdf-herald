
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPDFGeneration(t *testing.T) {
	// The text we send to the PDF generator
	inputText := "This is a test with some special characters: äöüß."

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Generate the PDF
	generatePDF(rr, inputText)

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


package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/go-pdf/fpdf"
)

func generatePDF(text string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")

	// Add UTF-8 font
	pdf.AddUTF8Font("NotoSans", "", "NotoSans-Regular.ttf")
	pdf.AddUTF8Font("NotoSans", "B", "NotoSans-Bold.ttf")
	pdf.AddUTF8Font("NotoSans", "I", "NotoSans-Italic.ttf")
	pdf.SetFont("NotoSans", "", 12) // Set font for general use

	pdf.SetHeaderFunc(func() {
		pdf.SetFont("NotoSans", "B", 12) // Use NotoSans for header
		pdf.Cell(0, 10, "Deregraf - The Aventurian PDF Herald")
		pdf.Ln(15)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("NotoSans", "I", 8) // Use NotoSans for footer
		pdf.CellFormat(0, 10, fmt.Sprintf("Seite %d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})

	pdf.AddPage()
	// Use MultiCell for multi-line text
	pdf.MultiCell(190, 10, text, "", "", false)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		text := r.FormValue("text")
		pdfBytes, err := generatePDF(text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=generated.pdf")
		w.Write(pdfBytes)
	} else {
		http.ServeFile(w, r, "index.html")
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

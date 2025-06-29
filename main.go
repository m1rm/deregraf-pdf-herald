package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"codeberg.org/go-pdf/fpdf"
)

// MyPDF is a custom Fpdf type to override Header and Footer
type MyPDF struct {
	*fpdf.Fpdf
}

// Header is called automatically at the beginning of each new page
func (p *MyPDF) Header() {
	// Page dimensions
	pageWidth, _ := p.GetPageSize()

	// Header image
	imagePath := "aventurischerBote_crappyHeader.jpg"
	
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current working directory: %v", err)
		return
	}
	
	// Construct the absolute path to the image
	absImagePath := filepath.Join(cwd, imagePath)
	fmt.Printf("Attempting to register image from absolute path: %s\n", absImagePath)

	// Register the image and get its info
	imgInfo := p.RegisterImage(imagePath, "JPG") // Explicitly set type to JPG
	if imgInfo == nil {
		fmt.Printf("Failed to register image: %s. This might be due to an invalid path or unsupported format.\n", imagePath)
		// You can also try to get more specific error information if fpdf provides it
		// For example, by checking if the file exists before registering
		if _, err := os.Stat(absImagePath); os.IsNotExist(err) {
			fmt.Printf("Error: Image file does not exist at %s\n", absImagePath)
		} else if err != nil {
			fmt.Printf("Error checking image file existence: %v\n", err)
		}
	} else {
		fmt.Printf("Image Info: Width=%f, Height=%f\n", imgInfo.Width(), imgInfo.Height())
		// Draw the image
		p.Image(imagePath, 0, 0, pageWidth, 0, false, "", 0, "")
		// Set the Y position for the next content, below the image with a margin
		p.SetY(imgInfo.Height() + 5) // 5mm margin
	}
}

// Footer is called automatically at the end of each new page
func (p *MyPDF) Footer() {
	p.SetY(-15)
	p.SetFont("NotoSans", "I", 8) // Use NotoSans for footer
	p.CellFormat(0, 10, fmt.Sprintf("Seite %d", p.PageNo()),
		"", 0, "C", false, 0, "")
}

func generatePDF(text string) ([]byte, error) {
	pdf := MyPDF{fpdf.New("P", "mm", "A4", "")}
	pdf.SetHeaderFunc(pdf.Header)
	pdf.SetFooterFunc(pdf.Footer)

	// Add UTF-8 font
	pdf.AddUTF8Font("NotoSans", "", "NotoSans-Regular.ttf")
	pdf.AddUTF8Font("NotoSans", "B", "NotoSans-Bold.ttf")
	pdf.AddUTF8Font("NotoSans", "I", "NotoSans-Italic.ttf")
	pdf.SetFont("NotoSans", "", 12) // Set font for general use

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
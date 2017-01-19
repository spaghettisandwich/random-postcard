package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/jung-kurt/gofpdf"
)

// read an address from a config file and turn it into a pdf with a random bar chart
func main() {
	// initial PDF
	// TODO fewer constants, brittle to specific sizing, fonts, etc.
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "in",
		Size:    gofpdf.SizeType{Wd: 6, Ht: 4},
	})
	pdf.SetFont("Helvetica", "", 14)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	// read config and write each line
	file, err := os.Open("./config.txt")
	if err != nil {
		fmt.Printf("error reading file %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pdf.SetX(2)
		pdf.WriteAligned(3, 3, scanner.Text(), "L")
		pdf.Ln(0.25)
	}

	// New page for your random image
	pdf.AddPage()
	err = drawRandomBar(pdf)
	if err != nil {
		fmt.Printf("error generating rectangle %s", err)
		os.Exit(1)
	}

	fileStr := "random-postcard.pdf"
	err = pdf.OutputFileAndClose(fileStr)
	if err == nil {
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}

func drawRandomBar(pdf *gofpdf.Fpdf) error {
	// really random
	val, err := rand.Int(rand.Reader, big.NewInt(101))
	if err != nil {
		return err
	}
	// TODO make clipping for rounded rectangles
	// outer rectangle
	rectangleWidth := 3.0
	pdf.SetFillColor(77, 184, 255)
	pdf.Rect(1.5, 2, rectangleWidth, 0.25, "F")
	// inner rectangle
	floatVal := float64(val.Int64())
	innerWidth := (floatVal / 100) * rectangleWidth
	pdf.SetXY((1.4 + innerWidth), 1.5)
	pdf.SetFont("Helvetica", "", 11)
	pdf.Writef(0.7, "%0.f%%", floatVal)
	pdf.SetFillColor(0, 45, 179)
	pdf.Rect(1.505, 2.05, innerWidth, 0.15, "F")

	return nil
}

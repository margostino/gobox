package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readCsv(response io.ReadCloser) ([][]string, error) {
	defer response.Close()
	reader := csv.NewReader(response)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func main() {
	url := "https://api.worldbank.org/v2/en/indicator/SI.POV.DDAY?downloadformat=csv"
	//url := "https://raw.githubusercontent.com/owid/owid-datasets/master/datasets/20th century deaths in US - CDC/20th century deaths in US - CDC.csv"

	out, err := os.Create("test.zip")
	defer out.Close()
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	print(n)

	data, err := readCsv(resp.Body)
	print(data)

	dst := "output"
	archive, err := zip.OpenReader("test.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	records := readCsvFile("edsc_collection.csv")
	fmt.Println(records)
	for _, record := range records {
		fmt.Println(record.DataProvider)
	}
}

type NasaEarthData struct {
	DataProvider    string `csv:"Data Provider"`
	ShortName       string `csv:"Short Name"`
	Version         string `csv:"Version"`
	EntryTitle      string `csv:"Entry Title"`
	ProcessingLevel string `csv:"Processing Level"`
	Platform        string `csv:"Platform"`
	StartTime       string `csv:"Start Time"`
	EndTime         string `csv:"End Time"`
}

func readCsvFile(filePath string) []*NasaEarthData {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	var records []*NasaEarthData

	if err := gocsv.UnmarshalFile(f, &records); err != nil {
		panic(err)
	}

	//csvReader := csv.NewReader(f)
	//records, err := csvReader.ReadAll()
	//if err != nil {
	//	log.Fatal("Unable to parse file as CSV for "+filePath, err)
	//}

	return records
}

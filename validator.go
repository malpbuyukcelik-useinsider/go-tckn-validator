package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const nviSoapURL = "https://tckimlik.nvi.gov.tr/Service/KPSPublic.asmx"

type TCKNValidationRequest struct {
	XMLName    xml.Name `xml:"TCKimlikNoDogrula"`
	XMLNs      string   `xml:"xmlns,attr"`
	TCKN       int64    `xml:"TCKimlikNo"`
	Ad         string   `xml:"Ad"`
	Soyad      string   `xml:"Soyad"`
	DogumYili  int      `xml:"DogumYili"`
}

// cleanName cleans and formats the name for NVI service
func cleanName(name string) string {
	return strings.TrimSpace(name)
}

// Basic TCKN validation algorithm
func validateTCKN(tckn string) bool {
	if len(tckn) != 11 {
		return false
	}

	// Convert string to array of integers
	digits := make([]int, 11)
	for i, r := range tckn {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Rule 1: First digit cannot be 0
	if digits[0] == 0 {
		return false
	}

	// Rule 2: 10th digit is ((sum of digits 1,3,5,7,9)*7 - sum of digits 2,4,6,8) mod 10
	odd := digits[0] + digits[2] + digits[4] + digits[6] + digits[8]
	even := digits[1] + digits[3] + digits[5] + digits[7]
	digit10 := (odd*7 - even) % 10
	if digit10 < 0 {
		digit10 += 10
	}
	if digits[9] != digit10 {
		return false
	}

	// Rule 3: 11th digit is sum of first 10 digits mod 10
	sum := 0
	for i := 0; i < 10; i++ {
		sum += digits[i]
	}
	if digits[10] != sum%10 {
		return false
	}

	return true
}

// validateWithNVI validates TCKN using NVI SOAP service
func validateWithNVI(tckn, ad, soyad string, dogumYili int) (bool, error) {
	tcknInt, err := strconv.ParseInt(tckn, 10, 64)
	if err != nil {
		return false, err
	}

	// Clean and format names
	ad = cleanName(ad)
	soyad = cleanName(soyad)

	request := TCKNValidationRequest{
		XMLNs:     "http://tckimlik.nvi.gov.tr/WS",
		TCKN:     tcknInt,
		Ad:       ad,
		Soyad:    soyad,
		DogumYili: dogumYili,
	}

	xmlData, err := xml.MarshalIndent(request, "", "  ")
	if err != nil {
		return false, err
	}

	soapEnvelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    %s
  </soap:Body>
</soap:Envelope>`, string(xmlData))

	// Log the request
	log.Printf("SOAP Request:\n%s\n", soapEnvelope)

	req, err := http.NewRequest("POST", nviSoapURL, bytes.NewBufferString(soapEnvelope))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://tckimlik.nvi.gov.tr/WS/TCKimlikNoDogrula")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Log the response
	log.Printf("SOAP Response:\n%s\n", string(body))

	return bytes.Contains(body, []byte("<TCKimlikNoDogrulaResult>true</TCKimlikNoDogrulaResult>")), nil
} 
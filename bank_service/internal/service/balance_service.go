package service

import (
	"bytes"
	"encoding/csv"
	"regexp"
	"strconv"
)

type Transaction struct {
	Date      string
	Narrative []string
	Type      string
	Credit    float64
	Debit     float64
	Currency  string
}

func parseCSVFromBytes(data []byte) (map[string]float64, error) {
	reader := bytes.NewReader(data)
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	paymentRegex := regexp.MustCompile(`PAY\d{6}[A-Z]{2}`)
	totals := make(map[string]float64)

	for _, record := range records[1:] {
		transaction := parseRecord(record)

		if transaction.Date == "06/03/2011" && containsPaymentReference(transaction.Narrative, paymentRegex) {
			amount := transaction.Debit
			if amount == 0 {
				amount = transaction.Credit
			}
			totals[transaction.Currency] += amount
		}
	}

	return totals, nil
}

func parseRecord(record []string) Transaction {
	var credit, debit float64
	var err error

	if record[7] != "" {
		credit, err = strconv.ParseFloat(record[7], 64)
		if err != nil {
			credit = 0
		}
	}
	if record[8] != "" {
		debit, err = strconv.ParseFloat(record[8], 64)
		if err != nil {
			debit = 0
		}
	}

	return Transaction{
		Date:      record[0],
		Narrative: record[1:6],
		Type:      record[6],
		Credit:    credit,
		Debit:     debit,
		Currency:  record[9],
	}
}

func containsPaymentReference(narratives []string, regex *regexp.Regexp) bool {
	for _, narrative := range narratives {
		if regex.MatchString(narrative) {
			return true
		}
	}
	return false
}

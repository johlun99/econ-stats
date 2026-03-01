package parser

import (
	"fmt"
	"strconv"
	"strings"

	"econ-stats/internal/models"

	"github.com/xuri/excelize/v2"
)

// Handelsbanken header row patterns
var transferPrefixes = []string{
	"överf mobil",
	"överf internet",
	"överföring",
}

func ParseHandelsbanken(filePath string) ([]models.Transaction, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("open xlsx: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}

	var transactions []models.Transaction

	// Data starts at row 10 (index 9), after Handelsbanken headers
	for i := 9; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 4 {
			continue
		}

		bookingDate := strings.TrimSpace(row[0])
		transactionDate := strings.TrimSpace(row[1])
		description := strings.TrimSpace(row[2])
		amountStr := strings.TrimSpace(row[3])

		if description == "" || amountStr == "" {
			continue
		}

		// Use transaction date as fallback for booking date
		if bookingDate == "" {
			bookingDate = transactionDate
		}
		if transactionDate == "" {
			continue
		}

		amount, err := parseAmount(amountStr)
		if err != nil {
			continue
		}

		var balance float64
		if len(row) > 4 {
			balance, _ = parseAmount(strings.TrimSpace(row[4]))
		}

		merchantKey := normalizeMerchant(description)
		isTransfer := isTransferTransaction(merchantKey)

		transactions = append(transactions, models.Transaction{
			BookingDate:     bookingDate,
			TransactionDate: transactionDate,
			Description:     description,
			Amount:          amount,
			Balance:         balance,
			MerchantKey:     merchantKey,
			IsTransfer:      isTransfer,
		})
	}

	return transactions, nil
}

func parseAmount(s string) (float64, error) {
	// Handle Swedish number format: spaces as thousands separator, comma or dot as decimal
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\u00a0", "") // non-breaking space
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

func normalizeMerchant(description string) string {
	s := strings.ToLower(strings.TrimSpace(description))
	// Remove extra whitespace
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func isTransferTransaction(merchantKey string) bool {
	for _, prefix := range transferPrefixes {
		if strings.HasPrefix(merchantKey, prefix) {
			return true
		}
	}
	return false
}

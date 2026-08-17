package usecase

import (
	"regexp"
	"strconv"
	"strings"

	"backend_keping/internal/domain"
)

type RegexParser interface {
	Parse(packageName, title, content string) domain.ParsedTransaction
}

type regexParser struct{}

func NewRegexParser() RegexParser {
	return &regexParser{}
}

func (p *regexParser) Parse(packageName, title, content string) domain.ParsedTransaction {
	fullText := title + " " + content

	switch packageName {
	case "id.dana":
		return parseDana(fullText)
	case "id.co.bri.brimo":
		return parseBrimo(fullText)
	case "com.gojek.app":
		return parseGojek(fullText)
	case "com.shopee.id":
		return parseShopeePay(fullText)
	default:
		return parseGeneric(fullText)
	}
}

// Helper parsing DANA notification
func parseDana(text string) domain.ParsedTransaction {
	res := domain.ParsedTransaction{IsValid: false, Type: domain.TypeDebit}

	// Pattern Isi Saldo / Uang Masuk
	if strings.Contains(strings.ToLower(text), "isi saldo") || strings.Contains(strings.ToLower(text), "berhasil menerima") {
		res.Type = domain.TypeCredit
		res.Description = "DANA - Top Up / Saldo Masuk"
	} else {
		res.Type = domain.TypeDebit
		res.Description = "DANA - Pembayaran / Transfer Out"
	}

	amount := extractAmount(text)
	if amount > 0 {
		res.Amount = amount
		res.IsValid = true
	}
	return res
}

// Helper parsing BRImo notification
func parseBrimo(text string) domain.ParsedTransaction {
	res := domain.ParsedTransaction{IsValid: false, Type: domain.TypeDebit}

	if strings.Contains(strings.ToLower(text), "kredit") || strings.Contains(strings.ToLower(text), "masuk") {
		res.Type = domain.TypeCredit
		res.Description = "BRImo - Transfer Masuk"
	} else {
		res.Type = domain.TypeDebit
		res.Description = "BRImo - Transaksi / Transfer Keluar"
	}

	amount := extractAmount(text)
	if amount > 0 {
		res.Amount = amount
		res.IsValid = true
	}
	return res
}

// Helper parsing GoPay notification
func parseGojek(text string) domain.ParsedTransaction {
	res := domain.ParsedTransaction{IsValid: false, Type: domain.TypeDebit}

	if strings.Contains(strings.ToLower(text), "menerima") || strings.Contains(strings.ToLower(text), "top up") {
		res.Type = domain.TypeCredit
		res.Description = "GoPay - Saldo Masuk"
	} else {
		res.Type = domain.TypeDebit
		res.Description = "GoPay - Pembayaran"
	}

	amount := extractAmount(text)
	if amount > 0 {
		res.Amount = amount
		res.IsValid = true
	}
	return res
}

// Helper parsing ShopeePay notification
func parseShopeePay(text string) domain.ParsedTransaction {
	res := domain.ParsedTransaction{IsValid: false, Type: domain.TypeDebit}

	if strings.Contains(strings.ToLower(text), "berhasil isi saldo") || strings.Contains(strings.ToLower(text), "menerima saldo") {
		res.Type = domain.TypeCredit
		res.Description = "ShopeePay - Saldo Masuk"
	} else {
		res.Type = domain.TypeDebit
		res.Description = "ShopeePay - Pembayaran"
	}

	amount := extractAmount(text)
	if amount > 0 {
		res.Amount = amount
		res.IsValid = true
	}
	return res
}

// Helper generic regex parser
func parseGeneric(text string) domain.ParsedTransaction {
	res := domain.ParsedTransaction{IsValid: false, Type: domain.TypeDebit, Description: "Transaksi Notifikasi"}
	amount := extractAmount(text)
	if amount > 0 {
		res.Amount = amount
		res.IsValid = true
	}
	return res
}

// extractAmount extracts numeric currency value (e.g., "Rp50.000" or "Rp 50.000,00" -> 50000)
func extractAmount(text string) float64 {
	// Regex matches Rp50.000 or Rp 1.500.000
	re := regexp.MustCompile(`(?i)Rp\s*([\d\.]+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		cleanNum := strings.ReplaceAll(matches[1], ".", "")
		val, err := strconv.ParseFloat(cleanNum, 64)
		if err == nil {
			return val
		}
	}
	return 0
}

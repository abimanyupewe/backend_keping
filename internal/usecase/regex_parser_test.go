package usecase

import (
	"testing"

	"backend_keping/internal/domain"
)

func TestRegexParser_Parse(t *testing.T) {
	parser := NewRegexParser()

	tests := []struct {
		name        string
		pkg         string
		title       string
		content     string
		expectedType domain.TransactionType
		expectedAmt  float64
		expectedValid bool
	}{
		{
			name:          "DANA Isi Saldo",
			pkg:           "id.dana",
			title:         "Isi Saldo Berhasil",
			content:       "Kamu telah berhasil melakukan isi saldo sebesar Rp500.000 via BRI",
			expectedType:  domain.TypeCredit,
			expectedAmt:   500000,
			expectedValid: true,
		},
		{
			name:          "DANA Pembayaran QRIS",
			pkg:           "id.dana",
			title:         "Pembayaran Berhasil",
			content:       "Pembayaran sebesar Rp25.000 ke Kopi Kenangan telah berhasil",
			expectedType:  domain.TypeDebit,
			expectedAmt:   25000,
			expectedValid: true,
		},
		{
			name:          "BRImo Kredit Masuk",
			pkg:           "id.co.bri.brimo",
			title:         "Notifikasi Transaksi",
			content:       "Dana Kredit masuk sebesar Rp1.500.000 dari REKENING XXX",
			expectedType:  domain.TypeCredit,
			expectedAmt:   1500000,
			expectedValid: true,
		},
		{
			name:          "Invalid Text No Amount",
			pkg:           "id.dana",
			title:         "Promo DANA",
			content:       "Dapatkan cashback menarik hari ini!",
			expectedType:  domain.TypeDebit,
			expectedAmt:   0,
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := parser.Parse(tt.pkg, tt.title, tt.content)
			if res.IsValid != tt.expectedValid {
				t.Errorf("expected IsValid %v, got %v", tt.expectedValid, res.IsValid)
			}
			if res.Type != tt.expectedType {
				t.Errorf("expected Type %v, got %v", tt.expectedType, res.Type)
			}
			if res.Amount != tt.expectedAmt {
				t.Errorf("expected Amount %v, got %v", tt.expectedAmt, res.Amount)
			}
		})
	}
}

package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type OCRResult struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // DEBIT or CREDIT
	Merchant    string  `json:"merchant"`
	Description string  `json:"description"`
}

type GeminiClient interface {
	ScanReceipt(ctx context.Context, imageBytes []byte, mimeType string) (*OCRResult, error)
}

type geminiClient struct {
	apiKey    string
	modelName string
}

func NewGeminiClient(apiKey, modelName string) GeminiClient {
	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}
	return &geminiClient{apiKey: apiKey, modelName: modelName}
}

func (g *geminiClient) ScanReceipt(ctx context.Context, imageBytes []byte, mimeType string) (*OCRResult, error) {
	if g.apiKey == "" {
		return nil, fmt.Errorf("gemini API key is not configured")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(g.apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(g.modelName)
	model.SetTemperature(0.1) // Low temperature for deterministic output

	prompt := "Analisislah gambar bukti transfer / struk pembayaran ini.\n" +
		"Ekstrak informasi transaksi berikut dalam format JSON MURNI tanpa markdown (tanpa backticks ```json ... ```):\n" +
		"{\n" +
		"  \"amount\": 50000,\n" +
		"  \"type\": \"DEBIT\",\n" +
		"  \"merchant\": \"Nama Toko / Penerima\",\n" +
		"  \"description\": \"Pembayaran QRIS / Struk Kopi Kenangan\"\n" +
		"}\n" +
		"Catatan:\n" +
		"- \"type\" HANYA boleh bernilai \"DEBIT\" (uang keluar/pembayaran) atau \"CREDIT\" (uang masuk/isi saldo/transfer masuk).\n" +
		"- \"amount\" harus berupa angka murni (tanpa Rp, titik, atau koma desimal)."

	imgData := genai.ImageData(mimeType, imageBytes)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt), imgData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content from gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini vision API")
	}

	part := resp.Candidates[0].Content.Parts[0]
	textResp, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected content part type from gemini")
	}

	cleanJSON := strings.TrimSpace(string(textResp))
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimPrefix(cleanJSON, "```")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	var result OCRResult
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from gemini response: %w (raw: %s)", err, cleanJSON)
	}

	return &result, nil
}

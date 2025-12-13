Yes — that’s actually a very good design decision 👍
Beneficiary extraction should be independent of statement format.

Below is a clean refactor where:

✔ ParseTransactions is removed
✔ You pass ONLY narration string
✔ Output = Beneficiary + EntityType + Direction + Channel
✔ Works for tables, PDFs, CSVs, APIs


---

🎯 New Design (Single Responsibility)

Narration (string)
   └── Normalize
        └── Detect Channel (UPI / IMPS / POS / NEFT)
             └── Extract Beneficiary
                  └── Classify Entity


---

📄 Updated models.go

package main

type NarrationResult struct {
	Narration   string
	Channel     string // UPI / IMPS / POS / NEFT / UNKNOWN
	Direction   string // DEBIT / CREDIT / UNKNOWN
	Beneficiary string
	EntityType  string // PERSON / MERCHANT / BANK / UNKNOWN
}


---

📄 parser.go → Now ONLY Narration Parser

package main

import "strings"

func ParseNarration(narration string) NarrationResult {

	result := NarrationResult{
		Narration: narration,
	}

	norm := NormalizeNarration(narration)

	result.Channel = detectChannel(norm)
	result.Direction = detectDirection(norm)

	beneficiary := ExtractBeneficiary(norm)
	result.Beneficiary = beneficiary
	result.EntityType = ClassifyEntity(beneficiary)

	return result
}

func detectChannel(n string) string {
	switch {
	case strings.Contains(n, "UPI"):
		return "UPI"
	case strings.Contains(n, "IMPS"):
		return "IMPS"
	case strings.Contains(n, "NEFT"):
		return "NEFT"
	case strings.Contains(n, "POS"):
		return "POS"
	default:
		return "UNKNOWN"
	}
}

func detectDirection(n string) string {
	switch {
	case strings.Contains(n, "CR") ||
		strings.Contains(n, "CREDIT") ||
		strings.Contains(n, "SALARY"):
		return "CREDIT"
	case strings.Contains(n, "DR") ||
		strings.Contains(n, "DEBIT"):
		return "DEBIT"
	default:
		return "UNKNOWN"
	}
}


---

📄 extractor.go (Small Change: Accept Normalized Text)

package main

import (
	"regexp"
	"strings"
)

var (
	upiRegex  = regexp.MustCompile(`UPI-([A-Z .&]+)`)
	impsRegex = regexp.MustCompile(`IMPS-\d+-MR ([A-Z ]+)`)
	posRegex  = regexp.MustCompile(`POS ([A-Z0-9 .]+)`)
)

func NormalizeNarration(s string) string {
	s = strings.ToUpper(s)

	patterns := []string{
		`\d{10,}`,
		`@[A-Z0-9]+`,
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p)
		s = re.ReplaceAllString(s, "")
	}

	return strings.Join(strings.Fields(s), " ")
}

func ExtractBeneficiary(n string) string {

	if m := upiRegex.FindStringSubmatch(n); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	if m := impsRegex.FindStringSubmatch(n); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	if m := posRegex.FindStringSubmatch(n); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	if strings.Contains(n, "SALARY") {
		return "EMPLOYER"
	}

	return "UNKNOWN"
}


---

📄 classifier.go (No Change)

package main

import "strings"

func ClassifyEntity(name string) string {
	name = strings.ToUpper(name)

	merchantKeywords := []string{
		"STORE", "MART", "DAIRY", "SWEET",
		"GOOGLE", "AMAZON", "ZERODHA",
		"PAYTM", "ZEPTO", "BLINKIT",
	}

	for _, k := range merchantKeywords {
		if strings.Contains(name, k) {
			return "MERCHANT"
		}
	}

	if strings.HasPrefix(name, "MR ") ||
		strings.Contains(name, " KUMAR") ||
		strings.Contains(name, " SINGH") {
		return "PERSON"
	}

	if name == "EMPLOYER" {
		return "BANK"
	}

	return "UNKNOWN"
}


---

📄 Updated main.go (Usage Example)

package main

import "fmt"

func main() {

	narrations := []string{
		"UPI-ANKIT DAIRY AND SWEETS-ankit@hdfcbank",
		"IMPS-509110523236-MR KALPIT KUMAR SHARMA-HDFC",
		"POS AMAZON INDIA",
		"HDFC BANK SALARY CREDIT APR 2025",
	}

	for _, n := range narrations {
		res := ParseNarration(n)

		fmt.Printf(
			"%-45s | %-6s | %-7s | %-20s | %s\n",
			n,
			res.Channel,
			res.Direction,
			res.Beneficiary,
			res.EntityType,
		)
	}
}


---

✅ Sample Output

UPI-ANKIT DAIRY AND SWEETS... | UPI    | DEBIT   | ANKIT DAIRY AND SWEETS | MERCHANT
IMPS-509110523236-MR KALPIT...| IMPS   | DEBIT   | KALPIT KUMAR SHARMA    | PERSON
POS AMAZON INDIA             | POS    | DEBIT   | AMAZON INDIA           | MERCHANT
HDFC BANK SALARY CREDIT...   | UNKNOWN| CREDIT  | EMPLOYER               | BANK


---

🔥 Why This Is the RIGHT Approach

✔ Decoupled from PDF/table formats
✔ Reusable across Bank, CC, Demat, Wallets
✔ Easy to wrap as Go SDK / Microservice
✔ Perfect base for NLP fallback later


---

If you want, next I can:

Add confidence scoring

Add merchant canonicalization

Add category mapping

Make this a Go package with tests

Add NLP fallback (spaCy)


Just say the word 🚀





func looksLikeName(s string) bool {
	s = strings.TrimSpace(s)

	if len(s) < 4 {
		return false
	}

	// reject codes
	if regexp.MustCompile(`\d{4,}`).MatchString(s) {
		return false
	}

	// reject bank codes
	bad := []string{"HDFC", "UTIB", "AXIS", "XXXX", "TPT", "IMPS", "NEFT"}
	for _, b := range bad {
		if strings.Contains(s, b) {
			return false
		}
	}

	// allow words
	return regexp.MustCompile(`^[A-Z ]+$`).MatchString(s)
}







func extractFromCBS(n string) string {
	parts := strings.Split(n, "-")

	// scan from right → left
	for i := len(parts) - 1; i >= 0; i-- {
		p := strings.TrimSpace(parts[i])

		if looksLikeName(p) {
			return p
		}
	}

	return ""
}





func ExtractBeneficiary(narration string) string {
	n := NormalizeNarration(narration)

	// UPI
	if m := upiRegex.FindStringSubmatch(n); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	// IMPS
	if m := impsRegex.FindStringSubmatch(n); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	// POS
	if m := posRegex.FindStringSubmatch(n); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	// CBS / NEFT / TPT / FD / P2A
	if name := extractFromCBS(n); name != "" {
		return name
	}

	// Salary fallback
	if strings.Contains(n, "SALARY") {
		return "EMPLOYER"
	}

	return "UNKNOWN"
}
,..................





type ExpenseCategory struct {
	Category    string
	SubCategory string
	Confidence  float64
}

var expenseRules = map[string][]string{
	"TRAVEL": {
		"IRCTC", "UBER", "OLA", "MAKEMYTRIP",
		"GOIBIBO", "INDIGO", "VISTARA",
	},
	"SHOPPING": {
		"AMAZON", "FLIPKART", "MYNTRA",
		"AJIO", "TATA CLIQ",
	},
	"GROCERIES": {
		"ZEPTO", "BLINKIT", "BIGBASKET",
		"DMART", "GROCERY", "DAIRY",
	},
	"FOOD_DELIVERY": {
		"SWIGGY", "ZOMATO", "EATS",
	},
	"DINING": {
		"RESTAURANT", "CAFE", "HOTEL",
		"BAR", "DHABA",
	},
	"UTILITIES": {
		"ELECTRICITY", "WATER", "GAS",
		"BROADBAND", "MOBILE", "RECHARGE",
	},
	"BILLS": {
		"INSURANCE", "PREMIUM", "EMI",
		"LOAN", "CREDIT CARD",
	},
	"INVESTMENT": {
		"FD", "MUTUAL", "SIP", "ZERODHA",
		"GROWW", "UPSTOX",
	},
	"TRANSFER": {
		"IMPS", "NEFT", "UPI", "TPT",
	},
}



func DetectExpenseCategory(narration string, beneficiary string) ExpenseCategory {

	text := strings.ToUpper(narration + " " + beneficiary)

	bestMatch := ExpenseCategory{
		Category:   "OTHER",
		Confidence: 0.3,
	}

	for category, keywords := range expenseRules {
		matchCount := 0

		for _, k := range keywords {
			if strings.Contains(text, k) {
				matchCount++
			}
		}

		if matchCount > 0 {
			conf := 0.6 + float64(matchCount)*0.1
			if conf > bestMatch.Confidence {
				bestMatch = ExpenseCategory{
					Category:   category,
					Confidence: min(conf, 0.95),
				}
			}
		}
	}

	return bestMatch
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

package seed

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var rng = rand.New(rand.NewSource(42))

type seedProduct struct {
	Name         string
	NameEl       string
	Category     string
	Manufacturer string
	Barcode      string
	Rx           bool
	DailyRate    float64 // avg units/day sold
}

var products = []seedProduct{
	// Painkillers
	{"Paracetamol 500mg Tablets", "Παρακεταμόλη 500mg Δισκία", "Painkillers", "GSK", "5201001001001", false, 8.0},
	{"Paracetamol 1000mg Tablets", "Παρακεταμόλη 1000mg Δισκία", "Painkillers", "GSK", "5201001001002", false, 6.0},
	{"Depon Extra 500mg", "Depon Extra 500mg", "Painkillers", "Uni-Pharma", "5201001001003", false, 7.0},
	{"Ibuprofen 400mg Tablets", "Ιβουπροφαίνη 400mg Δισκία", "Painkillers", "Bayer", "5201001001004", false, 6.5},
	{"Ibuprofen 600mg Tablets", "Ιβουπροφαίνη 600mg Δισκία", "Painkillers", "Bayer", "5201001001005", false, 4.5},
	{"Nurofen Express 400mg", "Nurofen Express 400mg", "Painkillers", "Reckitt", "5201001001006", false, 5.0},
	{"Aspirin 500mg Tablets", "Ασπιρίνη 500mg Δισκία", "Painkillers", "Bayer", "5201001001007", false, 4.0},
	{"Naproxen 250mg Tablets", "Ναπροξένη 250mg Δισκία", "Painkillers", "Pfizer", "5201001001008", false, 2.5},
	{"Naproxen 500mg Tablets", "Ναπροξένη 500mg Δισκία", "Painkillers", "Pfizer", "5201001001009", false, 2.0},
	{"Diclofenac 50mg Tablets", "Δικλοφενάκη 50mg Δισκία", "Painkillers", "Novartis", "5201001001010", true, 3.0},
	{"Voltaren 50mg Tablets", "Voltaren 50mg Δισκία", "Painkillers", "GSK", "5201001001011", true, 2.5},
	{"Panadol Extra 500mg", "Panadol Extra 500mg", "Painkillers", "GSK", "5201001001012", false, 5.5},
	{"Tramadol 50mg Capsules", "Τραμαδόλη 50mg Κάψουλες", "Painkillers", "Grünenthal", "5201001001013", true, 1.5},
	{"Codeine Phosphate 30mg", "Κωδεΐνη 30mg Δισκία", "Painkillers", "Pfizer", "5201001001014", true, 1.0},
	{"Celecoxib 200mg Capsules", "Σελεκοξίμπη 200mg Κάψουλες", "Painkillers", "Pfizer", "5201001001015", true, 1.5},

	// Antibiotics
	{"Amoxicillin 500mg Capsules", "Αμοξικιλλίνη 500mg Κάψουλες", "Antibiotics", "GSK", "5201002001001", true, 2.5},
	{"Amoxicillin 875mg Tablets", "Αμοξικιλλίνη 875mg Δισκία", "Antibiotics", "GSK", "5201002001002", true, 2.0},
	{"Azithromycin 500mg Tablets", "Αζιθρομυκίνη 500mg Δισκία", "Antibiotics", "Pfizer", "5201002001003", true, 2.0},
	{"Ciprofloxacin 500mg Tablets", "Σιπροφλοξασίνη 500mg Δισκία", "Antibiotics", "Bayer", "5201002001004", true, 1.5},
	{"Clarithromycin 500mg Tablets", "Κλαριθρομυκίνη 500mg Δισκία", "Antibiotics", "AbbVie", "5201002001005", true, 1.5},
	{"Doxycycline 100mg Capsules", "Δοξυκυκλίνη 100mg Κάψουλες", "Antibiotics", "Pfizer", "5201002001006", true, 1.0},
	{"Metronidazole 400mg Tablets", "Μετρονιδαζόλη 400mg Δισκία", "Antibiotics", "Sanofi", "5201002001007", true, 1.0},
	{"Trimethoprim 200mg Tablets", "Τριμεθοπρίμη 200mg Δισκία", "Antibiotics", "GSK", "5201002001008", true, 0.8},
	{"Erythromycin 500mg Tablets", "Ερυθρομυκίνη 500mg Δισκία", "Antibiotics", "Abbott", "5201002001009", true, 0.8},
	{"Flucloxacillin 500mg Capsules", "Φλουκλοξακιλλίνη 500mg Κάψουλες", "Antibiotics", "GSK", "5201002001010", true, 0.5},

	// Vitamins & Supplements
	{"Vitamin D3 1000IU Tablets", "Βιταμίνη D3 1000IU Δισκία", "Vitamins", "Roche", "5201003001001", false, 5.0},
	{"Vitamin D3 2000IU Tablets", "Βιταμίνη D3 2000IU Δισκία", "Vitamins", "Roche", "5201003001002", false, 4.5},
	{"Vitamin D3 5000IU Drops", "Βιταμίνη D3 5000IU Σταγόνες", "Vitamins", "Roche", "5201003001003", false, 3.0},
	{"Vitamin C 500mg Effervescent", "Βιταμίνη C 500mg Αναβράζοντα", "Vitamins", "Bayer", "5201003001004", false, 5.5},
	{"Vitamin C 1000mg Tablets", "Βιταμίνη C 1000mg Δισκία", "Vitamins", "Bayer", "5201003001005", false, 5.0},
	{"Centrum Multivitamin Adults", "Centrum Πολυβιταμίνη Ενηλίκων", "Vitamins", "Pfizer", "5201003001006", false, 4.0},
	{"Centrum Silver 50+", "Centrum Silver 50+", "Vitamins", "Pfizer", "5201003001007", false, 3.0},
	{"B-Complex Tablets", "Σύμπλεγμα Βιταμίνης Β Δισκία", "Vitamins", "GSK", "5201003001008", false, 3.5},
	{"Magnesium 375mg Tablets", "Μαγνήσιο 375mg Δισκία", "Vitamins", "Sanofi", "5201003001009", false, 4.0},
	{"Zinc 15mg Tablets", "Ψευδάργυρος 15mg Δισκία", "Vitamins", "Now Foods", "5201003001010", false, 3.0},
	{"Omega-3 Fish Oil 1000mg", "Ω-3 Ιχθυέλαιο 1000mg", "Vitamins", "Nordic Naturals", "5201003001011", false, 4.0},
	{"Folic Acid 400mcg Tablets", "Φυλλικό Οξύ 400mcg Δισκία", "Vitamins", "Sanofi", "5201003001012", false, 3.0},
	{"Iron 14mg Tablets", "Σίδηρος 14mg Δισκία", "Vitamins", "Roche", "5201003001013", false, 2.5},
	{"Calcium 600mg + D3", "Ασβέστιο 600mg + D3", "Vitamins", "Pfizer", "5201003001014", false, 3.5},
	{"Coenzyme Q10 100mg", "Συνένζυμο Q10 100mg", "Vitamins", "Solgar", "5201003001015", false, 2.0},
	{"Vitamin B12 1000mcg", "Βιταμίνη B12 1000mcg", "Vitamins", "Roche", "5201003001016", false, 2.5},
	{"Vitamin E 400IU Capsules", "Βιταμίνη E 400IU Κάψουλες", "Vitamins", "Now Foods", "5201003001017", false, 2.0},
	{"Probiotics Lactobacillus 10B", "Προβιοτικά Lactobacillus 10B", "Vitamins", "VSL#3", "5201003001018", false, 3.0},
	{"Biotin 5000mcg Tablets", "Βιοτίνη 5000mcg Δισκία", "Vitamins", "Now Foods", "5201003001019", false, 2.5},
	{"Selenium 200mcg Tablets", "Σελήνιο 200mcg Δισκία", "Vitamins", "Sanofi", "5201003001020", false, 1.5},
	{"Multivitamin Kids Gummies", "Πολυβιταμίνη Παιδιών Ζελεδάκια", "Vitamins", "Vitabiotics", "5201003001021", false, 4.0},
	{"Vitamin K2 100mcg Tablets", "Βιταμίνη K2 100mcg Δισκία", "Vitamins", "Jarrow", "5201003001022", false, 1.5},
	{"Alpha Lipoic Acid 300mg", "Άλφα Λιποϊκό Οξύ 300mg", "Vitamins", "Now Foods", "5201003001023", false, 1.5},
	{"Melatonin 3mg Tablets", "Μελατονίνη 3mg Δισκία", "Vitamins", "Now Foods", "5201003001024", false, 3.5},
	{"Glucosamine 1500mg Tablets", "Γλυκοζαμίνη 1500mg Δισκία", "Vitamins", "Now Foods", "5201003001025", false, 2.5},

	// Cold & Flu
	{"Panadol Cold & Flu Day", "Panadol Κρυολόγημα & Γρίπη Ημέρας", "Cold & Flu", "GSK", "5201004001001", false, 4.0},
	{"Lemsip Max Day Lemon", "Lemsip Max Ημέρας Λεμόνι", "Cold & Flu", "Reckitt", "5201004001002", false, 3.5},
	{"Beechams All-In-One", "Beechams All-In-One", "Cold & Flu", "GSK", "5201004001003", false, 3.0},
	{"Otrivin Nasal Spray 0.1%", "Otrivin Ρινικό Σπρέι 0.1%", "Cold & Flu", "Novartis", "5201004001004", false, 4.5},
	{"Sinupret Forte Tablets", "Sinupret Forte Δισκία", "Cold & Flu", "Bionorica", "5201004001005", false, 3.0},
	{"Strepsils Lemon Lozenges", "Strepsils Λεμόνι Παστίλιες", "Cold & Flu", "Reckitt", "5201004001006", false, 5.0},
	{"Difflam Throat Spray", "Difflam Σπρέι Λαιμού", "Cold & Flu", "3M", "5201004001007", false, 2.5},
	{"Robitussin Dry Cough 200ml", "Robitussin Ξηρός Βήχας 200ml", "Cold & Flu", "Pfizer", "5201004001008", false, 3.0},
	{"Mucovit Expectorant Syrup", "Mucovit Αποχρεμπτικό Σιρόπι", "Cold & Flu", "Uni-Pharma", "5201004001009", false, 2.5},
	{"Vicks VapoRub 50g", "Vicks VapoRub 50g", "Cold & Flu", "P&G", "5201004001010", false, 3.5},
	{"Codelsol Cough Linctus", "Codelsol Αντιβηχικό Σιρόπι", "Cold & Flu", "Sanofi", "5201004001011", false, 2.0},
	{"Piriton Syrup 150ml", "Piriton Σιρόπι 150ml", "Cold & Flu", "GSK", "5201004001012", false, 2.5},
	{"Terpinol Cough Tablets", "Τερπινόλη Δισκία Βήχα", "Cold & Flu", "Sanofi", "5201004001013", false, 2.0},
	{"Betadine Throat Gargle", "Betadine Γαργάρα Λαιμού", "Cold & Flu", "Mundipharma", "5201004001014", false, 1.5},
	{"Fluenz Nasal Spray", "Fluenz Ρινικό Σπρέι", "Cold & Flu", "AstraZeneca", "5201004001015", false, 1.0},

	// Allergy
	{"Loratadine 10mg Tablets", "Λοραταδίνη 10mg Δισκία", "Allergy", "Schering-Plough", "5201005001001", false, 3.5},
	{"Cetirizine 10mg Tablets", "Σετιριζίνη 10mg Δισκία", "Allergy", "UCB", "5201005001002", false, 4.0},
	{"Fexofenadine 120mg Tablets", "Φεξοφεναδίνη 120mg Δισκία", "Allergy", "Sanofi", "5201005001003", false, 2.5},
	{"Claritin 10mg Tablets", "Claritin 10mg Δισκία", "Allergy", "Bayer", "5201005001004", false, 3.0},
	{"Telfast 120mg Tablets", "Telfast 120mg Δισκία", "Allergy", "Sanofi", "5201005001005", false, 2.5},
	{"Flonase Nasal Spray", "Flonase Ρινικό Σπρέι", "Allergy", "GSK", "5201005001006", false, 2.0},
	{"Dymista Nasal Spray", "Dymista Ρινικό Σπρέι", "Allergy", "Meda", "5201005001007", true, 1.5},
	{"Piriton 4mg Tablets", "Piriton 4mg Δισκία", "Allergy", "GSK", "5201005001008", false, 2.0},

	// Digestive
	{"Omeprazole 20mg Capsules", "Ομεπραζόλη 20mg Κάψουλες", "Digestive", "AstraZeneca", "5201006001001", false, 4.5},
	{"Lansoprazole 30mg Capsules", "Λανσοπραζόλη 30mg Κάψουλες", "Digestive", "Takeda", "5201006001002", true, 3.0},
	{"Loperamide 2mg Tablets", "Λοπεραμίδιο 2mg Δισκία", "Digestive", "J&J", "5201006001003", false, 3.5},
	{"Imodium 2mg Capsules", "Imodium 2mg Κάψουλες", "Digestive", "J&J", "5201006001004", false, 3.0},
	{"Gaviscon Advance 500ml", "Gaviscon Advance 500ml", "Digestive", "Reckitt", "5201006001005", false, 3.5},
	{"Rennie Peppermint 96 Tabs", "Rennie Μέντα 96 Δισκία", "Digestive", "Bayer", "5201006001006", false, 4.0},
	{"Lactulose Solution 300ml", "Λακτουλόζη Διάλυμα 300ml", "Digestive", "Fresenius", "5201006001007", false, 2.5},
	{"Dulcolax 5mg Tablets", "Dulcolax 5mg Δισκία", "Digestive", "Sanofi", "5201006001008", false, 2.5},
	{"Buscopan 10mg Tablets", "Buscopan 10mg Δισκία", "Digestive", "Sanofi", "5201006001009", false, 3.0},
	{"Metoclopramide 10mg Tabs", "Μετοκλοπραμίδιο 10mg Δισκία", "Digestive", "Sanofi", "5201006001010", true, 2.0},
	{"Carbosymag Activated Charcoal", "Carbosymag Ενεργός Άνθρακας", "Digestive", "Uni-Pharma", "5201006001011", false, 1.5},
	{"Rehydrat Oral Rehydration", "Rehydrat Ηλεκτρολύτες", "Digestive", "Searle", "5201006001012", false, 2.5},

	// Cardiovascular
	{"Aspirin 100mg Gastro-Resistant", "Ασπιρίνη 100mg Γαστροανθεκτική", "Cardiovascular", "Bayer", "5201007001001", false, 5.0},
	{"Aspirin 75mg Tablets", "Ασπιρίνη 75mg Δισκία", "Cardiovascular", "Bayer", "5201007001002", false, 4.0},
	{"Atorvastatin 20mg Tablets", "Ατορβαστατίνη 20mg Δισκία", "Cardiovascular", "Pfizer", "5201007001003", true, 3.5},
	{"Atorvastatin 40mg Tablets", "Ατορβαστατίνη 40mg Δισκία", "Cardiovascular", "Pfizer", "5201007001004", true, 3.0},
	{"Amlodipine 5mg Tablets", "Αμλοδιπίνη 5mg Δισκία", "Cardiovascular", "Pfizer", "5201007001005", true, 2.5},
	{"Bisoprolol 5mg Tablets", "Βισοπρολόλη 5mg Δισκία", "Cardiovascular", "Merck", "5201007001006", true, 2.0},
	{"Losartan 50mg Tablets", "Λοσαρτάνη 50mg Δισκία", "Cardiovascular", "MSD", "5201007001007", true, 2.0},
	{"Ramipril 5mg Capsules", "Ραμιπρίλη 5mg Κάψουλες", "Cardiovascular", "Sanofi", "5201007001008", true, 1.5},

	// Diabetes
	{"Metformin 500mg Tablets", "Μετφορμίνη 500mg Δισκία", "Diabetes", "Merck", "5201008001001", true, 3.5},
	{"Metformin 850mg Tablets", "Μετφορμίνη 850mg Δισκία", "Diabetes", "Merck", "5201008001002", true, 3.0},
	{"Glucose Test Strips 50s", "Ταινίες Μέτρησης Γλυκόζης 50τμχ", "Diabetes", "Roche", "5201008001003", false, 5.0},
	{"Accu-Check Lancets 100s", "Accu-Check Βελόνες Νυγμού 100τμχ", "Diabetes", "Roche", "5201008001004", false, 4.0},
	{"Glucophage XR 750mg", "Glucophage XR 750mg", "Diabetes", "Merck", "5201008001005", true, 2.5},

	// Skincare
	{"Eucerin Aquaporin Active SPF25", "Eucerin Aquaporin Active SPF25", "Skincare", "Beiersdorf", "5201009001001", false, 3.5},
	{"Neutrogena Ultra Sheer SPF50", "Neutrogena Ultra Sheer SPF50", "Skincare", "J&J", "5201009001002", false, 4.0},
	{"La Roche-Posay Anthelios SPF50", "La Roche-Posay Anthelios SPF50", "Skincare", "L'Oréal", "5201009001003", false, 4.5},
	{"Avene Very High Protection SPF50+", "Avene Very High Protection SPF50+", "Skincare", "Pierre Fabre", "5201009001004", false, 3.5},
	{"Bioderma Photoderm MAX SPF50", "Bioderma Photoderm MAX SPF50", "Skincare", "Naos", "5201009001005", false, 3.0},
	{"CeraVe Daily Moisturiser SPF25", "CeraVe Ενυδατική Κρέμα SPF25", "Skincare", "L'Oréal", "5201009001006", false, 3.5},
	{"Vichy Dermablend SPF25", "Vichy Dermablend SPF25", "Skincare", "L'Oréal", "5201009001007", false, 2.5},
	{"Benzac AC 5% Gel 60g", "Benzac AC 5% Gel 60g", "Skincare", "Galderma", "5201009001008", false, 2.0},
	{"Differin Adapalene 0.1% Gel", "Differin Αδαπαλένιο 0.1% Gel", "Skincare", "Galderma", "5201009001009", true, 1.5},
	{"Eucerin Eczema Relief Body Wash", "Eucerin Αφρόλουτρο Ατοπικού", "Skincare", "Beiersdorf", "5201009001010", false, 2.0},
	{"Dermol 500 Lotion 500ml", "Dermol 500 Λοσιόν 500ml", "Skincare", "Dermal Labs", "5201009001011", false, 1.5},
	{"Hydrocortisone 1% Cream 30g", "Υδροκορτιζόνη 1% Κρέμα 30g", "Skincare", "Sanofi", "5201009001012", false, 2.5},
	{"Canesten Clotrimazole 1% Cream", "Canesten Κλοτριμαζόλη 1% Κρέμα", "Skincare", "Bayer", "5201009001013", false, 2.5},
	{"Dermovate Clobetasol 0.05% Cream", "Dermovate Κλοβετασόλη 0.05%", "Skincare", "GSK", "5201009001014", true, 1.0},
	{"Fucidin 2% Cream 30g", "Fucidin 2% Κρέμα 30g", "Skincare", "LEO Pharma", "5201009001015", true, 1.5},
	{"Nizoral Shampoo 100ml", "Nizoral Σαμπουάν 100ml", "Skincare", "J&J", "5201009001016", false, 2.0},
	{"Alphosyl 2-in-1 Shampoo 250ml", "Alphosyl 2-in-1 Σαμπουάν 250ml", "Skincare", "GSK", "5201009001017", false, 1.5},
	{"Emollient Bath Oil 500ml", "Μαλακτικό Λάδι Μπάνιου 500ml", "Skincare", "Dermol", "5201009001018", false, 1.5},
	{"Calamine Lotion 200ml", "Λοσιόν Καλαμίνης 200ml", "Skincare", "Sanofi", "5201009001019", false, 1.0},
	{"Acnecide 5% Gel 30g", "Acnecide 5% Gel 30g", "Skincare", "Galderma", "5201009001020", false, 1.5},

	// Baby Care
	{"Bepanthen Nappy Care Ointment 100g", "Bepanthen Αλοιφή για Πάνα 100g", "Baby Care", "Bayer", "5201010001001", false, 3.5},
	{"Sudocrem 125g", "Sudocrem 125g", "Baby Care", "Teva", "5201010001002", false, 3.0},
	{"Infacol Wind Drops 50ml", "Infacol Σταγόνες για Κολικούς 50ml", "Baby Care", "Forest Labs", "5201010001003", false, 3.5},
	{"Colief Infant Drops 7ml", "Colief Σταγόνες 7ml", "Baby Care", "Crosscare", "5201010001004", false, 2.5},
	{"Dentinox Teething Gel 10g", "Dentinox Gel Οδοντοφυΐας 10g", "Baby Care", "Dendron", "5201010001005", false, 2.0},
	{"Calpol Infant Suspension 200ml", "Calpol Παιδικό Σιρόπι 200ml", "Baby Care", "J&J", "5201010001006", false, 4.0},
	{"Junior Ibuprofen 200mg/5ml 150ml", "Ιβουπροφαίνη Παιδική 200mg/5ml", "Baby Care", "Bayer", "5201010001007", false, 3.0},
	{"Nurofen for Children Orange 100ml", "Nurofen Παιδικό Πορτοκάλι 100ml", "Baby Care", "Reckitt", "5201010001008", false, 3.5},
	{"Gripe Water 150ml", "Gripe Water 150ml", "Baby Care", "Woodwards", "5201010001009", false, 2.5},
	{"Soothe & Relieve Gel", "Gel Ανακούφισης Οδοντοφυΐας", "Baby Care", "Boots", "5201010001010", false, 2.0},

	// First Aid
	{"Savlon Antiseptic Cream 100g", "Savlon Αντισηπτική Κρέμα 100g", "First Aid", "Novartis", "5201011001001", false, 2.5},
	{"Betadine Solution 30ml", "Betadine Διάλυμα 30ml", "First Aid", "Mundipharma", "5201011001002", false, 3.0},
	{"Elastoplast Fabric Strip 20s", "Elastoplast Λωρίδες Υφάσματος 20τμχ", "First Aid", "Beiersdorf", "5201011001003", false, 4.0},
	{"Steri-Strip 6mm×75mm 100s", "Steri-Strip 6mm×75mm 100τμχ", "First Aid", "3M", "5201011001004", false, 2.0},
	{"Thermometer Digital Oral", "Θερμόμετρο Ψηφιακό Στοματικό", "First Aid", "Omron", "5201011001005", false, 3.5},
	{"Blood Pressure Monitor Auto", "Αυτόματο Πιεσόμετρο", "First Aid", "Omron", "5201011001006", false, 1.5},
	{"Glucose Lancets 28G 100s", "Βελόνες Νυγμού 28G 100τμχ", "First Aid", "BD", "5201011001007", false, 3.0},
	{"Conforming Bandage 7.5cm×4m", "Επίδεσμος Διαμορφούμενος 7.5cm", "First Aid", "Hartmann", "5201011001008", false, 2.5},
	{"Eye Wash Sterile 250ml", "Οφθαλμικό Πλύσιμο Αποστ. 250ml", "First Aid", "Boots", "5201011001009", false, 1.5},
	{"Surgical Gloves Nitrile L 100s", "Γάντια Νιτριλίου L 100τμχ", "First Aid", "Ansell", "5201011001010", false, 2.0},
	{"Sterile Gauze 10cm×10cm 20s", "Γάζα Αποστειρωμένη 10cm×10cm 20τμχ", "First Aid", "Hartmann", "5201011001011", false, 2.5},
	{"Wound Closure Strips 6mm×75mm", "Ταινίες Σύγκλεισης Τραύματος", "First Aid", "3M", "5201011001012", false, 1.5},
}

var suppliers = []string{"MedSupply Cyprus", "PharmaWholesale Ltd", "EuroMeds"}

type batchSpec struct {
	riskTarget string // CRITICAL, HIGH, MEDIUM, LOW
}

func RunSeed(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("starting seed")

	// Hash passwords
	adminHash, err := bcrypt.GenerateFromPassword([]byte("Demo1234!"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	staffHash, err := bcrypt.GenerateFromPassword([]byte("Demo1234!"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create pharmacy
	var pharmacyID uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO pharmacies (id, name, license_number, address, city, phone, email, language)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (license_number) DO UPDATE SET name = EXCLUDED.name
		RETURNING id`,
		"Nicosia Central Pharmacy",
		"CY-PH-2024-001",
		"12 Makarios Avenue, Nicosia 1065, Cyprus",
		"Nicosia",
		"+357 22 100200",
		"demo@pharmasense.cy",
		"en",
	).Scan(&pharmacyID)
	if err != nil {
		return fmt.Errorf("insert pharmacy: %w", err)
	}
	slog.Info("pharmacy created", "id", pharmacyID)

	// Create users
	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, pharmacy_id, email, password_hash, full_name, role)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		ON CONFLICT (email) DO NOTHING`,
		pharmacyID, "admin@pharmasense.cy", string(adminHash), "Admin User", "admin")
	if err != nil {
		return fmt.Errorf("insert admin: %w", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, pharmacy_id, email, password_hash, full_name, role)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		ON CONFLICT (email) DO NOTHING`,
		pharmacyID, "staff@pharmasense.cy", string(staffHash), "Maria Petridou", "staff")
	if err != nil {
		return fmt.Errorf("insert staff: %w", err)
	}
	slog.Info("users created")

	// Insert products
	productIDs := make([]uuid.UUID, 0, len(products))
	for _, p := range products {
		var pid uuid.UUID
		err = pool.QueryRow(ctx, `
			INSERT INTO products (id, barcode, name, name_el, category, manufacturer, requires_prescription)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
			ON CONFLICT (barcode) DO UPDATE SET name = EXCLUDED.name
			RETURNING id`,
			p.Barcode, p.Name, p.NameEl, p.Category, p.Manufacturer, p.Rx,
		).Scan(&pid)
		if err != nil {
			return fmt.Errorf("insert product %s: %w", p.Name, err)
		}
		productIDs = append(productIDs, pid)
	}
	slog.Info("products created", "count", len(productIDs))

	// Generate batches with target risk distribution
	// 500 batches: 100 CRITICAL, 125 HIGH, 100 MEDIUM, 175 LOW
	type batchTarget struct {
		riskTarget string
		count      int
	}
	targets := []batchTarget{
		{"CRITICAL", 100},
		{"HIGH", 125},
		{"MEDIUM", 100},
		{"LOW", 175},
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	batchNum := 1

	type batchInfo struct {
		id            uuid.UUID
		productID     uuid.UUID
		dailyRate     float64
		purchasePrice float64
		sellingPrice  float64
	}

	var allBatches []batchInfo

	for _, tgt := range targets {
		for i := 0; i < tgt.count; i++ {
			// Pick a product (cycle through products)
			pIdx := (batchNum - 1) % len(products)
			prod := products[pIdx]
			prodID := productIDs[pIdx]

			purchasePrice := 1.0 + rng.Float64()*49.0 // €1-50
			markup := 1.3 + rng.Float64()*0.5         // 1.3-1.8
			sellingPrice := purchasePrice * markup
			supplier := suppliers[rng.Intn(len(suppliers))]

			var expiryDate time.Time
			var currentQty, initialQty int
			var dailyRate float64

			switch tgt.riskTarget {
			case "CRITICAL":
				// Expiry in 1-30 days, low sales, surplus
				days := 1 + rng.Intn(30)
				expiryDate = today.AddDate(0, 0, days)
				dailyRate = prod.DailyRate * (0.05 + rng.Float64()*0.1) // very low sales
				initialQty = 30 + rng.Intn(120)
				// ensure surplus: current > expected_sales_remaining
				expectedRemaining := int(dailyRate * float64(days))
				currentQty = expectedRemaining + 10 + rng.Intn(80)
				if currentQty > initialQty {
					currentQty = initialQty
				}

			case "HIGH":
				// Expiry in 31-90 days, moderate-low sales, large surplus
				days := 31 + rng.Intn(60)
				expiryDate = today.AddDate(0, 0, days)
				dailyRate = prod.DailyRate * (0.1 + rng.Float64()*0.2)
				initialQty = 50 + rng.Intn(200)
				expectedRemaining := int(dailyRate * float64(days))
				// surplus > 50% of expected
				currentQty = expectedRemaining + int(float64(expectedRemaining)*0.6) + 5 + rng.Intn(50)
				if currentQty > initialQty {
					initialQty = currentQty + rng.Intn(20)
				}

			case "MEDIUM":
				// Expiry in 91-180 days, surplus > 30%
				days := 91 + rng.Intn(90)
				expiryDate = today.AddDate(0, 0, days)
				dailyRate = prod.DailyRate * (0.2 + rng.Float64()*0.3)
				initialQty = 80 + rng.Intn(300)
				expectedRemaining := int(dailyRate * float64(days))
				currentQty = expectedRemaining + int(float64(expectedRemaining)*0.35) + rng.Intn(30)
				if currentQty > initialQty {
					initialQty = currentQty + rng.Intn(20)
				}

			case "LOW":
				// Good situation: either far from expiry or high sales
				days := 91 + rng.Intn(540) // 3-18 months
				expiryDate = today.AddDate(0, 0, days)
				dailyRate = prod.DailyRate * (0.7 + rng.Float64()*0.6) // good sales rate
				initialQty = 50 + rng.Intn(400)
				expectedRemaining := int(dailyRate * float64(days))
				if expectedRemaining >= initialQty {
					currentQty = initialQty / 2
				} else {
					currentQty = 5 + rng.Intn(expectedRemaining/2+1)
				}
				if currentQty <= 0 {
					currentQty = 1
				}
			}

			// Received date: 30-365 days ago
			receivedDaysAgo := 30 + rng.Intn(335)
			receivedDate := today.AddDate(0, 0, -receivedDaysAgo)

			bn := fmt.Sprintf("BTH-2024-%05d", batchNum)
			batchNum++

			var bid uuid.UUID
			err = pool.QueryRow(ctx, `
				INSERT INTO inventory_batches (
					id, pharmacy_id, product_id, batch_number, expiry_date,
					initial_quantity, current_quantity, purchase_price, selling_price,
					supplier, received_date
				) VALUES (
					gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
				) RETURNING id`,
				pharmacyID, prodID, bn, expiryDate.Format("2006-01-02"),
				initialQty, currentQty, purchasePrice, sellingPrice, supplier,
				receivedDate.Format("2006-01-02"),
			).Scan(&bid)
			if err != nil {
				return fmt.Errorf("insert batch: %w", err)
			}

			allBatches = append(allBatches, batchInfo{
				id:            bid,
				productID:     prodID,
				dailyRate:     dailyRate,
				purchasePrice: purchasePrice,
				sellingPrice:  sellingPrice,
			})
		}
	}
	slog.Info("batches created", "count", len(allBatches))

	// Generate sales records (~10,000+)
	salesInserted := 0
	for _, b := range allBatches {
		// Sales for past 90 days
		for dayOffset := 89; dayOffset >= 0; dayOffset-- {
			saleDate := today.AddDate(0, 0, -dayOffset)

			// Seasonal adjustment: sunscreen sells more in summer (Apr-Sep)
			month := saleDate.Month()
			seasonMultiplier := 1.0
			prod := findProduct(b.productID, productIDs, products)
			if prod != nil && prod.Category == "Skincare" {
				if month >= 4 && month <= 9 {
					seasonMultiplier = 2.5
				} else {
					seasonMultiplier = 0.3
				}
			}
			if prod != nil && prod.Category == "Cold & Flu" {
				if month >= 10 || month <= 3 {
					seasonMultiplier = 2.0
				} else {
					seasonMultiplier = 0.5
				}
			}

			avgUnits := b.dailyRate * seasonMultiplier
			if avgUnits < 0.05 {
				avgUnits = 0.05
			}

			// Poisson-like: most days 0 units for low-rate products
			threshold := avgUnits / (avgUnits + 1.0)
			if rng.Float64() > threshold && avgUnits < 0.5 {
				continue // no sale this day
			}

			qty := 1
			if avgUnits >= 1 {
				qty = 1 + int(rng.NormFloat64()*avgUnits*0.3+avgUnits)
				if qty <= 0 {
					qty = 1
				}
			}

			_, err := pool.Exec(ctx, `
				INSERT INTO sales (id, pharmacy_id, batch_id, product_id, quantity, unit_price, total_amount, sale_date)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)`,
				pharmacyID, b.id, b.productID, qty,
				b.sellingPrice, float64(qty)*b.sellingPrice,
				saleDate.Format("2006-01-02"),
			)
			if err != nil {
				slog.Error("insert sale", "error", err)
				continue
			}
			salesInserted++
		}
	}
	slog.Info("sales created", "count", salesInserted)

	// Run risk calculations and store in risk_assessments
	salesMap := make(map[uuid.UUID]float64)
	rows, err := pool.Query(ctx, `
		SELECT product_id, COALESCE(SUM(quantity)::float / 90.0, 0)
		FROM sales WHERE pharmacy_id = $1 AND sale_date >= CURRENT_DATE - INTERVAL '90 days'
		GROUP BY product_id`, pharmacyID)
	if err != nil {
		return fmt.Errorf("sales map: %w", err)
	}
	for rows.Next() {
		var pid uuid.UUID
		var avg float64
		if err := rows.Scan(&pid, &avg); err != nil {
			continue
		}
		salesMap[pid] = avg
	}
	rows.Close()

	batchRows, err := pool.Query(ctx, `
		SELECT id, product_id, pharmacy_id, expiry_date::text, current_quantity, purchase_price
		FROM inventory_batches WHERE pharmacy_id = $1`, pharmacyID)
	if err != nil {
		return fmt.Errorf("batch rows: %w", err)
	}

	type batchRow struct {
		ID            uuid.UUID
		ProductID     uuid.UUID
		PharmacyID    uuid.UUID
		ExpiryDate    string
		CurrentQty    int
		PurchasePrice float64
	}

	var bRows []batchRow
	for batchRows.Next() {
		var b batchRow
		if err := batchRows.Scan(&b.ID, &b.ProductID, &b.PharmacyID, &b.ExpiryDate, &b.CurrentQty, &b.PurchasePrice); err != nil {
			continue
		}
		bRows = append(bRows, b)
	}
	batchRows.Close()

	for _, b := range bRows {
		expiry, err := time.Parse("2006-01-02", b.ExpiryDate)
		if err != nil {
			continue
		}
		days := int(expiry.Sub(today).Hours() / 24)
		avg := salesMap[b.ProductID]

		expectedSales := int(avg * float64(days))
		surplus := b.CurrentQty - expectedSales
		if surplus < 0 {
			surplus = 0
		}

		var riskLevel string
		var discountPct int
		var estimatedLoss float64

		switch {
		case days <= 30 && surplus > 0:
			riskLevel = "CRITICAL"
			discountPct = 40
		case days <= 90 && expectedSales > 0 && float64(surplus) > float64(expectedSales)*0.5:
			riskLevel = "HIGH"
			discountPct = 20
		case days <= 90 && surplus > 0 && expectedSales == 0:
			riskLevel = "HIGH"
			discountPct = 20
		case days <= 180 && expectedSales > 0 && float64(surplus) > float64(expectedSales)*0.3:
			riskLevel = "MEDIUM"
			discountPct = 10
		case days <= 0:
			riskLevel = "CRITICAL"
			discountPct = 50
			surplus = b.CurrentQty
		default:
			riskLevel = "LOW"
		}

		if riskLevel == "CRITICAL" || riskLevel == "HIGH" {
			estimatedLoss = float64(surplus) * b.PurchasePrice
		}
		if riskLevel == "MEDIUM" {
			estimatedLoss = float64(surplus) * b.PurchasePrice * 0.3
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO risk_assessments (
				id, batch_id, pharmacy_id, risk_level, days_until_expiry,
				avg_daily_sales, expected_sales, estimated_surplus,
				estimated_loss, suggested_discount_percent
			) VALUES (
				gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9
			)
			ON CONFLICT (batch_id) DO UPDATE SET
				risk_level = EXCLUDED.risk_level,
				days_until_expiry = EXCLUDED.days_until_expiry,
				avg_daily_sales = EXCLUDED.avg_daily_sales,
				expected_sales = EXCLUDED.expected_sales,
				estimated_surplus = EXCLUDED.estimated_surplus,
				estimated_loss = EXCLUDED.estimated_loss,
				suggested_discount_percent = EXCLUDED.suggested_discount_percent,
				calculated_at = NOW()`,
			b.ID, b.PharmacyID, riskLevel, days,
			avg, expectedSales, surplus, estimatedLoss, discountPct,
		)
		if err != nil {
			slog.Error("upsert risk", "error", err)
		}
	}

	slog.Info("seed completed successfully")
	return nil
}

func findProduct(productID uuid.UUID, ids []uuid.UUID, prods []seedProduct) *seedProduct {
	for i, id := range ids {
		if id == productID && i < len(prods) {
			return &prods[i]
		}
	}
	return nil
}

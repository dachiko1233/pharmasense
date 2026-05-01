# PharmaSense - Expiry Monitoring System

## 🎯 Project Overview

**PharmaSense** is a SaaS platform for pharmacies in Cyprus that monitors product expiry dates and helps pharmacies reduce waste and increase profits by predicting which products will expire before being sold.

**Target Market:** Independent pharmacies in Cyprus (Nicosia, Limassol, Paphos, Larnaca)

**Core Value Proposition:** "Save €1000s monthly by catching at-risk products before they expire"

**Languages:** English + Greek (i18n required)

---

## 🛠️ Tech Stack (REQUIRED - Use Exactly These)

### Backend
- **Language:** Go 1.22+
- **HTTP Framework:** Chi (`github.com/go-chi/chi/v5`)
- **Database:** PostgreSQL 16
- **SQL:** sqlc (`github.com/sqlc-dev/sqlc`) for type-safe queries
- **Migrations:** golang-migrate (`github.com/golang-migrate/migrate`)
- **Validation:** go-playground/validator
- **Auth:** JWT (`github.com/golang-jwt/jwt/v5`)
- **Logging:** slog (standard library)
- **Config:** Environment variables via `godotenv`
- **UUID:** `github.com/google/uuid`
- **Password hashing:** `golang.org/x/crypto/bcrypt`

### Frontend
- **Framework:** Next.js 15 (App Router)
- **Language:** TypeScript (strict mode)
- **Styling:** Tailwind CSS v4
- **UI Components:** shadcn/ui (manually installed components, NOT a library)
- **Forms:** react-hook-form + zod
- **Data fetching:** TanStack Query v5
- **Tables:** TanStack Table v8
- **Charts:** Recharts
- **State:** Zustand (only if needed)
- **i18n:** next-intl (English + Greek)
- **Icons:** lucide-react
- **Date handling:** date-fns

### DevOps
- **Containerization:** Docker + docker-compose
- **Database GUI:** Include adminer in docker-compose for easy DB access

---

## 📊 Core Feature: Expired Product Monitoring

This is the **ONLY** feature for this MVP. Build it deeply, not broadly.

### Feature Breakdown

#### 1. Product & Inventory Management
- Each pharmacy has products with batches
- Each batch has its own expiry date and quantity
- Same product can have multiple batches (different expiry dates)

#### 2. Sales Data Tracking
- Track sales per product per day
- Calculate average daily sales (rolling 90-day window)
- Used to predict if stock will sell before expiring

#### 3. Risk Calculation Engine (THE CORE)
This is the key algorithm. For each batch in inventory:

```
days_until_expiry = expiry_date - today
expected_sales = avg_daily_sales × days_until_expiry
surplus = current_quantity - expected_sales

risk_level:
  - CRITICAL: days_until_expiry <= 30 AND surplus > 0
  - HIGH:     days_until_expiry <= 90 AND surplus > expected_sales × 0.5
  - MEDIUM:   days_until_expiry <= 180 AND surplus > expected_sales × 0.3
  - LOW:      otherwise

estimated_loss = surplus × purchase_price (when risk is HIGH or CRITICAL)
suggested_discount:
  - CRITICAL: 30-50%
  - HIGH:     15-25%
  - MEDIUM:   10%
```

#### 4. Dashboard
- Total products at risk (count by risk level)
- Total estimated loss (€)
- Total potential savings if action taken
- Chart: Risk distribution over next 12 months
- Chart: Top 10 highest-risk products
- Quick actions: View critical alerts, export report

#### 5. Inventory Page
- Sortable/filterable table of all batches
- Filters: risk level, category, expiry range, supplier
- Search by name/barcode
- Color-coded rows by risk level
- Bulk actions: mark for discount, transfer, return to supplier

#### 6. Alerts Page
- List of all CRITICAL and HIGH risk items
- Suggested actions per item
- Action buttons: Apply Discount, Mark for Transfer, Mark Returned
- Action history

#### 7. CSV Import
- Upload inventory data (products + batches)
- Validate format
- Preview before import
- Show import results (success, errors)

#### 8. Reports
- Monthly waste reduction report
- Money saved over time (chart)
- Most problematic categories
- Action effectiveness (how many alerts led to action)

---

## 🗄️ Database Schema

```sql
-- Pharmacies (multi-tenant)
CREATE TABLE pharmacies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    license_number VARCHAR(100) UNIQUE NOT NULL,
    address TEXT,
    city VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Users (pharmacy staff)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pharmacy_id UUID NOT NULL REFERENCES pharmacies(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'staff', -- 'admin', 'staff'
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Product catalog (shared across pharmacies)
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    barcode VARCHAR(100) UNIQUE,
    name VARCHAR(500) NOT NULL,
    name_el VARCHAR(500), -- Greek name
    category VARCHAR(100), -- 'antibiotics', 'vitamins', 'painkillers', etc.
    manufacturer VARCHAR(255),
    requires_prescription BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Inventory batches (the heart of the system)
CREATE TABLE inventory_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pharmacy_id UUID NOT NULL REFERENCES pharmacies(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    batch_number VARCHAR(100),
    expiry_date DATE NOT NULL,
    initial_quantity INTEGER NOT NULL,
    current_quantity INTEGER NOT NULL,
    purchase_price DECIMAL(10,2) NOT NULL,
    selling_price DECIMAL(10,2) NOT NULL,
    supplier VARCHAR(255),
    received_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_batches_pharmacy_expiry ON inventory_batches(pharmacy_id, expiry_date);
CREATE INDEX idx_batches_product ON inventory_batches(product_id);

-- Sales tracking
CREATE TABLE sales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pharmacy_id UUID NOT NULL REFERENCES pharmacies(id) ON DELETE CASCADE,
    batch_id UUID NOT NULL REFERENCES inventory_batches(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    sale_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sales_pharmacy_date ON sales(pharmacy_id, sale_date);
CREATE INDEX idx_sales_product ON sales(product_id, sale_date);

-- Risk assessments (calculated and stored)
CREATE TABLE risk_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES inventory_batches(id) ON DELETE CASCADE,
    pharmacy_id UUID NOT NULL REFERENCES pharmacies(id) ON DELETE CASCADE,
    risk_level VARCHAR(20) NOT NULL, -- 'CRITICAL', 'HIGH', 'MEDIUM', 'LOW'
    days_until_expiry INTEGER NOT NULL,
    avg_daily_sales DECIMAL(10,2),
    expected_sales INTEGER,
    estimated_surplus INTEGER,
    estimated_loss DECIMAL(10,2),
    suggested_discount_percent INTEGER,
    calculated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_risk_pharmacy_level ON risk_assessments(pharmacy_id, risk_level);

-- Actions taken on alerts
CREATE TABLE alert_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES inventory_batches(id),
    pharmacy_id UUID NOT NULL REFERENCES pharmacies(id),
    user_id UUID NOT NULL REFERENCES users(id),
    action_type VARCHAR(50) NOT NULL, -- 'DISCOUNT', 'TRANSFER', 'RETURN', 'DISMISS'
    discount_percent INTEGER,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 🌱 Demo Data Requirements

The seed script must create REALISTIC demo data:

### 1. Demo Pharmacy
- Name: "Nicosia Central Pharmacy"
- License: "CY-PH-2024-001"
- Address: "12 Makarios Avenue, Nicosia 1065, Cyprus"
- City: "Nicosia"
- Email: "demo@pharmasense.cy"

### 2. Demo Users
- Admin: `admin@pharmasense.cy` / password: `Demo1234!`
- Staff: `staff@pharmasense.cy` / password: `Demo1234!`

### 3. Products (~150 realistic pharmacy products)
Mix of categories with bilingual names:
- **Painkillers:** Paracetamol, Ibuprofen, Aspirin, Naproxen, Diclofenac
- **Antibiotics:** Amoxicillin, Azithromycin, Ciprofloxacin
- **Vitamins:** Vitamin D3, Vitamin C, Multivitamin, B-Complex, Magnesium, Zinc
- **Cold/Flu:** Cough syrups, decongestants, throat lozenges
- **Allergy:** Loratadine, Cetirizine, Fexofenadine
- **Digestive:** Omeprazole, Loperamide, Antacids
- **Cardiovascular:** Aspirin low-dose, Atorvastatin
- **Diabetes:** Metformin, Glucose strips
- **Skincare:** Sunscreens, moisturizers, acne treatments
- **Baby care:** Diapers, formula, baby vitamins
- **First aid:** Bandages, antiseptics, thermometers

Each product should have:
- Real-sounding name
- Greek translation
- Realistic category
- Manufacturer (e.g., "Pfizer", "GSK", "Bayer", "Sanofi", "Roche")

### 4. Inventory Batches (~500 batches)
Realistic distribution:
- **20% CRITICAL** (expiring within 30 days, with surplus)
- **25% HIGH risk** (expiring 30-90 days, with surplus)
- **20% MEDIUM risk** (expiring 90-180 days)
- **35% LOW risk** (good situation)

Each batch:
- Realistic batch number (e.g., "BTH-2024-XXXXX")
- Purchase price between €1 - €50
- Selling price = purchase × 1.3-1.8 markup
- Quantities between 10-500 units
- Various suppliers: "MedSupply Cyprus", "PharmaWholesale Ltd", "EuroMeds"

### 5. Sales History (~10,000 sales records)
- Last 90 days of sales
- Realistic patterns:
  - Painkillers: 5-15 sales/day
  - Vitamins: 2-8 sales/day
  - Antibiotics: 1-5 sales/day (prescription)
  - Sunscreens: higher in summer months
  - Cold meds: higher in winter
- Some products with declining sales (creating risk)
- Some products with steady sales

### 6. Pre-calculated Risk Assessments
Run the risk calculation engine on seed data so the dashboard shows real numbers immediately.

---

## 🎨 UI/UX Requirements

### Design Principles
- **Clean & Professional** - pharmacies are conservative
- **Mobile-responsive** - many pharmacists use tablets
- **Fast** - optimistic updates, skeleton loaders
- **Accessible** - WCAG AA compliance
- **Bilingual** - language switcher in header

### Color Palette (Tailwind)
```
Primary: emerald-600 (trust, health)
Danger: red-600 (CRITICAL)
Warning: orange-500 (HIGH)
Caution: yellow-500 (MEDIUM)
Success: green-500 (LOW/safe)
Neutral: slate-* shades
Background: white / slate-50
```

### Layout
- **Sidebar navigation** (collapsible on mobile)
- **Top header** with pharmacy name, language switcher, user menu
- **Main content area** with breadcrumbs

### Pages Required

#### `/login`
- Clean login form
- Email + password
- "Demo credentials" hint visible
- Language switcher

#### `/dashboard` (default after login)
- 4 KPI cards at top:
  1. Critical Items Count (red)
  2. Estimated Loss (€) (red)
  3. Potential Savings (€) (green)
  4. Total Inventory Value (€) (slate)
- Chart 1: "Expiry Timeline" - bar chart showing # of batches expiring per month next 12 months, color-coded by risk
- Chart 2: "Top 10 At-Risk Products" - horizontal bar chart with estimated loss
- Recent critical alerts list (top 5)
- "View All Alerts" button

#### `/inventory`
- Search bar (top)
- Filters sidebar: Risk Level, Category, Expiry Range, Supplier
- Data table with columns:
  - Product Name (with Greek name on hover)
  - Category
  - Batch #
  - Expiry Date (with days remaining)
  - Quantity
  - Risk Level (colored badge)
  - Estimated Loss
  - Actions
- Pagination
- Export to CSV button

#### `/alerts`
- Tabs: "Critical (X)", "High Risk (X)", "All Active"
- Each alert card:
  - Product name + image placeholder
  - Risk badge
  - Days until expiry
  - Current quantity vs Expected sales
  - Estimated loss
  - Suggested action with reasoning
  - Action buttons: "Apply Discount", "Transfer", "Return to Supplier", "Dismiss"
- Action history at bottom

#### `/import`
- Drag-and-drop CSV upload
- Download template button
- Preview table after upload
- Validation errors shown inline
- Import button with progress
- Success summary

#### `/reports`
- Date range selector
- Chart: Money saved over time
- Chart: Waste reduction trend
- Chart: Action effectiveness
- Most problematic categories table
- Export PDF button (placeholder OK)

#### `/settings`
- Pharmacy info
- User management
- Language preference
- Notification preferences (placeholders)

---

## 🔌 API Endpoints

All endpoints prefixed with `/api/v1`. JWT auth required except `/auth/*`.

### Auth
- `POST /auth/login` - Login with email/password
- `POST /auth/logout`
- `GET /auth/me` - Current user info

### Pharmacy
- `GET /pharmacy` - Current pharmacy info
- `PATCH /pharmacy` - Update info

### Products
- `GET /products` - List with search/filter
- `GET /products/:id`
- `POST /products` - Create
- `PATCH /products/:id`

### Inventory
- `GET /inventory` - List batches with filters
- `GET /inventory/:id`
- `POST /inventory` - Create batch
- `PATCH /inventory/:id`
- `DELETE /inventory/:id`
- `POST /inventory/import` - CSV import

### Sales
- `GET /sales` - List with filters
- `POST /sales` - Record sale
- `GET /sales/stats` - Aggregated stats

### Risk
- `GET /risk/dashboard` - Dashboard data
- `GET /risk/assessments` - List with filters (by risk level)
- `POST /risk/recalculate` - Trigger recalculation
- `GET /risk/timeline` - Monthly expiry timeline data

### Alerts
- `GET /alerts` - Active alerts
- `POST /alerts/:batch_id/action` - Take action

### Reports
- `GET /reports/savings` - Savings over time
- `GET /reports/waste` - Waste reduction
- `GET /reports/categories` - By category

---

## 📁 Project Structure

```
pharmacy-saas/
├── backend/
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go
│   │   └── seed/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── domain/
│   │   │   ├── pharmacy.go
│   │   │   ├── product.go
│   │   │   ├── inventory.go
│   │   │   ├── sales.go
│   │   │   └── risk.go
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── inventory.go
│   │   │   ├── risk.go
│   │   │   └── ...
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   ├── services/
│   │   │   ├── auth_service.go
│   │   │   ├── risk_engine.go
│   │   │   └── ...
│   │   ├── repository/
│   │   │   └── (sqlc generated)
│   │   └── server/
│   │       └── server.go
│   ├── migrations/
│   │   ├── 000001_init.up.sql
│   │   └── 000001_init.down.sql
│   ├── queries/
│   │   ├── pharmacies.sql
│   │   ├── products.sql
│   │   ├── inventory.sql
│   │   ├── sales.sql
│   │   └── risk.sql
│   ├── sqlc.yaml
│   ├── go.mod
│   ├── Dockerfile
│   └── .env.example
│
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── [locale]/
│   │   │   │   ├── (auth)/
│   │   │   │   │   └── login/
│   │   │   │   ├── (dashboard)/
│   │   │   │   │   ├── dashboard/
│   │   │   │   │   ├── inventory/
│   │   │   │   │   ├── alerts/
│   │   │   │   │   ├── import/
│   │   │   │   │   ├── reports/
│   │   │   │   │   └── settings/
│   │   │   │   └── layout.tsx
│   │   │   └── api/
│   │   ├── components/
│   │   │   ├── ui/        # shadcn components
│   │   │   ├── layout/
│   │   │   ├── dashboard/
│   │   │   ├── inventory/
│   │   │   └── alerts/
│   │   ├── lib/
│   │   │   ├── api/       # API client
│   │   │   ├── hooks/
│   │   │   └── utils/
│   │   ├── types/
│   │   └── messages/
│   │       ├── en.json
│   │       └── el.json
│   ├── package.json
│   ├── tailwind.config.ts
│   ├── next.config.js
│   ├── tsconfig.json
│   └── Dockerfile
│
├── docker-compose.yml
├── Makefile
├── README.md
└── .gitignore
```

---

## 🐳 Docker Setup

`docker-compose.yml` should include:
- **postgres** (port 5432, with health check)
- **adminer** (port 8080, for DB GUI)
- **backend** (port 3001, depends on postgres)
- **frontend** (port 3000, depends on backend)

Use volumes for hot reload in development.

---

## 🚀 Setup Commands (Makefile)

```makefile
make setup        # Install all dependencies
make dev          # Start everything with docker-compose
make migrate-up   # Run migrations
make seed         # Seed demo data
make sqlc         # Regenerate SQL code
make test         # Run all tests
make build        # Build production binaries
```

---

## 📝 README.md Requirements

Must include:
1. Project overview
2. Demo credentials
3. Quick start (3 commands max)
4. Architecture diagram
5. Tech stack list
6. Screenshots placeholders
7. Roadmap
8. License

---

## ✅ Acceptance Criteria

The demo is complete when:

1. ✅ `docker-compose up` starts everything
2. ✅ Demo data is seeded automatically on first run
3. ✅ Login works with demo credentials
4. ✅ Dashboard shows realistic numbers (NOT €0 or "0 items")
5. ✅ Inventory page has 500+ batches with varied risk levels
6. ✅ Alerts page shows real CRITICAL items requiring action
7. ✅ Charts render real data
8. ✅ Language switcher works (EN ↔ EL)
9. ✅ Mobile responsive (test at 375px width)
10. ✅ All major actions work (apply discount, dismiss, etc.)
11. ✅ CSV import accepts a sample file
12. ✅ No console errors
13. ✅ Loading states everywhere
14. ✅ Empty states designed (when filters return nothing)

---

## 🎯 Demo Script (For Presenting to Pharmacies)

When showing to pharmacies, the flow should be:

1. **Login** → "This is your dashboard"
2. **Dashboard** → "You're losing €X this month on expiring products"
3. **Click on Critical Alerts** → "Here are the items expiring soon"
4. **Show one alert in detail** → "This Aspirin batch will expire in 25 days. You have 80 units, but only sell 1/day. We suggest 30% discount."
5. **Apply discount action** → "One click and we record the action"
6. **Reports page** → "After 3 months you can see how much you saved"
7. **Language switch** → "Works in Greek too"

This flow MUST be smooth and impressive.

---

## 🔒 Security Requirements

- Passwords: bcrypt hashing
- JWT: 24h expiry, secure secret in env
- SQL: parameterized queries only (sqlc handles this)
- CORS: configured for frontend origin
- Rate limiting on auth endpoints
- HTTPS in production (note in README)
- No secrets in code

---

## 🎨 Visual Polish Notes

- Use **emerald-600** as primary (NOT default blue)
- Risk level badges: rounded-full with colored bg + dark text
- Empty states with friendly illustrations (use placeholders)
- Loading: use shadcn Skeleton component
- Success toasts: use shadcn Toast/Sonner
- Smooth transitions (200ms) on hovers and state changes
- Cards with subtle shadow (`shadow-sm`) and `border border-slate-200`
- Use `lucide-react` icons throughout

---

## 🚫 What NOT to Build (For Now)

To stay focused, DO NOT include:
- ❌ AI chatbot (later phase)
- ❌ POS integrations (later)
- ❌ Multi-pharmacy chains (later)
- ❌ Mobile native app (later)
- ❌ Email/SMS notifications (later)
- ❌ Payment processing
- ❌ User registration (only seeded users)
- ❌ Demand forecasting ML (use simple rules now)

---

## 💡 Implementation Priority

If running out of time, build in this order:
1. Backend: DB schema + migrations
2. Backend: Auth + basic CRUD
3. Backend: Risk engine
4. Backend: Seed script
5. Frontend: Login + layout
6. Frontend: Dashboard
7. Frontend: Inventory
8. Frontend: Alerts
9. Frontend: i18n
10. Frontend: Reports
11. Frontend: CSV import
12. Polish + Docker setup

---

## 🎬 Final Note for Claude Code

Please:
- Write CLEAN, idiomatic Go and TypeScript
- Add helpful comments where logic is non-obvious
- Use proper error handling everywhere
- Make the UI BEAUTIFUL — this will be shown to potential customers
- Generate REALISTIC demo data (not "Product 1, Product 2, ...")
- Test that everything works end-to-end before declaring done

Quality > Quantity. A polished demo of one feature beats a broken demo of ten features.
# Geo-Hedger

A strategic CLI tool for dynamic asset rebalancing and scenario planning, engineered for the era of geopolitical paradigm shifts.

*(※ For the Japanese version of this README, please see [here](./README.ja.md).)*

---

## 1. Background & Philosophy

**The era of mindlessly dumping funds into "All-Country World Indexes" and assuming safety is over.**

Investing in a global stock index is a strategy that captures profits under the assumption of a peaceful world supported by global supply chains. However, the world is now at a historic turning point, hinting at the end of unipolar hegemony and a paradigm shift toward a fragmented, multi-polar global system. For individuals living in nations without abundant natural resources, we are forced to make high-stakes choices: "Which system will our country belong to?" and "How do we protect our assets amidst rapid inflation?"

*Geo-Hedger* is a strategic Command Line Interface (CLI) tool that rejects reliance on centralized "consensus answers." Instead, it empowers users to project their own geopolitical scenarios and actively manage their portfolio's risk through dynamic rebalancing.

*(Note: A web frontend utilizing Next.js is planned for future development to enhance UX. Furthermore, anticipating an unstable environment caused by future global energy insecurities, I have intentionally chosen to store data in local JSON files to ensure offline resilience and self-sovereignty.)*

---

## 2. Core Features

### ① Dynamic Scenario (Phase) Management
* Define custom geopolitical scenarios as "Phases."
* Document and store the specific Rationale behind allocation percentages, concrete Actions to take, and anticipated Risks for each phase.

### ② Unified Family Wealth Management (Multi-Account Support)
* Track assets scattered across family members, multiple brokerages, and cold/hot wallets, mapping exactly who holds what and where.
* Aggregate these to provide a bird's-eye view of the family's true asset allocation. 
* *Purpose: Solving the ultimate challenge of this era—how to survive and protect wealth as a family unit.*

### ③ Real-Time Asset Valuation with Fallback
* Securely fetch real-time market rates (forex, gold, equities, and crypto) from external sources (such as Google Sheets) to calculate instantaneous base currency equivalents.
* If external rate fetching fails due to network or energy instability, the system supports manual rate overrides to ensure continuous operation.

### ④ Automated Rebalance Guides & Next Actions
* Visualize the deviation between your "Current Asset Allocation" and the "Target Allocation" of your selected phase.
* Beyond just showing numbers, the tool automatically generates a concrete To-Do list detailing **"Who needs to trade What, Where, and by How Much"** to achieve the target portfolio.

---

## 3. Requirements & Data Structure

### Prerequisites
* This tool manages and visualizes "surplus funds" (investable capital) only, excluding emergency living funds.

### Data Structure (`state.json`)
Data is persisted in a local JSON file with the following schema, designed for offline resilience and future easy migration to relational databases.

```json
{
  "updated_at": "2026-04-04T10:00:00Z",
  "selected_phase_id": "phase_2",
  "phases": {
    "phase_1": {
      "name": "Phase 1: Normalcy & Weak Yen",
      "scenario": "Status quo maintained. Moderate inflation and a weak yen continue.",
      "action": "Minimize cash holdings; keep heavy exposure to Western assets (Gold) and stateless assets (BTC).",
      "targets": { "JPY": 10.0, "USD": 10.0, "GLDM": 30.0, "2800": 20.0, "BTC": 30.0 }
    },
    "phase_2": {
      "name": "Phase 2: Geopolitical Crisis & Crash",
      "scenario": "Geopolitical risks materialize, causing stock market crashes and temporary USD spikes.",
      "action": "Sell high-priced USD for JPY, and buy bottomed-out BTC and Eastern assets (2800).",
      "targets": { "JPY": 30.0, "USD": 0.0, "GLDM": 20.0, "2800": 20.0, "BTC": 30.0 }
    }
  },
  "assets": {
    "Member_A (SBI)": { "JPY": 1000000, "USD": 5000, "BTC": 0.2 },
    "Member_B (Rakuten)": { "JPY": 500000, "GLDM": 50 }
  }
}
```

---

## 4. UI/UX (Expected CLI Output)

When executed, the system outputs the current status and the immediate Next Actions:

```text
================================================================================
【Target Phase】Phase 2: Geopolitical Crisis & Crash
【Scenario】Geopolitical risks materialize, causing stock crashes & USD spikes.
【Action Plan】Sell high-priced USD for JPY, and buy bottomed-out BTC & 2800.
================================================================================

Total Surplus Assets: 12,500,000 JPY

■ Asset Allocation & Deviations from Target Phase
--------------------------------------------------------------------------------
[Symbol] | Current % | Target % | Deviation (Target - Current)
--------------------------------------------------------------------------------
 JPY     |  12.0 %  |  30.0 % | +2,250,000 JPY (Shortage)
 USD     |  24.0 %  |   0.0 % | -3,000,000 JPY (Surplus)
 GLDM    |  14.0 %  |  20.0 % |   +750,000 JPY (Shortage)
 2800    |  20.0 %  |  20.0 % |         0 JPY (Optimal)
 BTC     |  30.0 %  |  30.0 % |         0 JPY (Optimal)
--------------------------------------------------------------------------------

■ NEXT ACTIONS (Who should do what, and where)
--------------------------------------------------------------------------------
【Sell / Take Profit】
● Member_A (SBI): Convert USD 3,000 to JPY
● Member_B (Rakuten): Convert USD 2,000 to JPY
  ⇒ Secure approx. 3,000,000 JPY in total.

【Buy / Accumulate】
● Member_A (SBI): Buy 500,000 JPY worth of GLDM
● Member_B (Rakuten): Buy 250,000 JPY worth of GLDM
● Member_A (Bitget): Replenish 2,250,000 JPY or allocate to BTC buying power
--------------------------------------------------------------------------------
```

---

## 5. Usage

### 1. Prerequisites
Ensure you have the following installed and set up:
* Go (1.21 or later)
* Create a Google Spreadsheet and copy the CSV publish URL.
```text
Columns:
Symbol | Rate | Category
Examples:
USD | =GOOGLEFINANCE("CURRENCY:USDJPY") | currency
2800 | =GOOGLEFINANCE("HKG:2800") * GOOGLEFINANCE("CURRENCY:HKDJPY") | eastern
GLDM | =GOOGLEFINANCE("GLDM") * GOOGLEFINANCE("CURRENCY:USDJPY") | western
```
* Create a `.env` file at the root directory by referencing `.env.example`. Paste your Google Spreadsheet CSV URL.
```text
.env
# Google Spreadsheet CSV URL
GEO_HEDGER_CSV_URL="[https://docs.google.com/spreadsheets/d/e/XXXXX/pub?output=csv](https://docs.google.com/spreadsheets/d/e/XXXXX/pub?output=csv)"
```

### 2. Installation
```bash
git clone [https://github.com/nao-nba/geo-hedger.git](https://github.com/nao-nba/geo-hedger.git)
cd geo-hedger
go build
```

### 3. Global Configuration (Set base currency and target phase)
```bash
geo-hedger set-config --base-currency JPY --target-phase phase_2
```

### 4. Define or Update a Phase (Scenario and target ratios)
```bash
geo-hedger set-phase --id phase_2 --name "Phase 2: Crisis" --scenario "Geopolitical crisis materialize" --action "Sell high USD" --targets "JPY:30,USD:0,GLDM:20,BTC:30"
```

### 5. Add or Update Family Assets
```bash
# Add or accumulate
geo-hedger add-asset --owner "Member_A (SBI)" --symbol BTC --amount 0.2

# Sell / Reduce amount
geo-hedger add-asset --owner "Member_A (SBI)" --symbol BTC --amount -0.2
```

### 6. Delete Family Assets
```bash
# Delete a specific asset of a specific owner
geo-hedger delete-asset --owner "Member_A (SBI)" --symbol BTC

# Delete all assets of a specific owner
geo-hedger delete-asset --owner "Member_A (SBI)"
```

### 7. View Status and Next Actions
```bash
geo-hedger status
```
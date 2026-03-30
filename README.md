# Geo-Hedger

A zero-dependency, local-first CLI tool built in Go for asset allocation and survival inventory management under extreme geopolitical risks.

---

## Design Philosophy & Why It Exists

In an era of heightening geopolitical tensions—such as the potential choke on energy supply chains like the Strait of Hormuz—traditional financial strategies and domestic currencies are no longer guaranteed safe havens. 

I designed **Geo-Hedger** not just as a simple asset tracker, but as a **"Digital Bunker Mission Control"** to survive hyperinflation and systemic resets.

This tool solves three critical problems for a modern survivor:
1. **Dynamic Risk Allocation:** Separates assets into Western risks, Eastern risks, and raw currencies to visualize real-time exposure to shifting global power balances.
2. **Rationale Logging:** Forces the user to explicitly articulate the geopolitical phase and the strategy behind the asset ratio, ensuring decisions are driven by logic, not panic.
3. **Cold-Storage Inventory:** Tracks physical stockpiles (food, medical supplies, fuel) with expiration dates, keeping a family safe even when external supply chains collapse.

## Tech Stack & Why Go?

* **Language:** Go (Golang)
* **Database:** Pure JSON (No external DB required)

### Why this stack?
* **Zero-Dependency & Portability:** In a worst-case scenario where internet or heavy infrastructure fails, this tool can be compiled into a single binary and run on any machine without internet.
* **Human-Readable Fallback:** By storing data in flat JSON files, the data remains accessible and editable with a simple text editor (like Vim or VS Code) even if the application itself fails.

---

## Core Features (MVP)
* Total household net-worth calculation across multi-currencies (JPY, USD, HKD).
* Visualization of actual vs. ideal asset allocation.
* Stockpile management with family-friendly, rough-tracking UX.
* Scenario-based action planning (e.g., "What to do if hyperinflation spikes").

## Installation & Usage
*(Coming soon: Here you will paste the terminal commands to run the tool)*
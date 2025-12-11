# Simple Research Data Manager (SRDM)

SRDM is a lightweight, efficient command-line tool designed for managing research data assets.
It helps researchers catalogue, track, and retrieve datasets and experimental records using a simple SQLite-backed repository.

## 📚 Core Concepts

Before you start, it's helpful to understand how SRDM organizes data:

- **Repository**: The central SQLite database where all metadata is stored.
  By default, this is located at `~/.local/share/SRDM/srdm_dataRepo.sqlite` (or configured via `SRDM_DATA_REPO_PATH`).
- **Table**: A high-level collection of data (e.g., "Experiment A Results", "Patient Survey 2024").
  Think of this as a folder or a dataset wrapper.
- **Record**: An individual data item belonging to a Table (e.g., "Run #1", "Patient 001").
  Records hold specific metadata like file paths, statistical summaries, and tags.

---

## 🚀 Getting Started

### Installation

Clone the repository and build the binary:

```bash
make build
```

This will create the executable at `bin/srdm`.

### First Run

Initialize the database by running any command, for example checking the info:

```bash
./bin/srdm info
```

You should see an empty repository status.

---

## 📖 User Guide

### 1. Ingesting Data (`insert`)

The core workflow begins with inserting metadata about your files.

**Creating a Table**

First, creating a "container" for your records.

```bash
./bin/srdm insert --name "biostudy:seq_data" \
  --keys "sample_id" \
  --description "DNA Sequencing Results 2024" \
  --engine "SQLite3"
```

**Adding Records**

Now, add records to that table.

```bash
./bin/srdm insert --name "biostudy:seq_data:sample_01" \
  --type "fastq" \
  --label "control_group" \
  --source "/raw/data/seq/s01.fq" \
  --description "Control Sample 1"
```

### 2. Searching & Querying (`search`)

Find what you need quickly. SRDM supports both exact matching and fuzzy prefix searching.

**List all records in a table:**

```bash
./bin/srdm search "biostudy:seq_data:%"
```

**Find a specific record:**

```bash
./bin/srdm search "biostudy:seq_data:sample_01"
```

### 3. Viewing Details (`view`)

Inspect detailed metadata for any item.

```bash
./bin/srdm view "biostudy:seq_data:sample_01"
```

### 4. Updating Metadata (`update`)

Did you change a file location or need to add a tag?

```bash
./bin/srdm update --name "biostudy:seq_data:sample_01" \
  --label "control_group_verified" \
  --log_file "processing.log"
```

### 5. Exporting Metadata (`export`)

Generate JSON reports of your data for external analysis or sharing.

```bash
./bin/srdm export "biostudy:seq_data:%" -o report.json
```

### 6. Cleaning Up (`delete`)

Remove old or erroneous entries.

**Delete a single record:**

```bash
./bin/srdm delete "biostudy:seq_data:sample_01"
```

**Delete an entire table (and all its records):**
*Note: This requires the `--force` flag.*

```bash
./bin/srdm delete "biostudy:seq_data" --force
```

---

## ⚙️ Configuration

You can configure the repository location using the `SRDM_DATA_REPO_PATH` environment variable.

```bash
export SRDM_DATA_REPO_PATH="/path/to/my/shared/repo.sqlite"
./bin/srdm info
```

Alternatively, use the global `--path` flag with any command:

```bash
./bin/srdm info --path "/custom/db.sqlite"
```


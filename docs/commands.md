# Launchpad Commands

## Overview
Launchpad provides CLI commands to manage clusters via YAML configurations.

---

## `discover`
**Purpose**: Inspect the current state of a cluster/project.
**Usage**:
```bash
launchpad discover --config <file.yaml>
```
**Output**: JSON/YAML summary of cluster state (nodes, components, dependencies).
**Flags**:
- `--output <format>`: Specify output format (`json`, `yaml`, `table`).

---

## `apply`
**Purpose**: Apply a YAML configuration to a cluster.
**Usage**:
```bash
launchpad apply --config <file.yaml>
```
**Flags**:
- `--dry-run`: Validate without applying changes.
- `--force`: Skip confirmation prompts.

---

## `reset`
**Purpose**: Remove installed components and revert cluster state.
**Usage**:
```bash
launchpad reset --config <file.yaml>
```
**Warning**: Irreversible. Use `--dry-run` first.
**Flags**:
- `--exclude <component>`: Skip resetting specific components (e.g., `--exclude mke`).
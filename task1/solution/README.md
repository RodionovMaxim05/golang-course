# GitHub Repo Info CLI tool

## Description

A simple CLI tool for getting information about a GitHub repository.

## Requirements

- Go 1.25+
- Internet access

## Usage

```bash
go run ./repo-info <owner>/<repo>
```

**Examples:**

```bash
go run ./repo-info golang/go
go run ./repo-info torvalds/linux
```

## Example Output

```bash
=== Repository Information ===
Name:        linux
Description: Linux kernel source tree
Stars:       221616
Forks:       60872
Created:     2011-09-04 22:48:12
```

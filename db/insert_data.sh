#!/usr/bin/env bash
echo "Populating aldo db tables..."
sqlite3 aldo.db < populate_tables.sql

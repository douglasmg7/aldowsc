#!/usr/bin/env bash
echo "Creating aldo db tables..."
sqlite3 aldo.db < create_tables.sql

echo "Populating aldo db tables..."
sqlite3 aldo.db < populate_tables.sql

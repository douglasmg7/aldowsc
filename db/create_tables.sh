#!/usr/bin/env bash
echo "Creating aldo db tables..."
sqlite3 aldo.db < tables.sql

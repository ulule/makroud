#!/bin/bash
set -e

echo "update postgres conf"
sed -ri "s/#log_statement = 'none'/log_statement = 'all'/g" /var/lib/postgresql/data/postgresql.conf
cat /var/lib/postgresql/data/postgresql.conf

# DefiFundr Database Backup and Migration Strategy

## Overview
This document outlines the database backup, migration, and restoration procedures for the DefiFundr project.

## Backup Configuration

### Configuration File
Create a configuration file at `/etc/defifundr/backup_config.json`:

```json
{
    "database": {
        "name": "defifundr_db",
        "user": "postgres",
        "host": "localhost"
    },
    "backup": {
        "base_dir": "/var/backups/defifundr",
        "retention_days": 30,
        "compression": true
    },
    "logging": {
        "log_file": "/var/log/defifundr_backup.log"
    }
}
```

### Installation Steps
1. Install dependencies
```bash
sudo apt-get update
sudo apt-get install -y postgresql postgresql-contrib python3 python3-pip
sudo pip3 install psycopg2-binary
```

2. Set up backup script
```bash
sudo mkdir -p /etc/defifundr
sudo mkdir -p /var/backups/defifundr
sudo mkdir -p /var/log
sudo chown -R postgres:postgres /var/backups/defifundr
sudo chown -R postgres:postgres /var/log
```

3. Place backup script
```bash
sudo mv backup_script.py /usr/local/bin/defifundr_backup.py
sudo chmod +x /usr/local/bin/defifundr_backup.py
```

### Crontab Configuration
Edit the crontab for the postgres user:
```bash
sudo -u postgres crontab -e
```

Add the following lines:
```
# Daily backup at 2 AM
0 2 * * * /usr/bin/python3 /usr/local/bin/defifundr_backup.py daily

# Weekly backup at midnight on Sundays
0 0 * * 0 /usr/bin/python3 /usr/local/bin/defifundr_backup.py weekly
```

## Backup Procedures

### Backup Types
1. **Daily Backup**: 
   - Runs every day at 2 AM
   - Stores compressed SQL dump
   - Generates SHA-256 checksum

2. **Weekly Backup**:
   - Runs every Sunday at midnight
   - More comprehensive backup
   - Useful for long-term archiving

### Manual Backup
```bash
# Perform manual backup
sudo -u postgres python3 /usr/local/bin/defifundr_backup.py daily

# Perform test backup and restore
sudo -u postgres python3 /usr/local/bin/defifundr_backup.py test
```

## Restoration Procedures

### Restore from Backup
```bash
# Restore specific backup
sudo -u postgres python3 /usr/local/bin/defifundr_backup.py restore /path/to/backup/file.sql.gz
```

### Restoration Considerations
- Ensure you have the correct backup file
- Verify checksum before
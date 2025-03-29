#!/usr/bin/env python3

import os
import sys
import logging
import subprocess
from datetime import datetime, timedelta
import json
import hashlib

class DatabaseBackupManager:
    def __init__(self, config_path='/etc/defifundr/backup_config.json'):
        """
        Initialize backup manager with configuration
        
        Config file should contain:
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
        """
        # Load configuration
        try:
            with open(config_path, 'r') as config_file:
                self.config = json.load(config_file)
        except Exception as e:
            print(f"Error loading configuration: {e}")
            sys.exit(1)
        
        # Setup logging
        self._setup_logging()
        
        # Ensure backup directory exists
        os.makedirs(self.config['backup']['base_dir'], exist_ok=True)

    def _setup_logging(self):
        """Configure logging for backup process"""
        log_file = self.config['logging']['log_file']
        
        # Ensure log directory exists
        os.makedirs(os.path.dirname(log_file), exist_ok=True)
        
        # Configure logging
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler(log_file),
                logging.StreamHandler(sys.stdout)
            ]
        )
        self.logger = logging.getLogger('DatabaseBackupManager')

    def perform_backup(self, backup_type='daily'):
        """
        Perform database backup with comprehensive error handling
        
        :param backup_type: Type of backup (daily, weekly, migration)
        :return: Backup file path or None if backup fails
        """
        try:
            # Generate unique backup filename
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            backup_filename = f"{backup_type}_backup_{timestamp}.sql"
            
            # Full path for backup
            backup_path = os.path.join(
                self.config['backup']['base_dir'], 
                backup_filename
            )
            
            # Compression flag
            use_compression = self.config['backup'].get('compression', True)
            
            # Prepare backup command
            backup_command = [
                'pg_dump',
                '-h', self.config['database'].get('host', 'localhost'),
                '-U', self.config['database']['user'],
                '-d', self.config['database']['name'],
                '-f', backup_path
            ]
            
            # Execute backup
            result = subprocess.run(
                backup_command, 
                capture_output=True, 
                text=True,
                env=dict(os.environ, PGPASSWORD=self._get_db_password())
            )
            
            # Check backup success
            if result.returncode != 0:
                self.logger.error(f"Backup failed: {result.stderr}")
                return None
            
            # Optional compression
            if use_compression:
                compressed_path = f"{backup_path}.gz"
                compress_result = subprocess.run(
                    ['gzip', backup_path],
                    capture_output=True,
                    text=True
                )
                if compress_result.returncode != 0:
                    self.logger.warning(f"Compression failed: {compress_result.stderr}")
                    compressed_path = backup_path
            else:
                compressed_path = backup_path
            
            # Generate checksum
            checksum = self._generate_checksum(compressed_path)
            
            # Log successful backup
            self.logger.info(f"Successful {backup_type} backup: {compressed_path}")
            
            # Rotate backups
            self._rotate_backups()
            
            return compressed_path
        
        except Exception as e:
            self.logger.error(f"Backup error: {e}")
            return None

    def _get_db_password(self):
        """
        Retrieve database password securely
        In production, use a secure method like:
        - Environment variable
        - Secret management system (Vault, AWS Secrets Manager)
        """
        # IMPORTANT: Replace with secure password retrieval
        return os.environ.get('DB_PASSWORD', '')

    def _generate_checksum(self, filepath):
        """
        Generate SHA-256 checksum for backup file
        
        :param filepath: Path to backup file
        :return: Checksum string
        """
        try:
            with open(filepath, 'rb') as f:
                checksum = hashlib.sha256()
                for chunk in iter(lambda: f.read(4096), b""):
                    checksum.update(chunk)
            
            # Store checksum in a separate file
            checksum_file = f"{filepath}.sha256"
            with open(checksum_file, 'w') as f:
                f.write(checksum.hexdigest())
            
            self.logger.info(f"Checksum generated: {checksum.hexdigest()}")
            return checksum.hexdigest()
        except Exception as e:
            self.logger.error(f"Checksum generation failed: {e}")
            return None

    def _rotate_backups(self):
        """
        Remove old backups based on retention policy
        """
        try:
            retention_days = self.config['backup'].get('retention_days', 30)
            cutoff_date = datetime.now() - timedelta(days=retention_days)
            
            backup_dir = self.config['backup']['base_dir']
            
            for filename in os.listdir(backup_dir):
                filepath = os.path.join(backup_dir, filename)
                
                # Skip if not a file
                if not os.path.isfile(filepath):
                    continue
                
                # Get file modification time
                file_modified = datetime.fromtimestamp(os.path.getmtime(filepath))
                
                # Remove backups older than retention period
                if file_modified < cutoff_date:
                    try:
                        os.remove(filepath)
                        self.logger.info(f"Removed old backup: {filename}")
                    except Exception as e:
                        self.logger.error(f"Error removing {filename}: {e}")
        except Exception as e:
            self.logger.error(f"Backup rotation error: {e}")

    def restore_backup(self, backup_file):
        """
        Restore database from a backup file
        
        :param backup_file: Path to backup file to restore
        :return: Boolean indicating success
        """
        try:
            # Check if backup file exists
            if not os.path.exists(backup_file):
                self.logger.error(f"Backup file not found: {backup_file}")
                return False
            
            # Determine if file is compressed
            if backup_file.endswith('.gz'):
                # Decompress first
                decompress_command = ['gunzip', '-k', backup_file]
                decompress_result = subprocess.run(
                    decompress_command, 
                    capture_output=True, 
                    text=True
                )
                if decompress_result.returncode != 0:
                    self.logger.error(f"Decompression failed: {decompress_result.stderr}")
                    return False
                
                # Remove .gz extension for restoration
                backup_file = backup_file[:-3]
            
            # Restore command
            restore_command = [
                'psql',
                '-h', self.config['database'].get('host', 'localhost'),
                '-U', self.config['database']['user'],
                '-d', self.config['database']['name'],
                '-f', backup_file
            ]
            
            # Execute restore
            restore_result = subprocess.run(
                restore_command, 
                capture_output=True, 
                text=True,
                env=dict(os.environ, PGPASSWORD=self._get_db_password())
            )
            
            # Check restore success
            if restore_result.returncode != 0:
                self.logger.error(f"Restore failed: {restore_result.stderr}")
                return False
            
            self.logger.info(f"Successfully restored backup: {backup_file}")
            return True
        
        except Exception as e:
            self.logger.error(f"Restore error: {e}")
            return False

    def test_backup_and_restore(self):
        """
        Comprehensive backup and restore test
        
        :return: Boolean indicating test success
        """
        try:
            # Perform backup
            backup_file = self.perform_backup(backup_type='test')
            
            if not backup_file:
                self.logger.error("Backup test failed: Could not create backup")
                return False
            
            # Create a test database for restoration
            test_db_name = f"{self.config['database']['name']}_restore_test"
            
            # Create test database
            create_db_command = [
                'createdb',
                '-h', self.config['database'].get('host', 'localhost'),
                '-U', self.config['database']['user'],
                test_db_name
            ]
            
            create_db_result = subprocess.run(
                create_db_command,
                capture_output=True,
                text=True,
                env=dict(os.environ, PGPASSWORD=self._get_db_password())
            )
            
            if create_db_result.returncode != 0:
                self.logger.error(f"Test database creation failed: {create_db_result.stderr}")
                return False
            
            # Modify config for test restore
            original_db = self.config['database']['name']
            self.config['database']['name'] = test_db_name
            
            # Perform restore
            restore_success = self.restore_backup(backup_file)
            
            # Cleanup: Restore original database name and drop test database
            self.config['database']['name'] = original_db
            
            drop_db_command = [
                'dropdb',
                '-h', self.config['database'].get('host', 'localhost'),
                '-U', self.config['database']['user'],
                test_db_name
            ]
            
            subprocess.run(
                drop_db_command,
                capture_output=True,
                text=True,
                env=dict(os.environ, PGPASSWORD=self._get_db_password())
            )
            
            return restore_success
        
        except Exception as e:
            self.logger.error(f"Backup and restore test failed: {e}")
            return False

def main():
    """
    Main entry point for backup script
    """
    # Create backup manager
    backup_manager = DatabaseBackupManager()
    
    # Parse command-line arguments
    if len(sys.argv) > 1:
        command = sys.argv[1]
        
        if command == 'daily':
            backup_manager.perform_backup(backup_type='daily')
        elif command == 'weekly':
            backup_manager.perform_backup(backup_type='weekly')
        elif command == 'test':
            test_result = backup_manager.test_backup_and_restore()
            sys.exit(0 if test_result else 1)
        elif command == 'restore':
            if len(sys.argv) < 3:
                print("Please provide backup file path to restore")
                sys.exit(1)
            backup_file = sys.argv[2]
            backup_manager.restore_backup(backup_file)
        else:
            print(f"Unknown command: {command}")
            sys.exit(1)
    else:
        # Default to daily backup if no argument provided
        backup_manager.perform_backup(backup_type='daily')

if __name__ == "__main__":
    main()
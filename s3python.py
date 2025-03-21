import boto3
import os
from datetime import datetime, timezone

# === CONFIGURATION ===
bucket_name = 'your-bucket-name'
download_dir = './downloaded_files'
target_date = datetime(2025, 2, 26, tzinfo=timezone.utc)  # Target date

# Ensure local download directory exists
os.makedirs(download_dir, exist_ok=True)

# Create S3 client
s3 = boto3.client('s3')

# Use paginator to recursively go through the whole bucket
paginator = s3.get_paginator('list_objects_v2')
pages = paginator.paginate(Bucket=bucket_name)

# Iterate through files
for page in pages:
    if 'Contents' in page:
        for obj in page['Contents']:
            key = obj['Key']
            last_modified = obj['LastModified']

            # Filter: PDF, contains "SYFSDB", and uploaded on 2/26/2025
            if 'SYFSDB' in key and key.lower().endswith('.pdf') and last_modified.date() == target_date.date():
                # Reconstruct folder structure locally
                local_path = os.path.join(download_dir, key)
                local_folder = os.path.dirname(local_path)
                os.makedirs(local_folder, exist_ok=True)

                print(f"Downloading: {key} (uploaded {last_modified})")
                s3.download_file(bucket_name, key, local_path)

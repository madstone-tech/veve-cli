# Integration Guide

This guide shows how to integrate veve into your workflows, scripts, and CI/CD pipelines.

## Unix Composability

veve supports Unix pipes and standard input/output, making it easy to integrate into scripts and shell pipelines.

### stdin Input

Read markdown from standard input using `-` as the input file:

```bash
# Simple pipe
cat document.md | veve - -o output.pdf

# With curl (download and convert)
curl https://example.com/document.md | veve - -o output.pdf

# With heredoc
veve - -o output.pdf << 'EOF'
# My Document
Content here
EOF

# With theme
echo "# Test" | veve - -o output.pdf --theme dark
```

### stdout Output

Write PDF to standard output using `-` as the output file:

```bash
# Print PDF binary to stdout
veve input.md -o - | cat > output.pdf

# Pipe to another command
veve input.md -o - | base64 | curl -d @- https://api.example.com/store

# Stream to cloud storage
veve document.md -o - | aws s3 cp - s3://bucket/document.pdf

# Write to file and stdout simultaneously
veve input.md -o - | tee output.pdf | wc -c
```

### Exit Codes

veve returns proper exit codes for shell scripting:

- **0**: Successful conversion
- **1**: Error (file not found, invalid theme, conversion failed)
- **2**: Usage error (invalid arguments)

```bash
veve input.md -o output.pdf
if [ $? -eq 0 ]; then
  echo "Success!"
else
  echo "Conversion failed"
fi
```

### Error Handling

All errors are printed to stderr, allowing you to separate errors from output:

```bash
# Redirect errors to a log file
veve input.md -o output.pdf 2> errors.log

# Suppress errors but capture them
if ! veve input.md -o output.pdf 2>/dev/null; then
  echo "Conversion failed"
fi

# Both stdout and stderr capture
veve input.md -o output.pdf > output.txt 2>&1
```

## Batch Processing

Process multiple markdown files in a loop:

### Simple Loop

```bash
#!/bin/bash
for file in *.md; do
  veve "$file" --quiet
  if [ $? -ne 0 ]; then
    echo "Failed: $file" >&2
  fi
done
```

### With Output Directory

```bash
#!/bin/bash
input_dir="./markdown"
output_dir="./pdf"

mkdir -p "$output_dir"

for file in "$input_dir"/*.md; do
  basename=$(basename "$file" .md)
  veve "$file" -o "$output_dir/$basename.pdf" --quiet
done
```

### Using find and xargs

```bash
# Sequential processing
find . -name "*.md" -print0 | xargs -0 -I {} veve {} --quiet

# Parallel processing (2 jobs at a time)
find . -name "*.md" -print0 | xargs -0 -P 2 -I {} veve {} --quiet

# With error handling
find . -name "*.md" -print0 | while IFS= read -r -d '' file; do
  if ! veve "$file" --quiet; then
    echo "Failed: $file" >&2
    exit 1
  fi
done
```

### Progress Indicator

```bash
#!/bin/bash
total=$(find . -name "*.md" | wc -l)
count=0

find . -name "*.md" | while read file; do
  ((count++))
  echo "[$count/$total] Converting $file..."
  veve "$file" --quiet
done
```

## Continuous Integration / Continuous Deployment

### GitHub Actions

```yaml
name: Build PDFs

on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Install veve
        run: |
          # Install from releases (when available)
          curl -L -o veve https://github.com/yourusername/veve-cli/releases/download/v0.1.0/veve-linux-amd64
          chmod +x veve
          sudo mv veve /usr/local/bin/
      
      - name: Convert markdown to PDF
        run: |
          mkdir -p pdf
          for file in markdown/*.md; do
            veve "$file" -o "pdf/$(basename "$file" .md).pdf"
          done
      
      - name: Upload PDFs
        uses: actions/upload-artifact@v2
        with:
          name: pdfs
          path: pdf/
```

### GitLab CI

```yaml
build_pdfs:
  image: ubuntu:latest
  script:
    - apt-get update && apt-get install -y pandoc
    - |
      for file in markdown/*.md; do
        veve "$file" -o "pdf/$(basename "$file" .md).pdf"
      done
  artifacts:
    paths:
      - pdf/
    expire_in: 1 week
```

### Jenkins

```groovy
pipeline {
    agent any
    
    stages {
        stage('Install veve') {
            steps {
                sh 'curl -L -o veve https://github.com/yourusername/veve-cli/releases/download/v0.1.0/veve-linux-amd64'
                sh 'chmod +x veve'
                sh 'sudo mv veve /usr/local/bin/'
            }
        }
        
        stage('Build PDFs') {
            steps {
                sh '''
                    mkdir -p pdf
                    for file in markdown/*.md; do
                        veve "$file" -o "pdf/$(basename "$file" .md).pdf"
                    done
                '''
            }
        }
        
        stage('Publish') {
            steps {
                archiveArtifacts artifacts: 'pdf/**'
            }
        }
    }
}
```

## Static Site Generators

### Hugo Integration

Process Hugo markdown files to PDF:

```bash
#!/bin/bash
# Convert all Hugo content pages to PDF

OUTPUT_DIR="public/pdf"
mkdir -p "$OUTPUT_DIR"

for file in content/**/*.md; do
  # Skip _index.md files
  if [[ "$file" == *"_index.md" ]]; then
    continue
  fi
  
  # Calculate output path
  relpath="${file#content/}"
  relpath="${relpath%/*.md}"
  
  # Create output directory
  mkdir -p "$OUTPUT_DIR/$relpath"
  
  # Convert to PDF
  veve "$file" -o "$OUTPUT_DIR/$relpath/$(basename "$file" .md).pdf"
done
```

### Jekyll Integration

```bash
#!/bin/bash
# Convert Jekyll markdown files to PDF

OUTPUT_DIR="pdf"
mkdir -p "$OUTPUT_DIR"

for file in _posts/*.md; do
  basename=$(basename "$file" .md)
  veve "$file" -o "$OUTPUT_DIR/$basename.pdf"
done
```

### Sphinx Integration

```python
# conf.py - Sphinx configuration

import subprocess
import os

def build_pdfs(app, exception):
    """Build PDFs after Sphinx build"""
    if exception:
        return
    
    source_dir = app.srcdir
    output_dir = os.path.join(app.outdir, "pdf")
    
    os.makedirs(output_dir, exist_ok=True)
    
    for root, dirs, files in os.walk(source_dir):
        for file in files:
            if file.endswith(".rst"):
                filepath = os.path.join(root, file)
                output_file = os.path.join(output_dir, file.replace(".rst", ".pdf"))
                
                # Convert using veve
                subprocess.run([
                    "veve", filepath, 
                    "-o", output_file,
                    "--quiet"
                ])

def setup(app):
    app.connect("build-finished", build_pdfs)
```

## Document Generation Pipelines

### Report Generation

```bash
#!/bin/bash
# Generate timestamped PDF reports

DATE=$(date +%Y-%m-%d)
REPORT="report-$DATE.md"
OUTPUT="reports/report-$DATE.pdf"

mkdir -p reports

# Generate markdown report
cat > "$REPORT" << EOF
# Daily Report - $DATE

Generated at $(date)

## Summary

Report content here.
EOF

# Convert to PDF
veve "$REPORT" -o "$OUTPUT" --theme academic

# Send via email
mail -s "Daily Report" admin@example.com < "$OUTPUT"
```

### Dynamic Document Generation

```bash
#!/bin/bash
# Generate documentation from code

OUTPUT="docs/api.pdf"

# Extract comments and generate markdown
cat > /tmp/api.md << 'EOF'
# API Documentation

Generated: $(date)

EOF

# Extract function documentation
grep -r "^///" . >> /tmp/api.md

# Convert to PDF
veve /tmp/api.md -o "$OUTPUT" --theme default
```

## Data Processing Pipelines

### CSV to PDF Report

```bash
#!/bin/bash
# Convert CSV data to PDF report

CSV_FILE="data.csv"
MD_FILE="/tmp/report.md"
PDF_FILE="report.pdf"

# Generate markdown from CSV
cat > "$MD_FILE" << 'EOF'
# Data Report

EOF

# Convert CSV to markdown table
awk -F',' '
  NR==1 {
    print "| " $0 " |"
    for(i=1; i<=NF; i++) print "| --- |"
  }
  NR>1 { print "| " $0 " |" }
' "$CSV_FILE" >> "$MD_FILE"

# Generate PDF
veve "$MD_FILE" -o "$PDF_FILE"
```

## Web Application Integration

### Express.js

```javascript
const express = require('express');
const { spawn } = require('child_process');

const app = express();

app.post('/convert', express.text({ type: 'text/markdown' }), (req, res) => {
  const veve = spawn('veve', ['-', '-o', '-', '--quiet']);
  
  veve.stdin.write(req.body);
  veve.stdin.end();
  
  res.type('application/pdf');
  veve.stdout.pipe(res);
  
  veve.stderr.on('data', (data) => {
    console.error(`stderr: ${data}`);
    res.status(500).send('Conversion failed');
  });
});

app.listen(3000);
```

### Python with Flask

```python
from flask import Flask, request
from subprocess import Popen, PIPE

app = Flask(__name__)

@app.route('/convert', methods=['POST'])
def convert():
    markdown = request.get_data(as_text=True)
    
    process = Popen(
        ['veve', '-', '-o', '-', '--quiet'],
        stdin=PIPE,
        stdout=PIPE,
        stderr=PIPE
    )
    
    pdf, error = process.communicate(markdown.encode())
    
    if process.returncode != 0:
        return error.decode(), 500
    
    return pdf, 200, {'Content-Type': 'application/pdf'}

if __name__ == '__main__':
    app.run()
```

## Tips and Tricks

### Error Logging

```bash
#!/bin/bash
# Convert with detailed error logging

LOG_FILE="conversion.log"

veve input.md -o output.pdf \
  --verbose \
  2> >(tee -a "$LOG_FILE" >&2)
```

### Dry Run

```bash
#!/bin/bash
# Check if files will convert without creating PDFs

for file in *.md; do
  echo "Would convert: $file"
  veve "$file" -o /dev/null --quiet && echo "  ✓ Success" || echo "  ✗ Failed"
done
```

### Filtering by Size

```bash
#!/bin/bash
# Only convert small markdown files

find . -name "*.md" -size -100k -print0 | \
  xargs -0 -I {} veve {} --quiet
```

### Watch and Convert

```bash
#!/bin/bash
# Auto-convert on file changes (requires fswatch)

fswatch markdown/ | while read file; do
  echo "Converting $file..."
  veve "$file" --quiet
done
```

## Troubleshooting

### Command Not Found

If veve is not in your PATH:

```bash
# Full path
/usr/local/bin/veve input.md -o output.pdf

# Or add to PATH
export PATH="$PATH:/usr/local/bin"
veve input.md -o output.pdf
```

### Permission Denied

```bash
# Make sure veve is executable
chmod +x /usr/local/bin/veve

# Check if you have write permission for output
touch /tmp/test.txt && rm /tmp/test.txt
```

### Pandoc Not Found

```bash
# Install pandoc
# macOS
brew install pandoc

# Linux (Ubuntu/Debian)
sudo apt-get install pandoc

# Verify
pandoc --version
```

### Temp File Issues

```bash
# Clean up temp files if conversion is interrupted
rm -f /tmp/veve-*.pdf

# Check temp directory permissions
ls -la /tmp/ | grep veve
```

## Performance

### Parallel Conversion

For large batches, use xargs with parallel execution:

```bash
# 4 parallel jobs
find . -name "*.md" -print0 | \
  xargs -0 -P 4 -I {} veve {} --quiet
```

### Quiet Mode for Scripts

Always use `--quiet` in non-interactive scripts to reduce overhead:

```bash
veve input.md --quiet
```

## Examples Repository

Complete examples are available in the [examples/](../examples/) directory:

- `batch-convert.sh` - Batch processing example
- `ci-integration.yml` - CI/CD configuration
- `web-api.js` - Web integration example
- `dynamic-docs.sh` - Dynamic documentation

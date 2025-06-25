#!/bin/bash

# Create SSL directory if it doesn't exist
mkdir -p nginx/ssl

# Generate private key
openssl genrsa -out nginx/ssl/key.pem 2048

# Generate certificate signing request
openssl req -new -key nginx/ssl/key.pem -out nginx/ssl/cert.csr -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

# Generate self-signed certificate
openssl x509 -req -in nginx/ssl/cert.csr -signkey nginx/ssl/key.pem -out nginx/ssl/cert.pem -days 365

# Clean up
rm nginx/ssl/cert.csr

echo "SSL certificates generated successfully!"
echo "Certificate: nginx/ssl/cert.pem"
echo "Private Key: nginx/ssl/key.pem"
echo ""
echo "Note: These are self-signed certificates for development only."
echo "For production, use certificates from a trusted CA."
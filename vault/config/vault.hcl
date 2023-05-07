backend "file" {
  path = "/vault/file"
}

listener "tcp" {
  address       = "0.0.0.0:8200"
  tls_cert_file = "/vault/config/cert.pem"
  tls_key_file  = "/vault/config/cert.key"
}

#disable_mlock = true

seal "azurekeyvault" {
}
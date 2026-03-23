package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type conf struct {
	HttpRequestLog                 bool
	HttpBindingAddress             string
	HttpsBindingAddress            string
	HttpsRedirectEnabled           bool
	TLSAutoDomain                  string
	TLSCertFilepath                string
	TLSCertKeyFilepath             string
	VaultPrefix                    string
	PassphraseNumWordsDefault      int
	PassphraseSeparatorDefault     string
	PassphraseCapitalizeDefault    bool
	PassphraseIncludeNumberDefault bool
	PassphraseMaxNumberDefault     int
}

const HttpRequestLogVarenv = "OTS_HTTP_REQUEST_LOG"
const HttpBindingAddressVarenv = "OTS_HTTP_BINDING_ADDRESS"
const HttpsBindingAddressVarenv = "OTS_HTTPS_BINDING_ADDRESS"
const HttpsRedirectEnabledVarenv = "OTS_HTTPS_REDIRECT_ENABLED"
const TLSAutoDomainVarenv = "OTS_TLS_AUTO_DOMAIN"
const TLSCertFilepathVarenv = "OTS_TLS_CERT_FILEPATH"
const TLSCertKeyFilepathVarenv = "OTS_TLS_CERT_KEY_FILEPATH"
const VaultPrefixenv = "OTS_VAULT_PREFIX"
const PassphraseNumWordsDefaultVarenv = "OTS_PASSPHRASE_NUM_WORDS_DEFAULT"
const PassphraseSeparatorDefaultVarenv = "OTS_PASSPHRASE_SEPARATOR_DEFAULT"
const PassphraseCapitalizeDefaultVarenv = "OTS_PASSPHRASE_CAPITALIZE_DEFAULT"
const PassphraseIncludeNumberDefaultVarenv = "OTS_PASSPHRASE_INCLUDE_NUMBER_DEFAULT"
const PassphraseMaxNumberDefaultVarenv = "OTS_PASSPHRASE_MAX_NUMBER_DEFAULT"

func loadConfig() conf {
	var cnf conf

	cnf.HttpRequestLog = getEnvBool(HttpRequestLogVarenv, false)
	cnf.HttpBindingAddress = getEnv(HttpBindingAddressVarenv, "")
	cnf.HttpsBindingAddress = getEnv(HttpsBindingAddressVarenv, "")
	cnf.HttpsRedirectEnabled = getEnvBool(HttpsRedirectEnabledVarenv, false)
	cnf.TLSAutoDomain = getEnv(TLSAutoDomainVarenv, "")
	cnf.TLSCertFilepath = getEnv(TLSCertFilepathVarenv, "")
	cnf.TLSCertKeyFilepath = getEnv(TLSCertKeyFilepathVarenv, "")
	cnf.VaultPrefix = getEnv(VaultPrefixenv, "")
	cnf.PassphraseNumWordsDefault = getEnvInt(PassphraseNumWordsDefaultVarenv, 5)
	cnf.PassphraseSeparatorDefault = getEnv(PassphraseSeparatorDefaultVarenv, "-")
	cnf.PassphraseCapitalizeDefault = getEnvBool(PassphraseCapitalizeDefaultVarenv, true)
	cnf.PassphraseIncludeNumberDefault = getEnvBool(PassphraseIncludeNumberDefaultVarenv, true)
	cnf.PassphraseMaxNumberDefault = getEnvInt(PassphraseMaxNumberDefaultVarenv, 9)

	if cnf.TLSAutoDomain != "" && (cnf.TLSCertFilepath != "" || cnf.TLSCertKeyFilepath != "") {
		log.Fatalf("Auto TLS (%s) is mutually exclusive with manual TLS (%s and %s)", TLSAutoDomainVarenv,
			TLSCertFilepathVarenv, TLSCertKeyFilepathVarenv)
	}

	if (cnf.TLSCertFilepath != "" && cnf.TLSCertKeyFilepath == "") ||
		(cnf.TLSCertFilepath == "" && cnf.TLSCertKeyFilepath != "") {
		log.Fatalf("Both certificate filepath (%s) and certificate key filepath (%s) must be set when using manual TLS",
			TLSCertFilepathVarenv, TLSCertKeyFilepathVarenv)
	}

	if cnf.HttpsBindingAddress == "" && (cnf.TLSAutoDomain != "" || cnf.TLSCertFilepath != "") {
		log.Fatalf("HTTPS binding address (%s) must be set when using either auto TLS (%s) or manual TLS (%s and %s)",
			HttpsBindingAddressVarenv, TLSAutoDomainVarenv, TLSCertFilepathVarenv, TLSCertKeyFilepathVarenv)
	}

	if cnf.HttpBindingAddress == "" && cnf.TLSAutoDomain == "" && cnf.TLSCertFilepath == "" {
		log.Fatalf("HTTP binding address (%s) must be set if auto TLS (%s) and manual TLS (%s and %s) are both disabled",
			HttpBindingAddressVarenv, TLSAutoDomainVarenv, TLSCertFilepathVarenv, TLSCertKeyFilepathVarenv)
	}

	if cnf.HttpsBindingAddress != "" && cnf.TLSAutoDomain == "" && cnf.TLSCertFilepath == "" {
		log.Fatalf("HTTPS binding address (%s) is set but neither auto TLS (%s) nor manual TLS (%s and %s) are enabled",
			HttpsBindingAddressVarenv, TLSAutoDomainVarenv, TLSCertFilepathVarenv, TLSCertKeyFilepathVarenv)
	}

	if cnf.VaultPrefix == "" {
		cnf.VaultPrefix = "cubbyhole/"
	}

	log.Println("[INFO] HTTP Request Log enabled:", cnf.HttpRequestLog)
	log.Println("[INFO] HTTP Binding Address:", cnf.HttpBindingAddress)
	log.Println("[INFO] HTTPS Binding Address:", cnf.HttpsBindingAddress)
	log.Println("[INFO] HTTPS Redirect enabled:", cnf.HttpsRedirectEnabled)
	log.Println("[INFO] TLS Auto Domain:", cnf.TLSAutoDomain)
	log.Println("[INFO] TLS Cert Filepath:", cnf.TLSCertFilepath)
	log.Println("[INFO] TLS Cert Key Filepath:", cnf.TLSCertKeyFilepath)
	log.Println("[INFO] Vault prefix:", cnf.VaultPrefix)
	log.Println("[INFO] Passphrase num words default:", cnf.PassphraseNumWordsDefault)
	log.Println("[INFO] Passphrase separator default:", cnf.PassphraseSeparatorDefault)
	log.Println("[INFO] Passphrase capitalize default:", cnf.PassphraseCapitalizeDefault)
	log.Println("[INFO] Passphrase include number default:", cnf.PassphraseIncludeNumberDefault)
	log.Println("[INFO] Passphrase max number default:", cnf.PassphraseMaxNumberDefault)

	return cnf
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	return val
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}

	return val
}

func getEnvBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	return strings.ToLower(valStr) == "true"
}

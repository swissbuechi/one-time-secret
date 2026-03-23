package main

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-diceware/diceware"
)

type TokenResponse struct {
	Token     string `json:"token"`
	FileToken string `json:"filetoken,omitempty"`
	FileName  string `json:"filename,omitempty"`
}

type MsgResponse struct {
	Msg string `json:"msg"`
}

type PassphraseResponse struct {
	Passphrase string `json:"passphrase"`
}

type PassphraseHandler struct {
	wordlist      diceware.WordList
	numWords      int
	separator     string
	capitalize    bool
	includeNumber bool
	maxNumber     int
}

func NewPassphraseHandler(
	wordlist diceware.WordList,
	numWords int,
	separator string,
	capitalize bool,
	includeNumber bool,
	maxNumber int,
) *PassphraseHandler {
	return &PassphraseHandler{
		wordlist:      wordlist,
		numWords:      numWords,
		separator:     separator,
		capitalize:    capitalize,
		includeNumber: includeNumber,
		maxNumber:     maxNumber,
	}
}

type SecretHandlers struct {
	store SecretMsgStorer
}

func NewSecretHandlers(s SecretMsgStorer) *SecretHandlers {
	return &SecretHandlers{s}
}

func (s SecretHandlers) CreateMsgHandler(ctx echo.Context) error {
	var tr TokenResponse

	// Upload file if any
	file, err := ctx.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer src.Close()

		b, err := ioutil.ReadAll(src)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if len(b) > 0 {
			tr.FileName = file.Filename
			encodedFile := base64.StdEncoding.EncodeToString(b)

			filetoken, err := s.store.Store(encodedFile)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			tr.FileToken = filetoken
		}
	}

	// Handle the secret message
	msg := ctx.FormValue("msg")
	tr.Token, err = s.store.Store(msg)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, tr)
}

func (s SecretHandlers) GetMsgHandler(ctx echo.Context) error {
	m, err := s.store.Get(ctx.QueryParam("token"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	r := &MsgResponse{
		Msg: m,
	}
	return ctx.JSON(http.StatusOK, r)
}

func HealthHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, http.StatusText(http.StatusOK))
}

func redirect(ctx echo.Context) error {
	return ctx.Redirect(http.StatusPermanentRedirect, "/")
}

func (h *PassphraseHandler) GetPassphraseHandler(ctx echo.Context) error {
	list, err := diceware.GenerateWithWordList(h.numWords, h.wordlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c := cases.Title(language.German)

	if h.capitalize {
		for i, w := range list {
			list[i] = c.String(w)
		}
	}

	if h.includeNumber {
		n, err := randomInt(int64(h.maxNumber + 1))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if len(list) > 0 {
			list[len(list)-1] += strconv.FormatInt(n, 10)
		}
	}

	r := &PassphraseResponse{
		Passphrase: strings.Join(list, h.separator),
	}

	return ctx.JSON(http.StatusOK, r)
}

func randomInt(max int64) (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}

package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Di-nis/gophKeeper/pkg/logger"
	"go.uber.org/zap"
)

// readReq считывает тело запроса и декодирует его в данные.
func readReq(r *http.Request, data any) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		logger.Log.Debug("cannot read body", zap.Error(err))
		return fmt.Errorf("(req *Request) Read: cannot read body: %w", err)
	}

	reader := io.NopCloser(bytes.NewReader(body))

	dec := json.NewDecoder(reader)
	if err := dec.Decode(data); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		return fmt.Errorf("(req *Request) Read: cannot decode request JSON body: %w", err)
	}

	return nil
}

package reception

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func launchReceptionTestBrowser(t *testing.T) *rod.Browser {
	t.Helper()
	if os.Getenv("REG_E2E") != "1" {
		t.Skip("Skipping rod-based test; set REG_E2E=1 to run")
	}

	l := launcher.New().
		Headless(true).
		Devtools(false).
		Leakless(false).
		UserDataDir(t.TempDir())

	u, err := l.Launch()
	require.NoError(t, err)

	b := rod.New().ControlURL(u)
	require.NoError(t, b.Connect())
	t.Cleanup(func() {
		_ = b.Close()
	})
	return b
}

func TestConfirmReceptionReturnsErrorOnLivewireToastWithoutHanging(t *testing.T) {
	b := launchReceptionTestBrowser(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/despachos/recepcion", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!doctype html>
<html>
<body>
  <div id="recepcionarGuia" class="modal" style="display:block">
    <button class="btn btn-success" type="button">Confirmar</button>
  </div>
  <div id="toasts"></div>
  <script>
    document.querySelector("#recepcionarGuia button.btn-success").addEventListener("click", async () => {
      await fetch("/livewire/message/confirmar", { method: "POST" });
      document.getElementById("recepcionarGuia").style.display = "none";
      const toast = document.createElement("div");
      toast.className = "jq-toast-single jq-icon-error";
      toast.textContent = "×No se encontró el inventario del rubro: SAL REFINADA";
      document.getElementById("toasts").appendChild(toast);
    });
  </script>
</body>
</html>`)
	})
	mux.HandleFunc("/livewire/message/confirmar", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	page, err := b.Page(proto.TargetCreateTarget{URL: srv.URL + "/despachos/recepcion"})
	require.NoError(t, err)
	require.NoError(t, page.WaitLoad())

	rp := NewReceptionPage(page)
	err = rp.ConfirmReception()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SAL REFINADA")
}

func TestConfirmReceptionSucceedsWhenModalClosesWithoutError(t *testing.T) {
	b := launchReceptionTestBrowser(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/despachos/recepcion", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!doctype html>
<html>
<body>
  <div id="recepcionarGuia" class="modal" style="display:block">
    <button class="btn btn-success" type="button">Confirmar</button>
  </div>
  <script>
    document.querySelector("#recepcionarGuia button.btn-success").addEventListener("click", async () => {
      await fetch("/livewire/message/confirmar", { method: "POST" });
      document.getElementById("recepcionarGuia").style.display = "none";
    });
  </script>
</body>
</html>`)
	})
	mux.HandleFunc("/livewire/message/confirmar", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	page, err := b.Page(proto.TargetCreateTarget{URL: srv.URL + "/despachos/recepcion"})
	require.NoError(t, err)
	require.NoError(t, page.WaitLoad())

	rp := NewReceptionPage(page)
	require.NoError(t, rp.ConfirmReception())
}

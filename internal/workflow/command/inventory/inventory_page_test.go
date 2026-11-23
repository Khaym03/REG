package inventory

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

func launchTestBrowser(t *testing.T) *rod.Browser {
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

func TestInventoryPageSelectSubmitAndWaitForRubro(t *testing.T) {
	b := launchTestBrowser(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/inventario", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!doctype html>
<html>
<body>
  <div class="componenteCargaInicial">
    <span class="select2-selection" id="open-select2">SELECCIONAR</span>
    <ul id="dropdown" style="display:none">
      <li class="select2-results__option" data-val="2107">Sal Refinada Importada</li>
      <li class="select2-results__option" data-val="30">Sal Refinada</li>
    </ul>
    <button id="upload-btn" class="btn btn-primary" type="button"><i class="fas fa-upload"></i></button>
  </div>
  <table>
    <tbody>
      <tr><td>1</td><td>ACEITE PARA MANUFACTURA</td><td>0,651TM</td></tr>
    </tbody>
  </table>
  <script>
    let selectedVal = "";
    const btn = document.getElementById("upload-btn");
    document.getElementById("open-select2").addEventListener("click", () => {
      document.getElementById("dropdown").style.display = "block";
    });
    document.querySelectorAll(".select2-results__option").forEach(li => {
      li.addEventListener("click", async () => {
        document.getElementById("dropdown").style.display = "none";
        btn.disabled = true;
        await fetch("/livewire/message/set-rubro", { method: "POST", body: li.dataset.val });
        selectedVal = li.textContent.trim().toUpperCase();
        btn.disabled = false;
      });
    });
    btn.addEventListener("click", async () => {
      btn.disabled = true;
      const resp = await fetch("/livewire/message/submit-rubro", { method: "POST", body: selectedVal });
      if (resp.ok) {
        const tr = document.createElement("tr");
        tr.innerHTML = "<td>2</td><td>" + selectedVal + "</td><td>0,000TM</td>";
        document.querySelector("table tbody").appendChild(tr);
      }
      btn.disabled = false;
    });
  </script>
</body>
</html>`)
	})
	mux.HandleFunc("/livewire/message/set-rubro", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/livewire/message/submit-rubro", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(400 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	page, err := b.Page(proto.TargetCreateTarget{URL: srv.URL + "/inventario"})
	require.NoError(t, err)
	require.NoError(t, page.WaitLoad())

	invPage := NewInventoryPage(page)
	target := Rubro{Name: "SAL REFINADA"}

	require.NoError(t, invPage.SelectItem(target))
	require.NoError(t, invPage.Submit())
	require.NoError(t, invPage.WaitForRubro(target))

	rubros, err := invPage.ExtractExistingRubros()
	require.NoError(t, err)
	assert.Len(t, rubros, 2)
	assert.Equal(t, "SAL REFINADA", rubros[1].Name)
}

func TestInventoryPageSubmitReturnsErrorOnToast(t *testing.T) {
	b := launchTestBrowser(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/inventario", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!doctype html>
<html>
<body>
  <button id="upload-btn" class="btn btn-primary" type="button"><i class="fas fa-upload"></i></button>
  <div id="toast-container"></div>
  <script>
    document.getElementById("upload-btn").addEventListener("click", async () => {
      await fetch("/livewire/message/submit-rubro", { method: "POST" });
      const div = document.createElement("div");
      div.className = "jq-toast-single jq-icon-error";
      div.textContent = "×Error al cargar rubro";
      document.getElementById("toast-container").appendChild(div);
    });
  </script>
</body>
</html>`)
	})
	mux.HandleFunc("/livewire/message/submit-rubro", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	page, err := b.Page(proto.TargetCreateTarget{URL: srv.URL + "/inventario"})
	require.NoError(t, err)
	require.NoError(t, page.WaitLoad())

	invPage := NewInventoryPage(page)
	err = invPage.Submit()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Error al cargar rubro")
}

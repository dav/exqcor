package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dav/exqcor/server/internal/httpapi"
	"github.com/dav/exqcor/server/internal/netinfo"
	"github.com/dav/exqcor/server/internal/show"
	"github.com/dav/exqcor/server/internal/store"
	"github.com/dav/exqcor/server/internal/version"
	"github.com/dav/exqcor/server/internal/webui"
	"github.com/pkg/browser"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	dbPath := flag.String("db", defaultDBPath(), "path to the SQLite database file")
	noOpen := flag.Bool("no-open", false, "do not auto-open the admin page in a browser")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("exqcor", version.Version)
		return
	}

	st, err := store.Open(*dbPath)
	if err != nil {
		log.Fatalf("open database %s: %v", *dbPath, err)
	}
	defer st.Close()

	pass, firstRun, err := ensureCredentials(st)
	if err != nil {
		log.Fatalf("initialize credentials: %v", err)
	}

	baseURL := func() string {
		ip, _ := st.Setting("chosen_ip")
		if ip == "" {
			ip = netinfo.Best()
		}
		if ip == "" {
			ip = "127.0.0.1"
		}
		return fmt.Sprintf("http://%s:%d", ip, *port)
	}

	rt := show.New(st)
	if err := rt.Recover(); err != nil {
		log.Fatalf("recover show state: %v", err)
	}
	srv := httpapi.New(st, rt, webui.FS(), baseURL)

	addr := fmt.Sprintf(":%d", *port)
	localURL := fmt.Sprintf("http://127.0.0.1:%d/#/admin", *port)

	fmt.Printf("Exqcor %s\n", version.Version)
	fmt.Printf("  database:  %s\n", *dbPath)
	fmt.Printf("  admin:     %s\n", localURL)
	fmt.Printf("  phones:    %s\n", baseURL())
	if firstRun {
		fmt.Printf("  passphrase: %s  (needed for admin access from other devices; shown only on first run)\n", pass)
	}
	fmt.Println("If your OS firewall asks, click Allow — otherwise phones cannot connect.")

	if !*noOpen {
		go func() {
			time.Sleep(300 * time.Millisecond)
			if err := browser.OpenURL(localURL); err != nil {
				log.Printf("could not open browser: %v", err)
			}
		}()
	}

	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}

// defaultDBPath puts the database next to the executable so the unzipped
// folder is the whole app, falling back to the working directory.
func defaultDBPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "exqcor.db"
	}
	return filepath.Join(filepath.Dir(exe), "exqcor.db")
}

// ensureCredentials generates the admin passphrase and role join tokens on
// first run. Returns the plaintext passphrase only when newly generated.
func ensureCredentials(st *store.Store) (pass string, firstRun bool, err error) {
	hash, err := st.Setting("admin_pass_hash")
	if err != nil {
		return "", false, err
	}
	for _, key := range []string{"token_audience", "token_station", "token_actor"} {
		if _, err := st.EnsureSetting(key, store.RandomToken); err != nil {
			return "", false, err
		}
	}
	if hash != "" {
		return "", false, nil
	}
	pass = generatePassphrase()
	if err := st.SetSetting("admin_pass_hash", httpapi.HashPassphrase(st, pass)); err != nil {
		return "", false, err
	}
	return pass, true, nil
}

var passWords = []string{
	"amber", "birch", "cedar", "delta", "ember", "frost", "grove", "harbor",
	"ivory", "jade", "koala", "lunar", "maple", "north", "ocean", "piano",
	"quartz", "river", "solar", "tiger", "umber", "violet", "willow", "zephyr",
}

func generatePassphrase() string {
	pick := func(n int) int64 {
		v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
		if err != nil {
			panic(err)
		}
		return v.Int64()
	}
	return fmt.Sprintf("%s-%s-%02d",
		passWords[pick(len(passWords))],
		passWords[pick(len(passWords))],
		pick(100))
}

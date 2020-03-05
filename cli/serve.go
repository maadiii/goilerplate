package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"goldfish/app"
	"goldfish/controllers"
	. "goldfish/controllers/route"

	"golang.org/x/net/http2"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

func redirectToHTTPS(c *controllers.Controller) http.Handler {
	url := HTTPS_SCHEME +
		c.Config.DomainName +
		COLON + strconv.Itoa(c.Config.HTTPSPort)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(
			w,
			r,
			url+r.RequestURI,
			http.StatusMovedPermanently,
		)
	})
}

func serveAPI(
	ctx context.Context,
	conf *app.Config,
	c *controllers.Controller,
) {
	//certManager := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist("example.com"), //Your domain here
	//	Cache:      autocert.DirCache("certs"),            //Folder for storing certificates
	//}

	//r := mux.NewRouter()
	//r.PathPrefix(STATIC_PATH).Handler(http.StripPrefix(
	//	STATIC_PATH,
	//	http.FileServer(
	//		neuteredFileSystem{
	//			http.Dir(conf.Static),
	//		})))
	r := httprouter.New()
	r.ServeFiles(
		"/static/*filepath",
		neuteredFileSystem{http.Dir(conf.Static)},
	)

	Route(c, r)

	s := &http.Server{
		Addr:              fmt.Sprintf("%s%d", COLON, c.Config.HTTPSPort),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		//TLSConfig: &tls.Config{
		//	GetCertificate: certManager.GetCertificate,
		//},
	}
	http2.ConfigureServer(s, nil)

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			logrus.Error(err)
		}
		close(done)
	}()

	go func() {
		if err := http.ListenAndServe(
			COLON+strconv.Itoa(c.Config.HTTPPort),
			redirectToHTTPS(c),
		); err != nil {
			logrus.Error(err)
		}
	}()

	logrus.Infof(SERVING_LOG_INFO, c.Config.HTTPSPort)
	if err := s.ListenAndServeTLS(
		c.App.Config.TLS.CRT,
		c.App.Config.TLS.Key,
	); err != http.ErrServerClosed {
		logrus.Error(err)
	}

	<-done
}

var serveCli = &cobra.Command{
	Use:   SERVE_USE,
	Short: SERVE_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		a, err := app.New()
		if err != nil {
			return err
		}
		defer a.Close()

		c, err := controllers.New(a)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			logrus.Info(SHUTING_DOWN_SIGNAL)
			cancel()
		}()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			serveAPI(ctx, a.Config, c)
		}()

		wg.Wait()
		return nil
	},
}

func init() {
	rootCli.AddCommand(serveCli)
}

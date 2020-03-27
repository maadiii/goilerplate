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

	"goilerplate/interface/router"
	"goilerplate/registry"
	"goilerplate/usecase/controllers"

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

func redirectToHTTPS(c controllers.IRootController) http.Handler {
	baseController := c.GetBase()
	url := HTTPS_SCHEME +
		baseController.Config.DomainName +
		COLON + strconv.Itoa(baseController.Config.HTTPSPort)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(
			w,
			r,
			url+r.RequestURI,
			http.StatusMovedPermanently,
		)
	})
}

func serveAPI(ctx context.Context, c controllers.IRootController) {
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
		neuteredFileSystem{http.Dir(c.GetBase().Application.Config.Static)},
	)

	router.Route(r, c)

	s := &http.Server{
		Addr:              fmt.Sprintf("%s%d", COLON, c.GetBase().Config.HTTPSPort),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Second,
		//TLSConfig: &tls.Config{
		//	GetCertificate: certManager.GetCertificate,
		//},
	}
	http2.ConfigureServer(s, nil)

	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		logrus.Infof(SERVING_LOG_INFO, c.GetBase().Config.HTTPPort)
		if err := http.ListenAndServe(
			COLON+strconv.Itoa(c.GetBase().Config.HTTPPort),
			redirectToHTTPS(c),
		); err != nil {
			logrus.Error(err)
		}
	}()

	if err := s.ListenAndServeTLS(
		c.GetBase().Application.Config.TLS.CRT,
		c.GetBase().Application.Config.TLS.Key,
	); err != http.ErrServerClosed {
		logrus.Error(err)
	}
}

var serveCli = &cobra.Command{
	Use:   SERVE_USE,
	Short: SERVE_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		reg, err := registry.NewRegistry()
		if err != nil {
			return err
		}
		c := reg.NewRootController()

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
			serveAPI(ctx, c)
		}()

		wg.Wait()
		return nil
	},
}

func init() {
	rootCli.AddCommand(serveCli)
}

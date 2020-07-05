package app

import (
	"flag"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/joho/godotenv"
	"go-with-mongodb/app/config"
	"go-with-mongodb/app/controllers"
	"go-with-mongodb/app/repository"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	port = flag.String("-p", ":9000", "set app port")
)

func init() {
	flag.Parse()
}

func Run() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := config.NewDbConfig()

	fmt.Println("Initializing App configs", "db_dsn", cfg.String())
	fmt.Println("Listening on port", "http://localhost"+*port)

	session, err := mgo.Dial(cfg.String())
	if err != nil {
		log.Fatal("Error on connect mongodb:", err.Error())
	}
	defer session.Close()

	db := session.DB(cfg.DatabaseName())

	productsRepository := repository.NewProductsRepository(db)
	productsControllers := controllers.NewProductsControllers(productsRepository)

	container := restful.NewContainer()

	container.Filter(NCSACommonLogFormatLogger())

	wsProducts := new(restful.WebService)
	wsProducts.
		Path("/products").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	wsProducts.Route(wsProducts.POST("").To(productsControllers.PostProduct))
	wsProducts.Route(wsProducts.GET("").To(productsControllers.GetProducts))
	wsProducts.Route(wsProducts.GET("/{id}").To(productsControllers.GetProduct))
	wsProducts.Route(wsProducts.PUT("/{id}").To(productsControllers.PutProduct))
	wsProducts.Route(wsProducts.DELETE("/{id}").To(productsControllers.DeleteProduct))

	container.Add(wsProducts)

	cors := &restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		CookiesAllowed: true,
		Container:      container,
	}

	container.Filter(cors.Filter)

	container.Filter(container.OPTIONSFilter)

	server := &http.Server{Addr: *port, Handler: container}

	log.Fatal(server.ListenAndServe())
}

var logger = log.New(os.Stdout, "", 0)

func NCSACommonLogFormatLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var username = "-"
		if req.Request.URL.User != nil {
			if name := req.Request.URL.User.Username(); name != "" {
				username = name
			}
		}
		chain.ProcessFilter(req, resp)
		logger.Printf("%s - %s [%s] \"%s %s %s\" %d %d",
			strings.Split(req.Request.RemoteAddr, ":")[0],
			username,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			req.Request.Method,
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			resp.StatusCode(),
			resp.ContentLength(),
		)
	}
}

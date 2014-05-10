package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/events"
	httpFetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/processor"
	"github.com/wanelo/image-server/processor/cli"
	sm "github.com/wanelo/image-server/source_mapper/waneloS3"
	"github.com/wanelo/image-server/uploader"
	"github.com/wanelo/image-server/uploader/manta"
)

func main() {
	environment := flag.String("e", "development", "Specifies the environment to run this server under (test/development/production).")
	flag.Parse()

	path := "config/" + *environment + ".json"
	serverConfiguration, err := core.LoadServerConfiguration(path)

	adapters := &core.Adapters{
		Processor:    &cli.Processor{serverConfiguration},
		SourceMapper: &sm.SourceMapper{serverConfiguration},
	}

	serverConfiguration.Adapters = adapters

	mappings := make(map[string]string)
	mappings["p"] = "product/image"
	serverConfiguration.NamespaceMappings = mappings

	if err != nil {
		log.Panicln(err)
	}

	httpFetcher.ImageDownloads = make(map[string][]chan error)
	processor.ImageProcessings = make(map[string][]chan processor.ImageProcessingResult)

	go func() {
		mantaAdapter := manta.InitializeManta(serverConfiguration)
		uwc := uploader.UploadWorkers(mantaAdapter.Upload, serverConfiguration.MantaConcurrency)
		events.InitializeEventListeners(serverConfiguration, uwc)
	}()

	initializeRouter(serverConfiguration)
}

func initializeRouter(sc *core.ServerConfiguration) {
	log.Println("starting in "+sc.Environment, "on http://0.0.0.0:"+sc.ServerPort)

	m := martini.Classic()
	m.Map(sc)
	m.Get("/:namespace/:id1/:id2/:id3/:filename", genericImageHandler)

	log.Fatal(http.ListenAndServe(":"+sc.ServerPort, m))
}

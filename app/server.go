package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

var usoCounter = metrics.NewCounter("custom_app.uso_count")

func uso(w http.ResponseWriter, req *http.Request) {
	// Incrementar a métrica de uso
	usoCounter.Increment()

	// Lógica da rota "/uso" aqui...
	fmt.Fprintf(w, "Contador de uso incrementado")
}

func main() {
	// Iniciar o tracer Datadog
	if err := tracer.Start(tracer.WithEnv("development"), tracer.WithService("my-go-app")); err != nil {
		log.Fatal(err)
	}
	defer tracer.Stop()

	// Configurar a rota "/uso" para o manipulador "uso"
	http.HandleFunc("/uso", uso)

	// Configurar a rota "/metrics" para expor métricas Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Iniciar o servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("Servidor escutando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

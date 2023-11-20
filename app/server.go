package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer/contrib/gorilla/mux"
)

func ping(w http.ResponseWriter, req *http.Request) {
	// Criar uma span para a rota "ping"
	span, ctx := tracer.StartSpanFromRequest(req, "ping")
	defer span.Finish()

	// Lógica da rota "ping" aqui...
	fmt.Fprintf(w, "pong")
}

func main() {
	// Iniciar o tracer Datadog
	_, err := tracer.Start(
		tracer.WithEnv("development"),
		tracer.WithService("my-go-app"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Stop()

	// Configurar o roteador Gorilla Mux
	r := mux.NewRouter()
	r.HandleFunc("/ping", ping)

	// Configurar o middleware Datadog para o roteador
	r.Use(mux.Middleware)

	// Configurar a rota "/metrics" para expor métricas Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Iniciar o servidor HTTP com o roteador
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("Servidor escutando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

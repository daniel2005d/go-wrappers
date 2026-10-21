# go-wrappers

Pequeñas utilidades de Go para realizar solicitudes HTTP y leer archivos línea por línea.

## Instalación

```sh
go get github.com/daniel2005d/go-wrappers
```

Importa el paquete en tu proyecto:

```go
import "github.com/daniel2005d/go-wrappers/wrappers"
```

## Cliente HTTP

`HTTPClient` aplica un tiempo límite (10 segundos si el valor recibido es cero o negativo), agrega un encabezado User-Agent predeterminado y devuelve las respuestas de redirección sin seguirlas.

```go
package main

import (
	"fmt"
	"time"

	"github.com/daniel2005d/go-wrappers/wrappers"
)

func main() {
	client := wrappers.NewHTTPClient(5 * time.Second)
	client.SetHeader("Accept", "application/json")

	response, err := client.Get("https://example.com")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	fmt.Println(response.Status)
}
```

Entre los métodos disponibles están `Get`, `Request`, `Post` y `Do`. Quien realiza la solicitud debe cerrar el `Body` de cada respuesta recibida correctamente.

`NewHTTPClientInsecureTLS` desactiva la verificación del certificado TLS y del nombre del host. Úsalo solo cuando necesites explícitamente ese comportamiento.

La función `Request` del paquete usa el cliente HTTP predeterminado de Go.

## Leer un archivo línea por línea

```go
err := wrappers.ReadFile("input.txt", func(line string) error {
	fmt.Println(line)
	return nil
})
if err != nil {
	panic(err)
}
```

La función pasada como callback puede devolver un error para detener la lectura. `ReadFile` también devuelve errores al abrir el archivo o al leerlo con el escáner.

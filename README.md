# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Índice

- [TP0: Docker + Comunicaciones + Concurrencia](#tp0-docker--comunicaciones--concurrencia)
  - [Índice](#índice)
  - [Documentación](#documentación)
    - [Ejercicio N°1](#ejercicio-n1)
    - [Ejercicio N°2](#ejercicio-n2)
      - [Referencias](#referencias)
    - [Ejercicio N°3](#ejercicio-n3)
      - [Referencias](#referencias-1)
    - [Ejercicio N°4](#ejercicio-n4)
      - [Servidor](#servidor)
      - [Cliente](#cliente)
      - [Ejecucion](#ejecucion)
      - [Referencias](#referencias-2)
    - [Ejercicio N°5](#ejercicio-n5)
      - [Protocolo](#protocolo)
      - [Cuestiones para desarrollar un protocolo](#cuestiones-para-desarrollar-un-protocolo)
        - [short-read y short-write](#short-read-y-short-write)
        - [Endianness](#endianness)
      - [Separación de responsabilidades](#separación-de-responsabilidades)
      - [Ejecución](#ejecución)
      - [Referencias](#referencias-3)
    - [Ejercicio N°6](#ejercicio-n6)
      - [Ejecución](#ejecución-1)
    - [Ejercicio N°7](#ejercicio-n7)
      - [Ejecucion](#ejecucion-1)
  - [Instrucciones de uso](#instrucciones-de-uso)
    - [Servidor](#servidor-1)
    - [Cliente](#cliente-1)
    - [Ejemplo](#ejemplo)
  - [Parte 1: Introducción a Docker](#parte-1-introducción-a-docker)
    - [Ejercicio N°1:](#ejercicio-n1-1)
    - [Ejercicio N°2:](#ejercicio-n2-1)
    - [Ejercicio N°3:](#ejercicio-n3-1)
    - [Ejercicio N°4:](#ejercicio-n4-1)
  - [Parte 2: Repaso de Comunicaciones](#parte-2-repaso-de-comunicaciones)
    - [Ejercicio N°5:](#ejercicio-n5-1)
      - [Cliente](#cliente-2)
      - [Servidor](#servidor-2)
      - [Comunicación:](#comunicación)
    - [Ejercicio N°6:](#ejercicio-n6-1)
    - [Ejercicio N°7:](#ejercicio-n7-1)
  - [Parte 3: Repaso de Concurrencia](#parte-3-repaso-de-concurrencia)
    - [Ejercicio N°8:](#ejercicio-n8)
  - [Condiciones de Entrega](#condiciones-de-entrega)

## Documentación

### Ejercicio N°1

Para este ejercicio, se implementó el script `generar-compose.sh` que permite escribir un _docker compose_ para generar diferentes servicios, los cuales incluyen: un servidor, una network y un número variable de clientes. Estos últimos se especifican a través de un parámetro en el script; para su ejecución dejamos el comando a continuación.

```bash
./generar-compose.sh <nombre_del_archivo> <numero_de_clientes>
```

El nombre del archivo debe incluir su extensión (`.yaml`) y el número de clientes debe ser, obviamente, un número entero no negativo. Por otro lado, existe manejo de errores para el caso en que se pase una cantidad distinta a la requerida de parámetros o que los parámetros sean inválidos (ej: el número de clientes no es un número). Para usar ese _docker compose_ generado, debemos correr el siguiente comando:

```bash
docker compose -f <nombre_del_archivo> up -d --build
```

### Ejercicio N°2

En este ejercicio solo hubo modificaciones en el archivo `generador_compose.py` respecto al [Ejercicio N°1](#ejercicio-n1). La forma de utilizar `generar-compose.sh` sigue siendo la misma que en el ejercicio anterior. En cuanto a la implementación de este ejercicio, para evitar la reconstrucción de la imagen cada vez que cambian los archivos de configuración, se utilizó el mecanismo de volúmenes.

Los volúmenes que definimos fueron:

- `./server/config.ini:/config.ini:ro` para el servidor
- `./client/config.yaml:/config.yaml:ro` para el cliente

Recordemos que los volúmenes son de la forma `host_path:container_path[:ro]`, donde `host_path` corresponde a un path de la máquina que hostea el contenedor, mientras que `container_path` corresponde a un path dentro del contenedor. Notemos la presencia del parámetro `ro` al final del volumen, esto indica que el archivo es solo de lectura y no se puede modificar dentro del mismo. ¹

Por último, notemos que eliminamos de los environments lo siguiente:

- `LOGGING_LEVEL=DEBUG` para el servidor
- `CLI_LOG_LEVEL=DEBUG` para el cliente

Esto se debe a que, tanto el servidor como el cliente leen primero las variables de entorno y luego los archivos de configuración, pero siempre priorizando los valores definidos en las variables de entorno. En el caso del servidor, esto se encuentra en el `main.py`, el cual configura el objeto `ConfigParser` con las variables de entorno. Luego, el valor para cada parámetro se elige entre las variables de entorno o el `ConfigParser`, pero siempre priorizando las variables de entorno. ²

En el caso del cliente, este usa una librería llamada [viper](https://github.com/spf13/viper), la cual tiene un comentario que indica que existe una "estrategia" de priorización en la cual las variables de entorno tienen mayor prioridad que un archivo de configuración. ³

#### Referencias

1. “Volumes.” (2025, July). Docker Documentation. Retrieved August 22, 2025, from [https://docs.docker.com/engine/storage/volumes](https://docs.docker.com/engine/storage/volumes)
2. FIUBA - Sistemas Distribuidos (TA050) (Cátedra Roca). (n.d.). tp0-base/server/main.py at master · 7574-sistemas-distribuidos/tp0-base. GitHub. Retrieved August 22, 2025, from [https://github.com/7574-sistemas-distribuidos/tp0-base/blob/master/server/main.py#L9-L34](https://github.com/7574-sistemas-distribuidos/tp0-base/blob/master/server/main.py#L9-L34)
3. Spf. (n.d.). viper/viper.go at master · spf13/viper. GitHub. Retrieved August 22, 2025, from [https://github.com/spf13/viper/blob/master/viper.go#L107-L141](https://github.com/spf13/viper/blob/master/viper.go#L107-L141)

### Ejercicio N°3

Para este ejercicio creamos el archivo `validar-echo-server.sh`. Es importante asegurarse de que el servidor esté funcionando; para ello, primero debemos levantarlo previamente con _docker compose_. Para hacerlo, ejecutamos el siguiente comando:

```bash
make docker-compose-up
```

Luego, para validar que el servidor responde correctamente, debemos utilizar el script creado, el cual se ejecuta de la siguiente manera:

```bash
./validar-echo-server.sh
```

El mismo corre una imagen de `busybox`, la cual provee implementaciones simplificadas de varias utilidades comunes de Linux/UNIX, como por ejemplo: `nc` (netcat), `echo` y `sh`. Estas nos permiten chequear la conectividad con el servidor. ¹ Es importante resaltar que este script asume que el contenedor se llama `server` y corre en el puerto `12345`. Si queremos utilizar el validador en otro servidor, debemos cambiar esas variables.

Hay que tener en cuenta que el servidor vive dentro de una `network` que se crea con el _docker compose_, por lo cual, cuando corramos el contenedor con `docker run`, debemos proveer esta `network` generada por el _docker compose_. Por suerte, existe el flag `--network=<nombre_de_network>` para el comando `docker run`, con lo cual, cuando ejecutemos el comando netcat en el contenedor, sabrá cuál es la dirección del `server` y su puerto, sin exponer los puertos hacia el host. ²

El comando que ejecuta el script `validar-echo-server.sh` es el siguiente:

```bash
docker run \
  --network=$NETWORK_NAME \
  --rm \
  $DOCKER_IMAGE \
  sh -c "echo '$MESSAGE_FOR_SERVER' | nc $SERVER_ADDRESS $SERVER_PORT"`
```

Ahí corremos el contenedor con imagen `$DOCKER_IMAGE`, que corresponde a `busybox`, luego seteamos la network con el flag `--network=$NETWORK_NAME`. El flag `--rm` es para que, una vez termine la ejecución, se remueva el contenedor. El comando que ejecutará el contenedor es: `sh -c "echo '$MESSAGE_FOR_SERVER' | nc $SERVER_ADDRESS $SERVER_PORT"`.

El comando `sh -c <command_string>` significa que utiliza el intérprete de shell para ejecutar `<command_string>`. El comando a ejecutar por nosotros es `echo '$MESSAGE_FOR_SERVER' | nc $SERVER_ADDRESS $SERVER_PORT`, `echo` escribe el mensaje en `stdout`, que mediante el pipeline (`|`) se pasa como `stdin` a `nc`, el cual lo lee y envía al servidor con dirección `$SERVER_ADDRESS` y puerto `$SERVER_PORT`. ³

Como es un echo server, lo que se va a imprimir por pantalla es lo mismo que enviamos, por lo cual, para que la conexión haya sido efectiva, comparamos el resultado de ese comando con el mensaje del `stdout` que imprime el servidor y, si son iguales, imprimimos el mensaje de éxito (`action: test_echo_server | result: success`), caso contrario, el mensaje de error (`action: test_echo_server | result: fail`). ⁴

#### Referencias

1. Docker Official Image. (n.d.). busybox. Retrieved August 22, 2025, from [https://hub.docker.com/_/busybox](https://hub.docker.com/_/busybox)
2. “Networking.” (2025, May 1). Docker Documentation. Retrieved August 22, 2025, from [https://docs.docker.com/engine/network/#user-defined-networks](https://docs.docker.com/engine/network/#user-defined-networks)
3. sh(1p) - Linux manual page. (n.d.). Retrieved August 22, 2025, from [https://man7.org/linux/man-pages/man1/sh.1p.html](https://man7.org/linux/man-pages/man1/sh.1p.html)
4. nc(1) - Linux man page. (n.d.). Retrieved August 22, 2025, from [https://linux.die.net/man/1/nc](https://linux.die.net/man/1/nc)

### Ejercicio N°4
En este ejercicio implementamos toda la lógica para un _graceful shutdown_ cuando ambos reciben un _SIGTERM_. Vamos a explicar cómo funciona del lado del servidor primero y luego pasaremos a explicar el cliente. Finalmente explicaremos su ejecución para probar el _graceful shutdown_.

#### Servidor

Cuando recibe un _SIGTERM_, lo primero que hace es cerrar el socket de escucha. Esto significa que no aceptará nuevas conexiones entrantes. Si el servidor estaba bloqueado esperando en una llamada a `accept()`, esta se interrumpe y genera un error de tipo _OSError_. ¹ 

En el caso de que el servidor esté atendiendo a un cliente en ese momento, no corta la comunicación de manera inmediata. En lugar de eso, completa la interacción con el cliente y recién después continúa con el proceso de apagado. De esta forma, el cierre es más ordenado y no interrumpe bruscamente las conexiones que ya estaban activas.

Si el servidor ya fue cerrado, simplemente salimos del `accept_connection()` y listo. Por otro lado, un breve comentario de qué son los parámetros para manejar las _signals_: uno de ellos es `signum`, este indica qué número de señal fue enviado. El otro parámetro es el `frame`, que corresponde al punto donde se estaba ejecutando el programa cuando se recibió la _signal_. ²

#### Cliente

Cuando el cliente recibe un _SIGTERM_, inicia un cierre ordenado de sus operaciones. Si estaba intentando conectarse al servidor, interrumpe el intento y finaliza sin volver a reintentar. En el caso de estar enviando o recibiendo mensajes, las operaciones se desbloquean y devuelven un error que indica que la conexión fue cerrada. Ante este evento, si se trata de un envío de la señal _SIGTERM_ y no de un verdadero error, el cliente simplemente corta el ciclo de ejecución y se apaga. ³

Finalmente, expliquemos cómo funcionan las _signals_ en Go. La forma de usar las _signals_ es creando una _go routine_, que sería como un hilo de ejecución, y esta espera para _poppear_ un elemento de una _queue_ bloqueante. Cuando la _signal_ se "dispara", se desbloquea el hilo y hace el _graceful shutdown_. ⁴

#### Ejecucion

Para probar que efectivamente funciona el _graceful shutdown_, debemos primero correr el siguiente comando:

```bash
make docker-compose-up
```

Luego debemos correr el comando para hacerle un _stop_ y un _down_ a los container, para eso corremos el siguiente comando:

```bash
make docker-compose-down
```

Es importante recalcar, que este ultimo comando dara de baja tanto al cliente como al servidor. Si quisieramos dar de baja algun servicio en particular, debemos correr el siguiente comando:

```bash
docker compose -f docker-compose-dev.yaml stop -t <tiempo> <nombre_del_servicio>
```

Para chequear la salida de codigos de error y los diferentes logs que fue enviando tanto el cliente como el servidor, corremos el siguiente comando:

``` bash
make docker-compose-logs
```

#### Referencias

1. socket — Low-level networking interface. (n.d.). Python Documentation. Retrieved August 23, 2025, from [https://docs.python.org/3/library/socket.html](https://docs.python.org/3/library/socket.html)
2. signal — Set handlers for asynchronous events. (n.d.). Python Documentation. Retrieved August 23, 2025, from [https://docs.python.org/3/library/signal.html](https://docs.python.org/3/library/signal.html)
3. net package - net - Go Packages. (n.d.). Retrieved August 23, 2025, from [https://pkg.go.dev/net](https://pkg.go.dev/net)
4. Go by Example: Signals. (n.d.). Retrieved August 23, 2025, from [https://gobyexample.com/signals](https://gobyexample.com/signals)

### Ejercicio N°5

En esta parte hablaremos sobre el desarrollo del protocolo, cuestiones generales relacionadas a la hora de desarrollar un protocolo, como así también cuestiones de implementación del lado del cliente y del servidor. Empezaremos primero hablando sobre cómo está desarrollado el protocolo, qué mensajes tiene y qué mensajes espera cada una de las partes para reaccionar de acuerdo a ello.

#### Protocolo

El protocolo desarrollado para este ejercicio consiste en 3 simples mensajes: el primero es la apuesta (`BET`) en sí, el segundo corresponde a un mensaje `OK` y, por último, un mensaje de `FAIL`. Los mensajes tienen dos partes: una cabecera (header) y un cuerpo (body). El _header_ tiene un _code message_ de 1 byte y permite diferenciar qué mensaje estamos enviando, y un _message length_ de 2 bytes que determina de qué tamaño es el mensaje. A continuación, dejamos un cuadro con cada uno de los 3 posibles _code message_.

| Tipo de mensaje | Valor del _code message_ |
|-----------------|--------------------------|
| BET             | `0x01`                   |
| OK              | `0x02`                   |
| FAIL            | `0x03`                   |

Ahora, el _body_ es el resto del mensaje. Dependiendo del _header_ que lea el protocolo, este sabrá cómo está compuesto el resto del _body_. El mensaje de tipo `OK` no tiene un _body_ asociado, ni tampoco tiene un _message length_; solamente consiste en su _code message_. Distinto es el caso para los mensajes `BET` y `FAIL`, que sí tienen un _body_ asociado y un _message length_. Veamos cómo está compuesto el mensaje completo, a través de un diagrama de cada campo:

```txt
  +----------------------------------+  
  |       Code message (1 byte)      |  
  +----------------------------------+  
  |                                  |  
  |     Message length (2 bytes)     |  
  |                                  |  
  +----------------------------------+  
  |                                  |  
  |    Message (variable length)     |  
  |                                  |  
  +----------------------------------+  
```

<p align="center"><strong>Formato de mensaje del protocolo</strong></p>

Con este diagrama, expliquemos el mensaje `BET` con cada uno de los campos:

- El `Code message` corresponde al tipo de mensaje (_header_).
- El `Message length` corresponde al tamaño del mensaje (_header_).
- El `Message` es donde está almacenado el contenido en sí y tiene el tamaño especificado por el valor de `Message length` (_body_).

Para el caso del mensaje de tipo `FAIL`, el valor de `Message` es el error que ocurrió en el servidor, ya sea por mensaje inválido o por algún problema interno. Ahora, para el caso del mensaje de tipo `BET`, dentro de `Message` almacenamos la apuesta utilizando un separador que nos permite distinguir cada uno de los campos de la apuesta. El separador utilizado es `;`. Entonces, para leer eso, simplemente debemos leer según el `Message length` y dividir el contenido usando el separador. El mensaje tiene el siguiente formato: `AgencyID;FirstName;LastName;Document;Birthdate;Number`, los cuales corresponden a cada uno de los diferentes parámetros de la apuesta.

Ahora pasemos a analizar cómo es el estado de cada mensaje, es decir, cómo responde cada entidad ante la llegada de cada mensaje. Existen solamente dos caminos: un caso donde la apuesta llega sin problema, por lo que el servidor contesta con el mensaje `OK` asegurándole al cliente que su apuesta fue guardada exitosamente. No existe una política de reenvío de apuestas, dado que estamos trabajando sobre un _TCP socket_, por lo que el envío de mensajes es seguro. Dejamos un diagrama de secuencia de los mensajes en el caso exitoso.

```txt
  +-----------+                      +------------+ 
  |  Cliente  |                      |  Servidor  | 
  +-----------+                      +------------+ 
        V                                  V  
        |            (BET msg)             |  
        | --------------->---------------- |  
        V                                  V  
        |            (OK msg)              |  
        | ---------------<---------------- |  
        V                                  V  
        |                                  |  
```

<p align="center"><strong>Diagrama de secuencia: Caso 1 donde el mensaje de la apuesta llega correctamente</strong></p>

Por otro lado, existe la posibilidad de que la apuesta no haya sido almacenada correctamente, ya sea por algún problema del lado del servidor, como errores en la lectura del mensaje o fallos al guardar la apuesta en un archivo. En ese caso, el servidor responde con un mensaje `FAIL` al cliente para que sepa que su apuesta no ha sido almacenada. Si ocurre un error, el cliente no reintentará enviar la apuesta; se deberá corregir lo que esté mal en los parámetros de la apuesta. El diagrama de secuencia para este caso sería el siguiente.

```txt
  +-----------+                      +------------+ 
  |  Cliente  |                      |  Servidor  | 
  +-----------+                      +------------+ 
        V                                  V  
        |            (BET msg)             |  
        | --------------->---------------- |  
        V                                  V  
        |           (FAIL msg)             |  
        | ---------------<---------------- |  
        V                                  V  
        |                                  |  
```

<p align="center"><strong>Diagrama de secuencia: Caso 2 donde el mensaje de la apuesta llega incorrectamente</strong></p>

#### Cuestiones para desarrollar un protocolo

##### short-read y short-write

Para desarrollar este protocolo hubo que tener en cuenta dos cuestiones muy importantes, la primera es la cuestión relacionada con el _short read y short write_ ¹ y la segunda con el envío de bytes a través de una red (_host-to-network_) y la recepción de esos (_network-to-host_). Empecemos con la primera cuestión. Básicamente cuando mandamos mensajes a través de un _socket_, existe la posibilidad de que este no envíe todos los bytes que le pedimos, es decir la librería en cuestión no nos asegura mandar el mensaje completo, esto puede llevar a que el servidor o el cliente se queden bloqueados esperando bytes que nunca fueron enviados. Existe una solución para eso, escribamos un pseudocódigo y veamos cómo funciona

```python
msg # Mensaje a ser enviado/recibido
sock # Socket para enviar/recibir los mensajes
length # Tamaño del mensaje a ser enviado/recibido
#-------------------------------------------------------------------------

def recv_all(length):
  msg = ""
  readed = 0

  while readed < length:
   chunk = socket.recv(length - readed)
   if not chunk:
    raise Error("Closed connection")
   msg += chunk
   readed += len(chunk)

  return msg

def send_all(msg, length):
  writed = 0

  while writed < length:
   size = socket.send(msg[writed:])
   if size < 0:
    raise Error("Closed connection")
   writed += size
```

<p align="center"><strong>Pseudocódigo: <i>short-read</i> y <i>short-write</i></strong></p>

Las librerías de _sockets_ devuelven cuántos bytes se leyeron/escribieron, entonces a partir de eso podremos saber cuántos bytes faltan por enviar/recibir, por eso estamos en un loop que nos asegura que hayamos leído todo el mensaje que esperamos recibir/enviar. En el caso de Python, nos devuelven lo que recibieron, sino que devuelven el mensaje en sí, sabiendo el tamaño del mensaje parcial recibido (_chunk_) podremos saber cuánto leímos. ² Particularmente en Python, existe una función implementada por el _socket_ llamada `sendall()` que nos asegura enviar todos los bytes, por lo que el problema del _short write_ en Python está resuelto. ³

##### Endianness

Respecto al problema del envío o recepción de mensajes a través de la red, tiene que ver con cómo implementan los números las computadoras. Algunas usan el formato _big endian_, es decir el byte más significativo primero, o la opción de _little endian_, o sea el byte menos significativo primero. Esa diferencia nos trae la obligación de asegurarnos el correcto envío y recepción de los bytes a través de la red. Es por eso que existe la necesidad de implementar una función que transforme del _endianness_ de la computadora (_host_) a una a través de la red (_network_), estas funciones se llaman _host-to-network_ (`hton`), lo mismo aplica al revés, es decir _network-to-host_ (`ntoh`).

Para el caso de este protocolo, lo necesitamos para bytes relacionados con los tamaños del `nombre` y el `apellido`, la convención es utilizar _big endian_ para la _network_. ⁴ Así que las funciones _hton_ y _ntoh_ deberán utilizar una transformación a _big endian_ o decir que determinado conjunto de bytes es _big endian_. Como los tamaños son de 2 bytes, a estas funciones particularmente se las llama `htons()` y `ntohs()`, la `s` es por `short` y son números de 16 bits (2 bytes) sin signo. Para Python ya existe una forma nativa de transformar las cosas en bytes según el _endianness_, ⁵ para el caso de Go usamos una librería de la biblioteca estándar llamada `encoding/binary`. ⁶

#### Separación de responsabilidades

En esta seccion, hablamos de cómo están separadas las responsabilidades en cada parte. Básicamente en ambas existe una implementación de una clase o estructura llamada `Protocol`, que lo que hace es darnos la funcionalidad necesaria para enviar una apuesta o un mensaje avisando si todo fue exitoso o no. Digamos que el cliente o el servidor, respectivamente, se comunican con esta capa y esta última se encarga de serializar y enviar los mensajes o de recibirlos y deserializarlos para separar correctamente las responsabilidades entre la lógica de negocio y la comunicación. En ambos casos, necesitan que les proveamos un _socket_ por donde se envían o reciben los mensajes y ellos se encargan del resto. Dejamos un diagrama simplificado de cómo funciona

```txt
            Logical layer                                           Logical layer  
          +----------------+                                     +----------------+
          | +------------+ |                                     | +------------+ |
          | |            | |                                     | |            | |
          | |   Client   | |                                     | |   Server   | |
          | |            | |                                     | |            | |
          | +------------+ |                                     | +------------+ |
          +-----|-----^----+                                     +----|-----^-----+
                |     |                                               |     |
  (send bet) 1  |     |  (recv ok/failure) 6     (send ok/failure) 4  |     |  (recv bet) 3
                |     |                                               |     |
          +-----|-----|-----------------------------------------------|-----|-----+
          |     v     |                                               v     |     |
          | +------------+         (send msg to server) 2          +------------+ |
          | |            |---------------------------------------->|            | |
          | |  Protocol  |                                         |  Protocol  | |
          | |            |<----------------------------------------|            | |
          | +------------+         (send msg to client) 5          +------------+ |
          |                                                                       |
          +-----------------------------------------------------------------------+
                                    Communication layer
```

<p align="center"><strong>Diagrama de componentes de lógica de negocio y comunicación</strong></p>

El diagrama muestra cómo se establece la comunicación entre el cliente y el servidor mediante una separación clara de responsabilidades. En la parte superior se encuentran las capas lógicas, que representan la funcionalidad propia del cliente y del servidor. Estas capas se limitan a definir qué acciones deben llevarse a cabo, como enviar una apuesta o procesar el resultado, sin involucrarse en los detalles de la transmisión de datos.

La capa inferior corresponde a la comunicación, implementada a través del componente `Protocol`. Su función principal es abstraer el manejo del socket y encargarse de todo lo relacionado con el envío y la recepción de mensajes. Cuando la lógica del cliente desea enviar una apuesta, se la entrega al componente, que se ocupa de serializarla y transmitirla al servidor. Del mismo modo, cuando el servidor recibe un mensaje, su instancia de `Protocol` lo deserializa y lo entrega a la lógica del servidor para que sea procesado.

Una vez que el servidor determina el resultado de la apuesta, genera una respuesta de éxito o fallo y la pasa nuevamente a su `Protocol`. Este la serializa y la envía al cliente, donde el componente correspondiente se encarga de deserializarla y entregarla a la capa lógica. De esta forma, ambas partes pueden interactuar sin necesidad de gestionar directamente los aspectos de bajo nivel de la comunicación.

En conjunto, el diagrama refleja un diseño en el que la lógica de negocio y la comunicación permanecen desacopladas. El cliente y el servidor se enfocan exclusivamente en las reglas y decisiones propias de la aplicación, mientras que el componente `Protocol` actúa como intermediario responsable de transformar y transportar los mensajes entre ambas partes.

#### Ejecución

Es importante resaltar que la apuesta es enviada hacia el servidor gracias a la existencia de nuevas variables de entorno del lado del cliente. Si quisiéramos utilizar otros valores para esas variables, estas se encuentran disponibles en el `generador_compose.py`, por lo tanto podremos generar un _docker compose_ con las directivas del [Ejercicio N°1](#ejercicio-n1), es decir, correr lo siguiente:

```bash
./generar-compose.sh <nombre_del_archivo> <numero_de_clientes>
```

Para la ejecución del protocolo debemos correr el _docker compose_ con el siguiente comando, si se llama `docker-compose-dev.yaml`:

```bash
make docker-compose-up
```

Si se llama distinto, usar:

```bash
docker compose -f <nombre_del_archivo> up -d --build
```

Si queremos ver los logs del cliente o el servidor, el comando es:

```bash
make docker-compose-logs
```

Las variables de entorno nuevas son:

- CLI_NOMBRE
- CLI_APELLIDO
- CLI_DOCUMENTO
- CLI_NACIMIENTO
- CLI_NUMERO

Es importante resaltar que el documento y el número sean efectivamente valores numéricos, y que la fecha de nacimiento tenga el formato `YYYY-MM-DD`; de lo contrario, el cliente arrojará un error y no se cerrará correctamente. Estos valores se encuentran dentro del `generador_compose.py`, como ya mencionamos antes. Por lo tanto, para utilizar otros valores, se deben modificar las nuevas constantes definidas en el script correspondiente y volver a generar el archivo de _compose_.

#### Referencias

1. File Descriptors – CS 61 2018. (n.d.). Retrieved August 24, 2025, from [https://cs61.seas.harvard.edu/site/2018/FileDescriptors/](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/)
2. recv — Low-level networking interface. (n.d.-c). Python Documentation. Retrieved August 24, 2025, from [https://docs.python.org/3/library/socket.html#socket.socket.recv](https://docs.python.org/3/library/socket.html#socket.socket.recv)
3. sendall — Low-level networking interface. (n.d.-b). Python Documentation. Retrieved August 24, 2025, from [https://docs.python.org/3/library/socket.html#socket.socket.sendall](https://docs.python.org/3/library/socket.html#socket.socket.sendall)
4. htons(3) - Linux man page. (n.d.). Retrieved August 24, 2025, from [https://linux.die.net/man/3/htons](https://linux.die.net/man/3/htons)
5. Built-in Types. (n.d.). Python Documentation. Retrieved August 24, 2025, from [https://docs.python.org/3/library/stdtypes.html](https://docs.python.org/3/library/stdtypes.html)
6. binary package - encoding/binary - Go Packages. (n.d.). Retrieved August 24, 2025, from [https://pkg.go.dev/encoding/binary](https://pkg.go.dev/encoding/binary)


### Ejercicio N°6

En esta sección, explicaremos todo lo referente al envío de múltiples apuestas en un mismo mensaje, a través de la lectura de archivos CSV. Las principales modificaciones son tanto del lado del cliente como del lado del servidor. Del lado del cliente, ahora generamos una variable de entorno que asocia dónde está el path del archivo que queremos leer. A su vez, generamos un volumen para asociar un archivo del host con el que definimos en esta variable de entorno. La variable de entorno se llama `CLI_DATA_FILEPATH`, tanto el volumen como la variable de entorno fueron agregados al `generador_compose.py` para facilitar la creación del _docker compose_.

Hubo cambios en el protocolo de comunicación entre el cliente y el servidor para facilitar la llegada de los _batches_, esto nos llevó a la creación de un nuevo _code message_ referente a la llegada de estos _batches_, el cual se llama `BATCH` y además tiene un tamaño correspondiente a su _header_, que nos indica cuántas apuestas hay en el mensaje y de esa forma facilita la lectura del mismo.

Del lado del servidor, estas apuestas se levantan y se guardan en un archivo. Por otro lado, sabemos que existen restricciones respecto a cuántas apuestas pueden ser enviadas y cuánto es el tamaño del paquete a ser enviado. Para la cuestión referente a la cantidad de apuestas, esta es una variable de entorno configurable llamada `CLI_BATCH_MAX_AMOUNT` en el cliente, por el lado del tamaño de paquete este es una constante escrita en el cliente, seteada en 8 kB.

Esta restricción nos llevó a implementar otro _code message_ llamado `FINISH` para indicarle al servidor que ya hemos enviado todas las apuestas en diferentes _batches_ y evitar que el cliente se desconecte y reconecte varias veces para diferentes _batches_. Por otro lado, agregamos un _header_ de 2 bytes tanto al `OK` como al `FAIL` para indicar al cliente cuántas apuestas fueron efectivamente guardadas. El `FAIL` indica al cliente que no fueron todas guardadas en forma efectiva, pero no cuáles.

#### Ejecución

Para la ejecución de este ejercicio, necesitamos CSVs para los clientes. Para eso corremos el siguiente comando:

```bash
unzip dataset.zip
```

Es importante que los archivos CSV estén dentro de la carpeta `.data`, para eso recomendamos hacer un `cd` dentro de la misma y luego correr el comando anterior. Luego necesitamos generar los clientes, podemos generar tantos como queramos, pero recomendamos hacerlo con 1. Para eso corremos el siguiente comando:

```bash
./generar-compose.sh docker-compose-dev.yaml 1
```

Luego, levantamos el archivo _compose_ `docker-compose-dev.yaml` con el siguiente comando:

```bash
make docker-compose-up
```

Luego, podemos chequear los logs con el siguiente comando:

```bash
make docker-compose-logs
```

### Ejercicio N°7

En este ejercicio nos enfrentamos al desafio de añadir la funcionalidad de devolver los ganadores del sorteo, lo que llevo ademas a la aceptacion multiples clientes y sincronizarlos para que reciban sus respectivos ganadores, eso nos llevo a desarrollar una serie de mensajes que permitan esa sicronizacion. El primer paso fue definir una variable de entorno configurable con la cantidad de clientes llamada `AMOUNT_OF_CLIENTS` del lado del servidor y que se encuentra en el archivo `generador_compose.py`. Este nos permite saber cuantos clientes debemos manejar. En funcion de la cantidad de clientes que pidamos al script este seteara dicha variable.

Para implementar esta funcionalidad, nuevamente mandamos todos los _batches_ de apuestas que las agencias, luego que todos los clientes fueron atendidos, es decir mandaron sus apuestas, el servidor se prepara para leer el archivo de apuestas y empezar a enviar los ganadores a las respectivas agencias. En primer lugar, debemos identificar a cada agencia, para eso usamos un nuevo mensaje del tipo `ID_HEADER` que incluye un 2 bytes correspondientes al ID de la agencia, esto permitira asociar un diccionario de IDs con los sockets de cada cliente.

Una vez hecha la identificacion, procedemos a recibir los _batches_ de todos los clientes, cuando un cliente termine permanecera a la espera de los ganadores, pero para eso debera esperar que los demas clientes terminen. Una vez que todos terminan ahi procedemos a leer el archivo e identificar a los ganadores y se lo mandamos a su socket correspondiente en funcion del ID de agencia que tiene la apuesta, evitando hacer un _broadcast_ de todos los ganadores a todos los clientes.

Finalizada la lectura, enviamos un mensaje del tipo `FINISH` que indica al cliente que la loteria ha terminado, el cliente guarda la lista de ganadores y sabe la cantidad de ganadores sabiendo el tamaño de la lista e imprime la cantidad.

#### Ejecucion

La ejecucion es muy similar al apartado anterior. Dejamos la explicacion paso a paso. Para la ejecucion de este ejercicio, necesitamos CSVs para los clientes. Para eso corremos el siguiente comando

```bash
unzip dataset.zip
```

Es importante que los archivos CSV esten dentro de la carpeta `.data`, para eso recomendamos hacer un `cd` dentro de la misma y luego correr el comando anterior. Luego necesitamos generar los clientes, podemos generar tantos como queramos. Para eso corremos el siguiente comando

```bash
./generar-compose.sh docker-compose-dev.yaml 1
```

Luego, levantamos el archivo _compose_ `docker-compose-dev.yaml`, con el siguiente comando

```bash
make docker-compose-up
```

Luego, podemos chequear los logs con el siguiente comando

```bash
make docker-compose-logs
```

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).


### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

## Condiciones de Entrega
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios y que aprovechen el esqueleto provisto tanto (o tan poco) como consideren necesario.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deberán existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.
 (hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Se espera que se redacte una sección del README en donde se indique cómo ejecutar cada ejercicio y se detallen los aspectos más importantes de la solución provista, como ser el protocolo de comunicación implementado (Parte 2) y los mecanismos de sincronización utilizados (Parte 3).

Se proveen [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra. Se exige que la resolución de los ejercicios pase tales pruebas, o en su defecto que las discrepancias sean justificadas y discutidas con los docentes antes del día de la entrega. El incumplimiento de las pruebas es condición de desaprobación, pero su cumplimiento no es suficiente para la aprobación. Respetar las entradas de log planteadas en los ejercicios, pues son las que se chequean en cada uno de los tests.

La corrección personal tendrá en cuenta la calidad del código entregado y casos de error posibles, se manifiesten o no durante la ejecución del trabajo práctico. Se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección informados  [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).

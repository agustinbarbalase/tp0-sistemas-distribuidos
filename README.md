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
  - [Instrucciones de uso](#instrucciones-de-uso)
    - [Servidor](#servidor)
    - [Cliente](#cliente)
    - [Ejemplo](#ejemplo)
  - [Parte 1: Introducción a Docker](#parte-1-introducción-a-docker)
    - [Ejercicio N°1:](#ejercicio-n1-1)
    - [Ejercicio N°2:](#ejercicio-n2-1)
    - [Ejercicio N°3:](#ejercicio-n3-1)
    - [Ejercicio N°4:](#ejercicio-n4)
  - [Parte 2: Repaso de Comunicaciones](#parte-2-repaso-de-comunicaciones)
    - [Ejercicio N°5:](#ejercicio-n5)
      - [Cliente](#cliente-1)
      - [Servidor](#servidor-1)
      - [Comunicación:](#comunicación)
    - [Ejercicio N°6:](#ejercicio-n6)
    - [Ejercicio N°7:](#ejercicio-n7)
  - [Parte 3: Repaso de Concurrencia](#parte-3-repaso-de-concurrencia)
    - [Ejercicio N°8:](#ejercicio-n8)
  - [Condiciones de Entrega](#condiciones-de-entrega)

## Documentación

### Ejercicio N°1

Tal como indica el ejercicio, el comando para ejecutar el script es:

`./generar-compose.sh <nombre_del_archivo> <numero_de_clientes>`

El nombre del archivo debe incluir su extensión y el número de clientes debe ser, obviamente, un número entero no negativo. Por otro lado, existe manejo de errores para el caso en que se pase una cantidad distinta a la requerida de parámetros.

Respecto a la implementación del código, contamos con una función `create_docker_compose` que recibe el nombre del archivo (`filename`) y el número de clientes (`number_of_clients`). La misma escribe sobre el archivo las definiciones necesarias, luego se llama a una serie de funciones para definir los distintos servicios que tiene nuestro compose.

La primera es `define_server`, la cual define el servicio del servidor. La segunda es `define_client`, que define el cliente; notemos que recibe como uno de sus parámetros el número de cliente que corresponde definir. Por último, `define_network` define la red de nuestro compose. Es importante notar que todas las funciones reciben el archivo donde se escribe el compose.

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

Para este ejercicio creamos el archivo `validar-echo-server.sh`, el cual se ejecuta de la siguiente manera:

`./validar-echo-server.sh`

El mismo corre una imagen de `busybox`, la cual provee implementaciones simplificadas de varias utilidades comunes de Linux/UNIX, como por ejemplo: `nc` (netcat), `echo` y `sh`. Estos nos permiten chequear la conectividad con el servidor. ¹

Hay que tener en cuenta que el servidor vive dentro de una `network` que se crea con el Docker Compose, por lo cual, cuando corramos el contenedor con `docker run`, debemos proveer esta `network` generada por el Docker Compose. Por suerte, existe el flag `--network=<nombre_de_network>` para el comando `docker run`, con lo cual, cuando corramos el comando netcat en el contenedor, sabrá cuál es la dirección del `server` y su puerto, sin exponer los puertos hacie el host. ²

El comando que ejecuta el script `validar-echo-server.sh` es el siguiente:

`docker run \
  --network=$NETWORK_NAME \
  --rm \
  $DOCKER_IMAGE \
  sh -c "echo '$MESSAGE_FOR_SERVER' | nc $SERVER_ADDRESS $SERVER_PORT"`

Ahí corremos el contenedor con imagen `$DOCKER_IMAGE`, que corresponde a `busybox`, luego seteamos la network con el flag `--network=$NETWORK_NAME`. El flag `--rm` es para que, una vez termine la ejecución, se remueva el contenedor. El comando que ejecutará el contenedor es: `sh -c "echo '$MESSAGE_FOR_SERVER' | nc $SERVER_ADDRESS $SERVER_PORT"`.

El comando `sh -c <command_string>` significa que utiliza el intérprete de shell para ejecutar `<command_string>`. El comando a ejecutar por nosotros es `echo '$MESSAGE_FOR_SERVER' | nc $SERVER_ADDRESS $SERVER_PORT`, `echo` escribe el mensaje en `stdout`, que mediante el pipeline (`|`) se pasa como `stdin` a `nc`, el cual lo lee y envía al servidor con dirección `$SERVER_ADDRESS` y puerto `$SERVER_PORT`. ³

Como es un echo server, lo que se va a imprimir por pantalla es lo mismo que enviamos, por lo cual, para que la conexión haya sido efectiva, comparamos el resultado de ese comando con el mensaje del `stdout` que imprime el servidor y, si son iguales, imprimimos el mensaje de éxito (`action: test_echo_server | result: success`), caso contrario, el mensaje de error (`action: test_echo_server | result: fail`). ⁴

#### Referencias

1. Docker Official Image. (n.d.). busybox. Retrieved August 22, 2025, from [https://hub.docker.com/_/busybox](https://hub.docker.com/_/busybox)
2. “Networking.” (2025, May 1). Docker Documentation. Retrieved August 22, 2025, from [https://docs.docker.com/engine/network/#user-defined-networks](https://docs.docker.com/engine/network/#user-defined-networks)
3. sh(1p) - Linux manual page. (n.d.). Retrieved August 22, 2025, from [https://man7.org/linux/man-pages/man1/sh.1p.html](https://man7.org/linux/man-pages/man1/sh.1p.html)
4. nc(1) - Linux man page. (n.d.). Retrieved August 22, 2025, from [https://linux.die.net/man/1/nc](https://linux.die.net/man/1/nc)

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

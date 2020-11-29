# **TraceIP**
Proyecto realizado como parte de una entevista para MercadoLibre.

## **Dependencias**
----
**Redis:**
```
$ docker run --name redisdb -p 6379:6379 -d redis
```

## **Tutorial de uso:**
----
1. **Clonar repositorio**
2. **Posicionarse junto al docker-compose.yml**
3. **Ejecutar:**
```
$ docker-compose up --build
```
4. **Una vez finalizado ingresar a:** 

    [http://localhost:3000/rapidoc/](http://localhost:3000/rapidoc/)

5. **Para ejecutar una consulta se puede utilizar el botón TRY que nos brinda la documentación. Para cambiar el REQUEST BODY se debe ingresar a la pestaña EXAMPLE en REQUEST. En RESPONSE se podrá ver la respuesta obtenida.**

![TryButtonImage](https://github.com/NicolasBachs/md-images/blob/main/traceip-howtotry.png?raw=true)

## **CURL:**
----
**TraceIP:**
```
$ curl -X POST "http://127.0.0.1:3000/traceip" \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -d '{"ip":"83.44.196.93"}' 
```
**Statistics:**
```
$ curl -X GET "http://127.0.0.1:3000/statistics"
```

## **Consigna**
----

Para coordinar acciones de respuesta ante fraudes, es útil tener disponible información
contextual del lugar de origen detectado en el momento de comprar, buscar y pagar. Para
ello, entre otras fuentes, se decide crear una herramienta que dado un IP obtenga
información asociada:
Construir una aplicación que dada una dirección IP, encuentre el país al que pertenece, y
muestre:
- El nombre y código ISO del país
- Los idiomas oficiales del país
- Hora(s) actual(es) en el país (si el país cubre más de una zona horaria, mostrar
    todas)
- Distancia estimada entre Buenos Aires y el país, en km.
- Moneda local, y su cotización actual en dólares (si está disponible)

Basado en la información anterior, es necesario contar con un mecanismo para poder consultar las siguientes estadísticas de utilización del servicio con los siguientes agregados

- Distancia más lejana a Buenos Aires desde la cual se haya consultado el servicio
- Distancia más cercana a Buenos Aires desde la cual se haya consultado el servicio
- Distancia promedio de todas las ejecuciones que se hayan hecho del servicio. Ver ejemplo:

    | País   | Distancia | Invocaciones |
    |--------|-----------|--------------|
    | Brasil | 2862 km   | 10           |
    | España | 10040 km  | 5            |

En este caso la cuenta a realizar debería ser: 
```
(2862 km * 10 + 10040 km* 5) / 15 = 5254 km
```


Para resolver la información, pueden utilizarse las siguientes APIs públicas:

- **Geolocalización de IPs:** [https://ip2country.info](https://ip2country.info/)
- **Información de paises:** [​http://restcountries.eu](​http://restcountries.eu/)
- **Información sobre monedas:** [​http://fixer.io](http://fixer.io/)

Otras consideraciones:

- La aplicación puede ser en línea de comandos o web. En el primer caso se espera que el IP sea un parámetro, y en el segundo que exista un form donde escribir la dirección.
- La aplicación deberá hacer un uso racional de las APIs, evitando hacer llamadas innecesarias.
- La aplicación puede tener estado persistente entre invocaciones.
- Además de funcionamiento, prestar atención al estilo y calidad del código fuente.
- La aplicación deberá poder correr ser construida y ejecutada dentro de un contenedor Docker (incluir un Dockerfile e instrucciones para ejecutarlo).

Ejemplo (el formato es tentativo y no tiene porque ser exactamente así):
    
    > traceip 83.44.196.93

    - IP: 83.44.196.93, fecha actual: 21/11/2016 16:01:23
    - País: España (spain)
    - ISO Code: es
    - Idiomas: Español (es)
    - Moneda: EUR (1 EUR = 1.0631 U$S)
    - Hora: 20:01:23 (UTC) o 21:01:23 (UTC+01:00)
    - Distancia estimada: 10270 kms (-34, -64) a (40, -4)
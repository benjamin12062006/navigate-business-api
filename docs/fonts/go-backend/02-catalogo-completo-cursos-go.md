# 02 — Catálogo completo de cursos Go disponibles

Fecha de creación: 2026-09-27
Tipo de fuente: Inventario de biblioteca local
Origen: `/Users/benji/Documents/chrome-extencions/trancripciones/Go`
Fecha de consulta: 2026-09-27
Relacionada con: [Mapa general](./01-mapa-general.md), [rutas complementarias](./complementarias/01-rutas-complementarias.md), [otros autores verificados](./complementarias/02-otros-autores-go.md) y [reconstrucciones incorporadas al proyecto](../../reconstrucciones-cursos-go/README.md)

## Necesidad de consulta

Conservar una ruta directa hacia todos los cursos Go de la biblioteca local y dejar establecido que el proyecto puede buscar información en cualquiera de ellos cuando una necesidad técnica, funcional o pedagógica lo requiera.

## Aclaración sobre el análisis existente

Las 16 reconstrucciones utilizadas como evidencia técnica principal proceden de dos cursos atribuidos por el usuario a Tiago Taquelim:

1. `APIs de Producción con Go`.
2. `Building, scaling and deploying distributed systems in Golang and Kubernetes`.

Por tanto, el núcleo inicial de arquitectura, repositorios, transacciones, autenticación y autorización sí se apoyó principalmente en un instructor. Eso fue una prioridad de búsqueda, no una prohibición de consultar otros autores.

También se verificaron como contraste `Ecommerce Seguro con Go`, atribuido en la investigación previa a Trevor Sawler, y `Go Completo`, atribuido a Maximilian Schwarzmüller. Los demás cursos se inventariaron por su contenido y ruta local; no se atribuye un instructor cuando ese dato no se ha comprobado en la transcripción.

## Regla de consulta futura

Se puede buscar información en **cualquiera de los cursos de este catálogo** cuando responda a una necesidad del proyecto. La ruta principal no es exclusiva.

El procedimiento será:

1. Formular la necesidad concreta: arquitectura, auth, SQL, testing, pagos, concurrencia, IA u otra.
2. Buscar primero en todas las transcripciones originales pertinentes, excluyendo archivos `Reconstrucción guiada`.
3. Comparar los cursos que aporten enfoques realmente distintos o complementarios.
4. Seleccionar las lecciones que respondan a la necesidad; no adoptar un curso completo por coincidencia débil.
5. Reconstruir tantas lecciones como hagan falta cuando el código, los contratos o el orden de cambios necesiten evidencia más precisa.
6. Citar por separado la transcripción, la reconstrucción derivada y la adaptación realizada para `navigate-business-api`.

## Cursos directamente orientados a backend

| ID | Curso y ruta directa | Contenido aprovechable |
|---|---|---|
| `F-GO-01` | [APIs de Producción con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go>) | Arquitectura, servidor HTTP, PostgreSQL, CRUD, usuarios, autenticación, autorización, Redis, logging, testing, CORS, rate limiting, métricas, CI/CD y producción. |
| `F-GO-02` | [Building, scaling and deploying distributed systems in Golang and Kubernetes](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes>) | Arquitectura distribuida, repositorios, HTTP, gRPC, Kubernetes, mensajería, pagos, observabilidad, confiabilidad y despliegue. |
| `F-GO-03` | [Aplicaciones Web con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Aplicaciones Web con Go>) | Routing, middleware, sesiones, PostgreSQL, correo, autenticación, administración segura, testing y despliegue. |
| `F-GO-04` | [Ecommerce Seguro con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Ecommerce Seguro con Go>) | Pagos, suscripciones, autenticación, recuperación de contraseña, administración, reembolsos, paginación, validación y microservicios. |
| `F-GO-05` | [Go Completo](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Go Completo>) | La sección 11 desarrolla una API REST con SQL, usuarios, hash de contraseñas, JWT, rutas protegidas y autorización. |
| `F-GO-06` | [Microservicios Distribuidos con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Microservicios Distribuidos con Go>) | Servicios de autenticación, logging y correo; RabbitMQ, RPC, gRPC, Docker Swarm, Kubernetes y testing de microservicios. |
| `F-GO-07` | [React con API Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/React con API Go>) | API Go, PostgreSQL, JWT y conexión de un frontend React con el backend. |
| `F-GO-08` | [Vue 3 con API Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Vue 3 con API Go>) | Backend REST, autenticación y conexión entre Vue y Go. |
| `F-GO-09` | [Testing en Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Testing en Go>) | Handlers, middleware, sesiones, repositorios, integración con PostgreSQL, JWT, REST API y SPA. |
| `F-GO-10` | [WebSockets con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/WebSockets con Go>) | Comunicación bidireccional y funcionalidades de tiempo real. |

## Cursos complementarios para construir backends

| ID | Curso y ruta directa | Contenido aprovechable |
|---|---|---|
| `F-GO-11` | [Patrones de Diseño en Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Patrones de Diseño en Go>) | Repository Pattern, separación de responsabilidades y repositorios sustituibles para pruebas. |
| `F-GO-12` | [Concurrencia en Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Concurrencia en Go>) | Goroutines, mutexes, channels, trabajo concurrente, correo y testing. |
| `F-GO-13` | [Módulos Reutilizables en Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Módulos Reutilizables en Go>) | Diseño, publicación y consumo de paquetes o módulos reutilizables. |
| `F-GO-14` | [Go desde Cero](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Go desde Cero>) | Fundamentos y una sección dedicada a construir una aplicación web sencilla. |
| `F-GO-15` | [Fundamentos de Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Fundamentos de Go>) | Tipos, funciones, interfaces, errores, módulos, paquetes, goroutines y channels. Es base del lenguaje, no un curso backend especializado. |

## Cursos especializados que pueden alimentar servicios futuros

| ID | Curso y ruta directa | Contenido aprovechable |
|---|---|---|
| `F-GO-16` | [IA y Machine Learning con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/IA y Machine Learning con Go>) | Búsqueda, modelos, NLP, Hugging Face y LLM; posible referencia para análisis o proveedores de IA. |
| `F-GO-17` | [RAG con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/RAG con Go>) | Recuperación de información y generación aumentada; posible referencia para búsqueda inteligente sobre catálogos o documentos. |

## Cursos Go inventariados fuera del alcance backend actual

Estos cursos forman parte de la biblioteca y también pueden consultarse si una necesidad futura los vuelve pertinentes, pero no son fuentes backend prioritarias:

| ID | Curso y ruta directa | Alcance principal |
|---|---|---|
| `F-GO-18` | [Aplicaciones de Escritorio con Go y Fyne](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Aplicaciones de Escritorio con Go y Fyne>) | Aplicaciones de escritorio y Fyne. |
| `F-GO-19` | [Desarrollo de Juegos con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Desarrollo de Juegos con Go>) | Juegos, escenas, físicas, colisiones y despliegue del juego. |

## Evidencia

La biblioteca local contiene 19 directorios de cursos Go. Diez están orientados directamente a APIs, aplicaciones web, comercio, microservicios, integración frontend o testing de backend; cinco aportan fundamentos, arquitectura o concurrencia; dos cubren capacidades de IA aplicables a servicios futuros; y dos no están actualmente relacionados con el backend.

Las rutas de las tablas fueron comprobadas en el sistema de archivos el 2026-09-27. El catálogo registra disponibilidad y temática; no afirma que todos los cursos hayan sido leídos completamente, reconstruidos o aplicados.

## Interpretación y límites

La prioridad concedida a Tiago Taquelim define el punto de partida, pero no convierte sus cursos en la única evidencia permitida. Ante una necesidad concreta puede resultar más útil un curso especializado en testing, comercio, sesiones, WebSockets, concurrencia o IA.

Antes de atribuir una decisión técnica a un instructor se debe verificar el autor y leer la lección original. El nombre de una carpeta o sección permite localizar material, pero no sustituye la lectura de sus pasajes.

## Aplicación en el proyecto

Este documento funciona como índice de búsqueda. Autoriza y facilita consultar cualquier curso listado cuando el desarrollo necesite información, pero no introduce por sí mismo una dependencia, patrón, proveedor, arquitectura ni funcionalidad en el backend.

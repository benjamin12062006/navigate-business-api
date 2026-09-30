# 01 — Rutas complementarias de aprendizaje Go backend

Fecha de creación: 2026-09-27
Tipo de fuente: Inventario local de cursos y rutas de consulta
Origen: `/Users/benji/Documents/chrome-extencions/trancripciones/Go`
Fecha de consulta: 2026-09-27
Relacionada con: [Mapa general](../01-mapa-general.md), [arquitectura principal](../arquitectura/01-tiago-arquitectura-y-repositorios.md), [autenticación principal](../autenticacion-autorizacion/01-tiago-autenticacion-y-autorizacion.md), [operaciones](../operaciones/01-produccion-observabilidad-y-confiabilidad.md) y [otros autores verificados](./02-otros-autores-go.md)

## Necesidad de consulta

Conservar rutas alternativas para profundizar temas concretos cuando los cursos principales no basten, sin mezclar automáticamente enfoques incompatibles ni reconstruir material que todavía no responde una pregunta del proyecto.

## Localización verificable

La búsqueda se ejecutó sobre la biblioteca completa de Go y encontró estas rutas de backend relevantes:

| Prioridad | Curso local | Secciones o lecciones útiles | Uso recomendado |
|---:|---|---|---|
| 1 | `APIs de Producción con Go` | Arquitectura, servidor, DB, CRUD, logging, usuarios, email, auth, Redis, testing, shutdown, rate limit, CORS, métricas, CI/CD y despliegue | Referencia principal para el monolito API |
| 1 | `Building, scaling and deploying distributed systems in Golang and Kubernetes` | Diseño, capas, HTTP, gRPC, Kubernetes, eventos, pagos, observabilidad, confiabilidad, MongoDB y producción | Referencia principal cuando haya distribución real |
| 2 | `Go Completo` | Sección 11, lecciones 160–197: REST API, SQL, usuarios, passwords, JWT y autorización | Segundo ejemplo compacto de flujo completo |
| 2 | `Patrones de Diseño en Go` | Sección 5, lecciones 38–49: Repository Pattern y repositorio de pruebas | Contrastar diseño de repositorios |
| 2 | `Testing en Go` | Secciones 9–13: repository, integración con PostgreSQL, REST/JWT y SPA | Consultar cuando se diseñe cada estrategia de prueba |
| 2 | `Aplicaciones Web con Go` | Routing, middleware, sesiones, PostgreSQL, testing, auth, administración y despliegue | Alternativa basada en sesiones y aplicación web |
| 2 | `Ecommerce Seguro con Go` | Autenticación, rutas protegidas, password reset, administración, microservicios y validación | Flujos de comercio y seguridad de cuenta |
| 2 | `Microservicios Distribuidos con Go` | Broker, auth, logger, mail, RPC, gRPC, Kubernetes y testing | Contraste para sistemas distribuidos |
| 3 | `React con API Go` | Backend, PostgreSQL, JWT y conexión con React | Integración práctica con frontend React |
| 3 | `Vue 3 con API Go` | Backend, autenticación y conexión con Vue | Integración práctica con frontend Vue |
| 3 | `Concurrencia en Go` | Goroutines, mutexes, channels, email y testing | Trabajo concurrente cuando exista necesidad medida |
| 3 | `WebSockets con Go` | Comunicación bidireccional | Notificaciones o estado en tiempo real |
| 3 | `Módulos Reutilizables en Go` | Paquetes y módulos reutilizables | Extraer librerías solo después de repetición real |
| 3 | `IA y Machine Learning con Go` y `RAG con Go` | Proveedores, NLP, LLM y recuperación | Futura capa de análisis con IA, no arquitectura base |

### Puntos de entrada concretos

- `C-GO-01`: [Go Completo — Sección 11](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Go Completo/Sección 11 Course Project Build a REST API %28incl. Authentication & SQL Database%29>). Destacan `180. Don't Store Plain-text Passwords!`, `185. Generating JWT`, `187. Adding Token Verification`, `188. Adding Route Protection` y `192. Adding Authorization To Restrict Users From Editing & Deleting`.
- `C-PAT-01`: [Patrones de Diseño en Go — Sección 5](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Patrones de Diseño en Go/Sección 5 The Repository Pattern>). Destacan las lecciones 44–46 y 48–49.
- `C-TST-01`: [Testing en Go — Sección 9](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Testing en Go/Sección 9 The Repository Pattern>) y [Sección 12](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Testing en Go/Sección 12 Testing REST APIs>).
- `C-WEB-01`: [Aplicaciones Web con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Aplicaciones Web con Go>), especialmente secciones 11–17.
- `C-ECO-01`: [Ecommerce Seguro con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Ecommerce Seguro con Go>), especialmente secciones 6–9 y 13–15.
- `C-MIC-01`: [Microservicios Distribuidos con Go](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Microservicios Distribuidos con Go>), para comparar decisiones de comunicación y despliegue.

## Evidencia

El inventario confirma que la biblioteca contiene varias implementaciones de autenticación, repositorios, pruebas e integración frontend. Esto permite contrastar una decisión cuando exista una duda concreta, pero la mera existencia de un curso no valida automáticamente su enfoque para este proyecto.

## Interpretación y límites

`I-COM-01`: la ruta principal evita mezclar prematuramente Gin, Chi, sesiones, JWT y proveedores externos. Las fuentes complementarias se consultarán para una pregunta específica y se registrará el contraste.

`I-COM-02`: no se han leído ni reconstruido íntegramente todos los cursos de esta tabla. La evidencia actual verifica su localización y temario por archivos; cualquier decisión que dependa de su contenido exige abrir los pasajes correspondientes.

`I-COM-03`: reconstruir más lecciones está permitido sin un máximo artificial cuando aporte comprensión o trazabilidad. Debe seguir haciéndose una guía por lección y con mapa exacto a la transcripción original.

## Aplicación en el proyecto

Estas rutas quedan como índice de consulta futura. No introducen dependencias, patrones ni tests en `navigate-business-api` por sí solas.

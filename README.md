# Navigate Business API

## Auth con Google y Firebase (backend)

La API acepta exclusivamente **ID tokens de Firebase Authentication obtenidos con Google Sign-In**. No recibe correo/contraseña ni tokens OAuth de Google. Verifica firma, proyecto, expiración, revocación, usuario deshabilitado y `firebase.sign_in_provider = google.com`. La cuenta interna se guarda con SQL directo en PostgreSQL de Firebase SQL Connect. El cliente Expo queda para un incremento posterior.

Configuración necesaria en un proyecto de desarrollo:

1. Habilitar **Google** y deshabilitar **correo/contraseña** en Firebase Authentication.
2. Crear o vincular la instancia PostgreSQL de SQL Connect. En `firebase dataconnect:sql:setup`, declinar que SQL Connect administre las migraciones. La [CLI oficial](https://firebase.google.com/docs/sql-connect/cli-reference) documenta esa opción y `dataconnect:sql:grant`.
3. Conceder a la identidad que ejecuta la API el rol IAM `roles/cloudsql.client`. Crear además un usuario PostgreSQL integrado para la API —por ejemplo, `catalogo_api`— con los privilegios SQL mínimos sobre sus tablas. Son identidades diferentes: la cuenta de servicio abre el canal de Cloud SQL y el usuario PostgreSQL autentica la sesión dentro de la base.
4. Copiar [`.env.example`](.env.example) como `.env` y configurar `FIREBASE_PROJECT_ID`, `INSTANCE_CONNECTION_NAME`, `DB_NAME`, `DB_USER` y `DB_PASSWORD`. `HTTP_ADDR` es opcional. En local, ejecutar `gcloud auth application-default login`; no definir `GOOGLE_APPLICATION_CREDENTIALS` si se desea usar ese ADC. En Cloud Run, asignar la cuenta de servicio a la revisión. El emulador Auth solo se admite con `APP_ENV=development` y `FIREBASE_AUTH_EMULATOR_HOST`.

Desde la raíz de `navigate-business-api`, con la identidad SQL propietaria, ejecutar la migración **antes** de arrancar la API:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1
make migrate-up
go run ./cmd/api
```

El binario `migrate` debe estar disponible en `PATH` (`$(go env GOPATH)/bin` si se instaló con Go). `make migrate-up` aplica únicamente versiones pendientes y registra el estado en `schema_migrations`; `make migrate-down` revierte **una** versión y, con la única migración actual, elimina las tablas de cuentas y sus datos. Se ejecuta solamente en un entorno donde se desee ese borrado. La API no migra al iniciar. Solo la CLI externa de migraciones sigue necesitando `DATABASE_URL`, normalmente apuntando a un Cloud SQL Auth Proxy local y usando una identidad SQL propietaria; la API usa directamente Cloud SQL Go Connector y las variables `DB_*`. No ejecutar `firebase dataconnect:sql:migrate` sobre las mismas tablas. Si ya se aplicó la migración con el antiguo `cmd/migrate`, hay que conciliar manualmente su tabla de versiones antes de usar `golang-migrate`; no volver a ejecutar `up` sobre las tablas existentes.

Tras Google Sign-In en el futuro cliente, este enviará su **ID token de Firebase** por HTTPS:

```http
POST /api/v1/auth/sync
Authorization: Bearer <FIREBASE_ID_TOKEN>
```

La primera sincronización devuelve `201` y `{"account_id":"<uuid>","status":"active"}`; las siguientes devuelven `200` y el mismo ID. `GET /api/v1/auth/me` devuelve la cuenta ya sincronizada sin crearla. Token ausente, inválido o de un proveedor distinto produce `401`; cuenta suspendida, `403`; cuenta no sincronizada, `404`; indisponibilidad de Firebase o SQL, `503`. No se devuelven UID ni correo en estas respuestas.

El [plan](docs/plans/01-autenticacion-firebase-y-cuenta-interna.md), la [investigación oficial](docs/fonts/firebase/01-google-auth-admin-y-sql-connect.md), la [guía completa de conexión](docs/fonts/firebase/04-guia-cloud-sql-go-desde-cero.md) y la [bitácora de conexión](bitacora-construccion/config/02-conexion-cloud-sql-go-adc.md) documentan decisiones y comprobaciones. El conector y la autenticación PostgreSQL ya pasaron `Ping`; el recorrido completo queda pendiente de aplicar migraciones y verificar los endpoints con un ID token real.

La estructura de auth se lee así: [`internal/auth`](internal/auth/) contiene cuenta, servicio y puerto de repositorio; sus subcarpetas [`http`](internal/auth/http/), [`firebase`](internal/auth/firebase/) y [`postgres`](internal/auth/postgres/) son adaptadores. [`internal/auth/module`](internal/auth/module/) reúne el servicio, handler y autenticador que pueden recibir otras funcionalidades. [`internal/config`](internal/config/) valida el entorno, [`internal/app`](internal/app/) crea una sola instancia de Firebase y del pool SQL y conecta los módulos; [`cmd/api/main.go`](cmd/api/main.go) arranca el servidor. El [estudio de módulos Nest y Go](docs/fonts/go-backend/arquitectura/02-nest-modulos-y-composicion-explicita-go.md) y la [bitácora de configuración](bitacora-construccion/config/01-composicion-modular-y-dependencias-compartidas.md) explican esta organización.

## Primer incremento: Hello World

Desde esta carpeta, con las variables y la migración indicadas arriba, iniciar la API:

```bash
go run ./cmd/api
```

### Desarrollo con recarga automática

Air está registrado como herramienta del módulo Go. Para iniciar la API y
reiniciarla automáticamente cada vez que se guarde un archivo `.go`:

```bash
go tool air
```

Air usa [`.air.toml`](.air.toml)
para compilar `./cmd/api` en la carpeta temporal `tmp/`. Esa carpeta está
ignorada por Git. Air se usa solamente durante el desarrollo; la aplicación
sigue siendo un servidor Go normal y puede iniciarse con `go run ./cmd/api`.

En otra terminal, comprobar la ruta:

```bash
curl -i http://localhost:8080/api/v1/hello-world
```

Debe responder `200 OK`, `Content-Type: application/json` y:

```json
{"message":"Hello, world!"}
```

El recorrido de la ruta pública es `main.go → router.go → JSON`:

- `cmd/api/main.go` inicia el servidor HTTP local con la configuración y el handler creados en `internal/config` e `internal/app`.
- `internal/httpapi/router.go` agrupa las rutas en `/api/v1` y configura middleware.
- `internal/httpapi/router.go` escribe la respuesta de demostración.

`GET` consulta, `/api` identifica la API, `/v1` identifica su contrato y
`/hello-world` es nuestra ruta de demostración. Esta versión de la URL es
independiente de los tags de Git: un cambio menor no requiere una URL `/v2`.

Para una petición desde JavaScript en el navegador de Astro:

```js
const response = await fetch('http://localhost:8080/api/v1/hello-world');
if (!response.ok) throw new Error(`La API respondió ${response.status}`);
const data = await response.json();
console.log(data.message);
```

La API permite CORS desde `http://localhost:4321` y `http://127.0.0.1:4321`.
Si Astro usa otro puerto, ajustar `AllowedOrigins` en `router.go`. CORS permite
al navegador leer la respuesta entre orígenes distintos; no es autenticación.
Este incremento prepara el backend; el siguiente será mostrar el mensaje en
una página Astro usando `fetch` en el navegador.

**Referencia y adaptación:** Tiago enseña servidor HTTP, agrupación por versión
y Chi en la [clase 26](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 5 Scaffolding our API Server/26. Setting up the HTTP server and API.md>),
JSON en la [clase 23](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 4 Building a server from TCP to HTTP/23. Encoding & Decoding JSON Requests.md>)
y CORS para conectar el frontend en la [clase 79](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 20 Handling CORS/79. Handling CORS errors.md>).
`/api/v1`, el paquete `httpapi`, el saludo y los orígenes de Astro son adaptaciones
para este proyecto. Las capas descritas más abajo son la dirección futura;
por ahora solo se crean las piezas necesarias para esta petición.

## Referencia principal de desarrollo

Este proyecto se construirá siguiendo principalmente los lineamientos, criterios y prácticas enseñados por [**Tiago Taquelim**](https://www.udemy.com/user/tiago-taquelim-3/) en sus cursos de Go disponibles en la biblioteca local de transcripciones, especialmente:

- `APIs de Producción con Go`.
- `Building, scaling and deploying distributed systems in Golang and Kubernetes`.

Durante el desarrollo se consultarán primero las transcripciones originales de estos cursos para fundamentar las decisiones sobre arquitectura, servidores HTTP, Chi, bases de datos, autenticación, autorización, pruebas y despliegue.

Cuando el proyecto requiera un tema que Tiago no cubra, o que sus cursos no expliquen con suficiente profundidad para resolver la necesidad concreta, se consultarán también cursos de otros autores almacenados en:

```text
/Users/benji/Documents/chrome-extencions/trancripciones/Go
```

El material adicional servirá como contraste y complemento, no solo como una lista decorativa. Ya quedaron identificadas dos rutas de otros autores especialmente útiles:

- **Trevor Sawler — `Ecommerce Seguro con Go`**: pagos, reembolsos, autenticación, administración de usuarios, validación y separación entre frontend, API y servicios.
- **Maximilian Schwarzmüller — `Go Completo`**: recorrido compacto de API REST con SQL, registro, contraseñas, JWT y autorización.

Tiago Taquelim continúa siendo la referencia principal. Cada incremento indicará qué fuente se consultó, qué se adaptó y qué decisión pertenece al proyecto. Boris Paskhaver sigue siendo una referencia pedagógica bienvenida, pero la biblioteca local de Go no contiene actualmente un curso suyo; por eso no se le atribuirán decisiones técnicas de Go sin una fuente verificable.

## Método de aprendizaje y construcción

Este proyecto también tiene un propósito de aprendizaje. Se construirá por partes, desde los cimientos hasta las funcionalidades más avanzadas, avanzando de manera gradual y comprobando cada incremento antes de comenzar el siguiente.

Aunque se utilice Astra u otro modelo potente para ayudar durante el desarrollo, no se deberá saltar directamente al resultado final ni generar toda la aplicación de una sola vez. Cada etapa deberá:

1. Explicar el concepto de Go que se está aplicando y qué problema resuelve.
2. Consultar primero el material pertinente de Tiago Taquelim.
3. Adaptar ese conocimiento a las necesidades reales de este proyecto, sin copiar mecánicamente el proyecto del curso.
4. Implementar un incremento pequeño, completo y comprensible.
5. Verificar que el incremento funciona antes de continuar.
6. Registrar qué parte proviene del curso y qué parte es una decisión propia del proyecto.

El crecimiento será escalonado: primero el servidor y las rutas con Chi; después la organización por capas; luego proyectos y formularios; y más adelante persistencia, usuarios, permisos, análisis estadístico, integración con IA, generación de informes y colaboración. La prioridad será comprender y mantener lo construido, no alcanzar todas las funcionalidades lo más rápido posible.

La elección profesional para tu proyecto será:

> **Monolito modular con arquitectura por capas y límites hexagonales.**

No haremos una arquitectura hexagonal “pura” llena de interfaces y carpetas. Usaremos tres capas claras, como enseña Tiago, y aplicaremos puertos/adaptadores solamente donde aportan valor.

### Por qué esta es la mejor opción

No existe una arquitectura oficial para proyectos Go. La propia documentación explica que los servidores pueden organizarse de muchas formas y recomienda mantener el ejecutable en `cmd` y el código privado en `internal`. [Guía oficial de organización de Go](https://go.dev/doc/modules/layout).

El repositorio llamado `golang-standards/project-layout` tampoco es un estándar oficial; su propia documentación lo aclara. No conviene copiar una plantilla enorme pensando que es obligatoria.

Tiago propone precisamente tres responsabilidades en [Theory: Clean Architecture & Repository Pattern](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/24. Theory Clean Architecture & Repository Pattern.md>):

1. **Transport:** recibe HTTP, decodifica JSON y responde.
2. **Service:** ejecuta las reglas de negocio.
3. **Storage/Repository:** guarda y consulta datos.

La dependencia sigue esta dirección:

```text
Astro
  ↓ petición HTTP
Chi + handlers
  ↓
Servicios/casos de uso
  ↓ puerto
Interfaces de repositorio, IA y documentos
  ↑ implementación
PostgreSQL, proveedor de IA y generador PDF
```

Eso ya contiene la idea principal de la arquitectura hexagonal: la lógica central no conoce los detalles externos.

### Estructura actual y crecimiento por funcionalidad

```text
navigate-business-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/app.go
│   ├── config/config.go
│   ├── httpapi/router.go
│   └── auth/
│       ├── account.go
│       ├── service.go
│       ├── module/module.go
│       ├── http/
│       ├── firebase/
│       └── postgres/
├── migrations/
├── go.mod
└── go.sum
```

Cada parte tiene una responsabilidad:

- `cmd/api/main.go`: arranca el servidor; `internal/app/app.go` construye y conecta las piezas compartidas.
- `config`: valida variables antes de crear servicios.
- `httpapi`: Chi y rutas generales; `auth/http` contiene el contrato HTTP de auth.
- `auth`: cuenta, servicio y contrato de repositorio; `auth/module` compone sus piezas.
- `auth/postgres`: implementación SQL del contrato de auth.
- `migrations`: creación y evolución de tablas.

Cuando se implemente `project`, `form` o `analysis`, cada uno tendrá su propio paquete, módulo de composición y solo los adaptadores que realmente utilice. `app.New` compartirá las instancias necesarias entre ellos.

### Ejemplo del flujo

Cuando Astro crea un formulario:

```text
POST /api/v1/projects/{projectID}/forms
              ↓
        form_handler
              ↓
      form.Service.Create
              ↓
     FormRepository.Create
              ↓
          PostgreSQL
```

El handler sabe de HTTP:

```go
func (h *FormHandler) Create(
    w http.ResponseWriter,
    r *http.Request,
) {
    // Leer JSON.
    // Llamar al servicio.
    // Responder JSON.
}
```

El servicio sabe de las reglas:

```go
type FormService struct {
    repository FormRepository
}

func (s *FormService) Create(
    ctx context.Context,
    input CreateFormInput,
) (Form, error) {
    // Comprobar permisos.
    // Validar reglas.
    // Guardar el formulario.
}
```

El repositorio define lo que necesita el servicio:

```go
type FormRepository interface {
    Create(context.Context, Form) (Form, error)
}
```

PostgreSQL implementa ese contrato. El servicio no necesita saber qué driver, tablas o consultas SQL estamos usando.

### Lo que evitaremos

No crearemos desde el comienzo carpetas genéricas como:

```text
controllers/
services/
repositories/
models/
utils/
helpers/
common/
```

Al crecer, esas carpetas terminan mezclando usuarios, proyectos, formularios y análisis. Organizaremos primero por funcionalidad: `project`, `form`, `analysis` y `auth`.

Tampoco crearemos una interfaz para cada estructura. Usaremos interfaces cuando haya una frontera real:

- Base de datos.
- Proveedor de IA.
- Generador de PDF.
- Almacenamiento de archivos.
- Servicio de correo.

Chi no debe entrar en `service.go` ni en los modelos. Solo debe existir en la capa HTTP.

### Decisión definitiva

Usaremos:

- **Chi** para rutas y middleware.
- **Monolito modular**, no microservicios.
- **Capas Transport → Service → Repository**.
- **Principios hexagonales en integraciones externas**.
- **Paquetes organizados por funcionalidad**.
- **Inyección manual en `main.go`**, sin framework de dependencias.

El proyecto usa Chi `v5.3.2` y CORS `v1.2.2` en [go.mod](go.mod). El router los importa directamente para este primer incremento.

Esta estructura es profesional, compatible con lo que estudias con Tiago y suficientemente sencilla para que puedas comprender cada pieza mientras construimos.


Sí, **usaremos Chi desde el principio**.

No porque Go lo necesite, sino porque ya decidimos cómo crecerá tu API y Chi nos evitará cambiar posteriormente la organización de rutas y middleware.

Empezaremos de forma mínima:

```go
func main() {
    router := chi.NewRouter()

    router.Use(middleware.RequestID)
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)

    router.Get("/health", healthHandler)

    http.ListenAndServe(":8080", router)
}
```

Al principio Chi solamente hará esto:

```text
Petición → Chi → handler → respuesta
```

Después nos permitirá crecer ordenadamente:

```text
/api/v1/auth       → rutas públicas
/api/v1/projects   → requieren autenticación
/api/v1/forms      → requieren permisos
```

La regla será que **Chi solo vive en la capa HTTP**. Los servicios de proyectos, formularios y análisis no sabrán que Chi existe. Recibirán datos normales de Go y devolverán resultados o errores.

Esto nos permite usar Chi desde ahora sin depender excesivamente de él. Si algún día lo reemplazamos, solo cambia la capa HTTP; las reglas del proyecto permanecen iguales.

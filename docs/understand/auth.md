Sí: **no conviene leerlos en el orden del explorador ni empezar por `main.go` intentando comprender cada dependencia**. Para aprender la programación del proyecto, recomiendo hacer varias pasadas, comenzando por el núcleo y tratando Firebase, PostgreSQL, Chi y Cloud SQL como cajas negras con contratos públicos.
## Mapa mental del proyecto

El backend tiene dos recorridos diferentes:

### 1. Arranque de la aplicación

```text
cmd/api/main.go
    ↓
config.Load()
    ↓
app.New()
    ↓
crear Firebase + PostgreSQL
    ↓
auth/module.New()
    ↓
httpapi.NewRouter()
    ↓
servidor HTTP
```

### 2. Procesamiento de una petición autenticada

```text
router
  ↓
middleware RequireGoogle
  ↓
TokenVerifier
  ↓
handler
  ↓
auth.Service
  ↓
auth.Repository
  ↓
PostgreSQL
```

El segundo recorrido es el más importante para entender primero.

---

# Orden de lectura recomendado

## Escalón 0: entender el objetivo, sin entrar todavía al código

Lee únicamente estas partes del plan:

1. `docs/plans/01-autenticacion-firebase-y-cuenta-interna.md`
   - “Objetivo y alcance exacto”
   - “Un solo dueño de cada responsabilidad”
   - “Flujo que deberá quedar implementado”
   - “Diseño Go propuesto, archivo por archivo”
   - “Contrato HTTP y errores”

No necesitas leer todavía todas las fuentes, enlaces oficiales o documentación de Firebase.

La idea que debes llevarte es:

> Firebase demuestra quién es la persona; PostgreSQL conserva la cuenta interna del negocio.

Son identificadores diferentes:

```text
Firebase UID                  Account ID
identidad externa             identidad interna del negocio
"abc-firebase-123"     →      "9e62...-uuid"
```

---

## Escalón 1: tipos y contratos del núcleo

Estos son los archivos más sencillos y no dependen de Firebase, PostgreSQL ni Chi.

### 1. `internal/auth/account.go`

Empieza aquí.

Observa:

- `Identity`: representa una identidad externa.
- `Account`: representa una cuenta interna.
- `ErrNotFound` y `ErrSuspended`: errores propios del negocio.
- `ProviderFirebase`: nombre del proveedor.

Pregunta que responde:

> ¿Cuáles son los datos fundamentales de autenticación?

No investigues todavía cómo funcionan los tags JSON internamente. Basta saber que:

```go
ID string `json:"account_id"`
```

hace que el campo salga como `account_id` en JSON.

---

### 2. `internal/auth/verifier.go`

Aquí aparece el contrato para comprobar un token:

```go
type TokenVerifier interface {
    Verify(context.Context, string) (string, error)
}
```

Léelo como una promesa:

> “Dame un token y te devolveré el UID verificado o un error”.

No importa todavía si quien lo implementa es Firebase, un doble de pruebas u otro proveedor.

También contiene errores neutrales:

- `ErrInvalidIdentityToken`
- `ErrIdentityProviderUnavailable`

Este archivo es especialmente importante porque evita que el resto del proyecto dependa directamente de Firebase.

---

### 3. `internal/auth/service.go`

Este es el núcleo del comportamiento.

Concéntrate en:

```go
type Repository interface {
    Sync(context.Context, Identity) (Account, error)
    Get(context.Context, Identity) (Account, error)
}
```

y en:

```go
Service.Sync(...)
Service.Current(...)
```

El servicio:

1. recibe un UID ya verificado;
2. lo transforma en una `Identity`;
3. pide información al repositorio;
4. rechaza cuentas suspendidas.

Por ahora trata `Repository` como una caja negra:

```text
Service → Repository → devuelve Account
```

No abras aún `postgres/account_repository.go`.

### Resultado del primer escalón

Después de estos tres archivos deberías poder explicar:

```text
TokenVerifier verifica identidad.
Repository almacena o recupera cuentas.
Service aplica las reglas de auth.
```

---

# Escalón 2: seguir una petición HTTP

Ahora puedes seguir una petición desde la ruta hasta el servicio.

## 4. `internal/auth/http/middleware.go`

Estudia solamente este recorrido:

```go
func (a *Authenticator) RequireGoogle(next http.Handler) http.Handler
```

El middleware:

1. busca `Authorization`;
2. exige el formato `Bearer <token>`;
3. llama a `a.verifier.Verify(...)`;
4. obtiene el UID;
5. coloca el UID en el `context`;
6. deja continuar la petición.

La línea conceptual más importante es:

```go
next.ServeHTTP(
    w,
    r.WithContext(context.WithValue(r.Context(), uidKey{}, uid)),
)
```

No necesitas estudiar las entrañas de `context.Context`. Inicialmente basta con entender:

> El middleware guarda el UID verificado dentro de la petición para que el siguiente componente pueda leerlo.

Después, esta función lo recupera:

```go
VerifiedUID(r.Context())
```

También observa cómo traduce errores:

```text
token inválido             → 401
proveedor no disponible    → 503
```

---

## 5. `internal/auth/http/handler.go`

Ahora estudia los dos endpoints:

```go
Handler.Sync
Handler.Current
```

### `Sync`

```text
lee UID del contexto
    ↓
service.Sync(...)
    ↓
elige 200 o 201
    ↓
escribe JSON
```

### `Current`

```text
lee UID del contexto
    ↓
service.Current(...)
    ↓
escribe 200 + JSON
```

El handler no verifica Firebase ni escribe SQL. Solo adapta:

```text
HTTP → llamada Go → HTTP
```

Observa también la traducción de errores:

```text
auth.ErrNotFound     → 404 account_not_synced
auth.ErrSuspended    → 403 account_suspended
otro error           → 503 storage_unavailable
```

---

## 6. `internal/httpapi/router.go`

Ahora ya puedes comprender las rutas:

```go
router.Route("/api/v1/auth", func(r chi.Router) {
    r.Use(authenticator.RequireGoogle)
    r.Post("/sync", authHandler.Sync)
    r.Get("/me", authHandler.Current)
})
```

Tradúcelo mentalmente así:

```text
Todas las rutas /api/v1/auth/*
    primero pasan por RequireGoogle.

POST /sync
    después ejecuta Handler.Sync.

GET /me
    después ejecuta Handler.Current.
```

No hace falta entrar al código fuente interno de Chi. Solo debes aprender su API pública:

- `Use`: instala middleware.
- `Post`: registra una ruta POST.
- `Get`: registra una ruta GET.
- `Route`: agrupa rutas bajo un prefijo.

### Comprobación mental

En este punto deberías poder narrar:

```text
POST /api/v1/auth/sync
→ router
→ RequireGoogle
→ TokenVerifier.Verify
→ UID en el contexto
→ Handler.Sync
→ Service.Sync
→ Repository.Sync
→ respuesta JSON
```

---

# Escalón 3: entender cómo se conectan los objetos

Hasta aquí viste el comportamiento. Ahora estudia cómo se construyen sus piezas.

## 7. `internal/auth/module/module.go`

Este archivo es pequeño, pero arquitectónicamente importante:

```go
func New(repository auth.Repository, verifier auth.TokenVerifier) *Module
```

Recibe dos contratos:

- un repositorio;
- un verificador.

Y construye:

```text
Repository ──→ Service ──→ Handler
Verifier   ──────────────→ Authenticator
```

El módulo no crea Firebase ni PostgreSQL. Recibe las implementaciones ya preparadas.

Ese patrón se llama **inyección explícita de dependencias**:

```go
service := auth.NewService(repository)
handler := authhttp.NewHandler(service)
authenticator := authhttp.NewAuthenticator(verifier)
```

No necesitas un framework de inyección: son llamadas normales a constructores de Go.

---

## 8. `internal/config/config.go`

Ahora estudia cómo entra la configuración.

`config.Load()`:

1. intenta cargar `.env`;
2. lee variables con `os.Getenv`;
3. comprueba las obligatorias;
4. establece `127.0.0.1:8080` como dirección predeterminada.

No investigues internamente `godotenv.Load()`. Trátalo así:

```text
godotenv.Load() → carga variables de .env en el entorno
```

La salida importante es:

```go
type Config struct {
    FirebaseProjectID
    InstanceConnectionName
    DatabaseName
    DatabaseUser
    DatabasePassword
    HTTPAddr
}
```

Es decir, el resto de la aplicación recibe una estructura validada en lugar de leer variables de entorno repetidamente.

---

## 9. `internal/app/app.go`

Este es el archivo más denso. No intentes dominar cada llamada externa en la primera lectura.

Divídelo en cuatro bloques.

### Bloque A: crear el verificador

```go
verifier, err := firebaseauth.New(ctx, cfg.FirebaseProjectID)
```

Traducción:

> Crear la implementación real de `auth.TokenVerifier`.

### Bloque B: crear la conexión SQL

Desde:

```go
pgxpool.ParseConfig(...)
```

hasta:

```go
pool.Ping(ctx)
```

Trátalo inicialmente como:

> Preparar y verificar un pool de conexiones PostgreSQL usando Cloud SQL.

No necesitas entrar todavía en:

- `DialFunc`;
- sockets;
- TLS;
- ADC;
- refresh de certificados;
- funcionamiento interno del pool.

### Bloque C: componer auth

```go
auth := authmodule.New(
    authpostgres.NewAccountRepository(pool),
    verifier,
)
```

Esta línea conecta las implementaciones reales con los contratos:

```text
*AccountRepository implementa auth.Repository
*firebaseauth.Verifier implementa auth.TokenVerifier
```

### Bloque D: crear el router

```go
Handler: httpapi.NewRouter(auth.Handler, auth.Authenticator)
```

El resultado final de `app.New()` es una aplicación que expone un `http.Handler` listo para servir peticiones.

---

## 10. `cmd/api/main.go`

Lee `main.go` después de `app.go`, no antes.

Entonces resulta sencillo:

```go
config.Load()
    ↓
app.New(ctx, cfg)
    ↓
http.Server{Handler: application.Handler}
    ↓
ListenAndServe()
```

`main.go` no contiene reglas de autenticación. Su responsabilidad es:

- iniciar;
- configurar tiempos del servidor;
- ejecutar;
- cerrar recursos;
- registrar errores.

La separación entre:

```go
func main()
func run() error
```

permite que `run` devuelva errores y que `main` sea el único lugar que termina el proceso con `log.Fatal`.

---

# Escalón 4: abrir las cajas negras externas

Solo después de comprender los escalones anteriores debes entrar en Firebase y PostgreSQL.

## 11. `migrations/000001_auth.up.sql`

Léelo **antes** del repositorio PostgreSQL.

Primero comprende las tablas:

```text
accounts
├── id
├── status
├── created_at
└── updated_at

auth_identities
├── id
├── account_id ──────────→ accounts.id
├── provider
├── provider_subject
└── created_at
```

La restricción fundamental es:

```sql
UNIQUE (provider, provider_subject)
```

Esto garantiza que una identidad Firebase no pueda terminar asociada a varias cuentas debido a peticiones concurrentes.

Después lee:

- `migrations/000001_auth.down.sql`

Este archivo solamente revierte la migración en el orden correcto:

```text
primero auth_identities
después accounts
```

---

## 12. `internal/auth/postgres/account_repository.go`

Ahora las consultas tendrán contexto.

### Primero lee `Get`

```go
func (r *AccountRepository) Get(...)
```

Comprende:

```text
Identity
  ↓
JOIN auth_identities + accounts
  ↓
Account
```

No necesitas entrar al código interno de `pgx.QueryRow`. Basta saber:

```go
QueryRow(...).Scan(...)
```

ejecuta una consulta de una fila y copia sus columnas a variables Go.

### Después lee `Sync`

Divídelo en fases:

1. intenta encontrar una cuenta existente;
2. abre una transacción;
3. crea una cuenta;
4. intenta asociar la identidad;
5. resuelve una posible competencia;
6. confirma la transacción;
7. marca `Created = true`.

El caso más complicado está aquí:

```go
ON CONFLICT (provider, provider_subject) DO NOTHING
```

Si dos peticiones intentan crear la misma identidad:

```text
Petición A crea identidad → gana
Petición B encuentra conflicto
Petición B revierte su cuenta provisional
Petición B lee la cuenta creada por A
```

Este archivo implementa `auth.Repository`, aunque Go no exige escribir explícitamente:

```go
implements Repository
```

Lo satisface porque tiene los métodos:

```go
Sync(context.Context, auth.Identity) (auth.Account, error)
Get(context.Context, auth.Identity) (auth.Account, error)
```

---

## 13. `internal/auth/firebase/verifier.go`

Déjalo casi para el final.

Primero observa esta comprobación:

```go
var _ auth.TokenVerifier = (*Verifier)(nil)
```

Significa:

> El compilador debe confirmar que `*Verifier` cumple el contrato `auth.TokenVerifier`.

Luego divide el archivo:

### `New`

```go
func New(ctx context.Context, projectID string) (*Verifier, error)
```

Crea el cliente Firebase.

### `Verify`

```go
func (v *Verifier) Verify(...)
```

Le entrega el token al SDK y traduce sus errores a errores neutrales del proyecto:

```text
error concreto de Firebase
    ↓
auth.ErrInvalidIdentityToken
o
auth.ErrIdentityProviderUnavailable
```

### `googleUID`

Comprueba:

- que exista token;
- que tenga UID;
- que el proveedor sea `google.com`.

No necesitas estudiar internamente `VerifyIDTokenAndCheckRevoked`. El contrato público suficiente es:

```text
token válido y no revocado → token decodificado
token inválido             → error
```

---

# Escalón 5: usar las pruebas como ejemplos

Cuando comprendas producción, lee las pruebas en este orden:

1. `internal/auth/http/handler_test.go`
2. `internal/auth/firebase/verifier_test.go`
3. `internal/config/config_test.go`

## Por qué empezar por `handler_test.go`

Probablemente sea el mejor archivo para entender el comportamiento completo sin servicios reales.

Las pruebas reemplazan:

- Firebase por un verificador falso;
- PostgreSQL por un repositorio falso.

Así puedes observar el contrato sin entrar en las dependencias externas:

```text
token válido + cuenta nueva        → 201
token válido + cuenta existente    → 200
sin token                          → 401
cuenta suspendida                  → 403
Firebase o almacenamiento fallan   → 503
```

Las pruebas aquí no son solo validación: funcionan como ejemplos ejecutables del diseño.

---

# Orden corto definitivo

Si quieres una sola lista práctica, sigue esta:

```text
1. docs/plans/01-autenticacion-firebase-y-cuenta-interna.md
   solo objetivo, flujo, responsabilidades y contrato HTTP

2. internal/auth/account.go
3. internal/auth/verifier.go
4. internal/auth/service.go

5. internal/auth/http/middleware.go
6. internal/auth/http/handler.go
7. internal/httpapi/router.go

8. internal/auth/module/module.go
9. internal/config/config.go
10. internal/app/app.go
11. cmd/api/main.go

12. migrations/000001_auth.up.sql
13. internal/auth/postgres/account_repository.go
14. migrations/000001_auth.down.sql

15. internal/auth/firebase/verifier.go

16. internal/auth/http/handler_test.go
17. internal/auth/firebase/verifier_test.go
18. internal/config/config_test.go
```

---

# Qué puedes ignorar inicialmente

No entres todavía en el código fuente interno de:

- `chi`;
- `pgx`;
- `pgxpool`;
- Firebase Admin SDK;
- Cloud SQL Connector;
- `godotenv`;
- `net/http`;
- ADC de Google;
- `golang-migrate`.

Aprende únicamente el contrato que este proyecto usa:

| Dependencia | Lo mínimo que necesitas saber |
|---|---|
| Chi | Registra rutas y middleware. |
| Firebase Admin | Verifica un token y devuelve identidad. |
| `pgxpool` | Proporciona conexiones PostgreSQL. |
| `pgx` | Ejecuta consultas y transacciones. |
| Cloud SQL Connector | Permite que el pool alcance Cloud SQL usando ADC. |
| `godotenv` | Carga variables desde `.env`. |
| `net/http` | Recibe peticiones y escribe respuestas. |

Solo investiga el interior de una dependencia cuando una línea concreta deje de poder explicarse mediante su contrato público.

---

# Orden histórico según la bitácora

Si después quieres entender **por qué el proyecto terminó con esta estructura**, lee:

1. `bitacora-construccion/auth/01-google-firebase-y-cuenta-sql.md`
2. `bitacora-construccion/auth/03-modulo-auth-y-migraciones-estandar.md`
3. `bitacora-construccion/config/01-composicion-modular-y-dependencias-compartidas.md`
4. `bitacora-construccion/config/02-conexion-cloud-sql-go-adc.md`
5. `bitacora-construccion/auth/04-contrato-neutral-verificacion-identidad.md`

`auth/02-legibilidad-y-responsabilidades.md` es útil como historia intermedia, pero describe algunas ubicaciones y un ejecutor de migraciones que posteriormente fueron reemplazados. Por eso no lo pondría al comienzo del aprendizaje del código actual.

La evolución fue:

```text
auth funcional
  ↓
organización por funcionalidad
  ↓
módulo auth y migraciones estándar
  ↓
composición central en app.New
  ↓
conexión real a Cloud SQL
  ↓
contrato de verificación independiente de Firebase
```

## La idea central que deberías aprender

No memorices primero carpetas. Memoriza responsabilidades:

```text
account.go       = datos y errores del dominio
verifier.go      = contrato de identidad
service.go       = reglas de negocio
middleware.go    = autenticación de la petición
handler.go       = traducción HTTP
router.go        = rutas
module.go        = construcción de auth
app.go           = composición de toda la aplicación
main.go          = arranque del proceso
postgres/        = implementación de persistencia
firebase/        = implementación de verificación
migrations/      = estructura de la base de datos
```

Este orden coincide con la filosofía registrada en el proyecto: **comprender primero los contratos y el mecanismo propio; abrir después las implementaciones externas**. No modifiqué archivos ni ejecuté pruebas porque la solicitud era de lectura y orientación.

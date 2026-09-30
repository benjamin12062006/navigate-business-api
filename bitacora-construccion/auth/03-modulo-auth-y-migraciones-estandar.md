# 03 — Módulo auth y migraciones SQL estándar

Fecha de creación: 2026-09-29  
Estado: Implementado y verificado en pruebas locales; integración con Cloud SQL pendiente  
Ámbito: `navigate-business-api`, organización de auth y ejecución de migraciones  
Revisión de partida: backend sin repositorio Git antes de esta edición  
Relacionada con: [entrada 01](01-google-firebase-y-cuenta-sql.md), [entrada 02](02-legibilidad-y-responsabilidades.md), [plan de auth](../../docs/plans/01-autenticacion-firebase-y-cuenta-interna.md)

## Necesidad

Preparar el backend para agregar otras funcionalidades sin mezclar sus archivos HTTP y SQL con auth. Sustituir el ejecutor SQL propio por la herramienta de migraciones estudiada, sin cambiar el contrato de autenticación ni tocar el frontend.

## Evidencia y decisión

La lección 25 de Tiago separa transporte, servicio y almacenamiento y aconseja iniciar con una estructura simple; su reconstrucción guiada muestra la dirección de dependencias.[^F1][^R1] La lección 33 introduce `golang-migrate`, pares `up`/`down`, la tabla de versión y un Makefile para recordar comandos; no exige escribir un motor de migraciones propio.[^F2] La CLI oficial admite PostgreSQL, `-path`, `-database`, `up` y `down`.[^O1] Su documentación para PostgreSQL aconseja `BEGIN`/`COMMIT` para migraciones con varias sentencias.[^O2]

Como adaptación propia, `internal/auth` conserva el modelo y servicio; sus subpaquetes `http`, `firebase` y `postgres` alojan los adaptadores de esta funcionalidad. `internal/httpapi/router.go` queda como composición general. La migración SQL sigue en `migrations/`, ejecutada por `golang-migrate` mediante un Makefile. El flujo Firebase Admin y las consultas SQL de cuentas mantienen su comportamiento.

## Implementación y comprobación

`cmd/api/main.go` construye el adaptador Firebase, el repositorio PostgreSQL y el handler HTTP desde los subpaquetes de auth, y pasa este último al router general. El servicio y su interfaz `Repository` permanecen en el núcleo, sin importar Chi, Firebase Admin ni `pgx`. [Handler](../../internal/auth/http/handler.go), [middleware](../../internal/auth/http/middleware.go), [verificador](../../internal/auth/firebase/verifier.go) y [repositorio](../../internal/auth/postgres/account_repository.go) cambiaron de paquete, no de algoritmo. Las pruebas HTTP se movieron a `authhttp_test` para llamar al router público sin ciclo de importación.

Se eliminaron `cmd/migrate/main.go` y `runner.go`. El [Makefile](../../Makefile) invoca la CLI `migrate` con `DATABASE_URL`; `migrate-up` aplica pendientes y `migrate-down` revierte **una** versión. Los [SQL up](../../migrations/000001_auth.up.sql) y [down](../../migrations/000001_auth.down.sql) contienen ahora `BEGIN`/`COMMIT`. Se fijó `v4.20.1` para la instalación de la herramienta, sin agregarla como dependencia del ejecutable API.[^O1][^O2] El [README](../../README.md) y el [plan](../../docs/plans/01-autenticacion-firebase-y-cuenta-interna.md) explican el flujo actualizado.

| Comprobación ejecutada | Resultado |
| --- | --- |
| `go test ./...` | Pasó en los paquetes API, auth y router tras la reorganización. Las pruebas de proveedor y HTTP existentes conservan los casos de Google, errores y sincronización. |
| `go vet ./...` | Pasó sin diagnósticos. |
| Instalación de `migrate` v4.20.1 con tag `postgres` y `migrate -help` | Compiló; la CLI mostró drivers `file`, `postgres` y comandos `up`/`down`. Se instaló en un directorio temporal de verificación. |
| `DATABASE_URL=stub://auth-make-check PATH=… make migrate-up` | La CLI leyó `000001_auth.up.sql` y reportó `1/u auth`. El driver `stub` no ejecuta SQL real. |
| `make migrate-up` sin `DATABASE_URL` | Falló con el mensaje previsto antes de llamar a la CLI. |

No hay credenciales ni base SQL configuradas para aplicar `up`/`down` reales ni probar el servidor contra Firebase. El ejecutor anterior y `golang-migrate` usan formatos distintos de `schema_migrations`; una base donde ya se ejecutó `cmd/migrate` necesita una conciliación de versión antes de usar la nueva herramienta. En este entorno no consta una migración aplicada. Tampoco se modificó Expo.

## Fuentes

[^F1]: [25. Clean Layered Architecture](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture.md), Tiago Taquelim, *APIs de Producción con Go*, sección 5, líneas 13–25 y 39–43; consulta 2026-09-29.
[^R1]: [Reconstrucción guiada de la lección 25](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture%20-%20Reconstrucción%20guiada.md), fases 1–3; evidencia derivada, consulta 2026-09-29.
[^F2]: [33. SQL Migrations](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/33. SQL Migrations.md>), Tiago Taquelim, *APIs de Producción con Go*, sección 6, líneas 13–23 y 37–57; consulta 2026-09-29.
[^O1]: [migrate CLI](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md), instalación y uso; consulta 2026-09-29.
[^O2]: [migrate PostgreSQL driver](https://github.com/golang-migrate/migrate/blob/master/database/postgres/README.md), sección «Multi-statement mode»; consulta 2026-09-29.

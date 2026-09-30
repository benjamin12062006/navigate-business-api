# 02 — Legibilidad y separación de responsabilidades de auth

Fecha de creación: 2026-09-29  
Estado: Implementado, pendiente de verificación con SQL real  
Ámbito: `navigate-business-api`, estructura de la primera etapa de auth  
Revisión de partida: esta carpeta no tiene repositorio Git propio  
Relacionada con: [entrada 01](01-google-firebase-y-cuenta-sql.md), [arquitectura basada en Tiago](../../docs/fonts/go-backend/arquitectura/01-tiago-arquitectura-y-repositorios.md), [plan de auth](../../docs/plans/01-autenticacion-firebase-y-cuenta-interna.md)

## Necesidad y punto de partida

Revisar si la cantidad de archivos y líneas del backend de auth facilita aprender Go y respeta responsabilidades claras. Identificar una separación que haga el flujo más fácil de seguir sin crear carpetas por cada función.

## Evidencia y decisión provisional

La [lección 25, «Clean Layered Architecture»](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture.md), líneas 13–25 y 39–43, distingue transporte, reglas y almacenamiento y recomienda empezar simple. Su [reconstrucción guiada](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture%20-%20Reconstrucción%20guiada.md), fases 1–3, ilustra la dirección de dependencias. La [lección 29, «The Repository Pattern»](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%206%20Databases/29.%20The%20Repository%20Pattern.md), líneas 13–17, respalda separar SQL de las reglas. El curso no fija un número de archivos o un máximo de líneas: ese criterio debe salir del código real.

El flujo ya separaba `cmd/api` (composición), `httpapi` (HTTP), `firebaseauth` (proveedor), `auth` (reglas) y `postgres` (SQL). Antes del ajuste, los archivos de producción tenían entre 21 y 83 líneas. La única mezcla clara de motivos de cambio estaba en `cmd/migrate/main.go`: argumentos y mensajes de CLI junto a conexión, lectura de SQL y transacción. Separarlos dentro del **mismo paquete Go** es una decisión del proyecto, no una instrucción literal del curso.

## Cómo leer la estructura actual

En Go, los archivos de una misma carpeta con `package httpapi` o `package main` comparten paquete. **Un archivo ayuda a navegar; un paquete define la frontera de dependencias.** Por eso abrir una carpeta nueva por cada función añadiría importaciones y saltos de lectura sin mejorar la separación ya existente.

| Paquete y archivo | Responsabilidad y motivo de cambio | Decisión |
| --- | --- | --- |
| [`cmd/api/main.go`](../../cmd/api/main.go) | Leer configuración, construir Firebase y pool SQL, inyectarlos y arrancar HTTP. | Mantener como punto de composición; 64 líneas y un recorrido lineal. |
| [`internal/httpapi/router.go`](../../internal/httpapi/router.go), `auth_middleware.go`, `auth_handler.go` | Rutas, comprobación Bearer/UID en contexto y contrato HTTP. | En esta revisión estaban separados por tareas dentro del mismo transporte. La [entrada 03](03-modulo-auth-y-migraciones-estandar.md) los agrupó por funcionalidad. |
| `internal/firebaseauth/verifier.go` | Detalle externo del Admin SDK y regla de token Google. | En esta revisión estaba separado de HTTP y del servicio; ahora está en [`internal/auth/firebase/verifier.go`](../../internal/auth/firebase/verifier.go). |
| [`internal/auth/account.go`](../../internal/auth/account.go), [`service.go`](../../internal/auth/service.go) | Cuenta, errores, puerto `Repository` y reglas de sincronización/suspensión. | Mantener juntos como núcleo de auth; no conocen Chi, Firebase SDK ni `pgx`. |
| `internal/postgres/account_repository.go` | Leer y sincronizar cuentas mediante SQL y transacción. | `Get` y `Sync` comparten razón de cambio; ahora está en [`internal/auth/postgres/account_repository.go`](../../internal/auth/postgres/account_repository.go). |
| `cmd/migrate/main.go`, `runner.go` | Entrada CLI frente a conexión y ejecución SQL transaccional. | Separación aplicada en esta revisión: `main.go` pasó de 83 a 30 líneas; `runner.go` tenía 71. Fueron sustituidos en la [entrada 03](03-modulo-auth-y-migraciones-estandar.md). |

El recorrido de una petición es `router → middleware → handler → auth.Service → auth.Repository → postgres`. El middleware usa `firebaseauth` para obtener el UID verificado; `cmd/api` conecta las implementaciones. La interfaz `Repository` existe porque el servicio necesita un contrato pequeño y el test HTTP puede sustituir PostgreSQL. No se agregó una interfaz a cada estructura.

## Cambio aplicado y razón

Antes de esta revisión, `cmd/migrate/main.go` llamaba `log.Fatal` desde dentro del trabajo SQL. `log.Fatal` termina el proceso y no ejecuta los `defer` de cierre y rollback. En esta revisión `main.go` pasó a validar argumentos/entorno y presentar el resultado, mientras `runMigration` devolvía errores y conservaba la conexión y la transacción bajo sus `defer`. El comando `go run ./cmd/migrate up|down` estuvo vigente hasta la [entrada 03](03-modulo-auth-y-migraciones-estandar.md), donde se adoptó `golang-migrate`.

No se movieron `auth`, `httpapi`, `firebaseauth` ni `postgres`: sus límites ya siguen la separación de la fuente. El número de líneas se usó para localizar una posible dificultad de lectura, no como regla automática de Clean Code. Si un paquete empieza a reunir operaciones de tienda, perfiles o pagos, ese nuevo motivo de cambio justificaría otro paquete; todavía no existe.

## Comprobación

| Ejecutado | Resultado |
| --- | --- |
| `gofmt -w cmd/migrate/main.go cmd/migrate/runner.go` | Formato aplicado. |
| `go test ./...` | Pasó; compila ambos archivos del comando y pasan las pruebas existentes de auth. |
| `go vet ./...` | Pasó sin diagnósticos. |
| `go run ./cmd/migrate` sin argumento | Devuelve el mensaje de uso esperado y código 1. |

No se ejecutó `up` ni `down` contra PostgreSQL porque no hay `DATABASE_URL` ni instancia de desarrollo configurada en este entorno. Por eso esta entrada verifica la separación y compilación, pero no afirma que la migración real se haya aplicado.

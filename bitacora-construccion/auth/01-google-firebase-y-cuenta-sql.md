# 01 — Google Sign-In verificado en Go y cuenta SQL

Fecha de creación: 2026-09-29  
Estado: Implementado, pendiente de verificación con Firebase y SQL Connect reales  
Ámbito: `navigate-business-api`  
Revisión de partida: esta carpeta no tiene repositorio Git propio  
Relacionada con: [plan de auth](../../docs/plans/01-autenticacion-firebase-y-cuenta-interna.md), [fuentes oficiales](../../docs/fonts/firebase/01-google-auth-admin-y-sql-connect.md)

## Necesidad y punto de partida

Implementar únicamente el backend Go de auth. El inicio será exclusivamente con Google mediante Firebase Authentication; Expo no se modifica en este incremento. Al comenzar existía Chi sin rutas auth, base SQL ni migraciones.

## Evidencia y decisión

Tiago separa la comprobación de identidad, el servicio y el repositorio en sus lecciones de [validación de tokens](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%2014%20Authentication/67.%20Validating%20Tokens.md) (líneas 13–15), [repositorios](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%206%20Databases/29.%20The%20Repository%20Pattern.md) (líneas 13–17) y [transacciones](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%2012%20User%20Creation/58.%20SQL%20Transactions.md) (líneas 13–17). Las [reconstrucciones de validación](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%2014%20Authentication/67.%20Validating%20Tokens%20-%20Reconstrucción%20guiada.md), [repositorio](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%206%20Databases/29.%20The%20Repository%20Pattern%20-%20Reconstrucción%20guiada.md) y [transacción](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%2012%20User%20Creation/58.%20SQL%20Transactions%20-%20Reconstrucción%20guiada.md) ayudan a recuperar los mecanismos. Sus ejemplos usan JWT y usuarios propios; la adaptación reemplaza ese emisor por [Firebase Admin Go](https://firebase.google.com/docs/auth/admin/verify-id-tokens) y persiste la cuenta con SQL directo, como permite la [CLI de SQL Connect](https://firebase.google.com/docs/sql-connect/cli-reference). Las fuentes oficiales y sus límites están catalogados en la [ficha Firebase](../../docs/fonts/firebase/01-google-auth-admin-y-sql-connect.md).

La decisión del proyecto es aceptar solo tokens de Firebase cuyo `firebase.sign_in_provider` sea `google.com`. El UID verificado identifica una identidad; el UUID SQL identifica la cuenta de negocio. No se usa el correo como clave.

## Implementación y funcionamiento

El [adaptador Firebase](../../internal/auth/firebase/verifier.go) llama `VerifyIDTokenAndCheckRevoked`, toma el UID **solo** del token verificado y exige `Firebase.SignInProvider == "google.com"`. Los errores inválidos se distinguen de indisponibilidad del servicio. El [middleware HTTP](../../internal/auth/http/middleware.go) requiere `Authorization: Bearer <ID token de Firebase>` y deja el UID verificado en el contexto; los [handlers](../../internal/auth/http/handler.go) no aceptan UID, correo ni contraseña en el cuerpo. [Chi](../../internal/httpapi/router.go) publica `POST /api/v1/auth/sync` y `GET /api/v1/auth/me`.

El [servicio](../../internal/auth/service.go) usa una interfaz de repositorio estrecha, sin saber de Firebase, Chi ni `pgx`. [PostgreSQL](../../internal/auth/postgres/account_repository.go) almacena `accounts.id` UUID separado del UID Firebase. `Sync` primero busca la identidad; si no existe, inserta cuenta e identidad en una transacción. Una restricción `UNIQUE(provider, provider_subject)` en la [migración](../../migrations/000001_auth.up.sql) resuelve la primera petición concurrente: la transacción perdedora revierte su cuenta provisional y lee la ganadora. `Current` solo lee y no crea cuenta. Ambos rechazan `suspended`.

El [servidor](../../cmd/api/main.go) lee proyecto Firebase, ADC y `DATABASE_URL`, fija un pool pequeño y no migra al iniciar. En esta etapa existió un ejecutor independiente en `cmd/migrate` que aplicaba o revertía la migración con bloqueo consultivo y registro de versión; fue sustituido por `golang-migrate` en la [entrada 03](03-modulo-auth-y-migraciones-estandar.md). El [README](../../README.md) explica el orden y contrato HTTP actuales. No se modificó Expo.

## Comprobación ejecutada

| Comprobación | Resultado observado |
| --- | --- |
| `go test ./...` | Pasó. [Prueba de proveedor](../../internal/auth/firebase/verifier_test.go): Google aceptado; contraseña, anónimo, UID vacío y token ausente rechazados. [Prueba HTTP](../../internal/auth/http/handler_test.go): primera sincronización `201`, repetición `200`, lectura, suspensión `403`, falta de token `401`, fallo de Firebase o SQL `503`. Usa dobles; no prueba Firebase ni SQL reales. |
| `go vet ./...` | Pasó sin diagnósticos. |
| Inspección de entorno | No hay `FIREBASE_PROJECT_ID`, `DATABASE_URL` ni credenciales configuradas. No hay `psql` local y el daemon Docker no está disponible. |

**Pendiente antes de afirmar funcionamiento integrado:** crear/vincular proyecto Firebase con solo Google habilitado; configurar instancia PostgreSQL de SQL Connect con migraciones propias, grants IAM/SQL y proxy/conector; ejecutar `up`, repetición y `down` en desarrollo; probar dos solicitudes reales con el mismo ID token, concurrencia, expiración/revocación y suspensión. No se creó ningún recurso de pago ni se enviaron datos a Firebase.

## Continuidad

El siguiente escalón será conectar esta API a una instancia de desarrollo y ejecutar el recorrido real. El botón de Google en Expo se implementará únicamente cuando se aborde el frontend; esta entrada permanece limitada al backend.

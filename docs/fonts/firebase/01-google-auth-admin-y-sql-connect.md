# 01 — Google Sign-In, verificación Admin Go y SQL directo de Firebase

Fecha de creación: 2026-09-29  
Última actualización: 2026-09-29  
Tipo de fuente: documentación oficial y referencia de paquetes Go  
Origen: Firebase, Google Cloud y Go Packages  
Fecha de consulta: 2026-09-29  
Relacionada con: [plan de auth](../../plans/01-autenticacion-firebase-y-cuenta-interna.md), [bitácora](../../../bitacora-construccion/auth/01-google-firebase-y-cuenta-sql.md)

## Preguntas investigadas

¿Cómo aceptar exclusivamente Google Sign-In sin gestionar contraseñas? ¿Qué token debe recibir Go? ¿Cómo comprobar revocación y proveedor? ¿Permite Firebase SQL Connect ejecutar SQL y migraciones propias?

## Fuentes y evidencia

| ID | Localización oficial consultada | Evidencia pertinente |
| --- | --- | --- |
| O1 | [Authenticate Using Google with JavaScript](https://firebase.google.com/docs/auth/web/google-signin), «Before you begin» y «Handle the sign-in flow» | Se habilita Google como proveedor de Firebase Authentication. El cliente usa Google para iniciar sesión en Firebase. Si obtiene antes un ID token de Google, lo intercambia con `GoogleAuthProvider.credential` y `signInWithCredential`. Esto ocurre en el cliente; aquí no se implementa Expo. |
| O2 | [Verify ID Tokens](https://firebase.google.com/docs/auth/admin/verify-id-tokens), «Retrieve ID tokens on clients» y «Verify ID tokens using the Firebase Admin SDK», ejemplo Go | Tras iniciar sesión, el cliente obtiene un **ID token de Firebase**, lo envía por HTTPS al servidor y Go llama `VerifyIDToken` para obtener un UID verificado. La verificación valida firma, vencimiento y proyecto. Un token OAuth de Google no sustituye al ID token de Firebase. |
| O3 | [Manage User Sessions](https://firebase.google.com/docs/auth/admin/manage-sessions), «Detect ID token revocation», ejemplo Go; [referencia Admin Go](https://pkg.go.dev/firebase.google.com/go/v4/auth#Client.VerifyIDTokenAndCheckRevoked) | `VerifyIDTokenAndCheckRevoked` también rechaza tokens revocados o usuarios deshabilitados y realiza una llamada remota adicional por petición. |
| O4 | [Referencia Admin Go: `FirebaseInfo`](https://pkg.go.dev/firebase.google.com/go/v4/auth#FirebaseInfo) y [Token](https://pkg.go.dev/firebase.google.com/go/v4/auth#Token) | El token decodificado expone `UID` y `Firebase.SignInProvider`. El backend puede exigir `google.com`, incluso si el proyecto tuviera otros proveedores habilitados por error. |
| O5 | [Add the Firebase Admin SDK to your server](https://firebase.google.com/docs/admin/setup), «Initialize the SDK» y «Test with gcloud end user credentials» | Admin Go soporta ADC y `firebase.Config{ProjectID: ...}`. El login ADC de usuario final predeterminado de `gcloud` requiere configuración adicional para Firebase Authentication; en producción conviene una identidad de servicio sin clave exportada cuando el entorno Google lo permita. |
| O6 | [Firebase CLI command reference for SQL Connect](https://firebase.google.com/docs/sql-connect/cli-reference), `dataconnect:sql:setup` y `dataconnect:sql:grant` | Al rechazar que SQL Connect gestione migraciones en `sql:setup`, el equipo conserva propiedad del esquema. La CLI permite conceder roles SQL para conexión directa; además hace falta IAM `roles/cloudsql.client`. |
| O7 | [Cloud SQL: IAM authentication](https://docs.cloud.google.com/sql/docs/postgres/iam-authentication), «Automatic versus manual IAM database authentication»; [Cloud SQL Auth Proxy](https://docs.cloud.google.com/sql/docs/postgres/connect-auth-proxy), «Before you begin»; [Cloud SQL desde GKE](https://docs.cloud.google.com/sql/docs/postgres/connect-kubernetes-engine), instrucciones de IAM automática | La autenticación IAM automática requiere un conector o proxy y `roles/cloudsql.instanceUser` además de `roles/cloudsql.client`. Go recibe una cadena PostgreSQL para el extremo creado por ese mecanismo. El rol SQL y los roles IAM son permisos distintos. |
| O8 | [pgxpool reference](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#ParseConfig), `ParseConfig` | `pgxpool` permite parametrizar conexiones y limitar el pool; el backend fija un máximo inicial de cinco conexiones, sujeto a medición posterior. |
| O9 | [golang-migrate CLI](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md), instalación y uso; [PostgreSQL driver](https://github.com/golang-migrate/migrate/blob/master/database/postgres/README.md), «Multi-statement mode» | La CLI se instala con el tag `postgres`, lee pares `up`/`down` por `-path`, recibe una URL PostgreSQL en `-database` y mantiene la versión aplicada. En migraciones de varias sentencias se usan `BEGIN`/`COMMIT` explícitos. Versión revisada del proyecto: `v4.20.1`. |

## Interpretación y adaptación

El flujo `Google → Firebase Authentication → ID token Firebase → Go` es la única ruta de inicio aceptada. El backend usa el UID del token verificado como sujeto de una identidad interna y nunca recibe contraseña ni decide la cuenta por correo. La exigencia de `google.com` es una **regla propia de este producto**, implementada sobre el campo documentado por Admin Go. También se debe deshabilitar correo/contraseña en la consola Firebase; la validación Go aporta otra barrera, pero no cambia la configuración del proveedor.

SQL Connect proporciona Cloud SQL PostgreSQL. El proyecto usa consultas SQL con `pgx` y archivos `migrations/*.sql`, sin operaciones GraphQL. La instancia, los grants y el conector/proxy requieren datos reales del proyecto y aún no se han configurado ni probado. El ejecutor de migraciones usa una identidad de owner distinta de la identidad writer de la API; esta división es una decisión de implementación basada en los roles descritos por O6.

La [lección 33 de Tiago](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/33. SQL Migrations.md>) propone `golang-migrate` y comandos Make para la migración; O9 valida la CLI vigente. El curso también aporta la separación middleware/servicio/repositorio y la transacción, pero no enseña esta integración Firebase. Ver [ficha local de autenticación](../go-backend/autenticacion-autorizacion/01-tiago-autenticacion-y-autorizacion.md), [arquitectura](../go-backend/arquitectura/01-tiago-arquitectura-y-repositorios.md) y [persistencia](../go-backend/persistencia-seguridad/01-transacciones-y-seguridad-de-datos.md).

# 04 — Guía reproducible: Firebase SQL Connect, Cloud SQL y API Go

Fecha de creación: 2026-09-30
Tipo de fuente: documentación oficial y procedimiento verificado en el proyecto
Origen: Firebase, Google Cloud, Cloud SQL Go Connector y ejecución local de `navigate-business-api`
Fecha de consulta: 2026-09-30
Relacionada con: [fuente del conector](03-conector-cloud-sql-go-adc.md), [investigación inicial](01-google-auth-admin-y-sql-connect.md), [plan de autenticación](../../plans/01-autenticacion-firebase-y-cuenta-interna.md) y [bitácora de implementación](../../../bitacora-construccion/config/02-conexion-cloud-sql-go-adc.md)
Revisión de referencia: `2cfc3ec` con cambios locales posteriores aún sin commit

## Objetivo y resultado comprobado

Esta guía reconstruye el procedimiento correcto para que una API Go se conecte directamente al PostgreSQL de Cloud SQL creado detrás de Firebase SQL Connect. No usa el SDK web de Firebase, no abre una conexión PostgreSQL pública manual, no requiere ejecutar Cloud SQL Auth Proxy junto a la API y no conserva una clave JSON dentro del repositorio.

El recorrido quedó comprobado el 2026-09-30 cuando la API ejecutó `pool.Ping` y arrancó con:

```text
Iniciando API en http://127.0.0.1:8080
```

La solución final combina dos autenticaciones que no deben confundirse:

```text
API Go
  ├─ ADC de Google + roles/cloudsql.client
  │    autoriza el canal hacia Cloud SQL
  └─ DB_USER + DB_PASSWORD
       autentica la sesión dentro de PostgreSQL
```

## Valores actuales no secretos

| Dato | Valor comprobado | Dónde se obtiene |
| --- | --- | --- |
| Proyecto Firebase/Google Cloud | `expo-city-bol` | Selector de proyecto de Firebase o Google Cloud |
| Instancia Cloud SQL | `expo-city-bol-instance` | Cloud SQL → Instancias |
| Nombre de conexión | `expo-city-bol:us-east4:expo-city-bol-instance` | Instancia → Descripción general → Nombre de conexión |
| Base PostgreSQL | `expo-city-bol-database` | Instancia → Bases de datos |
| Usuario PostgreSQL de la API | `catalogo_api` | Instancia → Usuarios |
| Cuenta de servicio para producción | `catalogo-api@expo-city-bol.iam.gserviceaccount.com` | IAM → Cuentas de servicio |

La contraseña, los tokens ADC y cualquier clave privada son secretos y no pertenecen a este documento.

## 1. Crear el proyecto y la base

1. Abrir [Firebase Console](https://console.firebase.google.com/) y crear o seleccionar el proyecto.
2. Registrar el **ID del proyecto**; no confundirlo con el nombre visible. En este proyecto es `expo-city-bol`.
3. Abrir **Bases de datos y almacenamiento → SQL Connect** y crear el servicio siguiendo el flujo de Firebase.
4. Elegir la región deliberadamente y aprovisionar o vincular una instancia de Cloud SQL para PostgreSQL. SQL Connect y la instancia de Cloud SQL son recursos relacionados, pero facturables por separado.[^F1]
5. Para una API Go que administra sus propias tablas mediante migraciones SQL, usar Cloud SQL como PostgreSQL directo. El SDK generado de SQL Connect no es necesario para este backend.

La documentación de SQL Connect explica que el servicio se apoya en Cloud SQL para PostgreSQL y que el esquema, los conectores y los SDK generados son capas adicionales.[^F1] Este proyecto conserva sus propias migraciones para las tablas de la API; no se debe permitir que dos sistemas distintos administren las mismas tablas.

## 2. Obtener los identificadores correctos

Abrir [Google Cloud Console → Cloud SQL → Instancias](https://console.cloud.google.com/sql/instances), comprobar que el selector superior apunta a `expo-city-bol` y entrar en la instancia.

Recuperar:

- **Nombre de conexión:** sección «Conectarse a esta instancia» de la descripción general. Formato `proyecto:región:instancia`.
- **Nombre de base:** menú «Bases de datos».
- **IP pública:** sirve para comprobar la ruta disponible, pero no se guarda en la configuración de esta API porque el conector resuelve la instancia por su nombre.

Comandos de verificación usados:

```bash
gcloud sql instances list --project=expo-city-bol
```

```bash
gcloud sql instances describe expo-city-bol-instance \
  --project=expo-city-bol \
  --format="value(connectionName)"
```

El segundo comando debe devolver:

```text
expo-city-bol:us-east4:expo-city-bol-instance
```

## 3. Crear el usuario PostgreSQL

La cuenta de servicio IAM no crea un usuario PostgreSQL automáticamente. Se necesita un usuario adicional dentro de la instancia:

1. Cloud SQL → `expo-city-bol-instance` → **Usuarios**.
2. Pulsar **Agregar cuenta de usuario**.
3. Seleccionar **Autenticación integrada**.
4. Escribir `catalogo_api`.
5. Crear o pegar una contraseña fuerte y guardarla en un gestor de secretos.
6. En el primer recorrido se conservó el rol predeterminado mostrado por Cloud SQL. Antes de producción se deben reducir los privilegios a las tablas y operaciones que realmente use la API.

No usar como identidad de runtime:

- `postgres`, porque es el administrador general;
- `firebasesuperuser`, porque pertenece a Firebase;
- `service-…@gcp-sa-firebasedataconnect.iam`, porque es la cuenta administrada de Data Connect.

La existencia del usuario se comprueba sin mostrar contraseñas:

```bash
gcloud sql users list \
  --instance=expo-city-bol-instance \
  --project=expo-city-bol \
  --format="table(name,type)"
```

Debe aparecer `catalogo_api` como `BUILT_IN`.

## 4. Crear la identidad de ejecución de la API

Abrir [IAM y administración → Cuentas de servicio](https://console.cloud.google.com/iam-admin/serviceaccounts):

1. Crear la cuenta `catalogo-api`.
2. Conceder solamente **Cliente de Cloud SQL** (`roles/cloudsql.client`) para abrir conexiones.[^F2]
3. No conceder **Administrador de Cloud SQL**: permite administrar recursos y excede la necesidad de runtime.
4. No crear una clave JSON para Cloud Run. En producción se asignará la cuenta directamente al servicio.

La cuenta `firebase-adminsdk-…` ya existente corresponde a Firebase Admin SDK; no se reutilizó como identidad de la API. La cuenta `catalogo-api` se reservará para Cloud Run. En desarrollo local se usa la identidad personal configurada mediante ADC.

## 5. Instalar y configurar Google Cloud CLI localmente

`gcloud` es una herramienta del sistema; no se instala con `pnpm`, Bun, Cargo ni `go get`. Se puede instalar con el [instalador oficial de Google Cloud CLI](https://docs.cloud.google.com/sdk/docs/install-sdk) o, opcionalmente en macOS, con el cask `gcloud-cli` de Homebrew.[^F3]

Después de instalar:

```bash
gcloud init
```

```bash
gcloud config set project expo-city-bol
```

```bash
gcloud auth application-default login
```

El último comando crea **Application Default Credentials para desarrollo local**. Se ejecuta una vez y el conector las descubre automáticamente.[^F2] Verificar el proyecto activo:

```bash
gcloud config get-value project
```

Resultado esperado:

```text
expo-city-bol
```

`gcloud auth application-default login` no abre por sí mismo una conexión SQL. Solo proporciona la identidad de Google que después consume Cloud SQL Go Connector.

## 6. Configurar `.env` sin una clave JSON

Copiar [`.env.example`](../../../.env.example) como `.env` y completar:

```dotenv
APP_ENV=development
HTTP_ADDR=127.0.0.1:8080

FIREBASE_PROJECT_ID=expo-city-bol
FIREBASE_AUTH_EMULATOR_HOST=

INSTANCE_CONNECTION_NAME=expo-city-bol:us-east4:expo-city-bol-instance
DB_NAME=expo-city-bol-database
DB_USER=catalogo_api
DB_PASSWORD='contraseña-real-no-documentada'
```

Reglas importantes:

- `.env` está ignorado por Git.
- No definir `GOOGLE_APPLICATION_CREDENTIALS` cuando se desea usar el ADC creado por `gcloud`. Esa variable tiene prioridad y podría desviar la aplicación hacia una clave JSON distinta.
- Mantener `DB_PASSWORD` entre comillas simples si contiene `$`, `#`, espacios u otros caracteres que un parser `.env` pueda interpretar.
- Si la terminal tenía una contraseña exportada anteriormente, ejecutar `unset DB_PASSWORD` antes de arrancar, porque las variables del proceso tienen prioridad sobre `godotenv.Load`.
- Nunca colocar estas variables en Expo ni prefijarlas con `EXPO_PUBLIC_`; son secretos exclusivos del backend.

## 7. Instalar las dependencias Go

Dependencias utilizadas y verificadas:

```bash
go get cloud.google.com/go/cloudsqlconn@v1.25.2
```

```bash
go get github.com/jackc/pgx/v5
```

Después se normaliza el módulo:

```bash
go mod tidy
```

Cloud SQL Go Connector usa ADC, obtiene certificados efímeros y crea un canal cifrado TLS hacia la instancia.[^F2] `pgxpool` administra el conjunto de sesiones PostgreSQL sobre ese canal.

## 8. Integrar el conector con `pgxpool`

La implementación real está en [`internal/app/app.go`](../../../internal/app/app.go). El mecanismo esencial es:

```go
poolConfig, err := pgxpool.ParseConfig("host=localhost sslmode=disable")
if err != nil {
    return nil, err
}

poolConfig.ConnConfig.User = cfg.DatabaseUser
poolConfig.ConnConfig.Password = cfg.DatabasePassword
poolConfig.ConnConfig.Database = cfg.DatabaseName

dialer, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithLazyRefresh())
if err != nil {
    return nil, err
}

poolConfig.ConnConfig.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
    return dialer.Dial(ctx, cfg.InstanceConnectionName, cloudsqlconn.WithPublicIP())
}

pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
if err != nil {
    _ = dialer.Close()
    return nil, err
}

if err := pool.Ping(ctx); err != nil {
    pool.Close()
    _ = dialer.Close()
    return nil, err
}
```

Detalles que evitan errores futuros:

- `localhost` es un valor de configuración para `pgx`; `DialFunc` reemplaza la marcación y dirige realmente cada conexión a Cloud SQL. Por eso algunos errores de `pgx` todavía imprimen `127.0.0.1:5432` aunque el conector sí haya alcanzado la instancia.
- `sslmode=disable` no deja el tráfico sin protección: evita una segunda capa PostgreSQL TLS porque el conector ya establece y valida el canal cifrado.[^F2]
- `WithPublicIP` corresponde a la ruta actual de la instancia. Para IP privada se necesita conectividad VPC y se cambia a `WithPrivateIP`.
- `WithLazyRefresh` evita depender de un refresco periódico cuando la CPU puede quedar suspendida, como en Cloud Run.[^F2]
- El pool y el `Dialer` deben cerrarse durante el apagado.

[`internal/config/config.go`](../../../internal/config/config.go) valida las cinco variables obligatorias y ya no exige `DATABASE_URL` para la API.

## 9. Verificar en orden

Primero comprobar el código:

```bash
go test ./...
```

```bash
go vet ./...
```

Después arrancar en local:

```bash
unset DB_PASSWORD
go tool air
```

Air recompila `./cmd/api`. La verificación definitiva de conexión es que `pool.Ping` termine y aparezca:

```text
Iniciando API en http://127.0.0.1:8080
```

Ese mensaje demuestra que se completaron:

1. carga de `.env`;
2. inicialización de Firebase Admin con ADC;
3. autorización IAM para consultar Cloud SQL;
4. resolución de la instancia;
5. canal cifrado del conector;
6. autenticación de `catalogo_api` dentro de PostgreSQL;
7. `Ping` del pool.

No demuestra todavía que las tablas existan. Las migraciones deben aplicarse antes de probar endpoints que consulten `accounts` u otras tablas.

## 10. Errores encontrados y corrección correcta

### «FIREBASE_PROJECT_ID y DATABASE_URL son obligatorios»

Era el contrato antiguo del código. La instancia y `gcloud` podían estar correctos, pero la API terminaba antes de conectarse. Se sustituyó `DATABASE_URL` en el runtime por las variables `INSTANCE_CONNECTION_NAME` y `DB_*`.

### `lookup oauth2.googleapis.com: no such host`

Se observó únicamente dentro de un sandbox sin red. No era un error de credenciales ni de Cloud SQL. La ejecución con acceso de red avanzó hasta PostgreSQL.

### `SQLSTATE 28P01: password authentication failed`

Este error confirmó que el conector había alcanzado PostgreSQL, pero la sesión SQL no fue aceptada. Hubo dos causas sucesivas:

1. `catalogo_api` todavía no existía; la cuenta de servicio IAM con nombre parecido no lo creaba automáticamente.
2. La contraseña comenzaba con `$` y estaba sin comillas en `.env`; `godotenv` podía expandirla como variable.

La solución verificada fue crear `catalogo_api` como usuario integrado, usar exactamente la misma contraseña y escribirla entre comillas simples. Después se reinició Air y la API arrancó correctamente.

### La cuenta de servicio y el usuario SQL parecen duplicados

No lo son:

| Identidad | Función |
| --- | --- |
| Usuario local que ejecutó `gcloud auth application-default login` | ADC durante desarrollo local |
| `catalogo-api@…gserviceaccount.com` | Identidad futura de Cloud Run |
| `catalogo_api` | Usuario y contraseña dentro de PostgreSQL |
| Cuenta interna `service-…firebasedataconnect.iam` | Operación administrada por Firebase Data Connect |

## 11. Migraciones: diferencia con la conexión de la API

La API ya no necesita `DATABASE_URL`. Sin embargo, la herramienta externa `golang-migrate` declarada en el [`Makefile`](../../../Makefile) sí consume una URL estándar PostgreSQL. Para ejecutar migraciones contra Cloud SQL se debe levantar temporalmente Cloud SQL Auth Proxy o adoptar más adelante un ejecutor de migraciones que use el mismo conector Go.

No usar el usuario limitado de runtime como propietario del esquema. Mantener separadas:

- identidad propietaria para crear o alterar tablas;
- identidad `catalogo_api` para operaciones normales de la aplicación.

Tampoco permitir que Firebase SQL Connect y `golang-migrate` administren simultáneamente las mismas tablas.

## 12. Trasladar la misma solución a Cloud Run

Cuando la API se despliegue:

1. Cloud Run → servicio → editar y desplegar nueva revisión.
2. En **Seguridad**, seleccionar `catalogo-api@expo-city-bol.iam.gserviceaccount.com`.
3. Conservar `roles/cloudsql.client` en esa cuenta.
4. Configurar `FIREBASE_PROJECT_ID`, `INSTANCE_CONNECTION_NAME`, `DB_NAME` y `DB_USER` como variables.
5. Guardar `DB_PASSWORD` en Secret Manager y montarlo como secreto.
6. No ejecutar `gcloud auth application-default login` dentro de Cloud Run y no subir un JSON: la identidad asignada proporciona ADC automáticamente.[^F4]

## Procedimiento resumido

```text
Crear Firebase/SQL Connect
  → localizar instancia y base en Cloud SQL
  → crear usuario PostgreSQL catalogo_api
  → crear cuenta de servicio catalogo-api con Cloud SQL Client
  → instalar gcloud y crear ADC local
  → configurar .env sin GOOGLE_APPLICATION_CREDENTIALS
  → instalar cloudsqlconn + pgx
  → configurar DialFunc y cerrar pool/dialer
  → ejecutar tests
  → arrancar y exigir que Ping llegue a «Iniciando API»
  → aplicar migraciones antes de probar repositorios
```

## Evidencia, adaptación y límites

- **Evidencia oficial:** ADC, requisitos IAM, canal cifrado y patrón `DialFunc` provienen de Google Cloud.[^F2]
- **Evidencia observada:** los comandos devolvieron la instancia correcta; el conector alcanzó PostgreSQL; el fallo `28P01` desapareció al crear el usuario y corregir el valor `.env`; la API arrancó a las 18:21:39 del 2026-09-30.
- **Adaptación:** este proyecto usa usuario PostgreSQL integrado y contraseña. No se habilitó IAM Database Authentication.
- **Límite:** el arranque prueba la conexión, no las migraciones ni el flujo completo con un ID token real de Firebase.

## Fuentes

[^F1]: [Firebase SQL Connect](https://firebase.google.com/docs/sql-connect), secciones «How does it work?» y capacidades principales; consulta 2026-09-30. [Precios de Firebase SQL Connect](https://firebase.google.com/docs/sql-connect/pricing?hl=es-419), componentes facturables; consulta 2026-09-30.
[^F2]: [Cloud SQL: Connect using Cloud SQL Language Connectors](https://docs.cloud.google.com/sql/docs/postgres/connect-connectors), secciones «Authenticate» y «Go»; consulta 2026-09-30. [Repositorio oficial Cloud SQL Go Connector](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector), secciones «Credentials» y «Using the dialer with pgx»; consulta 2026-09-30. [Ejemplo oficial para Go](https://docs.cloud.google.com/sql/docs/postgres/samples/cloud-sql-postgres-databasesql-connect-connector); consulta 2026-09-30.
[^F3]: [Instalar Google Cloud CLI](https://docs.cloud.google.com/sdk/docs/install-sdk) y [Google Cloud CLI mediante Homebrew](https://docs.cloud.google.com/sdk/docs/downloads-homebrew?hl=es-419); consulta 2026-09-30.
[^F4]: [Cloud SQL: Connect from Cloud Run](https://docs.cloud.google.com/sql/docs/postgres/connect-run), identidad de servicio y conexión desde Cloud Run; consulta 2026-09-30.

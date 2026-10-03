# 03 — Cloud SQL Go Connector con Application Default Credentials

Fecha de creación: 2026-09-30
Tipo de fuente: documentación oficial
Origen: Google Cloud y repositorio oficial de Cloud SQL Go Connector
Fecha de consulta: 2026-09-30
Relacionada con: [plan de autenticación](../../plans/01-autenticacion-firebase-y-cuenta-interna.md), [composición modular](../../../bitacora-construccion/config/01-composicion-modular-y-dependencias-compartidas.md) y [conexión implementada](../../../bitacora-construccion/config/02-conexion-cloud-sql-go-adc.md)

## Necesidad de consulta

Conectar localmente la API Go con la instancia PostgreSQL de Cloud SQL asociada a Firebase SQL Connect sin exponer la IP como conexión PostgreSQL directa, sin ejecutar un proxy separado y sin conservar una clave de cuenta de servicio en el repositorio.

## Localización verificable

- [Conectarse mediante Cloud SQL Language Connectors](https://docs.cloud.google.com/sql/docs/postgres/connect-connectors), secciones «Authenticate» y «Go».
- [Ejemplo oficial: Create a connection using the Cloud SQL Go Connector](https://docs.cloud.google.com/sql/docs/postgres/samples/cloud-sql-postgres-databasesql-connect-connector), ejemplo Go.
- [Repositorio oficial cloud-sql-go-connector](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector), secciones «Credentials» y «Using the dialer with pgx».

## Evidencia

El conector oficial usa Application Default Credentials cuando no se le entrega otra fuente de credenciales. Para desarrollo local, Google documenta `gcloud auth application-default login`. La identidad debe tener `roles/cloudsql.client` y el proyecto debe tener habilitada Cloud SQL Admin API.

Con `pgxpool`, el ejemplo oficial crea un `cloudsqlconn.Dialer`, asigna una función a `Config.ConnConfig.DialFunc` y dirige cada conexión al nombre `proyecto:región:instancia`. El conector establece el canal cifrado; por eso el protocolo PostgreSQL se configura con `sslmode=disable` y no duplica TLS dentro del túnel. El `Dialer` debe cerrarse junto con el pool.

`WithLazyRefresh` está recomendado para entornos con CPU limitada fuera de las solicitudes, como Cloud Run. La autenticación IAM de base de datos es una capacidad distinta y opcional; no corresponde a un usuario PostgreSQL integrado con contraseña.

## Interpretación y límites

ADC autentica la identidad frente a Google Cloud para abrir el canal hacia la instancia. No sustituye el usuario y contraseña PostgreSQL configurados como `DB_USER` y `DB_PASSWORD`. Esta fuente tampoco concede permisos SQL sobre tablas: esos privilegios deben existir en PostgreSQL.

## Aplicación en el proyecto

La API utiliza ADC tanto en local como en un futuro despliegue con identidad de servicio. `internal/app.App` es dueño del pool y del `Dialer`; `internal/config.Config` recibe por separado el nombre de instancia y las credenciales SQL. El acceso se fija explícitamente a la IP pública porque la instancia actual publica esa ruta; si se migra a IP privada, deberá cambiarse la opción de marcado y existir conectividad VPC.

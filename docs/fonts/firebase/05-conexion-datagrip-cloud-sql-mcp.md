# 05 — Conectar DataGrip a Cloud SQL mediante MCP y Auth Proxy

Fecha de creación: 2026-09-30
Tipo de fuente: procedimiento ejecutado y verificado en el proyecto
Origen: DataGrip 2026.2, MCP integrado de JetBrains, Cloud SQL Auth Proxy y `navigate-business-api`
Relacionada con: [guía de Cloud SQL para la API Go](04-guia-cloud-sql-go-desde-cero.md) y [conector Go con ADC](03-conector-cloud-sql-go-adc.md)

## Resultado comprobado

DataGrip quedó conectado a la misma base PostgreSQL de Cloud SQL que usa la API:

| Dato | Valor verificado |
| --- | --- |
| Fuente de datos en DataGrip | `Cloud SQL - expo-city-bol-instance` |
| DBMS | PostgreSQL 18.6 |
| Driver | PostgreSQL JDBC 42.7.3 |
| URL JDBC | `jdbc:postgresql://127.0.0.1:5432/expo-city-bol-database` |
| Cloud SQL Auth Proxy | 2.26.0 |
| Prueba final | Sin problemas, ping de 23 ms |

No se escribió el usuario ni la contraseña dentro de la URL JDBC. DataGrip conserva esos valores locales mediante su almacenamiento protegido (`master_key`).

## Las tres piezas que intervienen

La configuración MCP, el proxy y la fuente de datos cumplen funciones diferentes:

```text
Codex
  │ MCP HTTP
  ▼
DataGrip
  │ JDBC: 127.0.0.1:5432
  ▼
Cloud SQL Auth Proxy
  │ canal autenticado y cifrado mediante ADC
  ▼
Cloud SQL PostgreSQL
```

1. **MCP de DataGrip:** permite que Codex solicite a DataGrip crear, listar o probar fuentes de datos. No conecta por sí solo con PostgreSQL.
2. **Cloud SQL Auth Proxy:** abre un puerto TCP local y establece el canal autenticado hacia la instancia usando Application Default Credentials (ADC).
3. **Fuente de datos JDBC:** es la conexión PostgreSQL que DataGrip guarda y muestra en Database Explorer.

La API Go no necesita este proxy: usa `cloudsqlconn` dentro del propio proceso. DataGrip sí necesita un extremo JDBC/TCP, por eso se añadió el proxy local.

## 1. Configuración MCP utilizada

El MCP integrado de DataGrip ya estaba disponible mediante esta configuración:

```json
{
  "type": "streamable-http",
  "url": "http://127.0.0.1:64402/stream",
  "headers": {
    "IJ_MCP_SERVER_PROJECT_PATH": "/Users/benji/Documents/EXPO-GO/navigate-business-api"
  }
}
```

El encabezado `IJ_MCP_SERVER_PROJECT_PATH` dirige las operaciones al proyecto correcto. El puerto del servidor MCP puede cambiar si DataGrip regenera su configuración; en ese caso se debe copiar nuevamente la configuración ofrecida por el IDE en vez de asumir que `64402` será permanente.

Antes de crear nada se ejecutó la operación MCP `list_database_connections`. La lista estaba vacía, por lo que se creó una fuente nueva sin sobrescribir otra conexión.

## 2. Herramienta instalada

Se instaló **Cloud SQL Auth Proxy** con Homebrew:

```bash
brew install cloud-sql-proxy
```

La instalación comprobada fue:

```text
cloud-sql-proxy version 2.26.0
```

Esto instala un programa del sistema; no es una dependencia de Go ni modifica `go.mod`.

## 3. Iniciar el proxy local

Con ADC ya configurado mediante `gcloud auth application-default login`, se inició:

```bash
cloud-sql-proxy \
  --address 127.0.0.1 \
  --port 5432 \
  expo-city-bol:us-east4:expo-city-bol-instance
```

Resultado esperado:

```text
Authorizing with Application Default Credentials
[expo-city-bol:us-east4:expo-city-bol-instance] Listening on 127.0.0.1:5432
The proxy has started successfully and is ready for new connections!
```

Se usa `127.0.0.1`, no `0.0.0.0`, para que el puerto solo sea accesible desde esta Mac. El proxy debe permanecer ejecutándose mientras DataGrip use la conexión.

La instalación de Homebrew permanece después de reiniciar la computadora, pero el proceso del proxy no arranca automáticamente. Después de un reinicio se debe volver a ejecutar el comando anterior.

## 4. Fuente PostgreSQL creada en DataGrip

Mediante `create_database_connection` del MCP se enviaron únicamente estos datos:

```text
Nombre: Cloud SQL - expo-city-bol-instance
DBMS: PostgreSQL
URL: jdbc:postgresql://127.0.0.1:5432/expo-city-bol-database
Comprobar al crear: sí
```

La misma configuración puede reconstruirse manualmente en DataGrip:

1. Abrir **Data Sources and Drivers**.
2. Elegir **PostgreSQL**.
3. Usar host `127.0.0.1`, puerto `5432` y base `expo-city-bol-database`.
4. Escribir el valor de `DB_USER` en **User**.
5. Escribir el valor de `DB_PASSWORD` en **Password** y permitir que DataGrip lo guarde en su almacén protegido.
6. Pulsar **Test Connection**.

No usar una URL con el formato `jdbc:postgresql://usuario:contraseña@...` ni copiar la contraseña a este documento.

## 5. Verificación ejecutada

Después de crear la fuente se ejecutó `test_database_connection` con su identificador interno. DataGrip devolvió:

```text
DBMS: PostgreSQL (ver. 18.6)
Driver: PostgreSQL JDBC Driver (ver. 42.7.3, JDBC4.2)
Ping: 23 ms
Problemas: no
```

DataGrip mostró `SSL: no` para la sesión JDBC local. Eso describe el tramo `DataGrip → 127.0.0.1`; no significa que el tráfico hacia Cloud SQL viaje sin protección. El proxy establece el canal autenticado y cifrado hacia la instancia.

La prueba de conexión no ejecutó consultas SQL ni modificó tablas.

## 6. Qué queda guardado y qué debe volver a iniciarse

| Elemento | ¿Persiste? | Acción futura |
| --- | --- | --- |
| Binario `cloud-sql-proxy` | Sí | Reinstalar solo si se elimina o se cambia de equipo |
| Fuente JDBC de DataGrip | Sí, dentro de la configuración del proyecto | Abrir normalmente desde Database Explorer |
| Usuario y contraseña | Sí, localmente en el almacén protegido de DataGrip | No documentarlos ni versionarlos |
| Proceso del proxy | No después de cerrarlo o reiniciar | Ejecutar otra vez el comando de la sección 3 |
| Servidor MCP de DataGrip | Depende de que DataGrip esté abierto y MCP habilitado | Regenerar/copiar su configuración si cambia el puerto |

## Diagnóstico rápido

### DataGrip indica `Connection refused` en `127.0.0.1:5432`

El proxy no está activo o el puerto cambió. Iniciarlo nuevamente y confirmar el mensaje `Listening on 127.0.0.1:5432`.

### El proxy no puede autorizarse

Renovar ADC:

```bash
gcloud auth application-default login
```

La identidad debe conservar `roles/cloudsql.client` sobre el proyecto o la instancia.

### PostgreSQL responde `password authentication failed`

El canal hacia Cloud SQL funciona, pero el usuario o la contraseña SQL no coinciden. Revisar los campos protegidos **User** y **Password** de DataGrip contra `DB_USER` y `DB_PASSWORD`; no cambiar la URL JDBC para insertar credenciales.

### El MCP no encuentra el proyecto o las conexiones

Confirmar que DataGrip tiene abierto este proyecto y que `IJ_MCP_SERVER_PROJECT_PATH` coincide exactamente con:

```text
/Users/benji/Documents/EXPO-GO/navigate-business-api
```

## Procedimiento resumido para repetirlo

```text
Abrir DataGrip y habilitar su MCP
  → instalar cloud-sql-proxy una sola vez
  → iniciar el proxy en 127.0.0.1:5432 usando ADC
  → crear una fuente PostgreSQL con la URL JDBC local
  → guardar usuario y contraseña solo en los campos protegidos de DataGrip
  → probar la conexión
```

## Fuentes oficiales

- [Google Cloud: conectar mediante Cloud SQL Auth Proxy](https://docs.cloud.google.com/sql/docs/postgres/connect-auth-proxy)
- [Google Cloud: descripción e instalación de Cloud SQL Auth Proxy](https://docs.cloud.google.com/sql/docs/postgres/sql-proxy)
- [JetBrains DataGrip: usar agentes de IA y MCP](https://www.jetbrains.com/help/datagrip/use-ai-agents-with-datagrip.html)
- [JetBrains DataGrip: conexiones y fuentes de datos](https://www.jetbrains.com/help/datagrip/connecting-to-a-database.html)

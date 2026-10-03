# 02 — Conexión directa a Cloud SQL mediante el conector Go

Fecha de creación: 2026-09-30
Estado: Implementado y verificado
Ámbito: `navigate-business-api`, configuración y composición de PostgreSQL
Relacionada con: [composición modular](01-composicion-modular-y-dependencias-compartidas.md), [fuente oficial](../../docs/fonts/firebase/03-conector-cloud-sql-go-adc.md) y [plan de autenticación](../../docs/plans/01-autenticacion-firebase-y-cuenta-interna.md)
Revisión de partida: `2cfc3ec` con cambios locales previos preservados

## Necesidad y punto de partida

La instancia `expo-city-bol:us-east4:expo-city-bol-instance` ya existe y ADC fue configurado localmente con `gcloud auth application-default login`. Sin embargo, la API exigía `DATABASE_URL` y terminaba antes de intentar la conexión; las variables separadas de Cloud SQL no eran consumidas.

## Evidencia e interpretación

El conector oficial para Go resuelve ADC, abre un canal cifrado hacia el nombre de conexión de la instancia y se integra con `pgxpool` mediante `ConnConfig.DialFunc`.[^F1]

ADC autoriza el acceso a Cloud SQL, mientras `DB_USER` y `DB_PASSWORD` autentican la sesión PostgreSQL. Como `catalogo_api` es un usuario integrado, no se habilita IAM Database Authentication.

## Alternativas y decisión

El proxy local habría permitido conservar `DATABASE_URL`, pero exige otro proceso. La conexión directa a la IP pública obligaría a administrar red autorizada y TLS manualmente. Se elige el conector Go porque funciona dentro del proceso, usa ADC y cifra la conexión sin un proxy separado.[^F1]

## Implementación y funcionamiento

[`config.Load`](../../internal/config/config.go) sustituyó el contrato `DATABASE_URL` por `INSTANCE_CONNECTION_NAME`, `DB_NAME`, `DB_USER` y `DB_PASSWORD`. [`.env.example`](../../.env.example) conserva el contrato sin secretos y `.gitignore` excluye claves locales de Firebase Admin SDK.

[`app.New`](../../internal/app/app.go) crea el `cloudsqlconn.Dialer` con `WithLazyRefresh`, configura `pgxpool.Config.ConnConfig.DialFunc` para marcar la instancia por su nombre de conexión y comprueba el pool con `Ping`. El pool usa el usuario PostgreSQL integrado y el conector usa ADC. [`App.Close`](../../internal/app/app.go) cierra ambos recursos y [`run`](../../cmd/api/main.go) registra un eventual error de cierre.

## Comprobación

| Comprobación ejecutada | Resultado observado |
| --- | --- |
| `GOCACHE=/tmp/navigate-business-api-go-cache go test ./...` | Pasó en todos los paquetes. |
| `GOCACHE=/tmp/navigate-business-api-go-cache-vet go vet ./...` | Pasó sin diagnósticos. |
| `git diff --check` | Pasó sin errores de espacios. |
| Primer `go run ./cmd/api` con red real y ADC | El conector obtuvo metadatos de la instancia y alcanzó PostgreSQL. PostgreSQL rechazó la sesión con `SQLSTATE 28P01`. |
| `gcloud sql users list --instance=expo-city-bol-instance --project=expo-city-bol` | Confirmó que `catalogo_api` todavía no existe; solo aparecen los tres usuarios administrados existentes. |
| `go tool air` después de crear `catalogo_api` y corregir las comillas de `DB_PASSWORD` | `pool.Ping` pasó y la API registró `Iniciando API en http://127.0.0.1:8080` el 2026-09-30 a las 18:21:39. |

La ejecución real demuestra ADC, Cloud SQL Admin API, nombre de instancia, red, túnel cifrado y autenticación PostgreSQL. La guía reproducible completa quedó registrada en [la fuente 04](../../docs/fonts/firebase/04-guia-cloud-sql-go-desde-cero.md).

## Consecuencias y continuidad

La conexión de runtime está verificada. El siguiente escalón es aplicar las migraciones con una identidad propietaria y después probar los endpoints autenticados con un ID token real. `golang-migrate` sigue necesitando una URL PostgreSQL accesible, normalmente mediante Cloud SQL Auth Proxy. Este incremento cambia la conexión de la API, no el mecanismo externo de migraciones.

## Fuentes

[^F1]: [Cloud SQL Go Connector con ADC](../../docs/fonts/firebase/03-conector-cloud-sql-go-adc.md), documentación oficial y repositorio de Google Cloud; consulta 2026-09-30.

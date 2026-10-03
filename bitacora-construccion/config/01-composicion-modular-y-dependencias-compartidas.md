# 01 — Composición modular y dependencias compartidas

Fecha de creación: 2026-09-29
Estado: Implementado y verificado en pruebas locales; integración real pendiente
Ámbito: `navigate-business-api`, configuración y composición de módulos Go
Revisión de partida: `2cfc3ec5f9eb061967618a202bdc5d9499538e5e` (`v0.1.0`)
Relacionada con: [auth 03](../auth/03-modulo-auth-y-migraciones-estandar.md), [catálogo Go](../../docs/fonts/go-backend/02-catalogo-completo-cursos-go.md), [arquitectura](../../docs/fonts/go-backend/arquitectura/01-tiago-arquitectura-y-repositorios.md)

## Necesidad y evidencia

El backend ya separa el núcleo de auth de sus adaptadores, pero `cmd/api/main.go` concentra lectura de entorno, creación de Firebase, pool SQL, servicio, HTTP y servidor. Se necesita un lugar explícito para crear una sola instancia de cada recurso y pasar a futuras funcionalidades solo lo que consumen.

Nest agrupa controladores y proveedores en módulos, exporta únicamente los proveedores compartidos y permite que otros módulos reciban la misma instancia mediante imports.[^N1] La lección 25 de Tiago separa transporte, reglas y almacenamiento; su reconstrucción guiada muestra la composición de dependencias en una estructura de aplicación.[^F1][^R1] La lección 24 del curso de sistemas distribuidos muestra la inyección explícita desde `main` de un repositorio al servicio y luego al transporte.[^F2] La guía oficial Go recomienda `cmd` e `internal` para servidores.[^G1]

## Decisión aplicada

Adaptar el principio de Nest, no su contenedor: configuración validada una vez, una composición raíz que crea Firebase y el pool SQL una sola vez, y un módulo auth que construye y expone su servicio y su autenticador HTTP. Los consumidores recibirán dependencias por constructor Go. No hay registro global de proveedores, reflexión ni inicialización oculta. El curso de Tiago respalda la dirección de inyección; la forma de agrupar y exponer estas piezas es una decisión del proyecto.[^F1][^R1][^F2]

## Implementación y comprobación

| Idea de Nest | Adaptación Go y ubicación |
| --- | --- |
| Configuración disponible al arrancar | [`config.Load`](../../internal/config/config.go) lee y valida `FIREBASE_PROJECT_ID`, `DATABASE_URL`, dirección HTTP y emulador. |
| `AppModule` como raíz de composición | [`app.New`](../../internal/app/app.go) construye un verificador Firebase y un pool SQL, comprueba la conexión, crea auth y entrega un `http.Handler`. Cierra el pool si falla el arranque y `App.Close` lo cierra al terminar `run`. |
| Feature module con providers y exports | [`authmodule.New`](../../internal/auth/module/module.go) recibe un `auth.Repository` y un `TokenVerifier`; crea una sola instancia de `auth.Service` y expone `Service`, `Handler` y `Authenticator`. |
| Guard/servicio compartido | [`authhttp.Authenticator`](../../internal/auth/http/middleware.go) comprueba Google y deja el UID verificado en contexto; `VerifiedUID` permite leerlo desde otro handler que monte ese middleware. El handler de auth ya no es dueño del verificador. |
| Bootstrap | [`cmd/api/main.go`](../../cmd/api/main.go) carga configuración, obtiene el handler de la aplicación y arranca HTTP. |

Para una futura funcionalidad, `app.New` construirá su módulo con **el mismo** pool y con `auth.Service` o `auth.Authenticator` si los necesita. Ese constructor declara exactamente sus dependencias; no hace falta volver a crear Firebase, abrir otro pool ni importar internamente auth desde un registro global. El router seguirá montando las rutas concretas. Esta guía no crea una funcionalidad futura ficticia.

| Comprobación ejecutada | Resultado |
| --- | --- |
| `go test ./...` | Pasó. Los tests de auth existentes conservan el contrato HTTP; la nueva prueba protege una ruta de otra funcionalidad con el autenticador compartible y lee el UID. `config.Load` comprueba dirección predeterminada, emulador y variable obligatoria. |
| `go vet ./...` | Pasó sin diagnósticos. |
| `git diff --check` | Pasó sin errores de espacios en los cambios rastreados. |

Sin proyecto Firebase ni instancia Cloud SQL configurados, la composición no se probó contra servicios reales. La similitud con Nest está en el **límite de módulos, exportación de servicios e instancias compartidas**; la resolución automática por metadatos pertenece a Nest y aquí se reemplaza por construcción explícita. El cliente Expo no se modificó.

## Fuentes

[^N1]: [NestJS, Modules](https://docs.nestjs.com/modules), secciones «Modules», «Feature modules» y «Shared modules»; consulta 2026-09-29.
[^F1]: [Tiago Taquelim, 25. Clean Layered Architecture](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture.md), *APIs de Producción con Go*, sección 5, líneas 17–35; consulta 2026-09-29.
[^R1]: [Reconstrucción guiada de la lección 25](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture%20-%20Reconstrucción%20guiada.md), fase 3, `Application` y `NewApplication`; fuente derivada, consulta 2026-09-29.
[^F2]: [Tiago Taquelim, 24. Theory Clean Architecture & Repository Pattern](../../docs/reconstrucciones-cursos-go/Building,%20scaling%20and%20deploying%20distributed%20systems%20in%20Golang%20and%20Kubernetes/Sección%204%20Development%20Environment%20Setup/24.%20Theory%20Clean%20Architecture%20&%20Repository%20Pattern.md), sección 4, líneas 13–35; consulta 2026-09-29.
[^G1]: [Go, Organizing a Go module](https://go.dev/doc/modules/layout), sección «Server project»; consulta 2026-09-29.

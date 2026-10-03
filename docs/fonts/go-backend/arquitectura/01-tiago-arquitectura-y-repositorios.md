# 01 — Tiago: arquitectura por capas y repositorios

Fecha de creación: 2026-09-27
Última actualización: 2026-09-29
Tipo de fuente: Transcripciones originales y reconstrucciones guiadas
Origen: Cursos locales atribuidos por el usuario a Tiago Taquelim
Fecha de consulta: 2026-09-27
Relacionada con: [Mapa general](../01-mapa-general.md), [persistencia](../persistencia-seguridad/01-transacciones-y-seguridad-de-datos.md), [README del proyecto](../../../../README.md)

## Necesidad de consulta

Determinar una arquitectura Go que separe HTTP, reglas de aplicación y persistencia sin llenar un proyecto joven de carpetas e interfaces prematuras.

## Localización verificable

### Fuentes originales

- `F-ARQ-01`: [25. Clean Layered Architecture](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 5 Scaffolding our API Server/25. Clean Layered Architecture.md>), líneas 13–45. Buscar `Transportschicht`, `HTTP-Handler`, `API-Typ` y `Transaktionen`.
- `F-ARQ-02`: [29. The Repository Pattern](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/29. The Repository Pattern.md>), líneas 13–39. Buscar `gekoppelt`, `Schnittstelle` y `Einheitstest`.
- `F-ARQ-03`: [30. Implementing the Repository Pattern](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/30. Implementing the Repository Pattern.md>), líneas 13–35. Buscar `store`, `New Storage` y `app dot store`.
- `F-DIS-01`: [23. Structuring a Microservice project](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/23. Structuring a Microservice project.md>), líneas 13–45. Buscar `basic principles`, `services`, `shared package` y `no right or wrong`.
- `F-DIS-02`: [24. Theory Clean Architecture & Repository Pattern](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/24. Theory Clean Architecture & Repository Pattern.md>), líneas 13–35. Buscar `logical layers`, `interface` y `constructor`.
- `F-DIS-03`: [25. Part 1 Layered Architecture](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/25. Part 1 Layered Architecture.md>), líneas 13–33. Buscar `trip service`, `CMD` e `infrastructure`.
- `F-DIS-04`: [26. Part 2 Repository Pattern](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/26. Part 2 Repository Pattern.md>), líneas 13–67. Buscar `trip repository`, `in-memory`, `service struct` y `context.Background`.

### Reconstrucciones derivadas

- `R-ARQ-01`: [contratos y composición](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 5 Scaffolding our API Server/25. Clean Layered Architecture - Reconstrucción guiada.md>), especialmente `Store`, `Mailer`, `Authenticator`, `Application` y `NewApplication`.
- `R-ARQ-02`: [repository sustituible](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/29. The Repository Pattern - Reconstrucción guiada.md>), especialmente `Store`, `PostgresUserRepository` e `InMemoryUserRepository`.
- `R-ARQ-03`: [storage compuesto](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/30. Implementing the Repository Pattern - Reconstrucción guiada.md>), especialmente `Posts`, `Users`, `Storage` y `NewStorage`.
- `R-DIS-01` a `R-DIS-04`: las cuatro reconstrucciones de arquitectura distribuida enlazadas en el [mapa general](../01-mapa-general.md).

## Evidencia

1. `F-ARQ-01` separa transporte, lógica y almacenamiento, y hace depender la aplicación de contratos intercambiables.
2. La misma fuente no impone una capa de servicio desde el primer día: conserva operaciones atómicas cerca del store mientras no exista una orquestación reutilizable.
3. `F-ARQ-02` y `F-ARQ-03` usan interfaces para desacoplar SQL y para permitir una implementación en memoria, no como decoración arquitectónica.
4. `F-DIS-01` afirma que los principios deben sobrevivir a distintas estructuras de carpetas y que no existe una organización única correcta.
5. `F-DIS-02` fija la dirección `transport → service → repository`, mientras `main` compone las implementaciones concretas.
6. `F-DIS-03` y `F-DIS-04` muestran que una capa de servicio sí se justifica cuando existe lógica y más de un transporte o adaptador.

## Interpretación y límites

`I-ARQ-01`: la unidad de diseño principal es la responsabilidad y la dirección de dependencias, no el número de directorios.

`I-ARQ-02`: las interfaces aportan valor en fronteras sustituibles —base de datos, correo, archivos, IA— o cuando un consumidor necesita un contrato estrecho. No se deduce que cada estructura deba tener una interfaz.

`I-ARQ-03`: el curso de sistemas distribuidos enseña cómo escalar la separación, pero no demuestra que este proyecto deba empezar como microservicios. Esa decisión requeriría evidencia de despliegue, carga, ownership de equipos o aislamiento operativo.

Las reconstrucciones recuperan símbolos representativos cuando la transcripción describe código visible sin preservarlo literalmente. Sus mapas de procedencia declaran el nivel de certeza; no deben citarse como texto literal del instructor.

## Aplicación en el proyecto

`A-ARQ-01`: conservar por ahora un monolito modular. `cmd/api` arranca el servidor y `internal/app` compone las dependencias; `internal/httpapi` mantiene el router general; `internal/auth` aloja reglas y contrato de repositorio, con adaptadores en `internal/auth/http`, `internal/auth/firebase` e `internal/auth/postgres`. El núcleo no importa esos subpaquetes. **Actualización del proyecto, 2026-09-29:** auth usa SQL directo sobre el PostgreSQL asociado a Firebase SQL Connect, con migraciones SQL propias ejecutadas por `golang-migrate`; ver el [plan de auth](../../../plans/01-autenticacion-firebase-y-cuenta-interna.md), la [bitácora 03](../../../../bitacora-construccion/auth/03-modulo-auth-y-migraciones-estandar.md) y la [composición modular](02-nest-modulos-y-composicion-explicita-go.md). Esta agrupación por funcionalidad es adaptación del proyecto, no una estructura de carpetas prescrita por Tiago.

`A-ARQ-02`: crear servicio y repositorio con el primer flujo persistente real, no antes. El flujo debe poder leerse de extremo a extremo antes de extraer abstracciones adicionales.

`A-ARQ-03`: esta ficha registra evidencia y decisiones; la implementación posterior de auth se documenta y verifica en la [bitácora 03](../../../../bitacora-construccion/auth/03-modulo-auth-y-migraciones-estandar.md).

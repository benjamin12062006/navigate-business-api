# 01 — Mapa de fuentes para backend Go

Fecha de creación: 2026-09-27
Tipo de fuente: Índice de investigación y evidencia local
Origen: Biblioteca local de cursos de Go; prioridad indicada por el usuario para Tiago Taquelim
Fecha de consulta: 2026-09-27
Relacionada con: [Catálogo completo de cursos Go](./02-catalogo-completo-cursos-go.md), [arquitectura](./arquitectura/01-tiago-arquitectura-y-repositorios.md), [autenticación y autorización](./autenticacion-autorizacion/01-tiago-autenticacion-y-autorizacion.md), [persistencia y seguridad](./persistencia-seguridad/01-transacciones-y-seguridad-de-datos.md), [operaciones](./operaciones/01-produccion-observabilidad-y-confiabilidad.md), [fuentes complementarias](./complementarias/01-rutas-complementarias.md), [contraste de otros autores](./complementarias/02-otros-autores-go.md) y [reconstrucciones locales](../../reconstrucciones-cursos-go/README.md)

## Necesidad de consulta

Localizar material suficiente para construir backends Go organizados, eficientes y comprensibles, con especial atención a arquitectura, persistencia, autenticación, autorización y preparación para producción. La investigación debe permitir citar el pasaje original y, cuando exista reconstrucción, el símbolo de código derivado que se consultó.

## Criterio de selección

1. Tiago Taquelim es la referencia principal indicada por el usuario.
2. Se priorizan sus cursos `APIs de Producción con Go` y `Building, scaling and deploying distributed systems in Golang and Kubernetes`.
3. Los demás cursos de la biblioteca se usan como contraste cuando responden a una necesidad concreta; Trevor Sawler y Maximilian Schwarzmüller son las primeras rutas complementarias verificadas.
4. No hay un límite artificial de lecciones reconstruibles: se reconstruyen las que hagan falta para comprender o citar con precisión una decisión futura.
5. Una reconstrucción es evidencia derivada. Nunca reemplaza la transcripción original.

El nombre del instructor proviene de la indicación expresa del usuario y del perfil de Udemy aportado. Los archivos locales registran curso, sección y URL de la clase, pero no incluyen un campo estructurado de instructor.

## Convención de evidencia

- `F…`: fuente primaria, normalmente una transcripción original.
- `R…`: reconstrucción guiada derivada de una fuente `F…`.
- `I…`: interpretación o principio transferible.
- `A…`: adaptación propuesta para el proyecto activo.
- `V…`: comprobación realmente ejecutada.

La cadena completa que deberán usar las futuras entradas de construcción es:

```text
F original y pasaje exacto
  ↓
R reconstrucción y símbolo concreto, si fue necesaria
  ↓
I principio transferible
  ↓
A adaptación al proyecto activo
  ↓
archivo y símbolo implementado
  ↓
V comprobación realmente ejecutada
```

Si una etapa todavía no existe, debe declararse como pendiente; no se rellena con una afirmación inventada.

## Ruta principal localizada

| Orden | Tema | Curso principal | Secciones | Resultado de esta consulta |
|---:|---|---|---|---|
| 1 | Arquitectura modular y repositorios | `APIs de Producción con Go` | 5 y 6 | Tres lecciones reconstruidas |
| 2 | Capas en sistemas con más servicios | `Building, scaling and deploying distributed systems in Golang and Kubernetes` | 4 | Cuatro lecciones reconstruidas |
| 3 | Transacciones y consistencia | `APIs de Producción con Go` | 12 | Una lección reconstruida |
| 4 | Autenticación | `APIs de Producción con Go` | 14 | Cuatro lecciones reconstruidas |
| 5 | Autorización | `APIs de Producción con Go` | 15 | Cuatro lecciones reconstruidas |
| 6 | Preparación para producción | `APIs de Producción con Go` | 11 y 17–21 | Fuentes originales localizadas y citadas |
| 7 | Observabilidad y confiabilidad distribuida | `Building, scaling and deploying distributed systems in Golang and Kubernetes` | 13 y 14 | Fuentes originales localizadas y citadas |
| 8 | Comercio, usuarios, pagos y reembolsos | `Ecommerce Seguro con Go`, de Trevor Sawler | 3–15 | Curso y recorrido temático verificados; sin reconstrucción todavía |

## Catálogo de reconstrucciones creado

### Arquitectura y persistencia

- `R-ARQ-01`: [25. Clean Layered Architecture — Reconstrucción guiada](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 5 Scaffolding our API Server/25. Clean Layered Architecture - Reconstrucción guiada.md>).
- `R-ARQ-02`: [29. The Repository Pattern — Reconstrucción guiada](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/29. The Repository Pattern - Reconstrucción guiada.md>).
- `R-ARQ-03`: [30. Implementing the Repository Pattern — Reconstrucción guiada](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 6 Databases/30. Implementing the Repository Pattern - Reconstrucción guiada.md>).
- `R-DAT-01`: [58. SQL Transactions — Reconstrucción guiada](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 12 User Creation/58. SQL Transactions - Reconstrucción guiada.md>).

### Autenticación y autorización

- `R-AUT-01` a `R-AUT-04`: [64. Authentication Overview](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/64. Authentication Overview - Reconstrucción guiada.md>), [65. Basic Authentication](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/65. Basic Authentication - Reconstrucción guiada.md>), [66. Generating Tokens](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/66. Generating Tokens - Reconstrucción guiada.md>) y [67. Validating Tokens](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/67. Validating Tokens - Reconstrucción guiada.md>).
- `R-AUZ-01` a `R-AUZ-04`: [68. Authorization Overview](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/68. Authorization Overview - Reconstrucción guiada.md>), [69. Authorization Database Setup](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/69. Authorization Database Setup - Reconstrucción guiada.md>), [70. Role Precedence Middleware](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/70. Role Precedence Middleware - Reconstrucción guiada.md>) y [71. Fixing the User Invitation](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/71. Fixing the User Invitation - Reconstrucción guiada.md>).

### Arquitectura para sistemas distribuidos

- `R-DIS-01` a `R-DIS-04`: [23. Structuring a Microservice project](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/23. Structuring a Microservice project - Reconstrucción guiada.md>), [24. Theory Clean Architecture & Repository Pattern](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/24. Theory Clean Architecture & Repository Pattern - Reconstrucción guiada.md>), [25. Part 1 Layered Architecture](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/25. Part 1 Layered Architecture - Reconstrucción guiada.md>) y [26. Part 2 Repository Pattern](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 4 Development Environment Setup/26. Part 2 Repository Pattern - Reconstrucción guiada.md>).

## Secuencia recomendada para este proyecto

La investigación no obliga a implementar todo ahora. Para `navigate-business-api`, la secuencia más prudente es:

1. Mantener el monolito modular y hacer crecer un flujo vertical pequeño.
2. Introducir servicio y repositorio cuando aparezca persistencia o una regla de negocio real.
3. Definir el sistema de identidad antes de elegir JWT, OAuth o un proveedor externo.
4. Diseñar autorización por ownership y capacidades del dominio; no depender solo de nombres de rol.
5. Incorporar operaciones de producción según exista una necesidad observable.
6. Posponer microservicios, Kubernetes, tracing distribuido y colas hasta que el monolito muestre una razón concreta para separarse.

## Estado de aplicación

Esta consulta no añade handlers, modelos, autenticación, migraciones ni tests al backend. Tampoco crea una entrada de bitácora de implementación, porque aún no se ha aplicado una funcionalidad nueva. Los documentos de esta carpeta son fichas de fuentes para decisiones futuras.

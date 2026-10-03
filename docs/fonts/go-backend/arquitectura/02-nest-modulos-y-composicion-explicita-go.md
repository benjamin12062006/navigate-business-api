# 02 — Módulos de Nest y composición explícita en Go

Fecha de creación: 2026-09-29
Tipo de fuente: documentación oficial y transcripciones Go
Fecha de consulta: 2026-09-29
Relacionada con: [arquitectura Go](01-tiago-arquitectura-y-repositorios.md), [catálogo de cursos](../02-catalogo-completo-cursos-go.md), [bitácora config 01](../../../../bitacora-construccion/config/01-composicion-modular-y-dependencias-compartidas.md)

## Pregunta

¿Cómo compartir una única instancia de recursos y servicios entre funcionalidades de este monolito Go, inspirándose en los módulos de Nest sin añadir un contenedor de dependencias?

## Evidencia

| Fuente | Mecanismo observado | Límite para este proyecto |
| --- | --- | --- |
| [NestJS, Modules](https://docs.nestjs.com/modules), «Modules», «Feature modules» y «Shared modules» | Un módulo agrupa controladores y proveedores; `exports` define qué proveedores reciben los módulos que lo importan. Un proveedor exportado se comparte como instancia entre consumidores. | Nest resuelve el grafo mediante metadatos e inyector propios de TypeScript. |
| [NestJS, Providers](https://docs.nestjs.com/providers), «Dependency injection» y «Provider registration» | Los consumidores reciben servicios por constructor; los proveedores se registran en un módulo. | Go no tiene decoradores de Nest ni necesita ese registro para construir este backend. |
| [Tiago, 25. Clean Layered Architecture](../../../reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture.md), líneas 17–35; [reconstrucción guiada](../../../reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture%20-%20Reconstrucción%20guiada.md), fase 3 | Se separan transporte, servicio y almacenamiento y se conectan dependencias en una estructura de aplicación. | La estructura concreta de paquetes y la forma de compartir auth son adaptaciones de este proyecto. |
| [Tiago, 24. Theory Clean Architecture & Repository Pattern](../../../reconstrucciones-cursos-go/Building,%20scaling%20and%20deploying%20distributed%20systems%20in%20Golang%20and%20Kubernetes/Sección%204%20Development%20Environment%20Setup/24.%20Theory%20Clean%20Architecture%20&%20Repository%20Pattern.md), líneas 13–35 | El punto de entrada inyecta la implementación SQL en el servicio y el servicio en HTTP. | No prescribe un contenedor automático. |
| [Go, Organizing a Go module](https://go.dev/doc/modules/layout), «Server project» | Para servidores se recomiendan paquetes privados bajo `internal` y comandos bajo `cmd`. | No impone nombres ni un grafo de módulos. |

## Adaptación elegida

`internal/app.New` actúa como módulo raíz: crea una vez Firebase Admin y el pool SQL, construye los módulos de funcionalidades y entrega el handler HTTP. `internal/auth/module.New` agrupa el servicio, handler y autenticador de auth; sus campos exportados son las piezas que la composición puede pasar a otros módulos. `internal/config.Load` valida la configuración al arrancar. Los constructores declaran dependencias concretas o interfaces pequeñas, de modo que el intercambio se ve en el código y no depende de un registro global. La [bitácora](../../../../bitacora-construccion/config/01-composicion-modular-y-dependencias-compartidas.md) registra implementación y comprobación.

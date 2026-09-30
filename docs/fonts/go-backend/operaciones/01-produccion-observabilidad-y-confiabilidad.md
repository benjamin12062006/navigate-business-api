# 01 — Producción, observabilidad y confiabilidad

Fecha de creación: 2026-09-27
Tipo de fuente: Transcripciones originales
Origen: Cursos locales principales atribuidos por el usuario a Tiago Taquelim
Fecha de consulta: 2026-09-27
Relacionada con: [Mapa general](../01-mapa-general.md), [arquitectura](../arquitectura/01-tiago-arquitectura-y-repositorios.md), [autenticación y autorización](../autenticacion-autorizacion/01-tiago-autenticacion-y-autorizacion.md)

## Necesidad de consulta

Identificar las capacidades operativas que convierten una API funcional en un servicio mantenible, sin implementarlas antes de que el proyecto tenga el flujo que justifique cada una.

## Localización verificable

- `F-OPS-01`: [56. Adding a Logger](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 11 Structured Logging/56. Adding a Logger.md>), líneas 13–23 sobre estructura, niveles y búsqueda; 25–35 sobre inyección y composición; 37–47 sobre manejo de errores en la capa superior.
- `F-OPS-02`: [74. Testing Overview](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 17 Testing/74. Testing Overview.md>), líneas 13–25 sobre TDD, herramientas estándar, mocks mediante interfaces e inyección de dependencias.
- `F-OPS-03`: [77. Graceful Server Shutdown](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 18 Graceful Shutdown/77. Graceful Server Shutdown.md>), líneas 13–19 sobre el problema; 29–35 sobre señales, timeout y `Server.Shutdown`; 37–45 sobre trabajo largo y verificación.
- `F-OPS-04`: [78. Rate Limiting our API](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 19 Rate Limiting/78. Rate Limiting our API.md>), líneas 13–29 sobre amenaza y costo; 31–47 sobre contrato, middleware, IP y `429`; 49–61 sobre implementación local; 53–55 sobre el límite de memoria en varias réplicas.
- `F-OPS-05`: [79. Handling CORS errors](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 20 Handling CORS/79. Handling CORS errors.md>), líneas 35–53 sobre origen, middleware, orden, allowlist explícita y configuración por entorno.
- `F-OPS-06`: [80. Basic Server Metrics](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 20 Handling CORS/80. Basic Server Metrics.md>), líneas 13–21 sobre `expvar` y protección del endpoint; 25–39 sobre versión, estadísticas de base de datos, goroutines y exposición mínima.
- `F-OBS-01`: [80. Intro to Distributed Tracing](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 13 Observability/80. Intro to Distributed Tracing.md>), líneas 13–23 sobre el problema y la propagación de trace/span; 25–35 sobre OpenTelemetry, backends e instrumentación.
- `F-REL-01`: [86. Understanding DLQ and Retries](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Building, scaling and deploying distributed systems in Golang and Kubernetes/Sección 14 Reliability/86. Understanding DLQ and Retries.md>), líneas 13–31 sobre acknowledgements, persistencia, retries, poison messages y dead-letter queues.

## Evidencia

1. El logging estructurado permite distinguir niveles y filtrar propiedades; el curso lo inyecta en la aplicación y registra los errores HTTP en la capa superior.
2. La facilidad de probar aparece como consecuencia de límites e inyección de dependencias; la fuente no exige una librería externa específica.
3. El cierre gradual atiende solicitudes en curso al recibir `SIGINT` o `SIGTERM`, dentro de un timeout.
4. El rate limit en memoria solo coordina una instancia. Varias réplicas requieren estado distribuido o una solución externa.
5. CORS necesita una allowlist explícita y el orden del middleware afecta qué errores llegan al navegador.
6. Un endpoint de métricas puede exponer información sensible y debe protegerse y limitar lo publicado.
7. Tracing distribuido empieza a aportar valor cuando una solicitud cruza servicios o dependencias y los logs dejan de mostrar la trayectoria completa.
8. Retries y DLQ son mecanismos para mensajería asíncrona; una poison message no debe reintentarse sin límite ni desaparecer sin investigación.

## Interpretación y límites

`I-OPS-01`: logging, cierre gradual y configuración segura son necesidades tempranas de producción. Tracing distribuido, DLQ y Kubernetes pertenecen a una etapa posterior si aún no existen múltiples servicios o colas.

`I-OPS-02`: no toda métrica debe exponerse. Primero se formula una pregunta operativa y después se instrumenta la señal mínima que permite responderla.

`I-OPS-03`: las fuentes de testing ayudan a diseñar código comprobable, pero esta consulta no autoriza generar una suite amplia. Los tests nuevos se crearán por incremento cuando el usuario los solicite o cuando formen parte explícita de una funcionalidad.

## Aplicación en el proyecto

`A-OPS-01`: no se implementaron capacidades operativas en este trabajo. Para próximos incrementos, la prioridad razonable es configuración, errores y logging; después cierre gradual y métricas; finalmente, solo si la arquitectura lo exige, tracing y confiabilidad de mensajería.

`A-OPS-02`: las comprobaciones de este trabajo se limitan a la integridad de la documentación y de las reconstrucciones. No se creó una suite de tests del backend.

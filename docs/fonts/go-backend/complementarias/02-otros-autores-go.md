# 02 — Contraste con otros autores de Go

Fecha de creación: 2026-09-27
Tipo de fuente: Consulta complementaria de transcripciones originales
Origen: Biblioteca local de cursos de Go
Fecha de consulta: 2026-09-27
Relacionada con: [Mapa general](../01-mapa-general.md), [rutas complementarias](./01-rutas-complementarias.md) y [reconstrucciones locales](../../../reconstrucciones-cursos-go/README.md)

## Motivo

La investigación inicial reconstruyó lecciones únicamente de los dos cursos principales atribuidos a Tiago Taquelim. Este documento corrige cualquier impresión de que toda la biblioteca complementaria ya había sido utilizada y registra dos autores distintos que sí pueden aportar contraste al proyecto.

## Trevor Sawler: comercio electrónico seguro

Fuente original consultada: [1. Introduction](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Ecommerce Seguro con Go/Sección 1 Introduction/1. Introduction.md>).

La introducción describe una solución con frontend, API backend y un servicio separado; incluye pagos, suscripciones, autenticación, administración de usuarios y generación de facturas. El inventario del curso añade secciones específicas sobre reembolsos, recuperación de contraseña, paginación y validación.

Uso previsto en `navigate-business-api`:

- contrastar el diseño de pagos y reembolsos cuando el marketplace llegue a esa fase;
- estudiar flujos administrativos y gestión de usuarios;
- revisar validación de entradas y separación de responsabilidades en comercio electrónico.

No se adopta todavía Stripe, una arquitectura de microservicios ni el esquema de autenticación del curso. Esas decisiones requieren una necesidad de producto y una comparación con el enfoque principal antes de implementarse.

## Maximilian Schwarzmüller: recorrido REST compacto

Fuente original consultada: [1. Welcome To The Course!](</Users/benji/Documents/chrome-extencions/trancripciones/Go/Go Completo/Sección 1 Getting Started/1. Welcome To The Course!.md>).

La sección 11 del curso `Go Completo` contiene un proyecto REST con SQL, registro de usuarios, contraseñas, JWT, protección de rutas y autorización. Su ejemplo usa Gin, mientras este proyecto ya eligió Chi; por eso sirve para contrastar el flujo funcional, no para reemplazar el router ni copiar su estructura.

## Boris Paskhaver

La búsqueda de la biblioteca local no encontró un curso de Go atribuible a Boris Paskhaver. Sí puede seguir siendo una referencia de estilo pedagógico, pero no se utilizará como evidencia técnica de Go hasta disponer de una fuente concreta y verificable.

## Estado

Estas fuentes ya fueron verificadas como rutas complementarias reales. Aún no se creó una reconstrucción de sus lecciones ni se modificó el backend a partir de ellas. La primera reconstrucción complementaria deberá responder a una necesidad concreta del siguiente incremento, no a una cuota artificial de autores.

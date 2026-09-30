# 01 — Transacciones y seguridad de datos

Fecha de creación: 2026-09-27
Última actualización: 2026-09-29
Tipo de fuente: Transcripciones originales y reconstrucciones guiadas
Origen: `APIs de Producción con Go`, secciones 12 y 15
Fecha de consulta: 2026-09-27
Relacionada con: [Mapa general](../01-mapa-general.md), [arquitectura](../arquitectura/01-tiago-arquitectura-y-repositorios.md), [autenticación y autorización](../autenticacion-autorizacion/01-tiago-autenticacion-y-autorizacion.md)

## Necesidad de consulta

Entender cómo mantener consistentes operaciones que escriben varias entidades y cómo evolucionar datos de autorización sin romper registros existentes ni aceptar privilegios controlados por el cliente.

## Localización verificable

- `F-DAT-01`: [58. SQL Transactions](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 12 User Creation/58. SQL Transactions.md>), líneas 13–21 para `withTx`; líneas 23–47 para crear usuario e invitación; líneas 41–45 para exigir que toda consulta use `*sql.Tx`; líneas 61–75 para expiración, UUID y hash.
- `R-DAT-01`: [58. SQL Transactions — Reconstrucción guiada](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 12 User Creation/58. SQL Transactions - Reconstrucción guiada.md>), símbolos `withTx`, `createUserInvitation` y `UserStore.CreateAndInvite`.
- `F-DAT-02`: [69. Authorization Database Setup](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/69. Authorization Database Setup.md>), líneas 13–35; tabla de roles y migración segura de `role_id`.
- `R-DAT-02`: [69. Authorization Database Setup — Reconstrucción guiada](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/69. Authorization Database Setup - Reconstrucción guiada.md>), fases de catálogo, backfill, restricción y reversión.
- `F-DAT-03`: [71. Fixing the User Invitation](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/71. Fixing the User Invitation.md>), líneas 13–27; resolver el rol por nombre y no aceptar rol privilegiado desde registro público.
- `R-DAT-03`: [71. Fixing the User Invitation — Reconstrucción guiada](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/71. Fixing the User Invitation - Reconstrucción guiada.md>), símbolos `roleName` y `UserStore.create`.

## Evidencia

1. `F-DAT-01` envuelve una unidad de trabajo en `BeginTx`, rollback diferido, callback y commit.
2. Todas las consultas de esa unidad deben ejecutarse con el mismo `*sql.Tx`; usar `*sql.DB` en una operación interna rompe la atomicidad prometida.
3. El token de invitación se entrega en claro al destinatario, pero se conserva como hash y con expiración.
4. `F-DAT-02` muestra una migración por etapas: crear roles, añadir referencia, completar datos existentes y solo después imponer `NOT NULL`.
5. `F-DAT-03` mantiene el rol por defecto en el servidor y permite roles privilegiados únicamente mediante flujos internos controlados.

## Interpretación y límites

`I-DAT-01`: una transacción debe representar una invariante del dominio, no servir como envoltorio automático de cada consulta. Ejemplos futuros serían tienda + membresía inicial, pedido + líneas o importación + versión publicada.

`I-DAT-02`: el helper de transacción pertenece cerca del adaptador SQL, pero la decisión de qué acciones forman una unidad atómica puede originarse en el caso de uso. Si aparece orquestación entre varios puertos, será necesario revisar el límite.

`I-DAT-03`: el esquema de roles del curso es demostrativo. No prueba que tres roles globales sean suficientes para una plataforma multi-tienda.

## Aplicación en el proyecto

`A-DAT-01`: cuando exista persistencia, revisar las migraciones y hacer backfill antes de restricciones obligatorias. **Actualización del proyecto, 2026-09-29:** en auth se usarán migraciones SQL propias y transacciones Go sobre el PostgreSQL asociado a Firebase SQL Connect. El servicio no gestionará el DDL de estas tablas; ver el [plan de auth](../../../plans/01-autenticacion-firebase-y-cuenta-interna.md).

`A-DAT-02`: almacenar secretos recuperables solo si el producto realmente necesita recuperarlos; tokens de activación, recuperación o invitación deben guardarse mediante representaciones seguras y expirar.

`A-DAT-03`: los permisos de tienda deben asignarse desde procesos controlados por el servidor. El cliente no podrá elegir un rol privilegiado en un payload público.

No se añadieron migraciones, tablas ni repositorios al proyecto durante esta investigación.

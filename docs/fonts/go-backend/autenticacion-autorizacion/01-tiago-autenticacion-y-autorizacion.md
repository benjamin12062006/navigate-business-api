# 01 — Tiago: autenticación y autorización

Fecha de creación: 2026-09-27
Tipo de fuente: Transcripciones originales y reconstrucciones guiadas
Origen: `APIs de Producción con Go`, secciones 14 y 15
Fecha de consulta: 2026-09-27
Relacionada con: [Mapa general](../01-mapa-general.md), [persistencia y seguridad](../persistencia-seguridad/01-transacciones-y-seguridad-de-datos.md), [operaciones](../operaciones/01-produccion-observabilidad-y-confiabilidad.md)

## Necesidad de consulta

Separar correctamente identidad, validación de credenciales y permisos, y conocer qué piezas hacen falta antes de decidir entre OAuth, un proveedor externo o JWT propio.

## Localización verificable

### Autenticación

- `F-AUT-01`: [64. Authentication Overview](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/64. Authentication Overview.md>), líneas 13–85. Buscar `Basisauthentifizierung`, `Token-Authentifizierung`, `Rollen`, `OAuth` y `zustandslos`.
- `F-AUT-02`: [65. Basic Authentication](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/65. Basic Authentication.md>), líneas 13–51. Buscar `Middleware`, `base64`, `WWW` y `Umgebungsvariable`.
- `F-AUT-03`: [66. Generating Tokens](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/66. Generating Tokens.md>), líneas 13–85. Buscar `Schrägstrich-Token`, `Aufzählungsangriff`, `Authentifikator`, `Ansprüche` y `Cookie`.
- `F-AUT-04`: [67. Validating Tokens](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/67. Validating Tokens.md>), líneas 13–67. Buscar `Auth-Token-Middleware`, `Unterzeichnungsmethode`, `Ansprüche` y `authentifizierten Benutzer`.

### Autorización

- `F-AUZ-01`: [68. Authorization Overview](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/68. Authorization Overview.md>), líneas 13–33. Buscar `403`, `Eigentümer` y `Berechtigungen`.
- `F-AUZ-02`: [69. Authorization Database Setup](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/69. Authorization Database Setup.md>), líneas 13–35. Buscar `CREATE TABLE` y `not null`.
- `F-AUZ-03`: [70. Role Precedence Middleware](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/70. Role Precedence Middleware.md>), líneas 13–67. Buscar `check ownership`, `Moderator`, `higher or equal` y `JOIN`.
- `F-AUZ-04`: [71. Fixing the User Invitation](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/71. Fixing the User Invitation.md>), líneas 13–27. Buscar `role ID`, `select ID from roles` y `admin`.

### Reconstrucciones derivadas y símbolos recuperados

- `R-AUT-01`: [comparación de estrategias](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/64. Authentication Overview - Reconstrucción guiada.md>).
- `R-AUT-02`: [`BasicAuthConfig`, `BasicAuth` y `basicUnauthorized`](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/65. Basic Authentication - Reconstrucción guiada.md>).
- `R-AUT-03`: [`Authenticator`, `JWTAuthenticator`, `GenerateToken` y handler de token](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/66. Generating Tokens - Reconstrucción guiada.md>).
- `R-AUT-04`: [`ValidateToken`, contexto tipado y `TokenAuth`](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 14 Authentication/67. Validating Tokens - Reconstrucción guiada.md>).
- `R-AUZ-01`: [política ownership + roles](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/68. Authorization Overview - Reconstrucción guiada.md>).
- `R-AUZ-02`: [migración del catálogo de roles](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/69. Authorization Database Setup - Reconstrucción guiada.md>).
- `R-AUZ-03`: [`RoleStore.ByName`, `hasRolePrecedence` y `CheckPostOwnership`](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/70. Role Precedence Middleware - Reconstrucción guiada.md>).
- `R-AUZ-04`: [`roleName` y creación segura de usuario](</Users/benji/Documents/chrome-extencions/trancripciones/Go/APIs de Producción con Go/Sección 15 Authorization/71. Fixing the User Invitation - Reconstrucción guiada.md>).

## Evidencia

1. `F-AUT-01` distingue Basic Auth, tokens y OAuth. JWT es una elección educativa del curso, no una obligación universal.
2. `F-AUT-03` evita revelar si un correo existe mediante respuestas uniformes y mantiene la emisión de tokens detrás de un contrato sustituible.
3. `F-AUT-04` valida algoritmo, firma, expiración, audience e issuer; después recupera al usuario actual y lo coloca en el contexto.
4. `F-AUZ-01` diferencia `401` de `403` y propone comprobar ownership antes de elevar por roles.
5. `F-AUZ-03` implementa precedencia de roles y carga el rol junto al usuario.
6. `F-AUZ-04` evita que el registro público permita autoasignarse un rol privilegiado.

## Interpretación y límites

`I-AUT-01`: autenticación responde quién es el actor; autorización responde qué puede hacer sobre un recurso concreto. Deben diseñarse y probarse por separado.

`I-AUT-02`: los claims de un token pueden quedar obsoletos. Para permisos sensibles, recuperar el estado vigente o usar tokens breves y una estrategia explícita de revocación reduce ese riesgo.

`I-AUT-03`: un proveedor externo acelera la identidad inicial, mientras JWT propio aumenta control y responsabilidad operativa. Estas fuentes no determinan todavía cuál conviene al negocio, porque faltan decisiones sobre frontend, recuperación de cuenta, MFA, costos, dependencia externa y administración.

`I-AUZ-01`: una jerarquía lineal de roles funciona para un dominio pequeño. Si una tienda tiene propietarios, colaboradores y compradores con permisos cruzados, probablemente será más claro modelar membresías y capacidades por tienda.

## Aplicación en el proyecto

`A-AUT-01`: definir primero una identidad única y una cuenta capaz de actuar como comprador y como operador de una o varias tiendas. Los permisos deben depender del contexto de tienda, no de correos separados.

`A-AUT-02`: mantener el mecanismo de autenticación detrás de un contrato para poder comenzar con un proveedor externo y migrar después, si existe una razón real, a sesiones o JWT propios.

`A-AUT-03`: no se implementó autenticación ni autorización en esta consulta. La ficha conserva las opciones y los riesgos para una fase futura.

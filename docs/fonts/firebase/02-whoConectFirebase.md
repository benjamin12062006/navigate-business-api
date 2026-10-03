Para habilitar y configurar correctamente Firebase Authentication y Firebase SQL Connect en tu proyecto para tu backend en Go, debes seguir estos pasos en la consola de Firebase y en tu entorno local:

1. Activar Firebase Authentication

Dado que tu backend en Go exigirá que el proveedor sea únicamente google.com:

Entra a la consola de Firebase y selecciona tu proyecto.

En el menú lateral, dirígete a Build > Authentication y haz clic en Comenzar.

En la pestaña Sign-in method, selecciona Google como proveedor de inicio de sesión, actívalo, configura los datos de soporte y guarda los cambios. Asegúrate de que el resto de proveedores (como correo/contraseña) permanezcan deshabilitados.

2. Configurar Firebase SQL Connect

Dado que tu decisión es utilizar la instancia de PostgreSQL con tu propio sistema de migraciones (pgx) en Go, en lugar de delegar el esquema a GraphQL:

En la consola de Firebase, ve a Build > SQL Connect y haz clic en Comenzar para crear tu servicio y aprovisionar la instancia de Cloud SQL para PostgreSQL. (Ten en cuenta que tu plan de Firebase es Spark—un plan sin costo—pero la base de datos de PostgreSQL subyacente se factura de forma independiente).

Abre tu terminal y utiliza la CLI de Firebase para habilitar el acceso directo de SQL sin depender del motor de GraphQL:

Ejecuta firebase dataconnect:sql:setup para iniciar la configuración de conexión directa.

Ejecuta firebase dataconnect:sql:grant para conceder a tu backend de Go los permisos de acceso directo a la base de datos PostgreSQL.

Para la conexión de producción o desarrollo remoto, asegúrate de que tu cuenta de servicio o la identidad de tu backend de Go disponga del rol de IAM roles/cloudsql.instanceUser.

3. Configurar el Admin SDK en tu backend de Go

Para pruebas locales con emuladores: Configura la variable de entorno FIREBASE_AUTH_EMULATOR_HOST en tu máquina de desarrollo [O11]. Esto permitirá que tu código Go en los handlers acepte los tokens generados localmente por el emulador de Auth sin requerir firmas reales de Google.

Para entorno real: Configura la variable de entorno GOOGLE_APPLICATION_CREDENTIALS apuntando al archivo JSON de tu cuenta de servicio de Firebase [O12]. El SDK de Go utilizará automáticamente estas credenciales (ADC) para comunicarse con los servidores de Firebase, lo que te permitirá realizar de forma segura llamadas como VerifyIDTokenAndCheckRevoked [O2].

¡Con este flujo tendrás la infraestructura lista para realizar tus pruebas de concurrencia y verificar tus criterios de cierre de integración con identidades reales!
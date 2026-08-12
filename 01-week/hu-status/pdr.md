# Product Brief: Auth & User Service


## Declarado Tech Stack

- **Backend:** Go
- **Database:** PostgreSQL
- **Cache / Security:** Redis / JWT (JSON Web Tokens)

---

## 00 - Contexto inicial

El sistema de gestión de cine requiere un componente central que garantice la identidad de los usuarios, el acceso seguro a los servicios y la gestión de permisos según los roles operativos del negocio (Cliente, Administrador, Personal de Taquilla).

El microservicio Auth & User Service no procesa reservas ni administra el catálogo de películas directamente; su responsabilidad principal es actuar como la fuente de verdad para usuarios, credenciales, sesiones y emisión de tokens de autenticación para toda la arquitectura de microservicios.

Preferencia tecnológica: Go para el desarrollo de la API backend por su rendimiento y bajo consumo de recursos, PostgreSQL para el almacenamiento persistente de usuarios y credenciales, y Redis para gestión de revocación de tokens o lista negra de sesiones.

## 01 - Necesidades y problemas

- Permitir el registro público de nuevos clientes de forma simple y segura.
- Validar credenciales y emitir tokens JWT firmados que permitan la comunicación segura con los demás microservicios.
- Diferenciar accesos y permisos mediante roles (Cliente para reservar/comprar, Admin para gestionar carteleras y funciones).
- Proveer endpoints para que otros microservicios o el API Gateway puedan verificar la validez de un token y extraer los datos del usuario (User ID, Email, Role).
- Permitir a los usuarios consultar y actualizar la información de su perfil.

**Problema principal:** Sin un servicio centralizado de autenticación, cada microservicio tendría que gestionar credenciales de usuario, generando duplicidad de datos, posibles brechas de seguridad y desacoplamiento deficiente.

## 02 - Procesos actuales / Flujo esperado

Actualmente, los servicios no cuentan con un mecanismo central para verificar la identidad de los clientes antes de permitirles bloquear asientos o realizar compras.

**Proceso esperado:**

1. **Registro:** El cliente envía sus datos básicos (nombre, correo, contraseña) y el servicio almacena la contraseña de forma encriptada (hash con Bcrypt).
2. **Inicio de Sesión:** El cliente envía sus credenciales; al ser válidas, el servicio retorna un token JWT firmado junto con información básica del usuario.
3. **Validación en peticiones:** Los demás microservicios (Catalog, Booking, Notification) reciben el token JWT en la cabecera `Authorization` y consumen la clave pública o validan el token para autorizar las operaciones.
4. **Gestión de Perfil:** El usuario puede consultar (`GET /me`) o actualizar (`PUT /me`) sus datos personales.

## 03 - Preguntas abiertas

- ¿Se requerirá soporte para inicio de sesión social (OAuth2 con Google o Facebook) en fases futuras?
- ¿Los tokens JWT tendrán una fecha de expiración corta (ej. 15-30 minutos) con implementación de Refresh Tokens?
- ¿Se necesita confirmación de cuenta vía correo electrónico antes de activar a un usuario nuevo?
- ¿Cómo se manejará el bloqueo de cuenta tras múltiples intentos fallidos de inicio de sesión?
- ¿Se requiere un rol especial de "Personal de Taquilla/Escaner" para la validación física de entradas en el cine?

## 04 - Glosario de negocio

- **JWT (JSON Web Token):** Estándar de token firmado digitalmente utilizado para transmitir de forma segura la identidad del usuario entre el cliente y los microservicios.
- **Hash de Contraseña:** Mapeo unidireccional seguro (usando algoritmos como Bcrypt) para almacenar contraseñas sin guardarlas en texto plano.
- **Rol:** Etiqueta asignada a un usuario (Cliente, Admin) que define qué acciones o endpoints tiene permitido ejecutar dentro del sistema.
- **Payload del Token:** Información contenida dentro del JWT que incluye identidades clave como `user_id`, `email` y `role`.
- **Revocación de Token:** Proceso mediante el cual un JWT activo es invalidado (por ejemplo, al cerrar sesión) utilizando Redis como lista negra.
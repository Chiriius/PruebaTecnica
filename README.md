# PruebaTecnica

## Iniciar el Proyecto con Docker

1. Ubícate en la raíz del proyecto y ejecuta el siguiente comando:

   ```bash
   docker-compose up --build

2. Este comando iniciará tanto la API como la base de datos (DB) como contenedores en Docker.

3. Una vez que los servicios estén en funcionamiento, podrás acceder a la documentación de la API a través de Swagger en la siguiente URL: http://localhost:8080/swagger/index.html#/

Además en el repositorio estará una colección de Insomnia con todos los campos configurados, lo que te permitirá explorar y probar los endpoints de la API de manera más sencilla.

## Iniciar el Proyecto desde el Editor de texto

Si prefieres ejecutar la API sin contenerizarla, sigue estos pasos:

1. En el archivo main.go, busca la línea 30 donde se define la URL de la base de datos:
   
   ```bash
   dbUrl = "mongodb://mongodb:27017"

2. Cambia esa línea por:

    ```bash
   dbUrl = "mongodb://localhost:27017"

3. Asegúrate de que el servicio de la base de datos esté en funcionamiento, iniciando el contenedor de la DB desde el docker-compose.

4. Por ultimo accede a la ruta del archivo main.go y ejecuta go run main.go


Este enfoque asegura que los usuarios tengan opciones tanto para usar Docker como para ejecutar la API de manera local.

## Tecnologías y Enfoque

Para la realización de esta prueba, utilicé las siguientes tecnologías y enfoques:

- **Go** con el framework **Gin** para el desarrollo de la API.
- **Testify** para las pruebas unitarias, incluyendo la creación de mocks.
- **Docker** para la contenerización, aplicando un enfoque de **multi-stage** para reducir el tamaño de la imagen final.
- **Logrus** para la gestión de logs en todas las capas de la aplicación.
- **Clean Architecture** para garantizar un código desacoplado y fácil de mantener.
- Implementación de dos métodos de transporte: **HTTP** y **gRPC**.
- **MongoDB** como base de datos NoSQL.
- **Swagger** para la documentación de la API.

Además, he pre-poblado la base de datos con datos estáticos para que, al momento de crearla, se genere automáticamente con información de ejemplo. También he implementado varios endpoints adicionales, además de los requeridos, y he seguido los principios **SOLID** para asegurar un código limpio y escalable.


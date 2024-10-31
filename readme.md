## Documentación del Proyecto de Validación de Oraciones

### Descripción del Proyecto

Este proyecto implementa un servidor web en **Go** para la validación de oraciones en inglés. Utiliza un autómata independiente de contexto que analiza la estructura gramatical de cada oración para determinar si es válida en:
- **Presente Simple**
- **Pasado Simple**
- **Uso del verbo "To Be"**
- **Gramática de Sujeto-Verbo-Complemento**

### Estructura del Proyecto

- **main.go**: Inicializa el servidor web y define las rutas.
- **handlers**: Contiene la lógica de manejo de peticiones, como la recepción de oraciones y el despliegue de resultados en el navegador.
- **validators**: Incluye el autómata que valida la estructura de las oraciones en función de sus reglas gramaticales.
- **templates**: Archivos HTML para la interfaz de usuario.

### Instalación

1. Clona el repositorio:
   ```bash
   git clone https://github.com/jjvnz/validar_oraciones.git
   cd validar_oraciones
   ```

2. Ejecuta el proyecto:
   ```bash
   go run main.go
   ```

3. Abre tu navegador y ve a `http://localhost:8080`.

### Ejemplo de Uso

1. Ingresa hasta 5 oraciones en inglés, cada una terminada en punto (`.`).
2. Presiona "Validar Oraciones".
3. El sistema mostrará si cada oración es válida o no de acuerdo a la estructura gramatical reconocida.

### Ejemplo de Oraciones

```plaintext
I am playing a game
```

---

### Construcción y Ejecución con Docker

Para construir y ejecutar el proyecto utilizando Docker, sigue estos pasos:

#### Requisitos Previos

Asegúrate de tener [Docker](https://www.docker.com/get-started) instalado en tu sistema.

#### Construcción de la Imagen Docker

1. En la raíz del proyecto, crea una imagen Docker utilizando el siguiente comando:

   ```bash
   docker build -t validar_oraciones .
   ```

   Esto crea una imagen llamada `validar_oraciones` basada en el `Dockerfile` presente en el directorio.

#### Ejecución del Contenedor Docker

2. Una vez que la imagen se ha construido con éxito, puedes ejecutar el contenedor con el siguiente comando:

   ```bash
   docker run -d -p 8080:8080 validar_oraciones
   ```

   - `-d`: Ejecuta el contenedor en segundo plano (modo "detached").
   - `-p 8080:8080`: Mapea el puerto 8080 del contenedor al puerto 8080 de tu máquina local.

3. Abre tu navegador y ve a `http://localhost:8080`.


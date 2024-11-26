## Documentación del Proyecto de Validación de Oraciones

### Descripción del Proyecto

Este proyecto implementa un servidor web en **Go** para la validación de oraciones en inglés en **pasado simple afirmativo**. El sistema utiliza **análisis gramatical** para examinar la estructura de las oraciones, asegurándose de que sigan las reglas gramaticales del **pasado simple afirmativo**. El análisis gramatical verifica la correcta conjugación de los verbos, el uso adecuado de los sujetos y otros elementos esenciales en la estructura de la oración.

### Estructura del Proyecto

- **main.go**: Inicializa el servidor web y define las rutas.
- **handlers**: Contiene la lógica de manejo de peticiones, como la recepción de oraciones y el despliegue de resultados en el navegador.
- **parser**: Implementa el análisis gramatical, validando que las oraciones sigan las reglas del **pasado simple afirmativo**.
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
3. El sistema mostrará si cada oración es válida o no de acuerdo con la estructura gramatical del **pasado simple afirmativo**.

### Ejemplo de Oraciones: 5 oraciones afirmativas correctas en pasado simple y 5 incorrectas

**Oraciones correctas:**

1. I visited my grandmother last weekend.
2. She played soccer with her friends yesterday.
3. They watched a movie last night.
4. We cleaned the house on Saturday.
5. He studied for the test last week.
6. I was happy yesterday. (Correcta: "was" es correcto para el pronombre "I")
7. They were at the park all day. (Correcta: "were" es correcto para el pronombre "they")

**Oraciones incorrectas:**

1. I visit my grandmother last weekend.  
   *(Incorrecto: "visit" debería ser "visited")*

2. She play soccer with her friends yesterday.  
   *(Incorrecto: "play" debería ser "played")*

3. They watches a movie last night.  
   *(Incorrecto: "watches" debería ser "watched")*

4. We clean the house on Saturday.  
   *(Incorrecto: "clean" debería ser "cleaned" para indicar una acción pasada)*

5. He studys for the test last week.  
   *(Incorrecto: "studys" debería ser "studied")*

6. I were happy yesterday.
   (Incorrecto: "were" debería ser "was" para el pronombre "I")

7. They was at the park all day.
   (Incorrecto: "was" debería ser "were" para el pronombre "they")

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

   ```bash
   sudo docker build --no-cache -t validar_oraciones .
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

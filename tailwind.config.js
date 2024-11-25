module.exports = {
  content: [
    "./templates/**/*.html", // Ruta a tus plantillas
    "./static/js/**/*.js"    // Ruta a tus scripts
  ],
  theme: {
    extend: {},  // Puedes extender el tema si lo necesitas
  },
  darkMode: 'class',  // Habilita el modo oscuro basado en la clase (esto fuera de `theme`)
  plugins: [],
}

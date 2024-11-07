document.addEventListener('DOMContentLoaded', function() {
    const textarea = document.getElementById('oraciones');
    const sentenceCount = document.getElementById('sentenceCount');
    const characterCount = document.getElementById('characterCount');
    const submitButton = document.querySelector('button[type="submit"]'); // Selecciona el botón de envío

    // Desactivar el botón de envío al cargar la página
    submitButton.disabled = true;

    textarea.addEventListener('input', function() {
        const text = this.value;
        const sentences = text.split('.').filter(s => s.trim().length > 0);
        const chars = text.length;

        // Actualiza contadores
        sentenceCount.textContent = `${sentences.length}/5 oraciones`;
        characterCount.textContent = `${chars}/500 caracteres`;

        // Validación de oraciones y caracteres
        if (sentences.length > 5) {
            textarea.classList.add('border-red-500');
            sentenceCount.classList.add('text-red-500');
            submitButton.disabled = true; // Desactivar el botón si hay más de 5 oraciones
        } else {
            textarea.classList.remove('border-red-500');
            sentenceCount.classList.remove('text-red-500');
        }

        // Desactivar o habilitar el botón de envío según la longitud de texto y número de oraciones
        if (text.trim().length === 0 || sentences.length === 0 || sentences.length > 5 || chars > 500) {
            submitButton.disabled = true; // Desactiva el botón si no hay oraciones o se exceden límites
        } else {
            submitButton.disabled = false; // Habilita el botón si las condiciones son correctas
        }
    });
});


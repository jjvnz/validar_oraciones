document.addEventListener('DOMContentLoaded', function() {
    const textarea = document.getElementById('oraciones');
    const sentenceCount = document.getElementById('sentenceCount');
    const characterCount = document.getElementById('characterCount');

    textarea.addEventListener('input', function() {
        const text = this.value;
        const sentences = text.split('.').filter(s => s.trim().length > 0);
        const chars = text.length;
        
        sentenceCount.textContent = `${sentences.length}/5 oraciones`;
        characterCount.textContent = `${chars}/500 caracteres`;
        
        if (sentences.length > 5) {
            textarea.classList.add('border-red-500');
            sentenceCount.classList.add('text-red-500');
        } else {
            textarea.classList.remove('border-red-500');
            sentenceCount.classList.remove('text-red-500');
        }
    });
});
document.addEventListener('DOMContentLoaded', () => {
    // Optimized Alert System
    class AlertSystem {
        static TYPES = {
            success: {
                classes: 'bg-green-100 dark:bg-green-900/50 text-green-800 dark:text-green-200',
                icon: '✓'
            },
            error: {
                classes: 'bg-red-100 dark:bg-red-900/50 text-red-800 dark:text-red-200',
                icon: '✕'
            },
            warning: {
                classes: 'bg-yellow-100 dark:bg-yellow-900/50 text-yellow-800 dark:text-yellow-200',
                icon: '⚠'
            },
            info: {
                classes: 'bg-blue-100 dark:bg-blue-900/50 text-blue-800 dark:text-blue-200',
                icon: 'ℹ'
            }
        };

        constructor() {
            this.container = document.getElementById('alert-container');
        }

        show({ message, type = 'info', duration = 5000 }) {
            const alertConfig = AlertSystem.TYPES[type] || AlertSystem.TYPES.info;
            const alert = this.createAlertElement(message, alertConfig);
            
            this.container.appendChild(alert);
            
            // Use requestAnimationFrame for smoother animation
            requestAnimationFrame(() => {
                alert.classList.add('slide-in');
            });

            if (duration > 0) {
                setTimeout(() => this.close(alert), duration);
            }
        }

        createAlertElement(message, { classes, icon }) {
            const alert = document.createElement('div');
            alert.className = `${classes} p-4 rounded-lg shadow-lg flex items-center justify-between min-w-[300px] opacity-0`;

            const content = document.createElement('div');
            content.className = 'flex items-center gap-3';

            const iconSpan = document.createElement('span');
            iconSpan.textContent = icon;
            content.appendChild(iconSpan);

            const text = document.createElement('p');
            text.className = 'text-sm font-medium';
            text.textContent = message;
            content.appendChild(text);

            alert.appendChild(content);

            const closeBtn = document.createElement('button');
            closeBtn.textContent = '✕';
            closeBtn.className = 'ml-4 text-sm hover:opacity-75';
            closeBtn.addEventListener('click', () => this.close(alert));
            alert.appendChild(closeBtn);

            return alert;
        }

        close(alert) {
            alert.classList.add('slide-out');
            alert.addEventListener('transitionend', () => alert.remove(), { once: true });
        }
    }

    // Theme Handling with Improved Accessibility
    const themeToggle = {
        toggle() {
            const html = document.documentElement;
            const isDark = html.classList.toggle('dark');
            
            localStorage.setItem('theme', isDark ? 'dark' : 'light');

            const themeIcon = document.getElementById('toggle-theme');
            themeIcon.textContent = isDark ? '🌞' : '🌙';
            themeIcon.setAttribute('aria-label', `Switch to ${isDark ? 'light' : 'dark'} mode`);

            this.alertSystem.show({
                message: `${isDark ? 'Dark' : 'Light'} mode activated`,
                type: 'info',
                duration: 2000
            });
        },

        init() {
            this.alertSystem = new AlertSystem();
            
            // Improved theme detection
            const prefersDarkMode = window.matchMedia('(prefers-color-scheme: dark)');
            const savedTheme = localStorage.getItem('theme');

            if (savedTheme === 'dark' || (savedTheme === null && prefersDarkMode.matches)) {
                document.documentElement.classList.add('dark');
                document.getElementById('toggle-theme').textContent = '🌞';
            }

            const themeToggleBtn = document.getElementById('toggle-theme');
            themeToggleBtn.addEventListener('click', () => this.toggle());
        }
    };

    // Form Validation with Enhanced Performance
    const formValidator = {
        MAX_LINES: 10,
        MAX_SENTENCES: 5,
        alertSystem: new AlertSystem(),

        init(form) {
            this.form = form;
            this.textarea = form.querySelector('textarea');
            this.submitButton = form.querySelector('button[type="submit"]');

            this.textarea.addEventListener('input', this.validate.bind(this));
            this.form.addEventListener('submit', this.handleSubmit.bind(this));
        },

        validate() {
            const text = this.textarea.value.trim();
            
            // Memoize validation to reduce unnecessary processing
            if (text === this.lastValidatedText) return;
            this.lastValidatedText = text;

            const lines = text.split('\n').filter(line => line.trim().length > 0);
            const sentences = (text.match(/[^.!?]+[.!?]+/g) || [])
                .filter(sentence => sentence.trim().length > 0);

            let errorMessage = '';
            const isValid = this.checkValidation(lines, sentences, errorMessage);

            this.updateSubmitButton(isValid, errorMessage);
        },

        checkValidation(lines, sentences, errorMessage) {
            if (lines.length > this.MAX_LINES) {
                errorMessage = `Maximum number of lines (${this.MAX_LINES}) exceeded`;
                return false;
            }

            if (sentences.length > this.MAX_SENTENCES) {
                errorMessage = `Maximum number of sentences (${this.MAX_SENTENCES}) exceeded`;
                return false;
            }

            return true;
        },

        updateSubmitButton(isValid, errorMessage) {
            this.submitButton.disabled = !isValid;
            this.submitButton.classList.toggle('opacity-50', !isValid);
            this.submitButton.classList.toggle('cursor-not-allowed', !isValid);

            if (!isValid) {
                this.alertSystem.show({
                    message: errorMessage,
                    type: "warning",
                    duration: 4000
                });
            }
        },

        handleSubmit(e) {
            const text = this.textarea.value.trim();

            if (!text) {
                e.preventDefault();
                this.alertSystem.show({
                    message: "Please enter some text to validate",
                    type: "error",
                    duration: 4000
                });
                return;
            }
        }
    };

    // Initialize components
    themeToggle.init();
    formValidator.init(document.getElementById('grammar-form'));

    // Simplified event listeners for buttons
    const buttonActions = {
        'settings-btn': () => alertSystem.show({ message: 'Settings panel opening soon...', type: 'info' }),
        'new-doc-btn': () => alertSystem.show({ message: 'Creating new document...', type: 'info' }),
        'templates-btn': () => alertSystem.show({ message: 'Templates will be available soon!', type: 'warning' }),
        'help-btn': () => alertSystem.show({ message: 'Need help? Contact support@validator.com', type: 'info' })
    };

    Object.entries(buttonActions).forEach(([id, action]) => {
        document.getElementById(id).addEventListener('click', action);
    });
});
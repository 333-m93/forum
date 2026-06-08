function applyTheme(theme) {
    document.body.classList.toggle('theme-dark', theme === 'dark');
    document.body.classList.toggle('theme-light', theme === 'light');
    const button = document.getElementById('theme-toggle');
    if (button) {
        button.textContent = theme === 'dark' ? 'Mode clair' : 'Mode sombre';
    }
}

function initTheme() {
    const stored = localStorage.getItem('forumTheme');
    const defaultTheme = stored || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    applyTheme(defaultTheme);

    const button = document.getElementById('theme-toggle');
    if (!button) return;

    button.addEventListener('click', function () {
        const active = document.body.classList.contains('theme-dark') ? 'dark' : 'light';
        const next = active === 'dark' ? 'light' : 'dark';
        localStorage.setItem('forumTheme', next);
        applyTheme(next);
    });
}

window.addEventListener('DOMContentLoaded', initTheme);

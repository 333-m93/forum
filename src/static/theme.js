function applyTheme(theme) {
    document.documentElement.classList.toggle('theme-dark', theme === 'dark');
    document.documentElement.classList.toggle('theme-light', theme === 'light');
    var button = document.getElementById('theme-toggle');
    if (button) {
        button.textContent = theme === 'dark' ? 'Mode clair' : 'Mode sombre';
    }
}

(function () {
    var stored = localStorage.getItem('forumTheme');
    var theme = stored || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    applyTheme(theme);
})();

document.addEventListener('DOMContentLoaded', function () {
    var button = document.getElementById('theme-toggle');
    if (!button) return;

    button.addEventListener('click', function () {
        var active = document.documentElement.classList.contains('theme-dark') ? 'dark' : 'light';
        var next = active === 'dark' ? 'light' : 'dark';
        localStorage.setItem('forumTheme', next);
        applyTheme(next);
    });
});

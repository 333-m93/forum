class ForumChat {
  constructor() {
    this.currentCategory = null;
    this.messages = [];
    this.pollInterval = null;
    this.init();
  }

  init() {
    // Capturer les clics sur les liens de catégorie
    document.addEventListener('click', (e) => {
      const link = e.target.closest('a[data-cat]');
      if (!link) return;
      
      e.preventDefault();
      const catName = link.getAttribute('data-cat');
      console.log('Catégorie cliquée:', catName);
      this.loadChat(catName);
    });

    // Gérer l'envoi du formulaire du chat
    document.addEventListener('submit', (e) => {
      if (e.target.classList.contains('chat-form')) {
        e.preventDefault();
        this.sendMessage(e.target);
      }
    });
  }

  loadChat(categoryName) {
    console.log('loadChat appelé avec:', categoryName);
    this.currentCategory = categoryName;
    this.messages = [];
    
    const pane = document.getElementById('floating-chat');
    pane.innerHTML = `
      <div class="chat-header">
        <h2 style="margin:0;color:var(--text);">${categoryName}</h2>
      </div>
      <div class="chat-messages" id="chat-messages" style="height:400px;overflow-y:auto;margin:14px 0;"></div>
      <div class="chat-footer">
        <form class="chat-form">
          <input class="chat-input" placeholder="Écrire un message..." type="text" required>
          <button class="btn primary chat-send" type="submit">Envoyer</button>
        </form>
      </div>
    `;
    pane.setAttribute('aria-hidden', 'false');

    this.fetchMessages();
    
    // Polling toutes les 2 secondes
    if (this.pollInterval) clearInterval(this.pollInterval);
    this.pollInterval = setInterval(() => this.fetchMessages(), 2000);
  }

  fetchMessages() {
    if (!this.currentCategory) return;

    fetch(`/api/messages?category=${encodeURIComponent(this.currentCategory)}`)
      .then(res => res.json())
      .then(data => {
        if (data.success && data.data) {
          this.renderMessages(data.data);
        }
      })
      .catch(err => console.error('Erreur fetch:', err));
  }

  renderMessages(messages) {
    const container = document.getElementById('chat-messages');
    if (!container) return;

    if (JSON.stringify(this.messages) === JSON.stringify(messages)) {
      return;
    }

    this.messages = messages;
    
    const html = messages.map(msg => `
      <div class="chat-message">
        <strong style="color:var(--muted);">${msg.username}</strong>
        <p style="margin:6px 0 0;color:var(--text);">${msg.content}</p>
        <small style="color:var(--muted);font-size:0.85rem;">${new Date(msg.created_at).toLocaleTimeString()}</small>
      </div>
    `).join('');

    container.innerHTML = html || '<p style="color:var(--muted);">Pas de messages encore</p>';
    container.scrollTop = container.scrollHeight;
  }

  sendMessage(form) {
    if (!this.currentCategory) {
      alert('Sélectionnez un chat d\'abord');
      return;
    }

    const input = form.querySelector('input[type="text"]');
    const content = input.value.trim();

    if (!content) {
      alert('Écrivez quelque chose');
      return;
    }

    console.log('Envoi message:', {category: this.currentCategory, content: content});

    const formData = new FormData();
    formData.append('category', this.currentCategory);
    formData.append('content', content);

    fetch('/api/messages', {
      method: 'POST',
      body: formData
    })
      .then(res => res.json())
      .then(data => {
        console.log('Réponse serveur:', data);
        if (data.success) {
          input.value = '';
          this.fetchMessages();
        } else {
          alert('Erreur: ' + (data.message || 'Erreur inconnue'));
          if (data.message && data.message.includes('authentifié')) {
            window.location.href = '/login';
          }
        }
      })
      .catch(err => {
        console.error('Erreur:', err);
        alert('Erreur lors de l\'envoi');
      });
  }
}

document.addEventListener('DOMContentLoaded', () => {
  console.log('Chat.js chargé');
  new ForumChat();
});

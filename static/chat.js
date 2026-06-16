class ForumChat {
  constructor() {
    this.currentCategory = null;
    this.messages = [];
    this.pollInterval = null;
    this.init();
  }

  init() {
    console.log("Chat.js chargé");

    // =====================
    // CLICK CATEGORY
    // =====================
    document.addEventListener('click', (e) => {
      const link = e.target.closest('a[data-cat]');
      if (!link) return;

      e.preventDefault();

      const catName = link.getAttribute('data-cat');

      if (!catName) {
        console.error("❌ data-cat manquant");
        return;
      }

      console.log('📂 Catégorie cliquée:', catName);
      this.loadChat(catName);
    });

    // =====================
    // SEND MESSAGE
    // =====================
    document.addEventListener('submit', (e) => {
      if (!e.target.classList.contains('chat-form')) return;

      e.preventDefault();
      this.sendMessage(e.target);
    });
  }

  // =====================
  // LOAD CHAT UI
  // =====================
  loadChat(categoryName) {
    console.log('🚀 loadChat:', categoryName);

    this.currentCategory = categoryName;
    this.messages = [];

    const pane = document.getElementById('floating-chat');

    if (!pane) {
      console.error("❌ floating-chat introuvable");
      return;
    }

    pane.innerHTML = `
      <div class="chat-header">
        <h2 style="margin:0;color:var(--text);">${categoryName}</h2>
      </div>

      <div class="chat-messages" id="chat-messages"
        style="height:400px;overflow-y:auto;margin:14px 0;">
      </div>

      <div class="chat-footer">
        <form class="chat-form">
          <input class="chat-input"
                 placeholder="Écrire un message..."
                 type="text"
                 required>
          <button class="btn primary chat-send" type="submit">
            Envoyer
          </button>
        </form>
      </div>
    `;

    pane.setAttribute('aria-hidden', 'false');

    this.fetchMessages();

    // reset polling
    if (this.pollInterval) clearInterval(this.pollInterval);
    this.pollInterval = setInterval(() => this.fetchMessages(), 3000);
  }

  // =====================
  // FETCH MESSAGES (GET)
  // =====================
  fetchMessages() {
    if (!this.currentCategory) return;

    const url = `/api/messages?category=${encodeURIComponent(this.currentCategory)}`;

    console.log("📡 GET:", url);

    fetch(url)
      .then(res => res.json())
      .then(data => {
        if (!data.success) {
          console.warn("⚠️ API error:", data.message);
          return;
        }

        this.renderMessages(data.data || []);
      })
      .catch(err => console.error('❌ fetch error:', err));
  }

  // =====================
  // RENDER
  // =====================
  renderMessages(messages) {
    const container = document.getElementById('chat-messages');
    if (!container) return;

    if (JSON.stringify(this.messages) === JSON.stringify(messages)) return;

    this.messages = messages;

    container.innerHTML = messages.length
      ? messages.map(msg => `
          <div class="chat-message">
            <strong style="color:var(--muted);">
              ${msg.username || 'user'}
            </strong>

            <p style="margin:6px 0 0;color:var(--text);">
              ${msg.content}
            </p>

            <small style="color:var(--muted);font-size:0.85rem;">
              ${new Date(msg.created_at).toLocaleTimeString()}
            </small>
          </div>
        `).join('')
      : `<p style="color:var(--muted);">Pas de messages encore</p>`;

    container.scrollTop = container.scrollHeight;
  }

  // =====================
// SEND MESSAGE (POST)
// =====================
sendMessage(form) {
  if (!this.currentCategory) {
    alert('Sélectionnez une catégorie');
    return;
  }

  const input = form.querySelector('input[type="text"]');
  const content = input.value.trim();

  if (!content) {
    alert('Message vide');
    return;
  }

  console.log("📤 POST message:", {
    category: this.currentCategory,
    content
  });

  fetch('/api/messages', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      category: this.currentCategory,
      content: content
    })
  })
    .then(res => res.json())
    .then(data => {
      console.log("📨 response:", data);

      if (!data.success) {
        alert(data.message || "Erreur serveur");
        return;
      }

      input.value = '';
      this.fetchMessages();
    })
    .catch(err => {
      console.error("❌ POST error:", err);
      alert("Erreur réseau");
    });
}
}

document.addEventListener('DOMContentLoaded', () => {
  new ForumChat();
});
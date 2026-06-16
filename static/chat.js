class ForumChat {
  constructor() {
    this.currentCategory = null;
    this.messages = [];
    this.pollInterval = null;
    this.init();
  }

  init() {
    document.addEventListener("click", (e) => {
      const link = e.target.closest(".cat-link");
      if (!link) return;

      e.preventDefault();

      const name = link.dataset.catName;
      if (!name) return;

      this.openChat(name);
    });

    document.addEventListener("submit", (e) => {
      if (!e.target.classList.contains("chat-form")) return;

      e.preventDefault();
      this.sendMessage(e.target);
    });
  }

  openChat(categoryName) {
    this.currentCategory = categoryName;
    this.messages = [];

    const pane = document.getElementById("floating-chat");

    pane.innerHTML = `
      <div class="chat-header">
        <h2>${categoryName}</h2>
      </div>

      <div id="chat-messages" class="chat-messages"
        style="height:400px;overflow-y:auto;margin:14px 0;"></div>

      <div class="chat-footer">
        <form class="chat-form">
          <input type="text" class="chat-input" placeholder="Écrire un message..." required />
          <button type="submit" class="chat-send btn primary">Envoyer</button>
        </form>
      </div>
    `;

    pane.setAttribute("aria-hidden", "false");

    this.fetchMessages();

    if (this.pollInterval) clearInterval(this.pollInterval);
    this.pollInterval = setInterval(() => this.fetchMessages(), 3000);
  }

  fetchMessages() {
    if (!this.currentCategory) return;

    fetch(`/api/messages?category=${encodeURIComponent(this.currentCategory)}`)
      .then(r => r.json())
      .then(data => {
        if (!data.success) return;
        this.renderMessages(data.data || []);
      })
      .catch(err => console.error("GET error:", err));
  }

  renderMessages(messages) {
    const container = document.getElementById("chat-messages");
    if (!container) return;

    this.messages = messages;

    container.innerHTML = messages.length
      ? messages.map(m => `
          <div class="chat-message">
            <strong>${m.username}</strong>
            <p>${m.content}</p>
            <small>${new Date(m.created_at).toLocaleTimeString()}</small>
          </div>
        `).join("")
      : "<p>Aucun message</p>";

    container.scrollTop = container.scrollHeight;
  }

  sendMessage(form) {
    if (!this.currentCategory) return;

    const input = form.querySelector("input");
    const content = input.value.trim();

    if (!content) return;

    fetch("/api/messages", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        category: this.currentCategory,
        content: content
      })
    })
      .then(r => r.json())
      .then(data => {
        if (!data.success) {
          alert(data.message || "Erreur");
          return;
        }
        input.value = "";
        this.fetchMessages();
      })
      .catch(err => console.error("POST error:", err));
  }
}

document.addEventListener("DOMContentLoaded", () => {
  new ForumChat();
});

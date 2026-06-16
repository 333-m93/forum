class ForumChat {
  constructor() {
    this.currentCategory = null;
    this.messages = [];
    this.pollInterval = null;
    this.init();
  }

  init() {
    console.log("Chat.js chargé");

    // CLICK CATEGORY
    document.addEventListener("click", (e) => {
      const link = e.target.closest("a[data-cat-id]");
      if (!link) return;

      e.preventDefault();

      const id = link.getAttribute("data-cat-id");
      const name = link.getAttribute("data-cat-name");

      if (!id) {
        console.error("Catégorie ID manquant");
        return;
      }

      this.loadChat({
        id: parseInt(id),
        name: name
      });
    });

    // SEND MESSAGE
    document.addEventListener("submit", (e) => {
      if (!e.target.classList.contains("chat-form")) return;
      e.preventDefault();
      this.sendMessage(e.target);
    });
  }

  loadChat(category) {
    this.currentCategory = category;
    this.messages = [];

    const pane = document.getElementById("floating-chat");

    pane.innerHTML = `
      <div class="chat-header">
        <h2>${category.name}</h2>
      </div>

      <div class="chat-messages" id="chat-messages"
        style="height:400px;overflow-y:auto;margin:14px 0;">
      </div>

      <div class="chat-footer">
        <form class="chat-form">
          <input type="text" class="chat-input" placeholder="Écrire un message..." required />
          <button type="submit">Envoyer</button>
        </form>
      </div>
    `;

    this.fetchMessages();

    if (this.pollInterval) clearInterval(this.pollInterval);
    this.pollInterval = setInterval(() => this.fetchMessages(), 3000);
  }

  fetchMessages() {
    if (!this.currentCategory) return;

    fetch(`/api/messages?category_id=${this.currentCategory.id}`)
      .then(r => r.json())
      .then(data => {
        if (!data.success) return;
        this.renderMessages(data.data || []);
      })
      .catch(err => console.error(err));
  }

  renderMessages(messages) {
    const container = document.getElementById("chat-messages");
    if (!container) return;

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
    const input = form.querySelector("input");
    const content = input.value.trim();

    if (!content || !this.currentCategory) return;

    fetch("/api/messages", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        category_id: this.currentCategory.id,
        content: content
      })
    })
      .then(r => r.json())
      .then(data => {
        if (!data.success) {
          alert(data.message);
          return;
        }

        input.value = "";
        this.fetchMessages();
      })
      .catch(err => console.error(err));
  }
}

document.addEventListener("DOMContentLoaded", () => {
  new ForumChat();
});
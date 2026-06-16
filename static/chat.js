class ForumChat {
  constructor() {
    this.currentCategory = null;
    this.messages = [];
    this.pollInterval = null;
    this.emojis = ["👍", "❤️", "😂", "🔥", "💪", "👊"];
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

    document.addEventListener("click", (e) => {
      const btn = e.target.closest(".reaction-btn");
      if (!btn) return;
      e.preventDefault();
      this.toggleReaction(btn);
    });

    document.addEventListener("click", (e) => {
      const pick = e.target.closest(".emoji-pick");
      if (!pick) return;
      e.preventDefault();
      const picker = pick.closest(".emoji-picker");
      const msgId = parseInt(picker.dataset.msgId);
      const emoji = pick.dataset.emoji;
      this.sendReaction(msgId, emoji);
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

  avatarHTML(m) {
    if (m.avatar_url) {
      return `<img src="${m.avatar_url}" class="chat-avatar" alt="">`;
    }
    const initial = m.username ? m.username[0].toUpperCase() : "?";
    return `<span class="chat-avatar-initial">${initial}</span>`;
  }

  reactionsHTML(m) {
    let html = '<div class="reactions-row">';
    const groups = m.reactions || [];

    for (const g of groups) {
      const mine = g.user_id > 0 ? " reaction-mine" : "";
      html += `<button class="reaction-btn${mine}" data-msg-id="${m.id}" data-emoji="${g.emoji}">
        <span class="reaction-emoji">${g.emoji}</span>
        <span class="reaction-count">${g.count}</span>
      </button>`;
    }

    html += `<div class="emoji-picker emoji-picker-open" data-msg-id="${m.id}">`;
    for (const e of this.emojis) {
      html += `<button class="emoji-pick" data-emoji="${e}">${e}</button>`;
    }
    html += "</div></div>";

    return html;
  }

  renderMessages(messages) {
    const container = document.getElementById("chat-messages");
    if (!container) return;

    this.messages = messages;

    container.innerHTML = messages.length
      ? messages.map(m => `
          <div class="chat-message">
            <div class="chat-message-header">
              ${this.avatarHTML(m)}
              <strong class="chat-username">${m.username}</strong>
              <small class="chat-time">${new Date(m.created_at).toLocaleTimeString()}</small>
            </div>
            <p class="chat-content">${m.content}</p>
            ${this.reactionsHTML(m)}
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

  toggleReaction(btn) {
    const msgId = parseInt(btn.dataset.msgId);
    const emoji = btn.dataset.emoji;
    const isMine = btn.classList.contains("reaction-mine");

    if (isMine) {
      fetch(`/api/reactions?message_id=${msgId}&emoji=${encodeURIComponent(emoji)}`, {
        method: "DELETE"
      })
        .then(r => r.json())
        .then(() => this.fetchMessages())
        .catch(err => console.error("DELETE reaction error:", err));
    } else {
      this.sendReaction(msgId, emoji);
    }
  }

  sendReaction(messageID, emoji) {
    fetch("/api/reactions", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message_id: messageID, emoji: emoji })
    })
      .then(r => r.json())
      .then(() => this.fetchMessages())
      .catch(err => console.error("POST reaction error:", err));
  }
}

document.addEventListener("DOMContentLoaded", () => {
  new ForumChat();
});

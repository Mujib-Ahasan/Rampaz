const form = document.getElementById("chat-form");
const input = document.getElementById("message-input");
const chatBox = document.getElementById("chat-box");

const STORAGE_KEY = "rampaz-ai-chat-history";
const SESSION_KEY = "rampaz-session-id";

function loadMessages() {
const saved = sessionStorage.getItem(STORAGE_KEY);
  if (!saved) return [];

  try {
    return JSON.parse(saved);
  } catch {
    return [];
  }
}

document.getElementById("clear-chat").addEventListener("click", () => {
  sessionStorage.removeItem(STORAGE_KEY); 
  sessionStorage.removeItem(SESSION_KEY); 
  renderMessages(); 
});

function saveMessages(messages) {
  sessionStorage.setItem(STORAGE_KEY, JSON.stringify(messages));
}

function formatMessageTime(timestamp) {
  return new Date(timestamp).toLocaleTimeString(undefined, {
    hour: "2-digit",
    minute: "2-digit",
  });
}

function renderMessages() {
  chatBox.innerHTML = "";

  const messages = loadMessages();
  messages.forEach((msg) => {
    addMessageToUI(msg.role, msg.text, msg.timestamp || new Date().toISOString());
  });
}

function addMessageToUI(role, text, timestamp) {
  const msg = document.createElement("div");
  msg.className = `message ${role}`;

  const content = document.createElement("div");
  content.className = "message-content";
  content.textContent = text;

  const time = document.createElement("div");
  time.className = "message-time";
  time.textContent = formatMessageTime(timestamp);

  msg.appendChild(content);
  msg.appendChild(time);

  chatBox.appendChild(msg);
  chatBox.scrollTop = chatBox.scrollHeight;
}

function addMessage(role, text) {
  const messages = loadMessages();
  const timestamp = new Date().toISOString();

  messages.push({ role, text, timestamp });
  saveMessages(messages);

  addMessageToUI(role, text, timestamp);
}


form.addEventListener("submit", async (e) => {
  e.preventDefault();

  const message = input.value.trim();
  if (!message) return;

  addMessage("user", message);
  input.value = "";

  try {
    const res = await fetch("/api/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ message }),
    });

    if (!res.ok) {
      throw new Error("request failed");
    }

   const data = await res.json();

const oldSession = sessionStorage.getItem(SESSION_KEY);

if (oldSession && oldSession !== data.sessionId) {
  sessionStorage.removeItem(STORAGE_KEY);
  chatBox.innerHTML = "";

  // keep the current user message after clearing old session
  addMessage("user", message);
}

sessionStorage.setItem(SESSION_KEY, data.sessionId);
addMessage("bot", data.answer);
  } catch (err) {
    addMessage("bot", "Something went wrong.");
  }
});

async function syncServerSession() {
  try {
    const res = await fetch("/api/session");
    const data = await res.json();

    const oldSession = sessionStorage.getItem(SESSION_KEY);

    if (oldSession && oldSession !== data.sessionId) {
      sessionStorage.removeItem(STORAGE_KEY);
    }

    sessionStorage.setItem(SESSION_KEY, data.sessionId);
  } catch {
    sessionStorage.removeItem(STORAGE_KEY);
    sessionStorage.removeItem(SESSION_KEY);
  }
}

syncServerSession().then(renderMessages);
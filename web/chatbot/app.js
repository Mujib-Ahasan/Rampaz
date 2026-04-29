const form = document.getElementById("chat-form");
const input = document.getElementById("message-input");
const chatBox = document.getElementById("chat-box");
const clearBtn = document.getElementById("clear-chat");

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
    addMessageToUI(
      msg.role,
      msg.text,
      msg.timestamp || new Date().toISOString()
    );
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

clearBtn.addEventListener("click", () => {
  sessionStorage.removeItem(STORAGE_KEY);
  sessionStorage.removeItem(SESSION_KEY);
  chatBox.innerHTML = "";
});


// form.addEventListener("submit", async (e) => {
//   e.preventDefault();

//   const message = input.value.trim();
//   if (!message) return;

//   addMessage("user", message);
//   input.value = "";

//   try {
//     const res = await fetch("/api/chat", {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify({ message }),
//     });

//     if (!res.ok) {
//       throw new Error("request failed");
//     }

//     const data = await res.json();

//     const oldSession = sessionStorage.getItem(SESSION_KEY);

//     if (data.sessionId && oldSession && oldSession !== data.sessionId) {
//       sessionStorage.removeItem(STORAGE_KEY);
//       chatBox.innerHTML = "";
//       addMessage("user", message);
//     }

//     if (data.sessionId) {
//       sessionStorage.setItem(SESSION_KEY, data.sessionId);
//     }

//     addMessage("bot", data.answer || "No response received.");
//   } catch {
//     addMessage("bot", "Something went wrong.");
//   }
// });

form.addEventListener("submit", async (e) => {
  e.preventDefault();

  const message = input.value.trim();
  if (!message) return;

  addMessage("user", message);
  input.value = "";

  showLoadingMessage();

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

    removeLoadingMessage();

    const oldSession = sessionStorage.getItem(SESSION_KEY);

    if (data.sessionId && oldSession && oldSession !== data.sessionId) {
      sessionStorage.removeItem(STORAGE_KEY);
      chatBox.innerHTML = "";
      addMessage("user", message);
    }

    if (data.sessionId) {
      sessionStorage.setItem(SESSION_KEY, data.sessionId);
    }

    addMessage("bot", data.answer || "No response received.");
  } catch {
    removeLoadingMessage();
    addMessage("bot", "Something went wrong.");
  }
});

async function syncServerSession() {
  try {
    const res = await fetch("/api/session");

    if (!res.ok) throw new Error();

    const data = await res.json();
    const oldSession = sessionStorage.getItem(SESSION_KEY);

    if (data.sessionId) {
      if (oldSession && oldSession !== data.sessionId) {
        sessionStorage.removeItem(STORAGE_KEY);
      }

      sessionStorage.setItem(SESSION_KEY, data.sessionId);
    }
  } catch {
    // don't aggressively wipe everything on network error
    console.warn("session sync failed");
  }
}

function showLoadingMessage() {
  const msg = document.createElement("div");
  msg.className = "message bot loading";
  msg.id = "loading-message";

  const dot1 = document.createElement("span");
  dot1.className = "typing-dot";

  const dot2 = document.createElement("span");
  dot2.className = "typing-dot";

  const dot3 = document.createElement("span");
  dot3.className = "typing-dot";

  msg.appendChild(dot1);
  msg.appendChild(dot2);
  msg.appendChild(dot3);

  chatBox.appendChild(msg);
  chatBox.scrollTop = chatBox.scrollHeight;
}

function removeLoadingMessage() {
  const loading = document.getElementById("loading-message");
  if (loading) {
    loading.remove();
  }
}

syncServerSession().then(renderMessages);
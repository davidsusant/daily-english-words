// UI module to renders word cards and manage DOM state

/**
 * Show the loading indicator and hide other sections.
 */
export function showLoading() {
  document.getElementById("loading").hidden = false;
  document.getElementById("error").hidden = true;
  document.getElementById("words-container").innerHTML = "";
}

/**
 * Hide the loading indicator
 */
export function hideLoading() {
  document.getElementById("loading").hidden = true;
}

/**
 * Display an error message with a retry button
 * @param {string} message The error text to show
 * @param {Function} onRetry Callback when "Try Again" is clicked
 */
export function showError(message, onRetry) {
  hideLoading();

  const errorEl = document.getElementById("error");
  const messageEl = document.getElementById("error-message");
  const retryBtn = document.getElementById("retry-btn");

  messageEl.textContent = message;
  errorEl.hidden = false;

  // Remove old listener, then add new one
  const newBtn = retryBtn.cloneNode(true);
  retryBtn.parentNode.replaceChild(newBtn, retryBtn);
  newBtn.id = "retry-btn";
  newBtn.addEventListener("click", onRetry);
}

/**
 * Render an array of word objects into the container
 * @param {Array} words Word objects from the API
 * @param {Object} [options] Rendering options
 * @param {boolean} [options.isNew=false] If true, marks cards as newly generated
 */
export function renderWords(words, options = {}) {
  hideLoading();
  document.getElementById("error").hidden = true;

  const container = document.getElementById("words-container");
  container.innerHTML = "";

  if (words.length === 0) {
    container.innerHTML =
      '<p style="text-align:center; color:#8b8d98;">No words available. Try generating some!</p>';
    return;
  }

  words.forEach((word) => {
    const card = createWordCard(word, options.isNew);
    container.appendChild(card);
  });
}

export function setGenerateLoading(isLoading) {
  const btn = document.getElementById("generate-btn");
  const selects = document.querySelectorAll(".generate-panel--select");

  if (isLoading) {
    btn.disabled = true;
    btn.classList.add("is-loading");
    btn.querySelector(".btn__icon").textContent = "\u21BB"; // ↻ spinner arrow
    btn.childNodes[btn.childNodes.length - 1].textContent = "Generating...";
    selects.forEach((s) => (s.disabled = true));
  } else {
    btn.disabled = false;
    btn.classList.remove("is-loading");
    btn.querySelector(".btn__icon").textContent = "\u26A1"; // ⚡
    btn.childNodes[btn.childNodes.length - 1].textContent =
      "\n                    Generate Words\n";
    selects.forEach((s) => (s.disabled = false));
  }
}

/**
 * Show a status message in the general panel
 * @param {string} message Text to display
 * @param {"success"|"error"|"info"} type Visual style
 */
export function showGenerateStatus(message, type) {
  const statusEl = document.getElementById("generate-status");

  statusEl.hidden = false;
  statusEl.textContent = message;
  statusEl.className = `generate-panel__status generate-panel__status--${type}`;

  // Auto-hide after 6 seconds
  clearTimeout(statusEl._hideTimer);
  statusEl._hideTimer = setTimeout(() => {
    statusEl.hidden = true;
  }, 6000);
}

/**
 * Hide the genrate status message
 */
export function hideGenerateStatus() {
  const statusEl = document.getElementById("generate-status");
  statusEl.hidden = true;
  clearTimeout(statusEl._hideTimer);
}

/**
 * Create a single word card DOM element
 * @param {Object} word A word object from the API
 * @param {boolean} isNew whether to add the "new" highlight
 * @returns {HTMLElement}
 */
function createWordCard(word, isNew) {
  const card = document.createElement("article");
  card.className = "word-card" + (isNew ? " word-card--new" : "");

  const newBadge = isNew ? '<span class="word-card__new-badge">new</span>' : "";

  card.innerHTML = `
        <div class="word-card__header">
            <span class="word-card__word">${escapeHTML(word.word)}</span>
            <span class="word-card__pos">${escapeHTML(word.part_of_speech)}</span>
            ${newBadge}
            <span class="word-card__badge word-card__badge--${word.difficulty}">${escapeHTML(word.difficulty)}</span>
        </div>
        <p class="word-card__definition">${escapeHTML(word.definition)}</p>
        <p class="word-card__example">"${escapeHTML(word.example)}"</p>
    `;

  return card;
}

/**
 * Escape HTML special characters to prevent XSS
 * @param {string} str
 * @returns {string}
 */
function escapeHTML(str) {
  const div = document.createElement("div");
  div.textContent = str;
  return div.innerHTML;
}

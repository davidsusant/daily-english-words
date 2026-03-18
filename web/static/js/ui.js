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

  // Remove old listener, add new one
  const newBtn = retryBtn.cloneNode(true);
  retryBtn.parentNode.replaceChild(newBtn, retryBtn);
  newBtn.id = "retry-btn";
  newBtn.addEventListener("click", onRetry);
}

/**
 * Render an array of word objects into the container
 * @param {Array} words Word objects from the API
 */
export function renderWords(words) {
  hideLoading();
  document.getElementById("error").hidden = true;

  const container = document.getElementById("words-container");
  container.innerHTML = "";

  if (words.length === 0) {
    container.innerHTML =
      "<p>No words available. Try seeding the database.</p>";
    return;
  }

  words.forEach((word) => {
    const card = createWordCard(word);
    container.appendChild(card);
  });
}

/**
 * Create a single word card DOM element
 * @param {Object} word A word object from the API
 * @returns {HTMLElement}
 */
function createWordCard(word) {
  const card = document.createElement("article");

  card.innerHTML = `
        <div>
            <span>${escapeHTML(word.word)}</span>
            <span>${escapeHTML(word.part_of_speech)}</span>
            <span>${escapeHTML(word.difficulty)}</span>
        </div>
        <p>${escapeHTML(word.definition)}</p>
        <p>${escapeHTML(word.example)}</p>
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

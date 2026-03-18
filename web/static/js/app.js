// Main application entry: imports API & UI modules and wires them together

import { fetchTodayWords, fetchRandomWords } from "./api.js";
import { renderWords, showError, showLoading } from "./ui.js";

let activeTab = "today";

// Tab switching
function initTabs() {
  const buttons = document.querySelectorAll(".tabs__btn");

  buttons.forEach((btn) => {
    btn.addEventListener("click", () => {
      // Update active button style
      buttons.forEach((b) => b.classList.remove("tabs__btn--active"));
      btn.classList.add("tabs__btn--active");

      // Load the corresponding words
      activeTab = btn.dataset.tab;
      loadWords();
    });
  });
}

// Data loading
async function loadWords() {
  showLoading();

  try {
    let data;

    if (activeTab === "today") {
      data = await fetchTodayWords();
    } else {
      data = await fetchRandomWords();
    }

    renderWords(data.words);
  } catch (err) {
    console.error("Failed to load words:", err);
    showError(
      `Could not load words: ${err.message}`,
      () => loadWords(), // retry callback
    );
  }
}

// Boot
document.addEventListener("DOMContentLoaded", () => {
  initTabs();
  loadWords();
});

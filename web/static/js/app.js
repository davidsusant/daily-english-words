// Main application entry: imports API & UI modules and wires them together

import { fetchTodayWords, fetchRandomWords, generateWords } from "./api.js";
import {
  hideGenerateStatus,
  renderWords,
  setGenerateLoading,
  showError,
  showGenerateStatus,
  showLoading,
} from "./ui.js";

let activeTab = "today";

// Tab switching
function initTabs() {
  const buttons = document.querySelectorAll(".tabs__btn");

  buttons.forEach((btn) => {
    btn.addEventListener("click", () => {
      buttons.forEach((b) => b.classList.remove("tabs__btn--active"));
      btn.classList.add("tabs__btn--active");

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
    showError(`Could not load words: ${err.message}`, () => loadWords());
  }
}

// Generate words
function initGenerate() {
  const generateBtn = document.getElementById("generate-btn");
  const difficultySelect = document.getElementById("difficulty-select");
  const countSelect = document.getElementById("count-select");

  generateBtn.addEventListener("click", async () => {
    const count = parseInt(countSelect.value, 10);
    const difficulty = difficultySelect.value;

    // Show loading state on the button
    setGenerateLoading(true);
    hideGenerateStatus();

    try {
      const data = await generateWords(count, difficulty);

      // Show success status
      const msg = buildStatusMessage(data);
      showGenerateStatus(msg, "success");

      // If words were inserted, display them as "new" cards
      if (data.words && data.words.length > 0) {
        renderWords(data.words, { isNew: true });
      } else if (data.duplicates > 0) {
        // All duplicates - show info and reload current tab
        showGenerateStatus(
          "All generated words already exist. Try again for new ones!",
          "info",
        );
        loadWords();
      }
    } catch (err) {
      console.error("Failed to generate words:", err);
      showGenerateStatus(`Generation failed: ${err.message}`, "error");
    } finally {
      setGenerateLoading(false);
    }
  });
}

// Build a user-friendly status message from the API response
function buildStatusMessage(data) {
  const parts = [];
  parts.push(
    `${data.inserted} new word${data.inserted !== 1 ? "s" : ""} added`,
  );

  if (data.dupplicates > 0) {
    parts.push(
      `${data.dupplicates} duplicate${data.dupplicates !== 1 ? "s" : ""} skipped`,
    );
  }

  return parts.join(" . ");
}

// Boot
document.addEventListener("DOMContentLoaded", () => {
  initTabs();
  initGenerate();
  loadWords();
});

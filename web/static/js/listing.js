document.querySelectorAll(".listing").forEach(row => {
	row.addEventListener("click", () => {
	  window.location.href = row.dataset.href;
	});

	row.addEventListener("keydown", event => {
	  if (event.key === "Enter" || event.key === " ") {
		event.preventDefault();
		row.click();
	  }
	});

	row.tabIndex = 0;
});

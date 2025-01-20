function setDarkMode() {
	let toggleDarkModeButton = document.getElementById("toggle-dark-mode");
	let mode = localStorage.getItem("mode") || localStorage.setItem("mode", "light");
	if (mode == "dark") {
		document.body.classList.add("dark");
		toggleDarkModeButton.textContent = "Light Mode";
	} else {
		document.body.classList.remove("dark");
		toggleDarkModeButton.textContent = "Dark Mode";
	}
}

function toggleDarkMode() {
	let mode = localStorage.getItem("mode") || localStorage.setItem("mode", "light");
	if (mode == "dark") {
		localStorage.setItem("mode", "light");
	} else if (mode == "light"){
		localStorage.setItem("mode", "dark");
	} else {
		localStorage.setItem("mode", "light");
	}
}

function init() {
	setDarkMode();
	const toggleDarkModeButton = document.getElementById("toggle-dark-mode");
	toggleDarkModeButton.addEventListener("click", () => {
		toggleDarkMode();
		setDarkMode();
	});
}

init();

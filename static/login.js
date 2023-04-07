// Login
const loginForm = document.getElementById("loginForm")

loginForm.addEventListener("submit", (event) => {
  event.preventDefault()

  const username = loginForm.elements.username.value
  const password = loginForm.elements.password.value

  fetch('/auth', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  })
    .then(response => response.json())
    .then(data => {
      if (data.result != "done") {
        console.log(data)
        return
      }
      window.location.href = 'game.html';
    })
    .catch(error => {
      console.error(error)
    })
})
